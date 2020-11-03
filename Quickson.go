package quickson

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//Marshal is used to convert a struct/interface into JSON. It outputs a string
//Usage: jsonText := quickson.Marshal(myInterface)
//Where the interface can be anything like a struct, map, int, slice....
//The result will be inline JSON with no indentation
func Marshal(x interface{}) string {
	v := reflect.ValueOf(x)
	vi := reflect.Indirect(v)
	var t string = ""
	if compareBytes(vi.Type().String()[:4], "map[") || compareBytes(vi.Type().String()[:1], "[") {
		t += marshalDeep(vi, vi.Type().String())
	} else {
		t = "{"
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
						if key.Type().String() != "string" {
							switch reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String() {
							case "int", "uint8", "bool":
								t += fmt.Sprint(key.Interface()) + ":" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + ","
							case "string":
								t += fmt.Sprint(key.Interface()) + ":\"" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + "\","
							default:
								t += fmt.Sprint(key.Interface()) + ":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String()) + ","
							}
						} else {
							switch reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String() {
							case "int", "uint8", "bool":
								t += "\"" + fmt.Sprint(key.Interface()) + "\":" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + ","
							case "string":
								t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + fmt.Sprint(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()) + "\","
							default:
								t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Field(i).Interface()).MapIndex(key).Type().String()) + ","
							}
						}
					}
					if len(mapTmpKeys) > 0 {
						t = t[:len(t)-1]
					}
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
					if reflect.ValueOf(vi.Field(i).Interface()).Len() > 0 {
						t = t[:len(t)-1]
					}
					t += "],"
				} else {
					t += Marshal(vi.Field(i).Interface())
				}
			}
		}
		if vi.NumField() > 0 {
			t = t[:len(t)-1]
		}
		t += "}"
	}
	return t
}

//to iterate deeply in the interface, can take map and slices/array. Interfaces are redirected to Marshal

func marshalDeep(vi reflect.Value, bytedType string) string {
	var t string = ""
	if compareBytes(bytedType[:4], "map[") {

		t += "{"
		mapTmpKeys := reflect.ValueOf(vi.Interface()).MapKeys()
		for _, key := range mapTmpKeys {
			if key.Type().String() != "string" {
				switch reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String() {
				case "int", "uint8", "bool":
					t += fmt.Sprint(key.Interface()) + ":" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + ","
				case "string":
					t += fmt.Sprint(key.Interface()) + ":\"" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + "\","
				default:
					t += fmt.Sprint(key.Interface()) + ":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String()) + ","
				}
			} else {
				switch reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String() {
				case "int", "uint8", "bool":
					t += "\"" + fmt.Sprint(key.Interface()) + "\":" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + ","
				case "string":
					t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + fmt.Sprint(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()) + "\","
				default:
					t += "\"" + fmt.Sprint(key.Interface()) + "\":\"" + marshalDeep(reflect.ValueOf(reflect.ValueOf(vi.Interface()).MapIndex(key).Interface()), reflect.ValueOf(vi.Interface()).MapIndex(key).Type().String()) + ","
				}
			}
		}
		if len(mapTmpKeys) > 0 {
			t = t[:len(t)-1]
		}
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
		if reflect.ValueOf(vi.Interface()).Len() > 0 {
			t = t[:len(t)-1]
		}
		t += "]"
	} else {
		t += Marshal(vi.Interface())
	}
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

//Unmarshal has two purposes:
//-It can parse JSON to a struct with a pointer Unmarshal(json, &myinterface)
//-Extract data to a map/slice/string/bool/int myMap := Unmarshal(json, false)
func Unmarshal(t string, i interface{}) interface{} {
	t = strings.Trim(t, " ")
	t = strings.Trim(t, "\n")
	if t[0] == '{' {
		a, _ := getMap(t, i, true)
		return a
	} else if t[0] == '[' {
		a, _ := getSlice(t)
		return a
	} else if t[0] == '"' {
		a, _ := getValue(t)
		return a
	} else if t[0] == 't' || t[0] == 'f' {
		a, _ := getBool(t)
		return a
	} else {
		a, _ := getInt(t)
		return a
	}
}

func getMap(t string, x interface{}, isFirst bool) (interface{}, string) {
	var is int = 0
	var antiSlash bool = false
	var memory []rune
	var maxI int = 0
	//get inner map to memory
	for i, c := range t {
		if c == '\\' {
			antiSlash = true
			if is >= 1 && antiSlash == true {
				memory = append(memory, c)
			}
		} else if antiSlash == true {
			antiSlash = false
			memory = append(memory, c)
		} else {
			if c == '{' {
				is++
			} else if c == '}' {
				is--
			}
			if is >= 1 {
				memory = append(memory, c)
			}

			if is == 0 && len(memory) > 0 {
				maxI = i + 1
				break
			}

		}
	}
	//analyze memory
	var isMade bool = false
	var tmpKey interface{}
	var tmpValue interface{}
	var mapValue reflect.Value
	var mapType reflect.Type
	isKey := true
	memory = memory[1:]
	for true {
		for len(memory) > 0 {
			if memory[0] == ' ' || memory[0] == '\n' || memory[0] == ',' {
				memory = memory[1:]
			} else {
				break
			}
		}
		if len(memory) == 0 {
			break
		}
		if isKey == true {
			if memory[0] == '{' {
				a, b := getMap(string(memory), false, false)
				memory = memory[len(b):]
				for i := 0; true; i++ {
					if memory[i] == ':' {
						memory = memory[:len(memory)-i]
						break
					}
				}
				tmpKey = a
			} else if memory[0] == '[' {

				a, b := getSlice(string(memory))
				memory = memory[len(b):]
				for i := 0; i < len(memory); i++ {
					if memory[i] == ',' {
						memory = memory[:len(memory)-i]
						break
					}
				}
				tmpKey = a
			} else if memory[0] == '"' {
				a, b := getValue(string(memory))
				memory = memory[len(b):]
				tmpKey = a
			} else if memory[0] == 't' || memory[0] == 'f' {
				a, b := getBool(string(memory))
				memory = memory[len(b):]
				tmpKey = a
			} else {
				a, b := getInt(string(memory))
				memory = memory[len(b):]
				tmpKey = a
			}
			if len(memory) >= 1 {
				memory = memory[1:]
			}
			isKey = false
		} else {
			if memory[0] == '{' {
				a, b := getMap(string(memory), "", false)
				memory = memory[len(b):]
				for i := 0; true; i++ {
					if memory[i] == ':' {
						memory = memory[:len(memory)-i]
						break
					}
				}
				tmpValue = a
			} else if memory[0] == '[' {
				a, b := getSlice(string(memory))
				memory = memory[len(b):]
				for i := 0; i < len(memory); i++ {
					if memory[i] == ',' {
						memory = memory[:len(memory)-i]
						break
					}
				}
				tmpValue = a
			} else if memory[0] == '"' {
				a, b := getValue(string(memory))
				memory = memory[len(b):]
				tmpValue = a
			} else if memory[0] == 't' || memory[0] == 'f' {
				a, b := getBool(string(memory))
				memory = memory[len(b):]
				tmpValue = a
			} else {
				a, b := getInt(string(memory))
				memory = memory[len(b):]
				tmpValue = a
			}
			if isMade == false && x == false {
				mapType = reflect.MapOf(reflect.TypeOf(tmpKey), reflect.TypeOf(tmpValue))
				mapValue = reflect.MakeMap(mapType)
				isMade = true
			}
			if isFirst && x != false {
				if reflect.ValueOf(tmpValue).Type().String() == "[]int" && reflect.Indirect(reflect.ValueOf(&x)).Elem().Elem().FieldByName(tmpKey.(string)).Type().String() == "[]uint8" {
					vi := reflect.ValueOf(tmpValue)
					var newTmpValue []uint8
					for ia := 0; ia < reflect.ValueOf(vi.Interface()).Len(); ia++ {
						iV, _ := strconv.Atoi(fmt.Sprint(reflect.ValueOf(vi.Interface()).Index(ia).Interface()))
						newTmpValue = append(newTmpValue, uint8(iV))
					}
					tmpValue = newTmpValue
				}
				reflect.Indirect(reflect.ValueOf(&x)).Elem().Elem().FieldByName(tmpKey.(string)).Set(reflect.ValueOf(tmpValue))
			}
			if x == false {
				go func() {
					if "map["+reflect.TypeOf(tmpKey).String()+"]"+reflect.TypeOf(tmpValue).String() != mapType.String() {
						panic("quickson: When exporting JSON to a map, every keys need to be the same type and evey values need to be the same type\n" + "map[" + reflect.TypeOf(tmpValue).String() + "]" + reflect.TypeOf(tmpValue).String() + " != " + mapType.String())
					}
				}()
				mapValue.SetMapIndex(reflect.ValueOf(tmpKey), reflect.ValueOf(tmpValue))
			}
			isKey = true
		}
		if len(memory) == 0 {
			break
		}
	}
	var returnMap interface{}
	if x == false {
		returnMap = mapValue.Interface()
	}
	return returnMap, t[:maxI]
}

func getSlice(t string) (interface{}, string) {
	var is int = 0
	var antiSlash bool = false
	var memory []rune
	var maxI int = 0
	for i, c := range t {
		if c == '\\' {
			antiSlash = true
			if is >= 1 && antiSlash == true {
				memory = append(memory, c)
			}
		} else if antiSlash == true {
			antiSlash = false
			memory = append(memory, c)
		} else {
			if c == '[' {
				is++
			} else if c == ']' {
				is--
			}
			if is >= 1 {
				memory = append(memory, c)
			}

			if is == 0 && len(memory) > 0 {
				maxI = i + 1
				break
			}

		}
	}
	//analyze memory
	var tmpValue interface{}
	var tmpSlice []interface{}
	var sliceValue reflect.Value
	var sliceType reflect.Type
	memory = memory[1:]
	for true {
		for len(memory) > 0 {
			if memory[0] == ' ' || memory[0] == '\n' || memory[0] == ',' {
				memory = memory[1:]
			} else {
				break
			}
		}
		if memory[0] == '{' {
			a, b := getMap(string(memory), false, false)
			if len(memory)-len(b) < 0 {
				break
			} else {
				memory = memory[len(b):]
			}
			if len(memory) > 0 {
				for i := 0; true; i++ {
					if memory[i] == ',' {
						memory = memory[:len(memory)-i]
						break
					}
				}
			}
			tmpValue = a
		} else if memory[0] == '[' {
			a, b := getSlice(string(memory))
			if len(memory)-len(b) < 0 {
				break
			} else {
				memory = memory[len(b):]
			}
			for i := 0; true; i++ {
				if memory[i] == ',' {
					memory = memory[:len(memory)-i]
					break
				}
			}
			tmpValue = a
		} else if memory[0] == '"' {
			a, b := getValue(string(memory))
			if len(memory)-len(b) < 0 {
				break
			} else {
				memory = memory[len(b):]
			}
			tmpValue = a
		} else if memory[0] == 't' || memory[0] == 'f' {
			a, b := getBool(string(memory))
			if len(memory)-len(b) < 0 {
				break
			} else {
				memory = memory[len(b):]
			}
			tmpValue = a
		} else {
			a, b := getInt(string(memory))
			if len(memory)-len(b) < 0 {
				break
			} else {
				memory = memory[len(b):]
			}
			tmpValue = a
		}
		tmpSlice = append(tmpSlice, tmpValue)
		Nmemory := strings.TrimPrefix(string(memory), ",")
		Nmemory = strings.TrimPrefix(Nmemory, " ")
		memory = []rune(Nmemory)
		if len(memory) == 0 {
			break
		}
	}
	sliceType = reflect.SliceOf(reflect.TypeOf(tmpValue))
	sliceValue = reflect.MakeSlice(sliceType, 0, len(tmpSlice))
	for _, in := range tmpSlice {
		sliceValue = reflect.Append(sliceValue, reflect.ValueOf(in).Convert(reflect.TypeOf(tmpValue)))
	}
	return sliceValue.Interface(), t[:maxI]
}

func getValue(t string) (string, string) {
	var is int = 0
	var antiSlash bool = false
	var memory []rune
	var maxI int = 0
	for i, c := range t {
		if c == '\\' {
			if antiSlash == true {
				memory = append(memory, c)
			}
			antiSlash = true
		} else if antiSlash == true {
			antiSlash = false
			memory = append(memory, c)
		} else {
			if c == '"' && is == 0 {
				is = 1
			} else if c == '"' && is == 1 {
				maxI = i
				break
			} else {
				memory = append(memory, c)
			}
		}
	}
	if len(t) >= maxI+2 {
		if t[maxI+1] == ',' {
			maxI++
		}
	}
	return string(memory), t[:maxI+1]
}

func getBool(t string) (bool, string) {
	for _, c := range t {
		if c == 'f' {
			if len(t) >= 6 {
				if t[5] == ',' {
					return false, "false,"
				}
			}
			return false, "false"
		} else if c == 't' {
			if len(t) >= 5 {
				if t[4] == ',' {
					return true, "true,"
				}
			}
			return true, "true"

		}
	}
	return false, ""
}
func getInt(t string) (int, string) {
	var memory []rune
	maxAi := 0
	for ai, c := range t {
		if c == ',' && len(memory) > 0 {
			i, _ := strconv.Atoi(string(memory))
			return i, t[:ai]
		} else if c != ':' {
			memory = append(memory, c)
		} else {
			i, _ := strconv.Atoi(string(memory))
			return i, t[:ai]
		}
		maxAi = ai
	}
	i, _ := strconv.Atoi(string(memory))
	return i, t[:maxAi+1]
}
