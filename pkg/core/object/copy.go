package object

import (
	"reflect"
)

func CopyObject[T Object](from, to T) {
	if reflect.ValueOf(to).Kind() != reflect.Ptr {
		panic("must copy to a pointer")
	}
	fromV := reflect.ValueOf(from)
	if fromV.Kind() == reflect.Ptr {
		fromElem := fromV.Elem()
		fromVal := reflect.New(fromElem.Type())
		fromVal.Elem().Set(fromElem)
		reflect.ValueOf(to).Elem().Set(fromVal.Elem())
	} else {
		to = fromV.Interface().(T)
	}
}
