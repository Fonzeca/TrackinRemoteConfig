package decoder

func DecodeLogin(data []byte) string {
	imeiData := data[7:14]
	return decodeImei(imeiData)
}
