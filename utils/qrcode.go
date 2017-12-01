package utils

import (
	goqrcode "github.com/skip2/go-qrcode"
	"os"
	"github.com/tuotoo/qrcode"
)

func QRCodePNG(str string) ([]byte , error) {
	return goqrcode.Encode(str, goqrcode.Highest,256)
}

func QRCodePNGWriteFile(str string,path string) error {
	return goqrcode.WriteFile(str, goqrcode.Highest, 256, path)
}

func QRCodeDecode(path string) (string,error) {
	fi, err := os.Open(path)
	if err != nil{
		return "", err
	}
	defer fi.Close()
	qrmatrix, err := qrcode.Decode(fi)
	if err != nil{
		return "", err
	}

	return qrmatrix.Content, nil
}

