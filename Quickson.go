package quickson

import (
	"fmt"
	"reflect"
	"strings"
)

//Marshal is used to convert a struct/interface into JSON. It outputs a string
func Marshal(x interface{}) string {
	v := reflect.ValueOf(x)
	vi := reflect.Indirect(v)
	var t = "{"
	for i := 0; i != vi.NumField(); i++ {
		osi := vi.Field(i).Type().String()
		switch osi {
		case "string":
			t += "\"" + vi.Type().Field(i).Name + "\":\"" + strings.Replace(fmt.Sprint(vi.Field(i).Interface()), "\"", "\\\"", -1) + "\","
		case "bool", "int", "uint8":
			t += "\"" + vi.Type().Field(i).Name + "\":\"" + fmt.Sprint(vi.Field(i).Interface()) + "\","
		default:
			if compareBytes(osi[:4], "map[") {
				t += "\"" + vi.Type().Field(i).Name + "\":{"
				mapTmpKeys := reflect.ValueOf(vi.Field(i).Interface()).MapKeys()
				for _, key := range mapTmpKeys {
					switch reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String() {
					case "int", "uint8", "bool":
						t += "\"" + fmt.Sprint(key.Interface()) + "\":" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + ","
					case "string":
						t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + "\","
					default:
						t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String()) + ","
					}
				}
				t = t[:len(t)-1]
				t += "},"
			} else if compareBytes(osi[:1], "[") {
				t += "\"" + vi.Type().Field(i).Name + "\":["
				for ia := 0; ia < reflect.ValueOf(vi.Field(i).Interface()).Len(); ia++ {
					switch reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Type().String() {
					case "int", "uint8", "bool":
						t += fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Interface()) + ","
					case "string":
						t += "\"" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Interface()) + "\","
					default:
						t += marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Interface()), reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Type().String()) + ","
					}
				}
				t = t[:len(t)-1]
				t += "],"
			} else {
				t += Marshal(vi.Field(i).Interface())
			}
		}
	}
	t = t[:len(t)-1]
	return t + "}"
}
func compareBytes(sa string, sb string) bool {
	a := []byte(sa)
	b := []byte(sb)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
func marshalDeep(vi reflect.Value, bytedType string) string {
	var t string = ""
	if compareBytes(bytedType[:4], "map[") {

		t += "{"
		mapTmpKeys := reflect.ValueOf(vi.Interface()).MapKeys()
		for _, key := range mapTmpKeys {
			switch reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String() {
			case "int", "uint8", "bool":
				t += "\"" + fmt.Sprint(key.Interface()) + "\":" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + ","
			case "string":
				t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + "\","
			default:
				t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String()) + ","
			}
		}
		t = t[:len(t)-1]
		t += "}"
	} else if compareBytes(bytedType[:1], "[") {
		t += "["
		for ia := 0; ia < reflect.ValueOf(vi.Interface()).Len(); ia++ {
			switch reflect.ValueOf(vi.Interface()).Index(ia).Type().String() {
			case "int", "uint8", "bool":
				t += fmt.Sprint(reflect.ValueOf(vi.Interface()).Index(ia).Interface()) + ","
			case "string":
				t += "\"" + fmt.Sprint(reflect.ValueOf(vi.Interface()).Index(ia).Interface()) + "\","
			default:
				t += marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Interface()).Index(ia).Interface()), reflect.ValueOf(vi.Interface()).Index(ia).Type().String()) + ","
			}
		}
		t = t[:len(t)-1]
		t += "]"
	} else {
		t += Marshal(vi.Interface())
	}
	return t
}
