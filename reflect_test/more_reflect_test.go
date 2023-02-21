package reflect_test_test

import (
	"fmt"
	"reflect"
	"testing"
)

func getType(arg any) {
	fmt.Printf("getType==> %v\n", reflect.TypeOf(arg))
}

func getKind(arg any) {
	typ := reflect.TypeOf(arg)
	fmt.Printf("getKind==>%v\n", typ.Kind())
}

func getValue(arg any) {
	typ := reflect.TypeOf(arg)
	val := reflect.ValueOf(arg)
	if typ.Kind() != reflect.Pointer {
		val = reflect.ValueOf(&arg)
	}
	val = val.Elem()
	fmt.Printf("getValue==>%v\n", val)
}

func setValue(arg any) {
	typ := reflect.TypeOf(arg)
	val := reflect.ValueOf(arg)
	if typ.Kind() != reflect.Pointer {
		val = reflect.ValueOf(&arg)
	}
	val = val.Elem()
	if val.CanSet() {
		val.Set(reflect.ValueOf(100))
	}
	fmt.Printf("setValue==>%v\n", arg)
}

func getStructFieldAndMethod(arg any) {
	typ := reflect.TypeOf(arg)
	val := reflect.ValueOf(arg)
	fieldCount := typ.NumField()
	for i := 0; i < fieldCount; i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		fieldName := field.Name
		fieldType := field.Type
		value := fieldVal.Interface()
		fmt.Printf("getStructFieldAndMethod==> fieldName: %s, fieldType: %s, fieldValue: %v\n", fieldName, fieldType, value)
	}

	methodCount := typ.NumMethod()
	for i := 0; i < methodCount; i++ {
		method := typ.Method(i)
		methodName := method.Name
		fmt.Printf("getStructFieldAndMethod==> methodName: %s, methodType:%s\n", methodName, method.Type)
	}
}

func callFunc(arg any, funcName string, funcArgs []any) {
	rv := reflect.ValueOf(arg)
	method := rv.MethodByName(funcName)
	rvArgs := make([]reflect.Value, 0)
	for _, v := range funcArgs {
		rvArgs = append(rvArgs, reflect.ValueOf(v))
	}
	ret := method.Call(rvArgs)
	for _, v := range ret {
		fmt.Printf("callFunc ret:%v\n", v.Interface())
	}
}

type User struct {
	Name string
	Age  int
}

func (u User) GetName() string {
	return u.Name
}
func TestReflect(t *testing.T) {
	i := 10
	getType(i)
	getKind(i)
	getValue(i)
	getValue(&i)
	setValue(i)
	getStructFieldAndMethod(User{
		Name: "Tom", Age: 18,
	})
	arg := make([]any, 0)
	callFunc(User{Name: "Tom", Age: 18}, "GetName", arg)
}
