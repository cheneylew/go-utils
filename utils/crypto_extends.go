package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(target interface{}) string {
	str := ""
	switch target.(type) {
	case string:
		str = target.(string)
	case int, int8, int32, int64:
		str = fmt.Sprintf("%d", target)
	default:
		str = fmt.Sprintf("%s", target)
	}

	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}


