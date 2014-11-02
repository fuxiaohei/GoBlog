package middle

import (
	"fmt"
	"github.com/fuxiaohei/GoBlog/gof"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Form() gof.RouterHandler {
	return nil
}

func FormBind() gof.RouterHandler {
	return nil
}

type FormBinderValue struct {
	UseTag     bool
	FieldTypes map[string]reflect.Type
	Fields     map[string]string
	//ValidRule map[string]string
}

type FormBinder struct {
	values map[reflect.Type]*FormBinderValue
}

type FormBinderError error

func (fb *FormBinder) toIntValue(v string, t reflect.Type) (reflect.Value, error) {
	intV, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return reflect.Zero(t), err
	}
	if t.Kind() == reflect.Int {
		return reflect.ValueOf(int(intV)), nil
	}
	if t.Kind() == reflect.Int8 {
		return reflect.ValueOf(int8(intV)), nil
	}
	if t.Kind() == reflect.Int16 {
		return reflect.ValueOf(int16(intV)), nil
	}
	if t.Kind() == reflect.Int32 {
		return reflect.ValueOf(int32(intV)), nil
	}
	if t.Kind() == reflect.Int64 {
		return reflect.ValueOf(intV), nil
	}
	return reflect.Zero(t), nil
}

func (fb *FormBinder) toUintValue(v string, t reflect.Type) (reflect.Value, error) {
	intV, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return reflect.Zero(t), err
	}
	if t.Kind() == reflect.Uint {
		return reflect.ValueOf(uint(intV)), nil
	}
	if t.Kind() == reflect.Uint8 {
		return reflect.ValueOf(uint8(intV)), nil
	}
	if t.Kind() == reflect.Uint16 {
		return reflect.ValueOf(uint16(intV)), nil
	}
	if t.Kind() == reflect.Uint32 {
		return reflect.ValueOf(uint32(intV)), nil
	}
	if t.Kind() == reflect.Uint64 {
		return reflect.ValueOf(intV), nil
	}
	return reflect.Zero(t), nil
}

func (fb *FormBinder) ToStruct(formValue url.Values, s interface{}, useTag bool) FormBinderError {
	fv, err := fb.parseStruct(s, useTag)
	if err != nil {
		return err
	}
	formValue = fb.parseFormValue(formValue)

	rv := reflect.ValueOf(s)
	for name, t := range fv.FieldTypes {
		if len(formValue[name]) == 0 {
			continue
		}
		value := formValue[name][0]
		typeName := t.String()
		rf := rv.Elem().FieldByName(fv.Fields[name])
		if typeName == "string" {
			rf.Set(reflect.ValueOf(value))
			continue
		}
		if strings.Contains(typeName, "int") {
			intValue, err := fb.toIntValue(value, t)
			if err != nil {
				return FormBinderError(err)
			}
			rf.Set(intValue)
			continue
		}
	}
	return nil
}

func (fb *FormBinder) parseStruct(s interface{}, useTag bool) (*FormBinderValue, FormBinderError) {
	rt := reflect.TypeOf(s).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, FormBinderError(fmt.Errorf("form binder need a struct pointer"))
	}
	fv := &FormBinderValue{
		UseTag:     useTag,
		Fields:     make(map[string]string),
		FieldTypes: make(map[string]reflect.Type),
	}
	fieldLength := rt.NumField()
	for i := 0; i < fieldLength; i++ {
		rf := rt.Field(i)
		name := strings.ToLower(rf.Name)
		if useTag {
			if t := rf.Tag.Get("form"); t != "" {
				name = t
			}
		}
		fv.FieldTypes[name] = rf.Type
		fv.Fields[name] = rf.Name
	}
	return fv, nil
}

func (fb *FormBinder) parseFormValue(formValue url.Values) map[string][]string {
	for k, v := range formValue {
		delete(formValue, k)
		formValue[strings.ToLower(k)] = v
	}
	return formValue
}

func TestFormBinder() *FormBinder {
	return &FormBinder{make(map[reflect.Type]*FormBinderValue)}
}
