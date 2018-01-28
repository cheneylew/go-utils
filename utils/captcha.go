package utils

/*
github.com/dchest/captcha/capexample
*/
import (
	"github.com/dchest/captcha"
	"image/jpeg"
	"os"
	path2 "path"
	"fmt"
	"encoding/base64"
	"bytes"
)

func NewCaptcha(dir string,digitLength, width, height int) (fpath string, fname string, value string, e error){
	fileName := RandomString(32)+".jpg"
	imagePath := path2.Join(dir, fileName)

	d := captcha.RandomDigits(digitLength)
	image := captcha.NewImage("anyid", d,width,height)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", "", "", err
	}

	var opt jpeg.Options
	opt.Quality = 80
	imageErr := jpeg.Encode(file,image,&opt)
	if e != nil {
		return "", "", "", imageErr
	}

	num := ""
	for _, v := range d {
		num += fmt.Sprintf("%d", v)
	}

	return imagePath, fname, num, nil
}

func NewCaptchaBase64Str(width, height int) (val, base64Str string, e error){
	d := captcha.RandomDigits(4)
	img := captcha.NewImage("anyid", d,width,height)

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		return "", "", err
	}

	num := ""
	for _, v := range d {
		num += fmt.Sprintf("%d", v)
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return num, "data:image/jpg;base64,"+str, nil
}
