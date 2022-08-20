package decoder

func DecodeLogin(data []byte) (string, error) {
	imeiData := data[7:15]
	return decodeImei(imeiData)
}
