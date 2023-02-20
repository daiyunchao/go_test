package reflect_test_test

import (
	"errors"
	"go_test/reflect_test"
	"reflect"
	"testing"
)

type MyUser struct {
	Name string
	age  string
}

func Test_getType(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "nil",
			args: args{
				val: nil,
			},
			want:    "",
			wantErr: errors.New("参数不能为nil"),
		},
		{
			name: "string",
			args: args{
				val: "Tom",
			},
			want:    "string",
			wantErr: nil,
		},
		{
			name: "struct",
			args: args{
				val: reflect_test.User{
					Name: "Tom",
				},
			},
			want:    "User",
			wantErr: nil,
		},
		{
			name: "pointer",
			args: args{
				val: &reflect_test.User{
					Name: "Tom",
				},
			},
			want:    "",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reflect_test.GetType(tt.args.val)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("getType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKind(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Kind
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "nil",
			args: args{
				val: nil,
			},
			want:    reflect.Invalid,
			wantErr: errors.New("参数不能为nil"),
		},
		{
			name: "string",
			args: args{
				val: "Tom",
			},
			want:    reflect.String,
			wantErr: nil,
		},
		{
			name: "Struct",
			args: args{
				val: reflect_test.User{
					Name: "Tom",
				},
			},
			want:    reflect.Struct,
			wantErr: nil,
		},
		{
			name: "Pointer",
			args: args{
				val: &reflect_test.User{
					Name: "Tom",
				},
			},
			want:    reflect.Pointer,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reflect_test.GetKind(tt.args.val)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFieldMap(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "nil",
			args: args{
				val: nil,
			},
			want:    nil,
			wantErr: errors.New("参数不能为nil"),
		},
		{
			name: "string",
			args: args{
				val: "Tom",
			},
			want:    nil,
			wantErr: errors.New("参数为非结构体或非结构体指针"),
		},
		{
			name: "User",
			args: args{
				val: reflect_test.User{
					Name: "Tom",
				},
			},
			want: map[string]any{
				"Name": "Tom",
			},
			wantErr: nil,
		},
		{
			name: "UserPointer",
			args: args{
				val: &reflect_test.User{
					Name: "Tom",
				},
			},
			want: map[string]any{
				"Name": "Tom",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reflect_test.GetFieldMap(tt.args.val)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetFieldMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldMap() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestSetField(t *testing.T) {
	type args struct {
		val       any
		fieldName string
		fieldVal  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "nil",
			args: args{
				val:       nil,
				fieldName: "",
				fieldVal:  "",
			},
			wantErr: errors.New("参数不能为nil"),
		},
		{
			name: "string",
			args: args{
				val:       "Tom",
				fieldName: "Tom",
				fieldVal:  "Tom",
			},
			wantErr: errors.New("参数必须是指针类型"),
		},
		{
			name: "MyUser",
			args: args{
				val: MyUser{
					Name: "Tom",
				},
				fieldName: "Name",
				fieldVal:  "Jerry",
			},
			wantErr: errors.New("参数必须是指针类型"),
		},
		{
			name: "UserPointer",
			args: args{
				val: &MyUser{
					Name: "Tom",
				},
				fieldName: "Name",
				fieldVal:  "Jerry",
			},
			wantErr: nil,
		},
		{
			name: "NotFoundField",
			args: args{
				val: &MyUser{
					Name: "Tom",
				},
				fieldName: "Score",
				fieldVal:  "10",
			},
			wantErr: errors.New("未找到指定字段"),
		},
		{
			name: "Can't Set",
			args: args{
				val: &MyUser{
					Name: "Tom",
				},
				fieldName: "age",
				fieldVal:  "10",
			},
			wantErr: errors.New("该字段不能设置"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := reflect_test.SetField(tt.args.val, tt.args.fieldName, tt.args.fieldVal); (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
