// Package bean /**
package bean

import (
	"fmt"
	"reflect"
)

func CopyFields(source interface{}, target interface{}, fields ...string) (err error) {
	if source == nil || target == nil {
		err = fmt.Errorf("source or target is nil")
		return
	}

	sourceType := reflect.TypeOf(source)
	sourceValue := reflect.ValueOf(source)
	targetType := reflect.TypeOf(target)
	targetValue := reflect.ValueOf(target)

	// 检查 source 是否为结构体或结构体指针
	if sourceType.Kind() != reflect.Struct {
		err = fmt.Errorf("source must be a struct")
		return
	}

	// 检查 target 是否为结构体指针
	if targetType.Kind() != reflect.Ptr {
		err = fmt.Errorf("target must be a struct pointer")
		return
	}

	// 要复制哪些字段
	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < sourceValue.NumField(); i++ {
			_fields = append(_fields, sourceType.Field(i).Name)
		}
	}

	if len(_fields) == 0 {
		fmt.Println("no fields to copy")
		return
	}

	// 复制
	elem := targetValue.Elem()
	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := elem.FieldByName(name)
		bValue := sourceValue.FieldByName(name)

		// a中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		} else {
			fmt.Printf("no such field or different kind, fieldName: %s\n", name)
		}
	}
	return
}
