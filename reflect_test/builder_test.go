package reflect_test_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

//通过结构体生成显示HTML+增删改查的例子

type Order struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Time int64  `json:"time" bson:"time"`
}

func TestBuilder(t *testing.T) {
	Builder(Order{})
}

func Builder(arg any) {
	buildHtml(arg)
	buildLogic(arg)
}

func buildHtml(arg any) error {
	argTyp := reflect.TypeOf(arg)
	argVal := reflect.ValueOf(arg)
	if argTyp.Kind() == reflect.Pointer {
		argTyp = argTyp.Elem()
		argVal = argVal.Elem()
	}
	if argTyp.Kind() != reflect.Struct {
		return errors.New("参数错误")
	}
	fieldCount := argTyp.NumField()
	htmlStr := ""
	for i := 0; i < fieldCount; i++ {
		fieldTyp := argTyp.Field(i)
		fieldType := fieldTyp.Type
		fieldJsonName := fieldTyp.Tag.Get("json")
		if fieldType.Kind() == reflect.String {
			htmlStr += `<div>
<spin>` + fieldJsonName + `:</spin>
<input type="text" id="input_` + fieldJsonName + `"/>
</div>`
		} else if (fieldType.Kind() == reflect.Int || fieldType.Kind() == reflect.Int32) || (fieldType.Kind() == reflect.Int32) || (fieldType.Kind() == reflect.Int64) {
			htmlStr += `<div>
<spin>` + fieldJsonName + `:</spin>
<input type="number" id="input_` + fieldJsonName + `"/>
</div>`
		}
	}
	fmt.Println(htmlStr)
	return nil
}

func buildLogic(arg any) {
	argTyp := reflect.TypeOf(arg)
	tableName := argTyp.Name()

	//查询全部
	queryAllStr := `SELECT * FROM ` + tableName
	fmt.Println("queryAllStr===>", queryAllStr)
	firstField := argTyp.Field(0)
	firstBson := firstField.Tag.Get("bson")

	//通过主键查询
	queryByPK := `SELECT * FROM ` + tableName + ` WHERE ` + firstBson + ` = ?`
	fmt.Println("queryByPK===>", queryByPK)

	//通过主键删除
	deleteByPK := `DELETE FROM ` + tableName + ` WHERE ` + firstBson + ` =?`
	fmt.Println("deleteByPK===>", deleteByPK)

	allCol := make([]string, 0)
	fieldCount := argTyp.NumField()
	for i := 0; i < fieldCount; i++ {
		field := argTyp.Field(i)
		allCol = append(allCol, field.Tag.Get("bson"))
	}
	allColUpdateStr := ``
	allColStr := ``
	allColVal := ``
	for i := 0; i < len(allCol); i++ {
		upperCol := strings.ToUpper(allCol[i])
		allColStr += upperCol + `, `
		allColVal += `?, `
		allColUpdateStr += `` + upperCol + ` = ? , `
	}
	//通过主键修改
	updateByPK := `UPDATE (` + allColUpdateStr + `) FROM ` + tableName + ` WHERE ` + firstBson + ` = ?`
	fmt.Println("updateByPK===>", updateByPK)

	//插入数据
	insert := `INSERT INTO (` + allColStr + `) VALUES (` + allColVal + `)`
	fmt.Println("insert===>", insert)
}
