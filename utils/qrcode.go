package utils

import "github.com/skip2/go-qrcode"

func QRCodePNG(str string) ([]byte , error) {
	return qrcode.Encode(str, qrcode.Highest,256)
}

func QRCodePNGWriteFile(str string,path string) error {
	return qrcode.WriteFile(str, qrcode.Highest, 256, path)
}

