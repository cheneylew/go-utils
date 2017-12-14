package models

import (
	"time"
	"reflect"
	"fmt"
	"errors"
	"github.com/cheneylew/goutil/utils"
)

type Response struct {
	Code int64
	Msg string
	Data map[string]map[string]interface{}
}

type KLineDay struct {
	OpeningPrice float64
	ClosingPrice float64
	MaxPrice float64
	MinPrice float64
	Date time.Time
}

type StockInfo struct {
	Qfqday [][]string
	Qt map[string][]string
	Prec string
	Version string
}


func (s *StockInfo) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetField(obj interface{}, name string, value interface{}) error {
	name = utils.UpperFirstChar(name)
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}