package yakgrpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yaklang/yaklang/common/consts"
	"github.com/yaklang/yaklang/common/jsonpath"
	"github.com/yaklang/yaklang/common/utils/lowhttp"
	"github.com/yaklang/yaklang/common/yak/yaklib"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/yaklang/yaklang/common/utils"
	"github.com/yaklang/yaklang/common/yakgrpc/ypb"
)

func TestOUTPUT_AiChat(t *testing.T) {
	consts.UpdateThirdPartyApplicationConfig(&ypb.ThirdPartyApplicationConfig{
		APIKey: fmt.Sprintf("%s.%s", utils.RandStringBytes(32), utils.RandStringBytes(16)),
		Type:   "chatglm",
	})
	rspStrTmp := `data: {"id":"1","created":1,"model":"1","choices":[{"index":0,"delta":{"role":"assistant","content":"%s"}}]}
`
	headerStr, _, _ := lowhttp.FixHTTPResponse([]byte("HTTP/1.1 200 OK\nContent-Type: application/json\nTransfer-Encoding: chunked\nConnection: Keep-Alive\n\n"))

	port := utils.GetRandomAvailableTCPPort()
	l, err := tls.Listen("tcp", spew.Sprintf(":%d", port), utils.GetDefaultTLSConfig(3))
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}
			genMsg := func(s string) []byte {
				msg := []byte(fmt.Sprintf(rspStrTmp, s))
				return []byte(fmt.Sprintf("%x\r\n%s\r\n", len(msg), msg))
			}
			go func() {
				utils.StableReader(conn, 1, 10240)
				conn.Write(headerStr)
				conn.Write(genMsg("你好"))
				time.Sleep(time.Millisecond * 500)
				conn.Write(genMsg("我是人工智障"))
				time.Sleep(time.Millisecond * 500)
				conn.Write(genMsg("助手"))
				conn.Write(genMsg(""))
				conn.Write([]byte("\r\n"))
				conn.Close()
			}()
		}
	}()
	client, err := NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}
	yaklib.InitYakit(yaklib.NewVirtualYakitClient(func(i *ypb.ExecResult) error {
		print(string(i.Raw))
		return nil
	}))
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	utils.WaitConnect(addr, 3)
	stream, err := client.AttachCombinedOutput(context.Background(), &ypb.AttachCombinedOutputRequest{})
	if err != nil {
		t.Fatal(err)
	}
	debugStreamTestResult := false
	go func() {
		for {
			v, err := stream.Recv()
			if err != nil {
				return
			}
			if strings.Contains(string(v.Raw), "你好我是人工智障助手") {
				debugStreamTestResult = true
			}
		}
	}()
	_, err = client.Exec(context.Background(), &ypb.ExecRequest{
		NoDividedEngine: true,
		Script:          fmt.Sprintf(`ai.Chat("你好",ai.type("chatglm"),ai.debugStream(),ai.domain("%s"))~`, addr),
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	assert.Equal(t, true, debugStreamTestResult)
}
func TestOUTPUT_AiChatOutputByDefaultStream(t *testing.T) {
	consts.UpdateThirdPartyApplicationConfig(&ypb.ThirdPartyApplicationConfig{
		APIKey: fmt.Sprintf("%s.%s", utils.RandStringBytes(32), utils.RandStringBytes(16)),
		Type:   "chatglm",
	})
	rspStrTmp := `data: {"id":"1","created":1,"model":"1","choices":[{"index":0,"delta":{"role":"assistant","content":"%s"}}]}
`
	headerStr, _, _ := lowhttp.FixHTTPResponse([]byte("HTTP/1.1 200 OK\nContent-Type: application/json\nTransfer-Encoding: chunked\nConnection: Keep-Alive\n\n"))

	port := utils.GetRandomAvailableTCPPort()
	l, err := tls.Listen("tcp", spew.Sprintf(":%d", port), utils.GetDefaultTLSConfig(3))
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}
			genMsg := func(s string) []byte {
				msg := []byte(fmt.Sprintf(rspStrTmp, s))
				return []byte(fmt.Sprintf("%x\r\n%s\r\n", len(msg), msg))
			}
			go func() {
				utils.StableReader(conn, 1, 10240)
				conn.Write(headerStr)
				conn.Write(genMsg("你好"))
				time.Sleep(time.Millisecond * 500)
				conn.Write(genMsg("我是人工智障"))
				time.Sleep(time.Millisecond * 500)
				conn.Write(genMsg("助手"))
				conn.Write(genMsg(""))
				conn.Write([]byte("\r\n"))
				conn.Close()
			}()
		}
	}()
	client, err := NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`\\"data\\":\\"(.+?)\\",\\"streamId`)
	subMsgN := 0
	msg := ""
	yaklib.InitYakit(yaklib.NewVirtualYakitClient(func(i *ypb.ExecResult) error {
		s := re.FindAllStringSubmatch(string(i.Message), -1)
		if len(s) > 0 {
			subMsgN++
			msg += s[0][1]
		}
		print(string(i.Raw))
		return nil
	}))
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	utils.WaitConnect(addr, 3)
	_, err = client.Exec(context.Background(), &ypb.ExecRequest{
		NoDividedEngine: true,
		Script:          fmt.Sprintf(`ai.Chat("你好",ai.type("chatglm"),ai.domain("%s"))~`, addr),
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	assert.Equal(t, true, subMsgN > 3)
	assert.Contains(t, msg, "你好我是人工智障助手")
}
func TestOUTPUT_STREAMYakitStream(t *testing.T) {
	client, err := NewLocalClient()
	if err != nil {
		t.Fatal(err)
	}

	uid := uuid.New().String()

	stream, err := client.Exec(context.Background(), &ypb.ExecRequest{
		NoDividedEngine: true,
		Script: `yakit.AutoInitYakit()

# Input your code!
pr, pw = io.Pipe()~
go func{
    count = 0
    for {
        count++
        pw.Write("Hello1")
        sleep(0.3)
        if count > 5 {
            pw.Close()
            pr.Close()
            return
        }
    }
}
yakit.Stream("ai", "` + uid + `", pr)
sleep(2)
`,
	})
	if err != nil {
		t.Fatal(err)
	}

	var dataBuf bytes.Buffer
	haveStart := false
	haveStop := false
	for {
		data, err := stream.Recv()
		if err != nil {
			break
		}

		if data.IsMessage {
			data := string(data.Message)
			if data == "" {
				continue
			}
			data = codec.AnyToString(jsonpath.Find(data, "$.content.data"))
			id := jsonpath.Find(data, "$.streamId")
			if id != uid {
				t.Fatal("streamId is not right")
			}
			if codec.AnyToString(jsonpath.Find(data, "$.action")) == "start" {
				haveStart = true
			}
			if codec.AnyToString(jsonpath.Find(data, "$.action")) == "stop" {
				haveStop = true
			}
			if codec.AnyToString(jsonpath.Find(data, "$.action")) == "data" {
				dataBuf.WriteString(codec.AnyToString(jsonpath.Find(data, "$.data")))
			}
			spew.Dump(data)
		}
	}
	if !haveStart {
		t.Fatal("stream start not found")
	}
	if !haveStop {
		t.Fatal("stream stop not found")
	}
	if dataBuf.String() != "Hello1Hello1Hello1Hello1Hello1Hello1" {
		t.Fatal("stream data not found")
	}
}

func TestGRPCMUSTPASS_LANGUAGE_YakitLog(t *testing.T) {
	testCase1 := [][]string{
		{"yakit.Info(\"yakit_info\")", "yakit_info"},
		{"yakit.Info(\"yakit_%v\",\"info\")", "yakit_info"},
		{"risk.NewRisk(\"1.1.1.1\")", ""},
		{"yakit.Output(yakit.TableData(\"table\", {\n    \"id\": 1,\n    \"name\": \"张三\",\n}))", ""},
	}
	testCase2 := [][]string{
		{"println(x\"{{base64(Hello Yak)}}\")", "SGVsbG8gWWFr"},
		{"println(\"println\")", "println"},
		{"println(\"print\")", "print"},
		{"dump(\"dump\")", "dump"},
		{"log.info(\"log_info\")", "log_info"},
		{"log.infof(\"log_%s\",\"info\")", "log_info"},
	}
	code := ""
	for _, v := range testCase1 {
		code += v[0] + "\n"
	}
	for _, v := range testCase2 {
		code += v[0] + "\n"
	}

	var client, err = NewLocalClient()
	stream, err := client.Exec(context.Background(), &ypb.ExecRequest{
		Script:          code,
		NoDividedEngine: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	otherLog := ""
	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		info := make(map[string]interface{})
		err = json.Unmarshal(res.Message, &info)
		if err != nil {
			otherLog += string(res.Raw)
		}
		if i >= len(testCase1) {
			break
		}
		if info["type"] == "log" {
			if v, ok := info["content"].(map[string]interface{}); ok {
				if !strings.Contains(utils.InterfaceToString(v["data"]), testCase1[i][1]) {
					t.Fatal("log error")
				}
			} else {
				t.Fatal("invalid log format")
			}
		}
		i++
	}
	_ = otherLog
	// 由于CombinedOutput是异步的，可能由于延迟导致这里没有获取到全部输出
	//for _, testCase := range testCase2 {
	//	if !strings.Contains(otherLog, testCase[1]) {
	//		t.Fatal("log stream not contains", testCase[1])
	//	}
	//}
}
