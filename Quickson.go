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
	var t string = ""
	if compareBytes(vi.Type().String()[:4], "map[") || compareBytes(vi.Type().String()[:1], "[") {
		go func() {
			t += marshalDeep(vi, vi.Type().String())
			fmt.Println("FT: " + t)
		}()
	} else {
		go func() {
			var Tt string = ""
			Tt += "{"
			for i := 0; i != vi.NumField(); i++ {
				osi := vi.Field(i).Type().String()
				switch osi {
				case "string":
					Tt += "\"" + vi.Type().Field(i).Name + "\":\"" + strings.Replace(fmt.Sprint(vi.Field(i).Interface()), "\"", "\\\"", -1) + "\","
				case "bool", "int", "uint8":
					Tt += "\"" + vi.Type().Field(i).Name + "\":\"" + fmt.Sprint(vi.Field(i).Interface()) + "\","
				default:
					if compareBytes(osi[:4], "map[") {
						Tt += "\"" + vi.Type().Field(i).Name + "\":{"
						mapTmpKeys := reflect.ValueOf(vi.Field(i).Interface()).MapKeys()
						for _, key := range mapTmpKeys {
							switch reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String() {
							case "int", "uint8", "bool":
								Tt += "\"" + fmt.Sprint(key.Interface()) + "\":" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + ","
							case "string":
								Tt += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + "\","
							default:
								Tt += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String()) + ","
							}
						}
						Tt = Tt[:len(Tt)-1]
						Tt += "},"
					} else if compareBytes(osi[:1], "[") {
						Tt += "\"" + vi.Type().Field(i).Name + "\":["
						for ia := 0; ia < reflect.ValueOf(vi.Field(i).Interface()).Len(); ia++ {
							switch reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Type().String() {
							case "int", "uint8", "bool":
								Tt += fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Interface()) + ","
							case "string":
								Tt += "\"" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Interface()) + "\","
							default:
								Tt += marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Interface()), reflect.ValueOf(vi.Field(i).Interface()).Index(ia).Type().String()) + ","
							}
						}
						Tt = Tt[:len(Tt)-1]
						Tt += "],"
					} else {
						Tt += Marshal(vi.Field(i).Interface())
					}
				}
			}
			Tt = Tt[:len(Tt)-1]
			Tt += "}"
			t += Tt
			fmt.Println("FT: " + t)
		}()
	}
	fmt.Println("LT: " + t)
	return t
}

//to compare types efficiently

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

//to iterate deeply in the interface, can take map and slices/array. Interfaces are redirected to Marshal

func marshalDeep(vi reflect.Value, bytedType string) string {
	var t string = ""
	go func() {
		var Tt string = ""
		if compareBytes(bytedType[:4], "map[") {

			Tt += "{"
			mapTmpKeys := reflect.ValueOf(vi.Interface()).MapKeys()
			for _, key := range mapTmpKeys {
				switch reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String() {
				case "int", "uint8", "bool":
					Tt += "\"" + fmt.Sprint(key.Interface()) + "\":" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + ","
				case "string":
					Tt += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + "\","
				default:
					Tt += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String()) + ","
				}
			}
			Tt = Tt[:len(Tt)-1]
			Tt += "}"
		} else if compareBytes(bytedType[:1], "[") {
			Tt += "["
			for ia := 0; ia < reflect.ValueOf(vi.Interface()).Len(); ia++ {
				switch reflect.ValueOf(vi.Interface()).Index(ia).Type().String() {
				case "int", "uint8", "bool":
					Tt += fmt.Sprint(reflect.ValueOf(vi.Interface()).Index(ia).Interface()) + ","
				case "string":
					Tt += "\"" + fmt.Sprint(reflect.ValueOf(vi.Interface()).Index(ia).Interface()) + "\","
				default:
					Tt += marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Interface()).Index(ia).Interface()), reflect.ValueOf(vi.Interface()).Index(ia).Type().String()) + ","
				}
			}
			Tt = Tt[:len(Tt)-1]
			Tt += "]"
		} else {
			Tt += Marshal(vi.Interface())
		}
		t += Tt
	}()
	return t
}

//[IN COMMING]Unmarshal is used to transform JSON into an interface. You will need to input the pointer of the interface to fill.
// func Unmarshal(j string, x *interface{}) {

// }
