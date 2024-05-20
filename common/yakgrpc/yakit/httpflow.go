package yakit

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/mutate"
	"github.com/yaklang/yaklang/common/schema"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/bizhelper"
	"github.com/yaklang/yaklang/common/utils/lowhttp"
	"github.com/yaklang/yaklang/common/utils/lowhttp/httpctx"
	"github.com/yaklang/yaklang/common/yakgrpc/model"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
)

func init() {
	RegisterPostInitDatabaseFunction(func() error {
		lowhttp.RegisterSaveHTTPFlowHandler(func(https bool, req []byte, rsp []byte, url string, remoteAddr string, reqSource string, runtimeId string, fromPlugin string, hiddenIndex string, payloads []string) {
			if rsp == nil || len(rsp) == 0 {
				return
			}
			db := consts.GetGormProjectDatabase()
			flow, err := CreateHTTPFlowFromHTTPWithBodySavedFromRaw(https, req, rsp, "scan", url, remoteAddr)
			if err != nil {
				log.Errorf("create httpflow from lowhttp failed: %s", err)
				return
			}
			if reqSource != "" {
				flow.AddTag(reqSource)
			} else {
				if reqSource := os.Getenv(consts.YAKIT_PLUGIN_ID); reqSource != "" {
					flow.AddTag(reqSource)
				}
			}
			switch ret := strings.ToLower(reqSource); ret {
			case "mitm":
				flow.SourceType = "mitm"
			case "basic-crawler", "crawler", "crawlerx":
				flow.SourceType = "basic-crawler"
			case "scan", "port-scan", "plugin":
				flow.SourceType = "scan"

			}
			flow.FromPlugin = fromPlugin
			flow.RuntimeId = runtimeId
			flow.HiddenIndex = hiddenIndex
			flow.Payload = strings.Join(payloads, ",")
			err = InsertHTTPFlow(db, flow)
			if err != nil {
				log.Errorf("save httpflow failed: %s", err)
			}
		})
		return nil
	})
}

type TagAndStatusCode struct {
	Value string
	Count int
}

type CreateHTTPFlowConfig struct {
	isHttps     bool
	reqRaw      []byte
	rspRaw      []byte
	fixRspRaw   []byte // 如果设置了，则不会再修复rspRaw
	source      string
	url         string
	remoteAddr  string
	reqIns      *http.Request // 如果设置了，则不会再解析reqRaw
	hiddenIndex string
}

type CreateHTTPFlowOptions func(c *CreateHTTPFlowConfig)

func CreateHTTPFlowWithHTTPS(isHttps bool) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.isHttps = isHttps
	}
}

func CreateHTTPFlowWithRequestRaw(reqRaw []byte) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.reqRaw = reqRaw
	}
}

func CreateHTTPFlowWithResponseRaw(rspRaw []byte) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.rspRaw = rspRaw
	}
}

// 如果传入了fixRspRaw，则不会再修复
func CreateHTTPFlowWithFixResponseRaw(fixRspRaw []byte) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.fixRspRaw = fixRspRaw
	}
}

func CreateHTTPFlowWithSource(source string) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.source = source
	}
}

func CreateHTTPFlowWithURL(url string) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.url = url
	}
}

func CreateHTTPFlowWithRemoteAddr(remoteAddr string) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.remoteAddr = remoteAddr
	}
}

// 如果传入了RequestIns，则优先使用这个作为NewFuzzRequest的参数
func CreateHTTPFlowWithRequestIns(reqIns *http.Request) CreateHTTPFlowOptions {
	return func(c *CreateHTTPFlowConfig) {
		c.reqIns = reqIns
	}
}

func FuzzerResponseToHTTPFlow(db *gorm.DB, rsp *ypb.FuzzerResponse) (*schema.HTTPFlow, error) {
	return SaveFromHTTPFromRaw(db, rsp.IsHTTPS, rsp.RequestRaw, rsp.GetResponseRaw(), "fuzzer", rsp.GetUrl(), rsp.Host)
}

func SaveFromHTTP(db *gorm.DB, isHttps bool, req *http.Request, rsp *http.Response, source string, url string, remoteAddr string) (*schema.HTTPFlow, error) {
	return SaveFromHTTPWithBodySaved(db, isHttps, req, rsp, source, url, remoteAddr)
}

func SaveFromHTTPFromRaw(db *gorm.DB, isHttps bool, req []byte, rsp []byte, source string, url string, remoteAddr string) (*schema.HTTPFlow, error) {
	flow, err := CreateHTTPFlowFromHTTPWithBodySavedFromRaw(isHttps, req, rsp, source, url, remoteAddr)
	if err != nil {
		return nil, utils.Errorf("create httpflow failed: %s", err)
	}
	err = CreateOrUpdateHTTPFlow(db, flow.CalcHash(), flow)
	if err != nil {
		return nil, err
	}
	return flow, nil
}

func SaveFromHTTPWithBodySaved(db *gorm.DB, isHttps bool, req *http.Request, rsp *http.Response, source string, url string, remoteAddr string) (*schema.HTTPFlow, error) {
	flow, err := CreateHTTPFlowFromHTTPWithBodySaved(isHttps, req, rsp, source, url, remoteAddr)
	if err != nil {
		return nil, utils.Errorf("create httpflow failed: %s", err)
	}
	err = CreateOrUpdateHTTPFlow(db, flow.CalcHash(), flow)
	if err != nil {
		return nil, err
	}
	return flow, nil
}

const maxBodyLength = 4 * 1024 * 1024

func CreateHTTPFlow(opts ...CreateHTTPFlowOptions) (*schema.HTTPFlow, error) {
	c := &CreateHTTPFlowConfig{}
	for _, opt := range opts {
		opt(c)
	}

	var (
		isHttps    = c.isHttps
		reqRaw     = c.reqRaw
		rspRaw     = c.rspRaw
		fixRspRaw  = c.fixRspRaw
		source     = c.source
		url        = c.url
		remoteAddr = c.remoteAddr
		reqIns     = c.reqIns
	)

	var (
		method     string
		requestUri string
		fReq       *mutate.FuzzHTTPRequest
	)

	header, body := lowhttp.SplitHTTPHeadersAndBodyFromPacketEx(reqRaw, func(m string, r string, proto string) error {
		method = m
		requestUri = r
		return nil
	})

	if false && len(body) > maxBodyLength {
		// Truncated by saver
		reqRaw = lowhttp.ReplaceHTTPPacketBody([]byte(header), body[:maxBodyLength], false)
	}
	requestRaw := strconv.Quote(string(reqRaw))
	if strings.HasPrefix(requestRaw, `"HTTP/1.`) {
		log.Errorf("[BUG] requestRaw is invalid: %s", requestRaw)
		log.Errorf("[BUG] requestRaw is invalid: %s", requestRaw)
		log.Errorf("[BUG] requestRaw is invalid: %s", requestRaw)
	}

	// 如果已经修复过响应，则不会再修复
	if len(fixRspRaw) == 0 {
		rawNoGzip, _, _ := lowhttp.FixHTTPResponse(rspRaw)
		if len(rawNoGzip) > 0 {
			rspRaw = rawNoGzip
		}
	} else {
		rspRaw = fixRspRaw
	}
	if rspRaw == nil {
		rspRaw = make([]byte, 0)
	}

	var rspContentType string
	header, body = lowhttp.SplitHTTPHeadersAndBodyFromPacket(rspRaw, func(line string) {
		k, v := lowhttp.SplitHTTPHeader(line)
		if strings.ToLower(k) == "content-type" {
			rspContentType = v
		}
	})
	responseRaw := strconv.Quote(string(rspRaw))
	// 如果设置了 reqIns，则不会再解析 reqRaw
	if reqIns != nil {
		fReq, _ = mutate.NewFuzzHTTPRequest(reqIns)
	} else {
		fReq, _ = mutate.NewFuzzHTTPRequest(reqRaw)
	}

	flow := &schema.HTTPFlow{
		IsHTTPS:     isHttps,
		Url:         url,
		Path:        requestUri,
		Method:      method,
		BodyLength:  int64(len(body)),
		ContentType: rspContentType,
		StatusCode:  int64(lowhttp.ExtractStatusCodeFromResponse(rspRaw)),
		SourceType:  source,
		Request:     requestRaw,
		Response:    responseRaw,
		RemoteAddr:  remoteAddr,
		HiddenIndex: uuid.NewString(),
	}
	ip, _, _ := utils.ParseStringToHostPort(remoteAddr)
	if ip != "" {
		flow.IPAddress = ip
		ipInt, _ := utils.IPv4ToUint64(ip)
		if ipInt > 0 {
			flow.IPInteger = int(ipInt)
		}
	}

	if fReq != nil {
		flow.GetParamsTotal = len(fReq.GetGetQueryParams())

		postParams := fReq.GetPostJsonParams()
		if len(postParams) <= 0 {
			postParams = fReq.GetPostXMLParams()
		}
		if len(postParams) <= 0 {
			postParams = fReq.GetPostParams()
		}
		flow.PostParamsTotal = len(postParams)

		flow.CookieParamsTotal = len(fReq.GetCookieParams())
	}

	flow.Hash = flow.CalcHash()
	return flow, nil
}

func CreateHTTPFlowFromHTTPWithBodySavedFromRaw(isHttps bool, reqRaw []byte, rspRaw []byte, source string, url string, remoteAddr string, opts ...CreateHTTPFlowOptions) (*schema.HTTPFlow, error) {
	extOpts := []CreateHTTPFlowOptions{
		CreateHTTPFlowWithHTTPS(isHttps), CreateHTTPFlowWithRequestRaw(reqRaw), CreateHTTPFlowWithResponseRaw(rspRaw), CreateHTTPFlowWithSource(source), CreateHTTPFlowWithURL(url), CreateHTTPFlowWithRemoteAddr(remoteAddr),
	}
	extOpts = append(extOpts, opts...)
	flow, err := CreateHTTPFlow(extOpts...)
	if err != nil {
		return nil, err
	}
	return flow, nil
}

func createHTTPFlowFromHTTP(isHttps bool, req *http.Request, rsp *http.Response, source string, url string, remoteAddr string, opts ...CreateHTTPFlowOptions) (*schema.HTTPFlow, error) {
	urlRaw := url
	if urlRaw == "" {
		u, err := lowhttp.ExtractURLFromHTTPRequest(req, isHttps)
		if err != nil {
			log.Warnf("extract url from request failed: %s", err)
		}
		if u != nil {
			urlRaw = u.String()
		} else {
			if isHttps {
				urlRaw = "https://" + remoteAddr
			} else {
				urlRaw = "http://" + remoteAddr
			}
		}
	}

	var (
		plainRequest  []byte
		plainResponse []byte
		err           error
	)
	// 为了此处的请求与mitm的请求保持一致，需要重新从httpctx中获取
	if httpctx.GetRequestIsModified(req) {
		plainRequest = httpctx.GetHijackedRequestBytes(req)
	} else {
		plainRequest = httpctx.GetPlainRequestBytes(req)
		if len(plainRequest) <= 0 {
			plainRequest = lowhttp.DeletePacketEncoding(httpctx.GetBareRequestBytes(req))
		}
	}
	if len(plainRequest) <= 0 {
		plainRequest, err = utils.HttpDumpWithBody(req, true)
		if err != nil {
			plainRequest, err = utils.HttpDumpWithBody(req, false)
			if err != nil {
				log.Errorf("dump request failed: %s", err)
			}
		}
	}

	// 为了此处的响应与mitm的响应保持一致，需要重新从httpctx中获取
	if rsp != nil {
		if httpctx.GetResponseIsModified(req) {
			plainResponse = httpctx.GetHijackedResponseBytes(req)
		} else {
			plainResponse = httpctx.GetPlainResponseBytes(req)
			if len(plainResponse) <= 0 {
				plainResponse = lowhttp.DeletePacketEncoding(httpctx.GetBareResponseBytes(req))
			}
		}
		if len(plainResponse) <= 0 {
			plainResponse, err = utils.HttpDumpWithBody(rsp, true)
			if err != nil {
				log.Errorf("dump response failed: %s", err)
			}
		}
	} else {
		plainResponse = make([]byte, 0)
	}

	return CreateHTTPFlowFromHTTPWithBodySavedFromRaw(isHttps, plainRequest, plainResponse, source, urlRaw, remoteAddr, opts...)
}

func CreateHTTPFlowFromHTTPWithNoRspSaved(isHttps bool, req *http.Request, source string, url string, remoteAddr string, opts ...CreateHTTPFlowOptions) (*schema.HTTPFlow, error) {
	return createHTTPFlowFromHTTP(isHttps, req, nil, source, url, remoteAddr, opts...)
}

func CreateHTTPFlowFromHTTPWithBodySaved(isHttps bool, req *http.Request, rsp *http.Response, source string, url string, remoteAddr string, opts ...CreateHTTPFlowOptions) (*schema.HTTPFlow, error) {
	return createHTTPFlowFromHTTP(isHttps, req, rsp, source, url, remoteAddr, opts...)
}

func UpdateHTTPFlowTags(db *gorm.DB, i *schema.HTTPFlow) error {
	if i == nil {
		return nil
	}
	db = db.Model(&schema.HTTPFlow{})

	if i.ID > 0 {
		if db = db.Where("id = ?", i.ID).Update("tags", i.Tags); db.Error != nil {
			log.Errorf("update tags(by id) failed: %s", db.Error)
			return db.Error
		}
	} else if i.HiddenIndex != "" {
		if db = db.Where("hidden_index = ?", i.HiddenIndex).Update("tags", i.Tags); db.Error != nil {
			log.Errorf("update tags(by hidden_index) failed: %s", db.Error)
			return db.Error
		}
	} else if i.Hash != "" {
		if db = db.Where("hash = ?", i.HiddenIndex).Update("tags", i.Tags); db.Error != nil {
			log.Errorf("update tags(by hash) failed: %s", db.Error)
			return db.Error
		}
	}
	return nil
}

func InsertHTTPFlow(db *gorm.DB, i *schema.HTTPFlow) (fErr error) {
	defer func() {
		if err := recover(); err != nil {
			fErr = utils.Errorf("met panic error: %v", err)
			debug.PrintStack()
		}
	}()

	i.ID = 0
	if db = db.Model(&schema.HTTPFlow{}).Save(i); db.Error != nil {
		return utils.Errorf("insert HTTPFlow failed: %s", db.Error)
	}

	return nil
}

func CreateOrUpdateHTTPFlow(db *gorm.DB, hash string, i *schema.HTTPFlow) (fErr error) {
	defer func() {
		if err := recover(); err != nil {
			fErr = utils.Errorf("met panic error: %v", err)
		}
	}()

	db = db.Model(&schema.HTTPFlow{})

	if db := db.Where("hash = ?", hash).Assign(i).FirstOrCreate(i); db.Error != nil {
		return utils.Errorf("create/update HTTPFlow failed: %s", db.Error)
	}
	return nil
}

func GetHTTPFlow(db *gorm.DB, id int64) (*schema.HTTPFlow, error) {
	var req schema.HTTPFlow
	if db := db.Model(&schema.HTTPFlow{}).Where("id = ?", id).First(&req); db.Error != nil {
		return nil, utils.Errorf("get HTTPFlow failed: %s", db.Error)
	}

	return &req, nil
}

func GetHTTPFlowByIDOrHash(db *gorm.DB, id int64, hash string) (*schema.HTTPFlow, error) {
	var req schema.HTTPFlow
	if db := db.Model(&schema.HTTPFlow{}).Where("id = ? OR hash = ?", id, hash).First(&req); db.Error != nil {
		return nil, utils.Errorf("get HTTPFlow failed: %s", db.Error)
	}

	return &req, nil
}

func GetHTTPFlowByHash(db *gorm.DB, hash string) (*schema.HTTPFlow, error) {
	var req schema.HTTPFlow
	if db := db.Model(&schema.HTTPFlow{}).Where("hash = ?", hash).First(&req); db.Error != nil {
		return nil, utils.Errorf("get HTTPFlow failed: %s", db.Error)
	}

	return &req, nil
}

func DeleteHTTPFlowByID(db *gorm.DB, id int64) error {
	if db := db.Model(&schema.HTTPFlow{}).Where(
		"id = ?", id,
	).Unscoped().Delete(&schema.HTTPFlow{}); db.Error != nil {
		return db.Error
	}
	return nil
}

func DeleteHTTPFlow(db *gorm.DB, req *ypb.DeleteHTTPFlowRequest) error {
	if req.GetDeleteAll() {
		db.DropTableIfExists(&schema.HTTPFlow{})
		if db := db.Exec(`UPDATE SQLITE_SEQUENCE SET SEQ=0 WHERE NAME='http_flows';`); db.Error != nil {
			log.Errorf("update sqlite sequence failed: %s", db.Error)
		}
		db.AutoMigrate(&schema.HTTPFlow{})
		DeleteProjectKeyBareRequestAndResponse(db)
		return nil
	}

	if len(req.GetId()) > 0 {
		db = db.Or("false")
		db = bizhelper.ExactQueryInt64ArrayOr(db, "id", req.GetId())
		// for _, id := range req.GetId() {
		// 	db = db.Or("id = ?", id)
		// }
		db.Unscoped().Delete(&schema.HTTPFlow{})
		return nil
	}

	if req.GetFilter() != nil {
		db = FilterHTTPFlow(db, req.GetFilter())
		db.Unscoped().Delete(&schema.HTTPFlow{})
		return nil
	}

	if req.GetURLPrefix() != "" {
		db = db.Model(&schema.HTTPFlow{})
		db = bizhelper.FuzzQueryLike(db, "url", req.GetURLPrefix()).Unscoped().Delete(&schema.HTTPFlow{})
		if db.Error != nil {
			return nil
		}
		return nil
	}

	if req.GetItemHash() != nil {
		db = db.Model(&schema.HTTPFlow{})
		db = bizhelper.ExactQueryStringArrayOr(db, "hash", req.GetItemHash())
		if db := db.Where("true").Unscoped().Delete(&schema.HTTPFlow{}); db.Error != nil {
			return db.Error
		}
		return nil
	}

	if len(req.GetURLPrefixBatch()) > 0 {
		db = db.Model(&schema.HTTPFlow{})
		db = bizhelper.FuzzQueryStringArrayOrLike(db, "url", req.GetURLPrefixBatch())
		db = db.Unscoped().Delete(&schema.HTTPFlow{})
		if db.Error != nil {
			return db.Error
		}
		return nil
	}
	return nil
}

func FilterHTTPFlow(db *gorm.DB, params *ypb.QueryHTTPFlowRequest) *gorm.DB {
	db = db.Model(&schema.HTTPFlow{}) //.Debug()
	if params == nil {
		params = &ypb.QueryHTTPFlowRequest{}
	}

	db = bizhelper.ExactQueryStringArrayOr(db, "source_type", utils.PrettifyListFromStringSplited(params.SourceType, ","))
	// 过滤 Methods
	if ms := utils.StringArrayFilterEmpty(utils.PrettifyListFromStringSplited(params.GetMethods(), ",")); ms != nil {
		db = bizhelper.ExactQueryStringArrayOr(db, "method", ms)
	}
	// 搜索 URL
	db = bizhelper.FuzzQueryLike(db, "url", params.GetSearchURL())
	// status code 这里可以支持范围搜索
	db = bizhelper.QueryBySpecificPorts(db, "status_code", params.GetStatusCode())
	if params.GetHaveBody() {
		db = db.Where("body_length > 0")
	}
	db = bizhelper.ExactQueryString(db, "runtime_id", params.GetRuntimeId())
	db = bizhelper.ExactOrQueryStringArrayOr(db, "runtime_id", params.GetRuntimeIDs())
	db = bizhelper.ExactQueryString(db, "from_plugin", params.GetFromPlugin())

	// 搜索是否有对应的参数
	if params.GetHaveCommonParams() {
		db = db.Where("((get_params_total > 0) OR (post_params_total > 0)) OR (cookie_params_total > 0)")
	}

	if params.GetHaveParamsTotal() == "true" {
		db = db.Where("((get_params_total > 0) OR (post_params_total > 0))")
	} else if params.GetHaveParamsTotal() == "false" {
		db = db.Where("((get_params_total = 0) and (post_params_total = 0))")
	}

	if len(params.GetTags()) > 0 {
		db = bizhelper.FuzzSearchWithStringArrayOrEx(db, []string{"tags"}, params.GetTags(), false)
	}

	if len(params.GetColor()) > 0 {
		db = bizhelper.FuzzSearchWithStringArrayOrAf(db, []string{"tags"}, params.GetColor(), false)
	}

	// 搜索 Content-Type
	db = bizhelper.FuzzQueryStringArrayOrLike(db, "content_type",
		utils.StringArrayFilterEmpty(utils.PrettifyListFromStringSplited(params.GetSearchContentType(), ",")))

	if len(params.GetIncludeInUrl()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLike(db, "url", params.GetIncludeInUrl())
	}

	if len(params.GetIncludeInIP()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLike(db, "ip_address", params.GetIncludeInIP())
	}

	if len(params.GetIncludeId()) > 0 {
		db = bizhelper.ExactQueryInt64ArrayOr(db, "id", params.GetIncludeId())
	}

	if len(params.GetExcludeInUrl()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLikeExclude(db, "url", params.GetExcludeInUrl())
	}

	if len(params.GetExcludeInIP()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLikeExclude(db, "ip_address", params.GetExcludeInIP())
	}

	if len(params.GetIncludePath()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLike(db, "path", params.GetIncludePath())
	}

	if len(params.GetExcludePath()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLikeExclude(db, "path", params.GetExcludePath())
	}

	if len(params.GetIncludeSuffix()) > 0 {
		var suffixes []string
		for _, suffix := range params.GetIncludeSuffix() {
			if !strings.HasPrefix(suffix, ".") {
				suffix = "." + suffix
			}
			suffixes = append(suffixes, suffix)
		}
		db = bizhelper.FuzzQueryStringArrayOrLike(db, "path", suffixes)
	}
	if len(params.GetExcludeSuffix()) > 0 {
		var suffixes []string
		for _, suffix := range params.GetIncludeSuffix() {
			if !strings.HasPrefix(suffix, ".") {
				suffix = "." + suffix
			}
			suffixes = append(suffixes, suffix)
		}
		db = bizhelper.FuzzQueryStringArrayOrLikeExclude(db, "path", suffixes)
	}
	if len(params.GetExcludeContentType()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLikeExclude(db, "content_type", params.GetExcludeContentType())
	}

	if len(params.GetExcludeId()) > 0 {
		db = bizhelper.ExactExcludeQueryInt64Array(db, "id", params.GetExcludeId())
	}

	if len(params.GetIncludeHash()) > 0 {
		db = bizhelper.ExactQueryStringArrayOr(db, "hash", params.GetIncludeHash())
	}

	if params.AfterBodyLength > 0 {
		db = db.Where("body_length >= ?", params.AfterBodyLength)
	}
	if params.BeforeBodyLength > 0 {
		db = db.Where("body_length <= ?", params.BeforeBodyLength)
	}

	return db
}

func QuickSearchHTTPFlowCount(token string) int {
	db := consts.GetGormProjectDatabase()
	var count int
	db.Model(&schema.HTTPFlow{}).Where(
		"(request like ?) OR (response like ?) OR (url like ?)",
		"%"+token+"%",
		"%"+token+"%",
		"%"+token+"%",
	).Count(&count)
	return count
}

func QuickSearchMITMHTTPFlowCount(token string) int {
	db := consts.GetGormProjectDatabase()
	var count int
	db.Model(&schema.HTTPFlow{}).Where(
		"(request like ?) OR (response like ?) OR (url like ?)",
		"%"+token+"%",
		"%"+token+"%",
		"%"+token+"%",
	).Where("source_type = ?", "mitm").Count(&count)
	return count
}

// BuildHTTPFlowQuery 构建带有过滤条件的查询
func BuildHTTPFlowQuery(db *gorm.DB, params *ypb.QueryHTTPFlowRequest) *gorm.DB {
	// 应用所有过滤条件
	if params == nil {
		params = &ypb.QueryHTTPFlowRequest{}
	}

	if !params.GetFull() {
		extraSelectField := ""
		if params.GetWithPayload() {
			extraSelectField = "payload,"
		}
		// 只查询部分字段，主要是为了处理大的 response 和 request 的情况，同时告诉用户
		// max request size is 200K -> 200 * 1024 -> 204800
		// max response size is 500K -> 500 * 1024 -> 512000
		db = db.Select(fmt.Sprintf(`id,created_at,updated_at,hidden_index,%s -- basic gorm fields
body_length, -- handle body length should be careful, if it's big, no return response

-- metainfo
is_http_s, -- legacy
no_fix_content_length, hash,url, path, method,
content_type, status_code, source_type,
get_params_total, post_params_total, cookie_params_total,
ip_address, remote_addr, ip_integer,
tags, is_websocket, websocket_hash, runtime_id, from_plugin,

-- request is larger than 200K, return empty string
LENGTH(request) > 204800 as is_request_oversize,
CASE WHEN LENGTH(request) > 204800 THEN '' ELSE request END as request,

-- response is larger than 500K, return empty string
LENGTH(response) > 512000 as is_response_oversize,
CASE WHEN LENGTH(response) > 512000 THEN '' ELSE response END as response,

-- is response too large
is_too_large_response, 
too_large_response_header_file, too_large_response_body_file
`, extraSelectField))
	}

	if params.Pagination == nil {
		params.Pagination = &ypb.Paging{
			Page:    1,
			Limit:   30,
			OrderBy: "updated_at",
			Order:   "desc",
		}
	}

	p := params.Pagination
	if p.OrderBy == "" {
		p.OrderBy = "id" // 如果 没有设置 orderby 则以ID排序
	}

	if params.GetAfterUpdatedAt() > 0 {
		db = bizhelper.QueryByTimeRangeWithTimestamp(db, "updated_at", params.GetAfterUpdatedAt(), time.Now().Add(10*time.Minute).Unix())
	}
	db = bizhelper.FuzzSearchEx(db, []string{
		"tags", "url", "path", "request",
		"response", "remote_addr",
	}, params.GetKeyword(), false)
	if params.GetAfterId() > 0 {
		db = db.Where("id > ?", params.GetAfterId())
	}
	if params.GetBeforeUpdatedAt() > 0 {
		db = bizhelper.QueryByTimeRangeWithTimestamp(db, "updated_at", 0, params.GetBeforeUpdatedAt())
	}
	if params.GetBeforeId() > 0 {
		db = db.Where("id < ?", params.GetBeforeId())
	}

	if params.GetOnlyWebsocket() {
		// log.Info("query websocket request flow")
		db = db.Where("(is_websocket = ?) AND (url LIKE 'ws%')", params.OnlyWebsocket)
	}
	switch params.GetIsWebsocket() {
	case "http/https":
		db = db.Where("is_websocket = false")
	case "websocket":
		db = db.Where("(is_websocket = true) AND (url LIKE 'ws%')")
	}

	db = bizhelper.QueryOrder(db, p.OrderBy, p.Order)

	db = FilterHTTPFlow(db, params)

	return db
}

func QueryHTTPFlow(db *gorm.DB, params *ypb.QueryHTTPFlowRequest) (paging *bizhelper.Paginator, httpflows []*schema.HTTPFlow, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
			utils.PrintCurrentGoroutineRuntimeStack()
		}
	}()
	queryDB := BuildHTTPFlowQuery(db.Model(&schema.HTTPFlow{}), params)

	return SelectHTTPFlowFromDB(queryDB, params)
}

func SelectHTTPFlowFromDB(queryDB *gorm.DB, params *ypb.QueryHTTPFlowRequest) (paging *bizhelper.Paginator, httpflows []*schema.HTTPFlow, err error) {
	var limitFlows, fullFlows []*schema.HTTPFlow

	if params.OffsetId > 0 {
		offsetDB := queryDB
		if params.Pagination.Order == "desc" {
			offsetDB = offsetDB.Where("id < ?", params.OffsetId)
		} else {
			offsetDB = offsetDB.Where("id > ?", params.OffsetId)
		}
		offsetDB.Limit(int(params.Pagination.Limit)).Offset(0).Scan(&limitFlows)
		paging, queryDB = bizhelper.Paging(queryDB, int(params.Pagination.Page), int(params.Pagination.Limit), &fullFlows)
	} else {
		paging, queryDB = bizhelper.Paging(queryDB, int(params.Pagination.Page), int(params.Pagination.Limit), &limitFlows)
	}

	if queryDB.Error != nil {
		return nil, nil, utils.Errorf("paging failed: %s", queryDB.Error)
	}

	return paging, limitFlows, nil
}

type HTTPFlowUrl struct {
	Url string `json:"url"`
}

func YieldHTTPUrl(db *gorm.DB, ctx context.Context) chan *HTTPFlowUrl {
	outC := make(chan *HTTPFlowUrl)
	go func() {
		defer close(outC)

		page := 1
		for {
			var items []*HTTPFlowUrl
			if _, b := bizhelper.NewPagination(&bizhelper.Param{
				DB:    db,
				Page:  page,
				Limit: 1000,
			}, &items); b.Error != nil {
				log.Errorf("paging failed: %s", b.Error)
				return
			}

			page++

			for _, d := range items {
				select {
				case <-ctx.Done():
					return
				case outC <- d:
				}
			}

			if len(items) < 1000 {
				return
			}
		}
	}()
	return outC
}

func YieldHTTPFlows(db *gorm.DB, ctx context.Context) chan *schema.HTTPFlow {
	outC := make(chan *schema.HTTPFlow)
	go func() {
		defer close(outC)

		page := 1
		for {
			var items []*schema.HTTPFlow
			if _, b := bizhelper.NewPagination(&bizhelper.Param{
				DB:    db,
				Page:  page,
				Limit: 1000,
			}, &items); b.Error != nil {
				log.Errorf("paging failed: %s", b.Error)
				return
			}

			page++

			for _, d := range items {
				select {
				case <-ctx.Done():
					return
				case outC <- d:
				}
			}

			if len(items) < 1000 {
				return
			}
		}
	}()
	return outC
}

const (
	HTTPFLOW_TAG        = "HTTPFLOW_TAG"
	HTTPFLOW_STATUSCODE = "HTTPFLOW_STATUSCODE"
)

/*func HTTPFlowStatusCode(refreshRequest bool) (req []*TagAndStatusCode, err error) {
	var db = consts.GetGormProjectDatabase()
	if db == nil {
		log.Error("cannot found database config")
		return nil, utils.Error("empty database")
	}
	if !refreshRequest {
		value := GetKey(db, HTTPFLOW_STATUSCODE)
		if value != "" {
			var statusCode []*TagAndStatusCode
			_ = json.Unmarshal([]byte(value), &statusCode)
			if len(statusCode) > 0 {
				return statusCode, nil
			}
		}
	}

	// log.Info("start to execute updating tags")
	db = db.Raw(`SELECT count(*) as count, status_code as value FROM http_flows GROUP BY status_code order by count desc;`)
	rows, err := db.Rows()
	if err != nil {
		return nil, utils.Errorf("rows failed: %s", err)
	}

	var statusCode = make([]*TagAndStatusCode, 0)
	for rows.Next() {
		var codeName string
		var count int
		err = rows.Scan(&count, &codeName)
		if err != nil {
			log.Errorf("scan code stats failed: %s", err)
			continue
		}
		statusCode = append(statusCode, &TagAndStatusCode{
			Value: codeName,
			Count: count,
		})
	}

	raw, _ := json.Marshal(statusCode)
	if len(raw) > 0 {
		log.Infof("start to cache statusCode[%v]", len(raw))
		SetKey(consts.GetGormProfileDatabase(), HTTPFLOW_STATUSCODE, string(raw))
	}
	return statusCode, nil
}*/

func HTTPFlowTags(refreshRequest bool) ([]*TagAndStatusCode, error) {
	tagCounts := make(map[string]int)
	for _, v := range model.GlobalHTTPFlowCache.GetAll() {
		for _, tag := range strings.Split(v.Tags, "|") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagCounts[tag]++
			}
		}
	}
	tags := make([]*TagAndStatusCode, 0)
	for k, v := range tagCounts {
		if !strings.HasPrefix(k, schema.COLORPREFIX) {
			tags = append(tags, &TagAndStatusCode{
				Value: k,
				Count: v,
			})
		}
	}
	return tags, nil
}

func QueryWebsocketFlowsByHTTPFlowHash(db *gorm.DB, req *ypb.DeleteHTTPFlowRequest) *gorm.DB {
	db = db.Model(&schema.HTTPFlow{})

	if len(req.GetId()) > 0 {
		db = db.Or("false")
		db = bizhelper.ExactQueryInt64ArrayOr(db, "id", req.GetId())
	}

	if req.GetFilter() != nil {
		db = FilterHTTPFlow(db, req.GetFilter())
	}

	if req.GetURLPrefix() != "" {
		db = bizhelper.FuzzQueryLike(db, "url", req.GetURLPrefix())
	}
	if len(req.GetURLPrefixBatch()) > 0 {
		db = bizhelper.FuzzQueryStringArrayOrLike(db, "url", req.GetURLPrefixBatch())
	}
	if req.GetItemHash() != nil {
		db = bizhelper.ExactQueryStringArrayOr(db, "hash", req.GetItemHash())
		db = db.Where("true")
	}
	return db
}

func ExportHTTPFlow(db *gorm.DB, params *ypb.ExportHTTPFlowsRequest) (paging *bizhelper.Paginator, ret []*schema.HTTPFlow, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
			utils.PrintCurrentGoroutineRuntimeStack()
		}
	}()
	queryParams := params.ExportWhere
	queryParams.IncludeId = params.Ids

	db = db.Model(&schema.HTTPFlow{})
	// overwrite Select Field, fix payloads
	for i, field := range params.FieldName {
		if field != "payloads" {
			continue
		}
		queryParams.WithPayload = true
		params.FieldName[i] = "payload"
	}

	queryDB := BuildHTTPFlowQuery(db, queryParams).Select(params.FieldName)
	return SelectHTTPFlowFromDB(queryDB, queryParams)
}

func HTTPFlowToOnline(db *gorm.DB, hash []string) error {
	db = db.Model(&schema.HTTPFlow{})
	db = bizhelper.ExactOrQueryStringArrayOr(db, "hash", hash)
	db = db.Update(map[string]interface{}{"upload_online": true})
	if db.Error != nil {
		return utils.Errorf("HTTPFlowToOnline failed %s", db.Error)
	}
	return nil
}
