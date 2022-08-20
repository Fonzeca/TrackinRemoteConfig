package decoder

func DecodeLogin(data []byte) (string, error) {
	imeiData := data[7:14]
	return decodeImei(imeiData)
}
