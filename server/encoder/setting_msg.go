package encoder

import (
	"encoding/hex"
	"strings"
)

var HEADER = []byte{0x25, 0x25}

var MessageType = []byte{0x81}

var PacketLenght = []byte{0x00, 0x10}

var SerialNumber = []byte{0x00, 0x01}

//-----------------

var FunctionType = []byte{0x01}

func Encode(imei string, content string) []byte {
	command := []byte{}
	command = append(command, HEADER...)
	command = append(command, MessageType...)
	command = append(command, PacketLenght...)
	command = append(command, SerialNumber...)
	command = append(command, EncodeIMEI(imei)...)
	command = append(command, FunctionType...)
	command = append(command, EncodeContent(content)...)

	return command
}

func EncodeIMEI(imei string) []byte {
	if len(imei) != 15 {
		return nil
	}
	imei = "0" + imei

	data, err := hex.DecodeString(imei)
	if err != nil {
		return nil
	}
	return data
}

func EncodeContent(content string) []byte {
	content = strings.ToValidUTF8(content, "")
	content = content + "#"
	return []byte(content)
}
