package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/yaklang/yaklang/common/go-funk"
	"github.com/yaklang/yaklang/common/yak/yaklib/codec"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func PrettifyListFromStringSplited(Raw string, sep string) (targets []string) {
	targetsRaw := strings.Split(Raw, sep)
	for _, tRaw := range targetsRaw {
		r := strings.TrimSpace(tRaw)
		if len(r) > 0 {
			targets = append(targets, r)
		}
	}
	return
}

func PrettifyListFromStringSplitEx(Raw string, sep ...string) (targets []string) {
	if len(sep) <= 0 {
		sep = []string{",", "|"}
	}
	var targetsRaw []string
	for _, s := range sep {
		targetsRaw = append(targetsRaw, strings.Split(Raw, s)...)
	}
	for _, tRaw := range targetsRaw {
		r := strings.TrimSpace(tRaw)
		if len(r) > 0 {
			targets = append(targets, r)
		}
	}
	return
}

func ToLowerAndStrip(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func StringSliceContain(s interface{}, raw string) (result bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	haveResult := false
	funk.ForEach(s, func(i interface{}) {
		if haveResult {
			return
		}
		if InterfaceToString(i) == raw {
			haveResult = true
		}
	})
	return haveResult
}

func StringContainsAnyOfSubString(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func StringContainsAllOfSubString(s string, subs []string) bool {
	if len(subs) <= 0 {
		return false
	}
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

func IStringContainsAnyOfSubString(s string, subs []string) bool {
	for _, sub := range subs {
		if IContains(s, sub) {
			return true
		}
	}
	return false
}

func ConvertToStringSlice(raw ...interface{}) (r []string) {
	for _, e := range raw {
		r = append(r, fmt.Sprintf("%v", e))
	}
	return
}

func ChanStringToSlice(c chan string) (result []string) {
	for l := range c {
		result = append(result, l)
	}
	return
}

var (
	cStyleCharPRegexp, _ = regexp.Compile(`\\((x[0-9abcdef]{2})|([0-9]{1,3}))`)
)

func ParseCStyleBinaryRawToBytes(raw []byte) []byte {
	// like "\\x12" => "\x12"
	return cStyleCharPRegexp.ReplaceAllFunc(raw, func(i []byte) []byte {
		if bytes.HasPrefix(i, []byte("\\x")) {
			if len(i) == 4 {
				rawChar := string(i[2:])
				charInt, err := strconv.ParseInt("0x"+string(rawChar), 0, 16)
				if err != nil {
					return i
				}
				return []byte{byte(charInt)}
			} else {
				return i
			}
		} else if bytes.HasPrefix(raw, []byte("\\")) {
			if len(i) > 1 && len(i) <= 4 {
				rawChar := string(i[1:])
				charInt, err := strconv.ParseInt(string(rawChar), 10, 8)
				if err != nil {
					return i
				}
				return []byte{byte(charInt)}
			} else {
				return i
			}
		}
		return i
	})
}

var GbkToUtf8 = codec.GbkToUtf8
var Utf8ToGbk = codec.Utf8ToGbk

func ParseStringToVisible(raw interface{}) string {
	var s = InterfaceToString(raw)
	s = EscapeInvalidUTF8Byte([]byte(s))
	//s = strings.ReplaceAll(s, "\x20", "\\x20")
	s = strings.ReplaceAll(s, "\x0b", "\\v")
	r, err := regexp.Compile(`\s`)
	if err != nil {
		return s
	}
	return r.ReplaceAllStringFunc(s, func(s string) string {
		var result = strconv.Quote(s)
		for strings.HasPrefix(result, "\"") {
			result = result[1:]
		}
		for strings.HasSuffix(result, "\"") {
			result = result[:len(result)-1]
		}
		return result
	})
}

func EscapeInvalidUTF8Byte(s []byte) string {
	// 这个操作返回的结果和原始字符串是非等价的
	ret := make([]rune, 0, len(s)+20)
	start := 0
	for {
		r, size := utf8.DecodeRune(s[start:])
		if r == utf8.RuneError {
			// 说明是空的
			if size == 0 {
				break
			} else {
				// 不是 rune
				ret = append(ret, []rune(fmt.Sprintf("\\x%02x", s[start]))...)
			}
		} else {
			// 不是换行之类的控制字符
			if unicode.IsControl(r) && !unicode.IsSpace(r) {
				ret = append(ret, []rune(fmt.Sprintf("\\x%02x", r))...)
			} else {
				// 正常字符
				ret = append(ret, r)
			}
		}
		start += size
	}
	return string(ret)
}

var GBKSafeString = codec.GBKSafeString

func LastLine(s []byte) []byte {
	s = bytes.TrimSpace(s)
	scanner := bufio.NewScanner(bytes.NewReader(s))
	scanner.Split(bufio.ScanLines)

	var lastLine = s
	for scanner.Scan() {
		lastLine = scanner.Bytes()
	}

	return lastLine
}

func RemoveUnprintableChars(raw string) string {
	scanner := bufio.NewScanner(bytes.NewBufferString(raw))
	scanner.Split(bufio.ScanBytes)

	var buf = bytes.NewBufferString("")
	for scanner.Scan() {
		c := scanner.Bytes()[0]

		if c <= 0x7e && c >= 0x20 {
			buf.WriteByte(c)
		} else {
			buf.WriteString(`\x` + fmt.Sprintf("%02x", c))
		}
	}

	return buf.String()
}

func RemoveUnprintableCharsWithReplace(raw string, handle func(i byte) string) string {
	scanner := bufio.NewScanner(bytes.NewBufferString(raw))
	scanner.Split(bufio.ScanBytes)

	var r []byte
	for scanner.Scan() {
		c := scanner.Bytes()[0]

		if c <= 0x7e && c >= 0x20 {
			r = append(r, c)
		} else {
			r = append(r, []byte(handle(c))...)
		}
	}

	return string(r)
}

func RemoveUnprintableCharsWithReplaceItem(raw string) string {
	return RemoveUnprintableCharsWithReplace(raw, func(i byte) string {
		return fmt.Sprintf("__HEX_%v__", codec.EncodeToHex([]byte{i}))
	})
}

func RemoveRepeatedWithStringSlice(slice []string) []string {
	r := map[string]interface{}{}
	for _, s := range slice {
		r[s] = 1
	}

	var r2 []string
	for k, _ := range r {
		r2 = append(r2, k)
	}
	return r2
}

var (
	titleRegexp = regexp.MustCompile(`(?is)\<title\>(.*?)\</?title\>`)
)

func ExtractTitleFromHTMLTitle(s string, defaultValue string) string {
	var title string
	l := titleRegexp.FindString(s)
	if len(l) > 15 {
		title = EscapeInvalidUTF8Byte([]byte(l))[7 : len(l)-8]
	}
	titleRunes := []rune(title)
	if len(titleRunes) > 128 {
		title = string(titleRunes[0:128]) + "..."
	}

	if title == "" {
		return defaultValue
	}

	return title
}
