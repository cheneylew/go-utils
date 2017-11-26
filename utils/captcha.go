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
)

func NewCaptcha(dir string,count, width, height int) (fpath string, fname string, value string, e error){
	fileName := RandomString(32)+".jpg"
	imagePath := path2.Join(dir, fileName)

	d := captcha.RandomDigits(count)
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
