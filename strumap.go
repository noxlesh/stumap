package strumap

import (
	"bytes"
	"errors"
	"reflect"
)

type (
	Map      map[string]interface{}
	iterFn   func(reflect.Value, iterFn) Map
	wordCase int
)

const (
	snakeCase wordCase = iota
	camelCase
)

var (
	ErrEmptyStruct = errors.New("an input argument represents an empty struct")
	ErrNotStruct   = errors.New("an input argument doesn't represent a struct")
)

func convert(s interface{}, c wordCase) (Map, error) {
	var iter = func(v reflect.Value, recFn iterFn) Map {
		a := map[string]interface{}{}
		for i := 0; i < v.NumField(); i++ {
			t := v.Type()
			fieldName := t.Field(i).Name
			key := convCase(fieldName, c)
			value := v.FieldByName(fieldName)
			if !value.CanInterface() {
				continue
			}
			switch value.Kind() {
			case reflect.Struct:
				a[key] = recFn(value, recFn)
			case reflect.Ptr:
				if value.Elem().Kind() == reflect.Struct {
					a[key] = recFn(value.Elem(), recFn)
				} else {
					if value.IsNil() {
						a[key] = value.Interface()
					} else {
						a[key] = value.Elem().Interface()
					}
				}

			default:
				a[key] = value.Interface()
			}
		}
		return a
	}
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	fieldCount := v.NumField()
	if fieldCount == 0 {
		return nil, ErrEmptyStruct
	}
	return iter(v, iter), nil
}

func Convert(s interface{}) (Map, error) {
	return convert(s, camelCase)
}

func ConvertSnakeCase(s interface{}) (Map, error) {
	return convert(s, snakeCase)
}

func convCase(i string, c wordCase) string {
	if c == camelCase {
		return i
	} else {
		var buf bytes.Buffer
		for _, c := range i {
			isCapitalLetter := 'A' <= c && c <= 'Z'
			if isCapitalLetter {
				if buf.Len() > 0 {
					buf.WriteRune('_')
				}
				lowCaseLetter := c - 'A' + 'a'
				buf.WriteRune(lowCaseLetter)
				continue
			}
			buf.WriteRune(c)
		}
		return buf.String()
	}
}
