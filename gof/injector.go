package gof

import (
	"reflect"
)

type Injector struct {
	values map[reflect.Type]reflect.Value
}

func (in *Injector) Into(v interface{}, t ...interface{}) {
	rt := reflect.TypeOf(v)
	if len(t) > 0 {
		rt = reflect.TypeOf(t[0]).Elem()
	}
	in.values[rt] = reflect.ValueOf(v)
}

func (in *Injector) Out(v interface{}) {
	rt := reflect.TypeOf(v)
	if iv, ok := in.values[rt]; ok {
		if rt.Kind() == reflect.Ptr {
			reflect.ValueOf(v).Elem().Set(iv.Elem())
			return
		}
		reflect.ValueOf(v).Set(iv)
		return
	}
	rv := reflect.ValueOf(v).Elem()
	rv.Set(reflect.Zero(rv.Type()))
}

func (in *Injector) Clone() *Injector {
	i := new(Injector)
	i.values = in.values
	return i
}

func NewInjector() *Injector {
	in := new(Injector)
	in.values = make(map[reflect.Type]reflect.Value)
	return in
}
