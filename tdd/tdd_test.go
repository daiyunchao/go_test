package tdd

import (
	"errors"
	"reflect"
	"testing"
)

func TestPrintType(t *testing.T) {
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
			name: "int",
			args: args{
				val: 123,
			},
			want:    "int",
			wantErr: nil,
		},
		{
			name: "nil",
			args: args{
				val: nil,
			},
			wantErr: errors.New("please input val arg"),
			want:    "",
		},
		{
			name: "user",
			args: args{
				val: User{
					Name: "zhangSan",
				},
			},
			want:    "User",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetType(tt.args.val)
			if err != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("PrintType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PrintType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type User struct {
	Name string
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
			wantErr: errors.New("please input val arg"),
			want:    reflect.Invalid,
		},
		{
			name: "string",
			args: args{
				val: "zhangSan",
			},
			want:    reflect.String,
			wantErr: nil,
		},
		{
			name: "int",
			args: args{
				val: 123,
			},
			wantErr: nil,
			want:    reflect.Int,
		},
		{
			name: "User",
			args: args{
				val: User{
					Name: "zhangSan",
				},
			},
			want:    reflect.Struct,
			wantErr: nil,
		},
		{
			name: "pointer",
			args: args{
				val: &User{Name: "zhangSan"},
			},
			want:    reflect.Pointer,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKind(tt.args.val)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}
