package reflect_test

import (
	"errors"
	"fmt"
	"reflect"
)

type User struct {
	Name string
}

func GetType(val any) (string, error) {
	if val == nil {
		return "", errors.New("参数不能为nil")
	}
	tpe := reflect.TypeOf(val)
	return tpe.Name(), nil
}

func GetKind(val any) (reflect.Kind, error) {
	if val == nil {
		return reflect.Invalid, errors.New("参数不能为nil")
	}
	tpe := reflect.TypeOf(val)
	return tpe.Kind(), nil
}

func GetFieldMap(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("参数不能为nil")
	}
	typ := reflect.TypeOf(val)
	typVal := reflect.ValueOf(val)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		typVal = typVal.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("参数为非结构体或非结构体指针")
	}
	fliedCount := typ.NumField()
	var retMap = make(map[string]any, fliedCount)
	for i := 0; i < fliedCount; i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldValue := typVal.FieldByName(fieldName)
		retMap[fieldName] = fieldValue.String()
	}
	return retMap, nil
}

func SetField(val any, fieldName string, fieldVal string) error {
	if val == nil || fieldName == "" || fieldVal == "" {
		return errors.New("参数不能为nil")
	}
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer {
		return errors.New("参数必须是指针类型")
	}
	_, ok := typ.Elem().FieldByName(fieldName)
	if !ok {
		return errors.New("未找到指定字段")
	}
	oldFieldVal := reflect.ValueOf(val).Elem().FieldByName(fieldName)
	if !oldFieldVal.CanSet() {
		return errors.New("该字段不能设置")
	}
	oldFieldVal.SetString("1000")
	fmt.Printf("newValue is :%s", reflect.ValueOf(val).Elem().FieldByName(fieldName))
	return nil
}
