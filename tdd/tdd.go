package tdd

import (
	"errors"
	"reflect"
)

// GetType shift command t 添加函数的测试用例
func GetType(val any) (string, error) {
	if val == nil {
		return "", errors.New("please input val arg")
	}
	typ := reflect.TypeOf(val)
	return typ.Name(), nil
}

func GetKind(val any) (reflect.Kind, error) {
	if val == nil {
		return reflect.Invalid, errors.New("please input val arg")
	}
	typ := reflect.TypeOf(val)
	return typ.Kind(), nil
}
