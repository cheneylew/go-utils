package utils

import "fmt"

func Equal(a, b interface{}) bool {
	var ta, tb float64

	strCount := 0
	switch a.(type) {
	case uint8,uint16,uint64,int8,int16,int32,int64,int:
		ta = float64(ToFloat64(fmt.Sprintf("%d", a)))
	case float32,float64:
		ta = float64(ToFloat64(fmt.Sprintf("%f", a)))
	case string:
		strCount += 1
		ta = ToFloat64(a.(string))
	}

	switch b.(type) {
	case uint8,uint16,uint64,int8,int16,int32,int64,int:
		tb = float64(ToFloat64(fmt.Sprintf("%d", b)))
	case float32,float64:
		tb = float64(ToFloat64(fmt.Sprintf("%f", b)))
	case string:
		strCount += 1
		tb = ToFloat64(a.(string))
	}

	if strCount == 2 {
		return fmt.Sprintf("%s", a) == fmt.Sprintf("%s", b)
	}

	return ta == tb
}

