package codec

import (
	"bytes"
	"fmt"
	"github.com/yaklang/yaklang/common/log"
	"github.com/yaklang/yaklang/common/mimetype"
	"strings"
)

func _mimeIsText(depth int, t *mimetype.MIME) bool {
	if depth > 20 || t == nil {
		return false
	}
	if strings.HasPrefix(t.String(), "text/") || strings.HasPrefix(t.String(), "Text/") {
		return true
	}
	return _mimeIsText(depth+1, t.Parent())
}

func mimeIsText(t *mimetype.MIME) bool {
	return _mimeIsText(0, t)
}

type MIMEResult struct {
	MIMEType    string
	IsText      bool
	NeedCharset bool
	Charset     string
}

func (t *MIMEResult) TryUTF8Convertor(raw []byte) []byte {
	if strings.Contains(t.MIMEType, "/html") || strings.Contains(t.MIMEType, "/xhtml") {
		result := raw
		// <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
		// <meta charset="UTF-8">
		// <meta http-equiv="Content-Type" content="text/html; charset=gb2312">
		// <meta charset="gb2312">
		newBuffer := new(bytes.Buffer)
		lastStart := -1
		var encodings []PrescanResult
		var set = make(map[string]struct{})
		enc, origin := HtmlCharsetPrescan(result, func(start, end int, match PrescanResult) {
			if _, ok := set[match.Name]; !ok {
				encodings = append(encodings, match)
				set[match.Name] = struct{}{}
			}
			if lastStart < 0 {
				newBuffer.Write(result[:start])
			} else {
				newBuffer.Write(result[lastStart:start])
			}
			newBuffer.WriteString("utf-8")
			lastStart = end
		})
		if strings.ToLower(origin) != "utf-8" {
			newBuffer.Write(result[lastStart:])
			result = newBuffer.Bytes()
		}

		if len(encodings) == 1 {
			decodedResult, err := enc.NewDecoder().Bytes(result)
			if err != nil {
				return result
			}
			return decodedResult
		}
		log.Warn("WARNING: ATTENTION multiple encodings, try the best")
		for _, v := range encodings {
			if v.Encoding != nil {
				decodeResult, err := v.Encoding.NewDecoder().Bytes(result)
				if err != nil {
					log.Infof("try encoding %#v failed: %v", v.Name, err)
					continue
				}
				return decodeResult
			}
		}
		return result
	}

	switch strings.ToLower(t.Charset) {
	case "gb18030", "gb-18030", "gbk", "gb2312", "gb-2312":
		result, err := GB18030ToUtf8(raw)
		if err != nil {
			return result
		}
		return result
	default:
		log.Warnf("TBD: charset %#v not supported yet, use origin raw input", t.Charset)
	}
	return raw
}

// MatchMIMEType will match via bytes
// note: if the raw input is overlarge, check the first n(4k) bytes to detect
// question: fix the tail, if the raw input is text (not structured file)
func MatchMIMEType(raw any) (*MIMEResult, error) {
	r := mimetype.Detect(interfaceToBytes(raw))
	if r == nil {
		return nil, fmt.Errorf("match(detect) mime type failed, check: %v", ShrinkString(fmt.Sprintf("%#v", raw), 64))
	}
	var result = &MIMEResult{
		MIMEType:    r.String(),
		IsText:      mimeIsText(r),
		NeedCharset: r.NeedCharset(),
		Charset:     r.Charset(),
	}
	return result, nil
}
