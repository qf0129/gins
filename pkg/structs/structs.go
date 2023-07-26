package structs

import (
	"encoding/json"
	"reflect"

	"github.com/sirupsen/logrus"
)

func GetFields(obj any) []string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		logrus.Warn("structs: InvalidStruct, " + t.Kind().String())
		return nil
	}

	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetJsonFields(obj any) []string {
	var objMap map[string]any

	b, err := json.Marshal(obj)
	if err != nil {
		logrus.Warn("structs: MarshalJsonFailed")
		return nil
	}
	err = json.Unmarshal(b, &objMap)
	if err != nil {
		logrus.Warn("structs: UnMarshalJsonFailed, " + string(b))
		return nil
	}

	var result []string
	for k, _ := range objMap {
		result = append(result, k)
	}
	return result
}
