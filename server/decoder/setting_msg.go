package decoder

import (
	"bytes"
	"encoding/hex"
	"errors"
	"unicode/utf16"
	"unicode/utf8"
)

func Decode(data []byte) (string, string, error) {
	if len(data) <= 16 {
		return "", "", errors.New("Data length corto")
	}
	imeiData := data[7:15]
	imei, err := decodeImei(imeiData)
	if err != nil {
		return "", "", err
	}

	contentData := data[16:]
	content, err := decodeContent(contentData)
	if err != nil {
		return "", "", err
	}

	return imei, content, nil
}

func decodeImei(data []byte) (string, error) {
	if len(data) != 8 {
		return "", errors.New("Imei tiene que ser 8 de lenght")
	}
	imei := hex.EncodeToString(data)

	imei = imei[1:]
	return imei, nil
}

func decodeContent(data []byte) (string, error) {

	if len(data)%2 != 0 {
		return "", errors.New("El content no es multiplo de 2")
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

	return ret.String(), nil
}
