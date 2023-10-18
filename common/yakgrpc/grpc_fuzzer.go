package yakgrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/samber/lo"
	"github.com/yaklang/yaklang/common/consts"
	filter2 "github.com/yaklang/yaklang/common/filter"
	"github.com/yaklang/yaklang/common/go-funk"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/mutate"
	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/utils/lowhttp"
	"github.com/yaklang/yaklang/common/yak"
	"github.com/yaklang/yaklang/common/yak/cartesian"
	"github.com/yaklang/yaklang/common/yak/httptpl"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"github.com/yaklang/yaklang/common/yakgrpc/yakit"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"

	"github.com/saintfish/chardet"
	uuid "github.com/satori/go.uuid"
)

var (
	taskIDBackTrackMap = make(map[int64]int64)
	taskIDBackTrackMu  sync.Mutex
)

func taskIDBackTrack(taskID int64) []int64 {
	var ok bool
	count := 0

	results := make([]int64, 0)
	for {
		count++
		taskID, ok = taskIDBackTrackMap[taskID]
		if !ok || taskID == 0 || count > 1000 {
			break
		}
		results = append(results, taskID)
	}
	return results
}

func Chardet(raw []byte) string {
	res, err := chardet.NewTextDetector().DetectBest(raw)
	if err != nil {
		return "utf-8"
	}
	return res.Charset
}

func (s *Server) ExtractUrl(ctx context.Context, req *ypb.FuzzerRequest) (*ypb.ExtractedUrl, error) {
	u, err := lowhttp.ExtractURLFromHTTPRequestRaw([]byte(req.Request), req.GetIsHTTPS())
	if err != nil {
		return nil, err
	}
	return &ypb.ExtractedUrl{Url: u.String()}, nil
}

func (s *Server) StringFuzzer(rootCtx context.Context, req *ypb.StringFuzzerRequest) (*ypb.StringFuzzerResponse, error) {
	max := req.GetLimit()
	timeoutSeconds := req.GetTimeoutSeconds()
	var ctx = rootCtx
	var cancel = func() {}
	if timeoutSeconds > 0 {
		ctx, cancel = context.WithTimeout(rootCtx, time.Duration(timeoutSeconds)*time.Second)
	}
	defer cancel()

	var res [][]byte
	var counter int64
	mutate.FuzzTagExec(
		req.GetTemplate(),
		mutate.Fuzz_WithResultHandler(func(origin string, payloads []string) bool {
			select {
			case <-ctx.Done():
				return false
			default:
				if max > 0 && counter >= max {
					return false
				}
			}
			counter++
			res = append(res, []byte(origin))
			return true
		}),
		yak.Fuzz_WithHotPatch(rootCtx, req.GetHotPatchCode()),
		mutate.Fuzz_WithEnableFiletag(),
	)
	return &ypb.StringFuzzerResponse{Results: res}, nil
}

func (s *Server) RedirectRequest(ctx context.Context, req *ypb.RedirectRequestParams) (*ypb.FuzzerResponse, error) {
	result := lowhttp.GetRedirectFromHTTPResponse([]byte(req.GetResponse()), false)
	if result == "" {
		return nil, utils.Error("cannot find redirect url")
	}

	isHttps := req.GetIsHttps()
	if strings.HasPrefix(result, "https://") {
		isHttps = true
	}
	_ = isHttps
	newUrl := lowhttp.MergeUrlFromHTTPRequest([]byte(req.GetRequest()), result, isHttps)
	resultRequest := lowhttp.UrlToGetRequestPacket(newUrl, []byte(req.GetRequest()), isHttps, lowhttp.ExtractCookieJarFromHTTPResponse([]byte(req.GetResponse()))...)
	if resultRequest == nil {
		return nil, utils.Errorf("cannot merge request packet. redirect url: %s", newUrl)
	}

	start := time.Now()
	host, port, _ := utils.ParseStringToHostPort(newUrl)
	rspIns, err := lowhttp.HTTPWithoutRedirect(
		lowhttp.WithHttps(isHttps),
		lowhttp.WithHost(host),
		lowhttp.WithPort(port),
		lowhttp.WithRequest(resultRequest),
		lowhttp.WithTimeoutFloat(req.GetPerRequestTimeoutSeconds()),
		lowhttp.WithGmTLS(req.GetIsGmTLS()),
		lowhttp.WithProxy(utils.PrettifyListFromStringSplited(req.GetProxy(), ",")...),
	)
	if err != nil {
		return nil, err
	}
	rspRaw := rspIns.RawPacket
	// 提取响应
	extractHTTPResponseResult, err := s.ExtractHTTPResponse(ctx, &ypb.ExtractHTTPResponseParams{
		HTTPResponse: string(rspRaw),
		Extractors:   req.GetExtractors(),
	})
	var extractResults []*ypb.KVPair
	if err == nil && extractHTTPResponseResult != nil && extractHTTPResponseResult.GetValues() != nil {
		for _, value := range extractHTTPResponseResult.GetValues() {
			extractResults = append(extractResults, &ypb.KVPair{
				Key:   value.GetKey(),
				Value: value.GetValue(),
			})
		}
	}
	// 匹配响应
	var httpTPLmatchersResult bool
	if len(req.GetMatchers()) != 0 {
		httpTplMatcher := make([]*httptpl.YakMatcher, 0)
		for _, matcher := range req.GetMatchers() {
			httpTplMatcher = append(httpTplMatcher, httptpl.NewMatcherFromGRPCModel(matcher))
		}
		cond := "and"
		switch ret := strings.ToLower(req.GetMatchersCondition()); ret {
		case "or", "and":
			cond = ret
		default:
		}
		ins := &httptpl.YakMatcher{
			SubMatcherCondition: cond,
			SubMatchers:         httpTplMatcher,
		}
		var mergedParams = make(map[string]interface{})
		renderedParams, err := s.RenderVariables(ctx, &ypb.RenderVariablesRequest{
			Params: funk.Map(req.GetParams(), func(i *ypb.FuzzerParamItem) *ypb.KVPair {
				return &ypb.KVPair{Key: i.GetKey(), Value: i.GetValue()}
			}).([]*ypb.KVPair),
			IsHTTPS: req.GetIsHttps(),
			IsGmTLS: req.GetIsGmTLS(),
		})
		if err != nil {
			return nil, utils.Errorf("render variables failed: %v", err)
		}
		for _, kv := range renderedParams.GetResults() {
			mergedParams[kv.GetKey()] = kv.GetValue()
		}

		matcherParams := utils.CopyMapInterface(mergedParams)
		httpTPLmatchersResult, err = ins.Execute(&lowhttp.LowhttpResponse{RawPacket: rspRaw}, matcherParams)
		if err != nil {
			log.Errorf("httptpl.YakMatcher execute failed: %s", err)
		}
	}

	var rsp = &ypb.FuzzerResponse{
		Method:                "GET",
		ResponseRaw:           rspRaw,
		GuessResponseEncoding: Chardet(rspRaw),
		RequestRaw:            resultRequest,
		ExtractedResults:      extractResults,
		MatchedByMatcher:      httpTPLmatchersResult,
		HitColor:              req.GetHitColor(),
	}
	rsp.UUID = uuid.NewV4().String()
	rsp.Timestamp = start.Unix()
	rsp.DurationMs = time.Now().Sub(start).Milliseconds()

	requestIns, err := lowhttp.ParseBytesToHttpRequest(resultRequest)
	if err != nil {
		return nil, err
	}
	rsp.Host = requestIns.Header.Get("Host")
	if rsp.Host == "" {
		rsp.Host = requestIns.Host
	}

	responseIns, err := lowhttp.ParseBytesToHTTPResponse(rspRaw)
	if responseIns != nil {
		rsp.Ok = true
		rsp.StatusCode = int32(responseIns.StatusCode)
		rsp.ContentType = responseIns.Header.Get("Content-Type")
		var bodyLen int64 = 0
		if responseIns.Body != nil {
			raw, _ := ioutil.ReadAll(responseIns.Body)
			bodyLen = int64(len(raw))
		}
		rsp.BodyLength = bodyLen

		// 解析 Headers
		for k, vs := range responseIns.Header {
			for _, v := range vs {
				rsp.Headers = append(rsp.Headers, &ypb.HTTPHeader{
					Header: k,
					Value:  v,
				})
			}
		}
	}
	return rsp, nil
}

func (s *Server) PreloadHTTPFuzzerParams(ctx context.Context, req *ypb.PreloadHTTPFuzzerParamsRequest) (*ypb.PreloadHTTPFuzzerParamsResponse, error) {
	vars := httptpl.NewVars()
	for _, k := range req.GetParams() {
		if k.GetType() == "raw" {
			vars.Set(k.GetKey(), k.GetValue())
			continue
		}
		vars.AutoSet(k.GetKey(), k.GetValue())
	}
	var results []*ypb.FuzzerParamItem
	for k, v := range vars.ToMap() {
		results = append(results, &ypb.FuzzerParamItem{
			Key:   k,
			Value: utils.InterfaceToString(v),
			Type:  "raw",
		})
	}
	return &ypb.PreloadHTTPFuzzerParamsResponse{Values: results}, nil
}

func (s *Server) HTTPFuzzer(req *ypb.FuzzerRequest, stream ypb.Yak_HTTPFuzzerServer) (finalError error) {
	defer func() {
		if err := recover(); err != nil {
			finalError = utils.Errorf("panic from httpfuzzer: %v", err)
			utils.PrintCurrentGoroutineRuntimeStack()
		}
	}()
	// retry
	isRetry := req.RetryTaskID > 0

	// hot code
	var extraOpt []mutate.FuzzConfigOpt
	if strings.TrimSpace(req.GetHotPatchCode()) != "" {
		extraOpt = append(extraOpt, yak.Fuzz_WithHotPatch(stream.Context(), req.GetHotPatchCode()))
	}

	/*
		Plugins
	*/
	var pocs []*yakit.YakScript
	for _, i := range req.GetYamlPoCNames() {
		poc, err := yakit.GetYakScriptByName(consts.GetGormProfileDatabase(), i)
		if err != nil {
			log.Errorf("get yaml poc[%v] failed: %s", i, err)
			continue
		}
		if poc.Type != "nuclei" {
			log.Errorf("poc[%s] is not yaml poc: %s", i, poc.Type)
			continue
		}
		pocs = append(pocs, poc)
	}

	var batchTarget string
	if req.GetBatchTargetFile() {
		if ret := utils.GetFirstExistedFile(string(req.BatchTarget)); ret != "" {
			fp, err := os.Open(ret)
			if err != nil {
				return utils.Errorf("open batch target file failed: %s", err)
			}
			raw, _ := io.ReadAll(fp)
			fp.Close()
			batchTarget = strings.TrimSpace(string(raw))
		} else {
			return utils.Errorf("batch target file not found: %s", req.GetBatchTarget())
		}
	} else {
		batchTarget = string(req.GetBatchTarget())
	}

	var swg = utils.NewSizedWaitGroup(int(req.GetConcurrent()))
	defer swg.Wait()
	var feedbackWg = new(sync.WaitGroup)
	defer func() {
		feedbackWg.Wait()
	}()
	var feedbackResponse = func(rsp *ypb.FuzzerResponse, skipPoC bool) error {
		err := stream.Send(rsp)
		if err != nil {
			return err
		}

		if skipPoC {
			return nil
		}

		feedbackWg.Add(1)
		defer func() {
			defer feedbackWg.Done()
			for _, p := range pocs {
				poc := p
				err := swg.AddWithContext(stream.Context())
				if err != nil {
					break
				}
				go func() {
					defer swg.Done()
					defer func() {
						if err := recover(); err != nil {
							spew.Dump(err)
							utils.PrintCurrentGoroutineRuntimeStack()
						}
					}()
					httptpl.ScanPacket(
						rsp.RequestRaw, lowhttp.WithHttps(rsp.IsHTTPS),
						httptpl.WithTemplateRaw(poc.Content),
						lowhttp.WithResponseCallback(func(i *lowhttp.LowhttpResponse) {
							err := stream.Send(ConvertLowhttpResponseToFuzzerResponseBase(i))
							if err != nil {
								log.Errorf("yaml poc send failed")
							}
						}),
						httptpl.WithOnRisk(rsp.Url, func(i *yakit.Risk) {
							log.Infof("found risk: %s", i.Title)
						}),
					)
				}()

			}
		}()
		return nil
	}

	if req.GetHistoryWebFuzzerId() > 0 {
		for resp := range yakit.YieldWebFuzzerResponses(s.GetProjectDatabase(), stream.Context(), int(req.GetHistoryWebFuzzerId())) {
			rsp, err := resp.ToGRPCModel()
			if err != nil {
				continue
			}
			err = feedbackResponse(rsp, true)
			if err != nil {
				log.Errorf("stream send failed: %s", err)
				continue
			}
		}
		return nil
	}
	if !isRetry && req.GetRequest() == "" && len(req.GetRequestRaw()) <= 0 {
		return utils.Errorf("empty request is not allowed")
	}

	var proxies = utils.StringArrayFilterEmpty(utils.PrettifyListFromStringSplited(req.GetProxy(), ","))
	var concurrent = req.GetConcurrent()
	if concurrent <= 0 {
		concurrent = 20
	}
	var timeoutSeconds = req.GetPerRequestTimeoutSeconds()
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}

	task, err := yakit.SaveWebFuzzerTask(s.GetProjectDatabase(), req, 0, false, "executing...")
	if err != nil {
		return utils.Errorf("save to web fuzzer to database failed: %s", err)
	}
	var taskId = task.ID
	task.FuzzerIndex = req.GetFuzzerIndex()
	task.FuzzerTabIndex = req.GetFuzzerTabIndex()
	defer func() {
		if db := s.GetProjectDatabase().Save(task); db.Error != nil {
			log.Errorf("update web fuzzer task failed: %s", db.Error)
		}
	}()

	taskIDBackTrackMu.Lock()
	taskIDBackTrackMap[int64(task.ID)] = req.RetryTaskID
	taskIDBackTrackMu.Unlock()

	/* 丢包过滤器 */
	includeStatusCodeFilter := utils.NewPortsFilter()
	var maxBody, minBody int64
	var regexps, keywords []string
	filter := req.GetFilter()
	if filter != nil {
		includeStatusCodeFilter.Add(filter.GetStatusCode()...)
		regexps = filter.GetRegexps()
		keywords = filter.GetKeywords()
		minBody = filter.GetMinBodySize()
		maxBody = filter.GetMaxBodySize()
	}

	var rawRequest []byte
	if !isRetry {
		if len(req.GetRequestRaw()) > 0 {
			rawRequest = req.GetRequestRaw()
		} else {
			rawRequest = []byte(req.GetRequest())
		}
	}

	// 保存 request 中 host/port
	defer func() {
		if req.GetActualAddr() != "" {
			task.Host = req.GetActualAddr()
		} else {
			results := extractHostRegexp.FindStringSubmatch(string(rawRequest))
			if len(results) > 1 {
				task.Host = results[1]
				if len(task.Host) > 40 {
					task.Host = task.Host[:40] + "..."
				}
			}
		}
		_, task.Port, _ = utils.ParseStringToHostPort(task.Host)
	}()

	var inStatusCode = utils.ParseStringToPorts(req.GetRetryInStatusCode())
	var notInStatusCode = utils.ParseStringToPorts(req.GetRetryNotInStatusCode())

	var httpTplMatcher = make([]*httptpl.YakMatcher, len(req.GetMatchers()))
	var httpTplExtractor = make([]*httptpl.YakExtractor, len(req.GetExtractors()))
	var haveHTTPTplMatcher = len(httpTplMatcher) > 0
	var haveHTTPTplExtractor = len(httpTplExtractor) > 0
	if haveHTTPTplExtractor {
		for i, e := range req.GetExtractors() {
			httpTplExtractor[i] = httptpl.NewExtractorFromGRPCModel(e)
		}
	}

	if haveHTTPTplMatcher {
		for i, m := range req.GetMatchers() {
			httpTplMatcher[i] = httptpl.NewMatcherFromGRPCModel(m)
		}
	}

	// 重试处理，通过taskid找到所有失败的发送包
	var iInput any
	httpPoolOpts := make([]mutate.HttpPoolConfigOption, 0)
	retryPayloadsMap := make(map[string][]string, 0) // key 是原始请求报文，value 是重试的payload，我们需要将重试的payload绑定回去
	// 这里可能会出现原始请求报文一样的情况，但是这样也是因为payload没有而导致的，例如{{repeat(10)}}

	if !isRetry {
		// 插入 {{repeat(n)}}的fuzz标签
		if req.GetRepeatTimes() > 0 {
			var buf bytes.Buffer
			buf.WriteString("{{repeat(" + fmt.Sprint(req.GetRepeatTimes()) + ")}}")
			buf.Write(rawRequest)
			rawRequest = buf.Bytes()
		}
		iInput = rawRequest
	} else {
		// 找到上次任务的包
		webFuzzerResponses, err := yakit.QueryWebFuzzerResponseWithoutPaging(s.GetProjectDatabase(), req.RetryTaskID)
		if err != nil {
			return err
		}
		failedResponses := make([]*yakit.WebFuzzerResponse, 0)
		for _, resp := range webFuzzerResponses {
			if !resp.OK {
				failedResponses = append(failedResponses, resp)
				retryPayloadsMap[resp.Request] = strings.Split(resp.Payload, ",")
			} else {
				respModel, err := resp.ToGRPCModel()
				if err != nil {
					log.Errorf("convert web fuzzer response to grpc model failed: %s", err)
					continue
				}
				feedbackResponse(respModel, true)
			}
		}

		// 回溯找到所有之前重试成功的包
		oldIDs := taskIDBackTrack(int64(req.RetryTaskID))
		oldSuccessResponses, err := yakit.QueryWebFuzzerResponseByTaskIDsWithOk(s.GetProjectDatabase(), oldIDs)
		if err != nil {
			log.Errorf("query old web fuzzer succes response failed: %s", err)
		}

		for _, resp := range oldSuccessResponses {
			respModel, err := resp.ToGRPCModel()
			if err != nil {
				log.Errorf("convert web fuzzer response to grpc model failed: %s", err)
				continue
			}
			feedbackResponse(respModel, true)
		}

		if len(failedResponses) == 0 {
			return utils.Errorf("no failed web fuzzer request found")
		}

		iInput = lo.Map(failedResponses, func(i *yakit.WebFuzzerResponse, _ int) []byte {
			return utils.UnsafeStringToBytes(i.Request)
		})
	}

	var requestCount = 0
	if req.GetForceOnlyOneResponse() {
		requestCount = 1
	}

	fuzzerRequestSwg := utils.NewSizedWaitGroup(int(concurrent))
	executeBatchRequestsWithParams := func(mergedParams map[string]any) (retErr error) {
		defer func() {
			if err := recover(); err != nil {
				retErr = utils.Errorf("panic from grpc.httpfuzzer executeBatchRequestsWithParams: %v", err)
				utils.Debug(func() {
					utils.PrintCurrentGoroutineRuntimeStack()
				})
			}
		}()

		httpPoolOpts = append(httpPoolOpts,
			mutate.WithPoolOpt_FuzzParams(mergedParams),
			mutate.WithPoolOpt_ForceFuzzfile(req.GetForceFuzz()),
			mutate.WithPoolOpt_ExtraFuzzOptions(extraOpt...),
			mutate.WithPoolOpt_Timeout(timeoutSeconds),
			mutate.WithPoolOpt_Proxy(proxies...),
			mutate.WithPoolOpt_BatchTarget(batchTarget),
			//mutate.WithPoolOpt_Concurrent(int(concurrent)),
			mutate.WithPoolOpt_SizedWaitGroup(fuzzerRequestSwg),
			mutate.WithPoolOpt_Addr(req.GetActualAddr(), req.GetIsHTTPS()),
			mutate.WithPoolOpt_RawMode(true),
			mutate.WithPoolOpt_Https(req.GetIsHTTPS()),
			mutate.WithPoolOpt_GmTLS(req.GetIsGmTLS()),
			mutate.WithPoolOpt_Context(stream.Context()),
			mutate.WithPoolOpt_NoFollowRedirect(req.GetNoFollowRedirect()),
			mutate.WithPoolOpt_FollowJSRedirect(req.GetFollowJSRedirect()),
			mutate.WithPoolOpt_RedirectTimes(int(req.GetRedirectTimes())),
			mutate.WithPoolOpt_noFixContentLength(req.GetNoFixContentLength()),
			//mutate.WithPoolOpt_ExtraMutateConditionGetter(yak.MutateWithParamsGetter(req.GetHotPatchCodeWithParamGetter())),
			//mutate.WithPoolOpt_ExtraMutateCondition(yak.MutateWithYaklang(req.GetHotPatchCode())),
			mutate.WithPoolOpt_DelayMinSeconds(req.GetDelayMinSeconds()),
			mutate.WithPoolOPt_DelayMaxSeconds(req.GetDelayMaxSeconds()),
			mutate.WithPoolOpt_HookCodeCaller(yak.MutateHookCaller(req.GetHotPatchCode())),
			mutate.WithPoolOpt_Source("webfuzzer"),
			mutate.WithPoolOpt_RetryTimes(int(req.GetMaxRetryTimes())),
			mutate.WithPoolOpt_RetryInStatusCode(inStatusCode),
			mutate.WithPoolOpt_RetryNotInStatusCode(notInStatusCode),
			mutate.WithPoolOpt_RetryWaitTime(req.GetRetryWaitSeconds()),
			mutate.WithPoolOpt_RetryMaxWaitTime(req.GetRetryMaxWaitSeconds()),
			mutate.WithPoolOpt_DNSServers(req.GetDNSServers()),
			mutate.WithPoolOpt_EtcHosts(req.GetEtcHosts()),
			mutate.WithPoolOpt_NoSystemProxy(req.GetNoSystemProxy()),
			mutate.WithPoolOpt_FuzzParams(mergedParams),
			mutate.WithPoolOpt_RequestCountLimiter(requestCount))

		if isRetry {
			// 重试的时候，不需要渲染fuzztag
			httpPoolOpts = append(httpPoolOpts, mutate.WithPoolOpt_ForceFuzz(false))
		} else {
			httpPoolOpts = append(httpPoolOpts, mutate.WithPoolOpt_ForceFuzz(req.GetForceFuzz()))
		}
		res, err := mutate.ExecPool(
			iInput,
			httpPoolOpts...,
		)
		if err != nil {
			task.Ok = false
			task.Reason = utils.Errorf("exec http pool failed: %s", err).Error()
			return err
		}
		// 可以用于计算相似度
		var firstHeader, firstBody []byte
		for result := range res {
			task.HTTPFlowTotal++
			var payloads []string
			if !isRetry {
				payloads = make([]string, len(result.Payloads))
				for i, payload := range result.Payloads {
					if len(payload) > 100 {
						payload = payload[:100] + "..."
					}
					payloads[i] = utils.ParseStringToVisible(payload)
				}
			} else {
				payloads, _ = retryPayloadsMap[utils.UnsafeBytesToString(result.RequestRaw)]
			}

			var extractorResults []*ypb.KVPair

			if result != nil && result.ExtraInfo != nil {
				for k, v := range result.ExtraInfo {
					extractorResults = append(extractorResults, &ypb.KVPair{Key: utils.EscapeInvalidUTF8Byte([]byte(k)), Value: utils.EscapeInvalidUTF8Byte([]byte(v))})
				}
			}

			if result.Error != nil {
				rsp := &ypb.FuzzerResponse{}
				rsp.RequestRaw = result.RequestRaw
				rsp.UUID = uuid.NewV4().String()
				rsp.Url = utils.EscapeInvalidUTF8Byte([]byte(result.Url))
				rsp.Ok = false
				rsp.Reason = result.Error.Error()
				rsp.TaskId = int64(taskId)
				rsp.Payloads = payloads
				if result.LowhttpResponse != nil && result.LowhttpResponse.TraceInfo != nil {
					rsp.TotalDurationMs = result.LowhttpResponse.TraceInfo.TotalTime.Milliseconds()
					rsp.DurationMs = result.LowhttpResponse.TraceInfo.ServerTime.Milliseconds()
					rsp.FirstByteDurationMs = result.LowhttpResponse.TraceInfo.ServerTime.Milliseconds()
					rsp.DNSDurationMs = result.LowhttpResponse.TraceInfo.DNSTime.Milliseconds()
					rsp.Proxy = result.LowhttpResponse.Proxy
					rsp.RemoteAddr = result.LowhttpResponse.RemoteAddr
				}

				task.HTTPFlowFailedCount++
				yakit.SaveWebFuzzerResponse(s.GetProjectDatabase(), int(task.ID), rsp)
				_ = feedbackResponse(rsp, false)
				continue
			}

			if haveHTTPTplExtractor {
				var params = make(map[string]any)
				for _, extractor := range httpTplExtractor {
					vars, err := extractor.Execute(result.ResponseRaw, params)
					if err != nil {
						log.Errorf("extractor execute failed: %s", err)
						continue
					}
					for k, v := range vars {
						params[k] = v
						extractorResults = append(extractorResults, &ypb.KVPair{Key: k, Value: httptpl.ExtractResultToString(v)})
					}
				}
			}
			extractorResultsOrigin := extractorResults
			for k, v := range mergedParams {
				extractorResults = append(extractorResults, &ypb.KVPair{
					Key: k, Value: utils.EscapeInvalidUTF8Byte(codec.AnyToBytes(v))},
				)
			}

			var httpTPLmatchersResult bool
			if haveHTTPTplMatcher && result.LowhttpResponse != nil {
				cond := "and"
				switch ret := strings.ToLower(req.GetMatchersCondition()); ret {
				case "or", "and":
					cond = ret
				default:
				}
				ins := &httptpl.YakMatcher{
					SubMatcherCondition: cond,
					SubMatchers:         httpTplMatcher,
				}
				matcherParams := utils.CopyMapInterface(mergedParams)
				for _, kv := range extractorResultsOrigin {
					matcherParams[kv.GetKey()] = kv.GetValue()
				}
				httpTPLmatchersResult, err = ins.Execute(result.LowhttpResponse, matcherParams)
				if finalError != nil {
					log.Errorf("httptpl.YakMatcher execute failed: %s", err)
				}
			}

			_, body := lowhttp.SplitHTTPHeadersAndBodyFromPacket(result.ResponseRaw)
			if len(body) > 2*1024*1024 {
				body = body[:2*1024*1024]
				body = append(body, []byte("...chunked by yakit web fuzzer")...)
			}

			if !req.GetNoFixContentLength() && (result.Request != nil && result.Request.ProtoMajor != 2) { // no fix for h2 rsp
				result.ResponseRaw = lowhttp.ReplaceHTTPPacketBody(result.ResponseRaw, body, false)
				result.Response, _ = lowhttp.ParseStringToHTTPResponse(string(result.ResponseRaw))
			}

			if len(result.RequestRaw) > 1*1024*1024 {
				result.RequestRaw = result.RequestRaw[:1*1024*1024]
				result.RequestRaw = append(result.RequestRaw, []byte("...chunked by yakit web fuzzer")...)
			}

			task.HTTPFlowSuccessCount++
			var rsp = &ypb.FuzzerResponse{
				Url:                   utils.EscapeInvalidUTF8Byte([]byte(result.Url)),
				Method:                utils.EscapeInvalidUTF8Byte([]byte(result.Request.Method)),
				ResponseRaw:           result.ResponseRaw,
				GuessResponseEncoding: Chardet(result.ResponseRaw),
				RequestRaw:            result.RequestRaw,
				Payloads:              payloads,
				IsHTTPS:               strings.HasPrefix(strings.ToLower(result.Url), "https://"),
				ExtractedResults:      extractorResults,
				MatchedByMatcher:      httpTPLmatchersResult,
				HitColor:              req.GetHitColor(),
			}

			redirectPacket := result.LowhttpResponse.RedirectRawPackets
			if result.LowhttpResponse != nil {
				// redirect
				for _, f := range redirectPacket {
					rsp.RedirectFlows = append(rsp.RedirectFlows, &ypb.RedirectHTTPFlow{
						IsHttps:  f.IsHttps,
						Request:  f.Request,
						Response: f.Response,
					})
				}
			}

			// 处理额外时间
			if result.LowhttpResponse != nil && result.LowhttpResponse.TraceInfo != nil {
				rsp.TotalDurationMs = result.LowhttpResponse.TraceInfo.TotalTime.Milliseconds()
				rsp.DurationMs = result.LowhttpResponse.TraceInfo.ServerTime.Milliseconds()
				rsp.FirstByteDurationMs = result.LowhttpResponse.TraceInfo.ServerTime.Milliseconds()
				rsp.DNSDurationMs = result.LowhttpResponse.TraceInfo.DNSTime.Milliseconds()
				rsp.Proxy = result.LowhttpResponse.Proxy
				rsp.RemoteAddr = result.LowhttpResponse.RemoteAddr
			}
			if rsp.ResponseRaw != nil {
				// 处理结果，相似度
				header, body := lowhttp.SplitHTTPHeadersAndBodyFromPacket(rsp.ResponseRaw)
				if firstHeader == nil {
					log.Debugf("start to set first header[%v]...", result.Url)
					firstHeader = []byte(header)
					rsp.HeaderSimilarity = 1.0
				} else {
					rsp.HeaderSimilarity = utils.CalcSimilarity(firstHeader, []byte(header))
				}

				if firstBody == nil {
					log.Debugf("start to set first body[%v]...", result.Url)
					firstBody = body
					rsp.BodySimilarity = 1.0
				} else {
					rsp.BodySimilarity = utils.CalcSimilarity(firstBody, body)
				}
			}

			rsp.UUID = uuid.NewV4().String()
			rsp.Timestamp = result.Timestamp
			rsp.DurationMs = result.DurationMs
			rsp.Host = utils.EscapeInvalidUTF8Byte([]byte(result.Request.Header.Get("Host")))
			if rsp.Host == "" {
				rsp.Host = result.Request.Host
			}
			rsp.Host = utils.EscapeInvalidUTF8Byte([]byte(utils.ParseStringToVisible(result.Request.Host)))

			if result.Response != nil {
				rsp.Ok = true
				rsp.StatusCode = int32(result.Response.StatusCode)
				rsp.ContentType = utils.ParseStringToVisible(result.Response.Header.Get("Content-Type"))
				var bodyLen int64 = 0
				if result.Response.Body != nil {
					raw, _ := ioutil.ReadAll(result.Response.Body)
					bodyLen = int64(len(raw))
				}
				rsp.BodyLength = bodyLen

				// 解析 Headers
				for k, vs := range result.Response.Header {
					for _, v := range vs {
						rsp.Headers = append(rsp.Headers, &ypb.HTTPHeader{
							Header: utils.ParseStringToVisible(k),
							Value:  utils.ParseStringToVisible(v),
						})
					}
				}
			}

			if rsp.StatusCode > 0 {
				// 通过长度过滤
				if minBody <= maxBody && (minBody > 0 || maxBody > 0) {
					if maxBody >= rsp.BodyLength && minBody <= rsp.BodyLength {
						rsp.MatchedByFilter = true
					}
				}

				// 通过 StatusCode 过滤
				if !rsp.MatchedByFilter {
					rsp.MatchedByFilter = includeStatusCodeFilter.Contains(int(rsp.StatusCode))
				}

				// rule
				if !rsp.MatchedByFilter && (len(regexps) > 0 || len(keywords) > 0) {
					if utils.MatchAnyOfRegexp(rsp.ResponseRaw, regexps...) {
						rsp.MatchedByFilter = true
					}
					if rsp.MatchedByFilter || utils.MatchAllOfRegexp(rsp.ResponseRaw, keywords...) {
						rsp.MatchedByFilter = true
					}
				}
			}
			// 自动重定向
			if !req.GetNoFollowRedirect() {

				for i := 0; i < len(redirectPacket)-1; i++ {
					redirectRes := redirectPacket[i].RespRecord
					method, _, _ := lowhttp.GetHTTPPacketFirstLine(redirectRes.RawRequest)
					var redirectRsp = &ypb.FuzzerResponse{
						Url:                   utils.EscapeInvalidUTF8Byte([]byte(redirectRes.Url)),
						Method:                utils.EscapeInvalidUTF8Byte([]byte(method)),
						ResponseRaw:           redirectRes.RawPacket,
						GuessResponseEncoding: Chardet(redirectRes.RawPacket),
						RequestRaw:            redirectRes.RawRequest,
						Payloads:              payloads,
						IsHTTPS:               redirectRes.Https,
						MatchedByMatcher:      httpTPLmatchersResult,
						HitColor:              req.GetHitColor(),
					}
					if redirectRes != nil && redirectRes.TraceInfo != nil {
						redirectRsp.TotalDurationMs = redirectRes.TraceInfo.TotalTime.Milliseconds()
						redirectRsp.DurationMs = redirectRes.TraceInfo.ServerTime.Milliseconds()
						redirectRsp.FirstByteDurationMs = redirectRes.TraceInfo.ServerTime.Milliseconds()
						redirectRsp.DNSDurationMs = redirectRes.TraceInfo.DNSTime.Milliseconds()
						redirectRsp.Proxy = redirectRes.Proxy
						redirectRsp.RemoteAddr = redirectRes.RemoteAddr
					}
					redirectRsp.UUID = uuid.NewV4().String()
					redirectRsp.Timestamp = result.Timestamp
					redirectRsp.DurationMs = result.DurationMs
					redirectRsp.Host = utils.EscapeInvalidUTF8Byte([]byte(lowhttp.GetHTTPPacketHeader(redirectRes.RawRequest, "Host")))

					if redirectRes.RawPacket != nil {
						redirectRsp.Ok = true
						redirectRsp.StatusCode = int32(lowhttp.GetStatusCodeFromResponse(redirectRes.RawPacket))
						redirectRsp.ContentType = utils.ParseStringToVisible(lowhttp.GetHTTPPacketHeader(redirectRes.RawPacket, "Content-Type"))
						var bodyLen int64 = 0
						if lowhttp.GetHTTPPacketBody(redirectRes.RawPacket) != nil {
							bodyLen = int64(len(lowhttp.GetHTTPPacketBody(redirectRes.RawPacket)))
						}
						redirectRsp.BodyLength = bodyLen

						// 解析 Headers
						for k, vs := range lowhttp.GetHTTPPacketHeaders(redirectRes.RawPacket) {
							for _, v := range vs {
								redirectRsp.Headers = append(redirectRsp.Headers, &ypb.HTTPHeader{
									Header: utils.ParseStringToVisible(k),
									Value:  utils.ParseStringToVisible(v),
								})
							}
						}
					}

					if redirectRsp.StatusCode > 0 {
						// 通过长度过滤
						if minBody <= maxBody && (minBody > 0 || maxBody > 0) {
							if maxBody >= redirectRsp.BodyLength && minBody <= redirectRsp.BodyLength {
								redirectRsp.MatchedByFilter = true
							}
						}

						// 通过 StatusCode 过滤
						if !redirectRsp.MatchedByFilter {
							redirectRsp.MatchedByFilter = includeStatusCodeFilter.Contains(int(redirectRsp.StatusCode))
						}

						// rule
						if !redirectRsp.MatchedByFilter && (len(regexps) > 0 || len(keywords) > 0) {
							if utils.MatchAnyOfRegexp(redirectRsp.ResponseRaw, regexps...) {
								redirectRsp.MatchedByFilter = true
							}
							if redirectRsp.MatchedByFilter || utils.MatchAllOfRegexp(redirectRsp.ResponseRaw, keywords...) {
								redirectRsp.MatchedByFilter = true
							}
						}
					}
					yakit.SaveWebFuzzerResponse(s.GetProjectDatabase(), int(task.ID), redirectRsp)
					rsp.TaskId = int64(taskId)
					err := feedbackResponse(redirectRsp, false)
					if err != nil {
						log.Errorf("send to client failed: %s", err)
						continue
					}
				}
				//如果重定向了,修正最后一个req
				if len(redirectPacket) > 0 {
					rsp.RequestRaw = redirectPacket[len(redirectPacket)-1].Request
				}
			}
			yakit.SaveWebFuzzerResponse(s.GetProjectDatabase(), int(task.ID), rsp)
			rsp.TaskId = int64(taskId)
			err := feedbackResponse(rsp, false)
			if err != nil {
				log.Errorf("send to client failed: %s", err)
				continue
			}
		}
		return nil
	}

	/*
		handle vars
	*/
	wg := new(sync.WaitGroup)
	var mergedErr = make(chan error)
	for _param := range s.PreRenderVariables(stream.Context(), req.GetParams(), req.GetIsHTTPS(), req.GetIsGmTLS()) {
		mergedParams := _param
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := executeBatchRequestsWithParams(mergedParams)
			if err != nil {
				mergedErr <- err
			}
		}()
	}
	go func() {
		wg.Wait()
		close(mergedErr)
	}()

	errFilter := filter2.NewFilter()
	var errBuf bytes.Buffer
	for retErr := range mergedErr {
		h := codec.Sha256(retErr.Error())
		if errFilter.Exist(h) {
			continue
		}
		errFilter.Insert(h)
		errBuf.WriteString(retErr.Error())
		errBuf.WriteString("\n")
	}

	if errBuf.Len() > 0 {
		task.Ok = false
		task.Reason = errBuf.String()
		return utils.Errorf("execute batch requests failed: %s", errBuf.String())
	}
	task.Ok = true
	task.Reason = "normal exit / user canceled"
	return nil
}

var requestToMutateResult = func(reqs []*http.Request, chunked bool) (*ypb.MutateResult, error) {
	var raws [][]byte
	for _, r := range reqs {
		if chunked {
			r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
			r.Header.Set("Accept-Encoding", "*/*")
			r.Close = true
			urlIns, err := lowhttp.ExtractURLFromHTTPRequest(r, false)
			if err != nil {
				log.Errorf("extract url from httprequest: %v", err)
			}

			if urlIns != nil {
				r.URL = urlIns
			} else {
				r.URL, err = url.Parse(fmt.Sprintf("http://%v", r.Header.Get("Host")))
				if err != nil {
					log.Errorf("fallback generate url failed: %s", err)
				}
			}
			reqRaw, err := utils.DumpHTTPRequest(r, true)
			if err != nil {
				log.Errorf("dump with transfer encoding failed: %s", err)
			}
			if len(reqRaw) > 0 {
				raws = append(raws, lowhttp.FixHTTPRequest(reqRaw))
			}
			continue
		}
		reqRaw, _ := utils.HttpDumpWithBody(r, true)
		if len(reqRaw) > 0 {
			raws = append(raws, reqRaw)
		}
	}

	if raws != nil && len(raws) > 1 {
		return &ypb.MutateResult{
			Result:       raws[0],
			ExtraResults: raws[1:],
		}, nil
	}

	if raws != nil && len(raws) == 1 {
		return &ypb.MutateResult{
			Result: raws[0],
		}, nil
	}

	return nil, utils.Errorf("empty result")
}

func (s *Server) HTTPRequestMutate(ctx context.Context, req *ypb.HTTPRequestMutateParams) (*ypb.MutateResult, error) {
	freq, err := mutate.NewFuzzHTTPRequest(lowhttp.TrimLeftHTTPPacket(req.Request))
	if err != nil {
		return nil, utils.Errorf("build fuzzer request failed: %s", err)
	}

	switch strings.Join(req.FuzzMethods, "") {
	case "POST":
		u, _ := lowhttp.ExtractURLFromHTTPRequestRaw(req.GetRequest(), true)
		if u != nil {
			reqs, _ := freq.FuzzMethod(
				"POST",
			).FuzzGetParamsRaw(
				"",
			).FuzzHTTPHeader(
				"Content-Type", "application/x-www-form-urlencoded",
			).FuzzHTTPHeader(
				"Transfer-Encoding", "",
			).FuzzHTTPHeader(
				"User-Agent", consts.DefaultUserAgent,
			).FuzzPostRaw(
				u.RawQuery,
			).Results()
			if len(reqs) > 0 {
				return requestToMutateResult(reqs, false)
			}
		}
	case "HEAD":
		fallthrough
	case "GET":
		u, _ := lowhttp.ExtractURLFromHTTPRequestRaw(req.GetRequest(), true)
		_, body := lowhttp.SplitHTTPHeadersAndBodyFromPacket(req.GetRequest())
		if u != nil {
			var params = make(url.Values)
			values, _ := url.ParseQuery(u.RawQuery)
			if values != nil {
				for k, v := range values {
					params[k] = v
				}
			}
			postValue, _ := url.ParseQuery(strings.TrimSpace(string(body)))
			if postValue != nil {
				for k, v := range postValue {
					params[k] = v
				}
			}

			reqs, _ := freq.FuzzMethod(
				strings.ToUpper(strings.Join(req.GetFuzzMethods(), "")),
			).FuzzPath(
				u.Path,
			).FuzzHTTPHeader(
				"Content-Type", "",
			).FuzzHTTPHeader(
				"Transfer-Encoding", "",
			).FuzzGetParamsRaw(params.Encode()).FuzzHTTPHeader(
				"User-Agent", consts.DefaultUserAgent,
			).FuzzPostRaw("").Results()
			if len(reqs) > 0 {
				return requestToMutateResult(reqs, false)
			}
		}
	}

	if len(req.FuzzMethods) > 0 {
		reqs, err := freq.FuzzMethod(req.FuzzMethods...).Results()
		if err != nil {
			return nil, err
		}
		return requestToMutateResult(reqs, false)
	}

	// 获取 body
	reqInstance, err := freq.GetOriginHTTPRequest()
	if err != nil {
		return nil, err
	}
	bodyRaw, err := ioutil.ReadAll(reqInstance.Body)
	if err != nil {
		return nil, err
	}
	if bodyRaw == nil {
		return nil, utils.Errorf("empty body")
	}

	// 获取 chunk encode
	if req.ChunkEncode {
		_, body := lowhttp.SplitHTTPHeadersAndBodyFromPacket(req.GetRequest())
		reqRaw := lowhttp.ReplaceHTTPPacketBody(req.GetRequest(), body, true)
		return &ypb.MutateResult{Result: reqRaw}, nil
	}

	if req.UploadEncode {
		freq, err := mutate.NewFuzzHTTPRequest(lowhttp.TrimLeftHTTPPacket(req.Request))
		if err != nil {
			return nil, utils.Errorf("build fuzz.HTTPRequest failed: %s", err)
		}
		kv := make(map[string]interface{})
		for _, getParam := range freq.GetGetQueryParams() {
			kv[getParam.Name()] = getParam.Value()
		}

		for _, postRawParam := range freq.GetPostParams() {
			kv[postRawParam.Name()] = postRawParam.Value()
		}

		currentPair := freq.FuzzMethod("POST")
		for k, v := range kv {
			currentPair = currentPair.FuzzUploadKVPair(k, v)
		}
		reqs, err := currentPair.Results()
		if err != nil {
			return nil, err
		}
		return requestToMutateResult(reqs, false)
	}

	return &ypb.MutateResult{
		Result:       []byte(req.Request),
		ExtraResults: nil,
	}, nil
}

func (s *Server) HTTPResponseMutate(ctx context.Context, req *ypb.HTTPResponseMutateParams) (*ypb.MutateResult, error) {
	return nil, nil
}

func (s *Server) QueryHistoryHTTPFuzzerTask(ctx context.Context, req *ypb.Empty) (*ypb.HistoryHTTPFuzzerTasks, error) {
	return &ypb.HistoryHTTPFuzzerTasks{Tasks: yakit.QueryFirst50WebFuzzerTask(s.GetProjectDatabase())}, nil
}

func (s *Server) QueryHistoryHTTPFuzzerTaskEx(ctx context.Context, req *ypb.QueryHistoryHTTPFuzzerTaskExParams) (*ypb.HistoryHTTPFuzzerTasksResponse, error) {
	paging, tasks, err := yakit.QueryFuzzerHistoryTasks(s.GetProjectDatabase(), req)
	if err != nil {
		return nil, err
	}
	newTasks := funk.Map(tasks, func(i *yakit.WebFuzzerTask) *ypb.HistoryHTTPFuzzerTaskDetail {
		return i.ToSwaggerModelDetail()
	}).([]*ypb.HistoryHTTPFuzzerTaskDetail)
	return &ypb.HistoryHTTPFuzzerTasksResponse{
		Data:       newTasks,
		Total:      int64(paging.TotalRecord),
		TotalPage:  int64(paging.TotalPage),
		Pagination: req.GetPagination(),
	}, nil
}

func (s *Server) GetHistoryHTTPFuzzerTask(ctx context.Context, req *ypb.GetHistoryHTTPFuzzerTaskRequest) (*ypb.HistoryHTTPFuzzerTaskDetail, error) {
	task, err := yakit.GetWebFuzzerTaskById(s.GetProjectDatabase(), int(req.GetId()))
	if err != nil {
		return nil, err
	}
	var reqRaw ypb.FuzzerRequest
	err = json.Unmarshal([]byte(task.RawFuzzTaskRequest), &reqRaw)
	if err != nil {
		return nil, err
	}
	return &ypb.HistoryHTTPFuzzerTaskDetail{
		BasicInfo:     task.ToSwaggerModel(),
		OriginRequest: &reqRaw,
	}, nil
}

func (s *Server) QueryHTTPFuzzerResponseByTaskIdRequest(ctx context.Context, req *ypb.QueryHTTPFuzzerResponseByTaskIdRequest) (*ypb.QueryHTTPFuzzerResponseByTaskIdResponse, error) {
	p, rets, err := yakit.QueryWebFuzzerResponse(s.GetProjectDatabase(), req)
	if err != nil {
		return nil, err
	}

	var results []*ypb.FuzzerResponse
	for _, i := range rets {
		r, err := i.ToGRPCModel()
		if err != nil {
			continue
		}
		results = append(results, r)
	}

	return &ypb.QueryHTTPFuzzerResponseByTaskIdResponse{
		Pagination: req.Pagination,
		Data:       results,
		Total:      int64(p.TotalRecord),
		TotalPage:  int64(p.TotalPage),
	}, nil
}

func (s *Server) ExtractHTTPResponse(ctx context.Context, req *ypb.ExtractHTTPResponseParams) (*ypb.ExtractHTTPResponseResult, error) {
	if req.GetHTTPResponse() == "" {
		return nil, utils.Error("empty http response")
	}

	if len(req.GetExtractors()) == 0 {
		return nil, utils.Error("empty extractors")
	}

	/*
		type YakExtractor struct {
			Name             string // name or index
			Type             string
			Scope            string // header body all
			Groups           []string
			RegexpMatchGroup []int
			XPathAttribute   string
		}
	*/
	extractors := funk.Map(req.GetExtractors(), func(i *ypb.HTTPResponseExtractor) *httptpl.YakExtractor {
		return httptpl.NewExtractorFromGRPCModel(i)
	}).([]*httptpl.YakExtractor)

	var params = make(map[string]interface{})
	for _, i := range extractors {
		p, err := i.Execute([]byte(req.GetHTTPResponse()), params)
		if err != nil {
			log.Errorf("extractor %s execute failed: %s", i.Name, err)
			continue
		}
		for k, v := range p {
			params[k] = httptpl.ExtractResultToString(v)
		}
	}

	var results []*ypb.FuzzerParamItem
	for k, v := range params {
		results = append(results, &ypb.FuzzerParamItem{
			Key:   k,
			Value: httptpl.ExtractResultToString(v),
		})
	}
	return &ypb.ExtractHTTPResponseResult{Values: results}, nil
}

func (s *Server) MatchHTTPResponse(ctx context.Context, req *ypb.MatchHTTPResponseParams) (*ypb.MatchHTTPResponseResult, error) {
	if req.GetHTTPResponse() == "" {
		return nil, utils.Error("empty http response")
	}

	if len(req.GetMatchers()) == 0 {
		return nil, utils.Error("empty matchers")
	}

	/*
		type YakMatcher struct {
			MatcherType string
			// just for expr
			ExprType string

			// groups
			Scope         string
			Condition     string
			Group         []string
			GroupEncoding string

			Negative bool

			// or / and
			SubMatcherCondition string
			SubMatchers         []*YakMatcher
		}
	*/
	matchers := funk.Map(req.GetMatchers(), func(i *ypb.HTTPResponseMatcher) *httptpl.YakMatcher {
		return &httptpl.YakMatcher{
			MatcherType:   i.GetMatcherType(),
			ExprType:      i.GetExprType(),
			Scope:         i.GetScope(),
			Condition:     i.GetCondition(),
			Group:         i.GetGroup(),
			GroupEncoding: i.GetGroupEncoding(),
			Negative:      i.GetNegative(),
		}
	}).([]*httptpl.YakMatcher)

	matcher := &httptpl.YakMatcher{
		SubMatcherCondition: req.GetMatcherCondition(),
		SubMatchers:         matchers,
	}
	if matcher.SubMatcherCondition == "" {
		matcher.SubMatcherCondition = "and"
	}

	result, err := matcher.ExecuteRawResponse([]byte(req.GetHTTPResponse()), nil)
	if err != nil {
		return nil, err
	}
	return &ypb.MatchHTTPResponseResult{Matched: result}, nil
}

func (s *Server) RenderVariables(ctx context.Context, req *ypb.RenderVariablesRequest) (*ypb.RenderVariablesResponse, error) {
	vars := httptpl.NewVars()
	for _, kv := range req.GetParams() {
		vars.AutoSet(kv.GetKey(), kv.GetValue())
	}
	var results = vars.ToMap()
	var finalResults []*ypb.KVPair
	for _, kv := range req.GetParams() {
		value, ok := results[kv.GetKey()]
		if !ok {
			continue
		}
		finalResults = append(finalResults, &ypb.KVPair{
			Key:   kv.GetKey(),
			Value: utils.EscapeInvalidUTF8Byte(utils.InterfaceToBytes(value)),
		})
	}

	var responseVars []*ypb.KVPair
	for k, v := range httptpl.LoadVarFromRawResponse(req.GetHTTPResponse(), 0) {
		responseVars = append(responseVars, &ypb.KVPair{
			Key:   k,
			Value: utils.EscapeInvalidUTF8Byte(utils.InterfaceToBytes(v)),
		})
	}
	sort.SliceStable(responseVars, func(i, j int) bool {
		return responseVars[i].Key < responseVars[j].Key
	})
	finalResults = append(finalResults, responseVars...)
	return &ypb.RenderVariablesResponse{Results: finalResults}, nil
}

func (s *Server) RenderVariablesWithTypedKV(ctx context.Context, kvs []*ypb.FuzzerParamItem) map[string]any {
	vars := httptpl.NewVars()
	for _, kv := range kvs {
		key, value := kv.GetKey(), kv.GetValue()
		if kv.GetType() == "nuclei-dsl" {
			vars.SetAsNucleiTags(key, value)
		} else {
			vars.Set(key, value)
		}
	}
	return vars.ToMap()
}

func (s *Server) PreRenderVariables(ctx context.Context, params []*ypb.FuzzerParamItem, https, gmtls bool) chan map[string]any {
	var resultsChan = make(chan map[string]any, 100)
	if len(params) <= 0 {
		resultsChan <- make(map[string]any)
		close(resultsChan)
		return resultsChan
	}

	l := make([][]string, len(params))
	idToParam := make(map[int]*ypb.FuzzerParamItem)
	hasNucleiTag := false

	for index, p := range params {
		_, value := p.GetKey(), p.GetValue()
		typ := strings.TrimSpace(strings.ToLower(p.GetType()))
		idToParam[index] = p

		if typ == "fuzztag" {
			rets, _ := mutate.FuzzTagExec(value, mutate.FuzzFileOptions()...)
			if len(rets) > 0 {
				l[index] = rets
				continue
			}
		} else if typ == "nuclei-dsl" {
			hasNucleiTag = true
		}

		l[index] = []string{value}
	}

	var count int64 = 0
	go func() {
		defer func() {
			if count <= 0 {
				resultsChan <- make(map[string]any)
				close(resultsChan)
				return
			}
			close(resultsChan)

			if err := recover(); err != nil {
				log.Errorf("cartesian to fuzztag vars failed: %v", err)
			}
		}()

		err := cartesian.ProductExContext(ctx, l, func(payloads []string) error {
			params := make([]*ypb.FuzzerParamItem, 0)
			resultMap := make(map[string]any)
			if hasNucleiTag {
				for index, v := range payloads {
					p := idToParam[index]
					key := p.GetKey()
					params = append(params, &ypb.FuzzerParamItem{Key: key, Value: v, Type: p.GetType()})
				}
				resultMap = s.RenderVariablesWithTypedKV(ctx, params)
			} else {
				for index, v := range payloads {
					p := idToParam[index]
					key := p.GetKey()
					resultMap[key] = v
				}
			}

			atomic.AddInt64(&count, 1)
			resultsChan <- resultMap
			return nil
		})
		if err != nil {
			log.Errorf("cartesian product failed: %s", err)
		}
	}()
	return resultsChan
}
func (s *Server) GetSystemDefaultDnsServers(ctx context.Context, req *ypb.Empty) (*ypb.DefaultDnsServerResponse, error) {
	servers, err := utils.GetSystemDnsServers()
	return &ypb.DefaultDnsServerResponse{DefaultDnsServer: servers}, err
}
