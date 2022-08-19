package decoder

import (
	"bytes"
	"encoding/hex"
	"unicode/utf16"
	"unicode/utf8"
)

func Decode(data []byte) (string, string) {
	if len(data) <= 16 {
		return "error", "error"
	}
	imeiData := data[7:14]
	imei := decodeImei(imeiData)

	contentData := data[16:]
	content := decodeContent(contentData)

	return imei, content
}

func decodeImei(data []byte) string {
	if len(data) != 8 {
		return ""
	}
	imei := hex.EncodeToString(data)

	imei = imei[1:]
	return imei
}

func decodeContent(data []byte) string {

	if len(data)%2 != 0 {
		return "error"
	}

	u16s := make([]uint16, 1)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(data)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(data[i]) + (uint16(data[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String()
}
