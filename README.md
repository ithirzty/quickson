# Quickson : Fast JSON Marshaller/Unmarshaller for golang. [![License](https://img.shields.io/badge/License-Apache--2.0-green)](https://github.com/ithirzty/quickson/blob/main/LICENSE)
![logo](https://github.com/ithirzty/quickson/blob/main/logo.png?raw=true)
##### What is a marshaller? Unmarshaller?
Marshal signifies parse. Unmarshal signifies parse but the other way (stringify)

## Installing
Go in your project directory, open a terminal and type the following
```bash
go get github.com/ithirzty/quickson
```
Then, in your `main.go` import the following
```golang
import(
	"github.com/ithirzty/quickson"
)
```
### How to update
```bash
go get -u github.com/ithirzty/quickson
```

## Why?
* It is up to 3x as fast as the native one (encoding/json), generaly 2x faster.
* It is really easy to use.

# How to use
* Converting a struct into JSON :
```golang
myConvertedJson := quickson.Marshal(MyInterface)
```
* Parsing JSON into a struct : 
```golang
data := myStruct{}
quickson.Unmarshal(json, &data)
```
* Paring JSON into a map/slice/string/bool/int
```golang
data := quickson.Unmarshal(json, false)
```

## Cons
* It is best you don't use this tool if you need to marshal complex interfaces, Quickson is not yet capable of converting big interfaces.

## Less performance under load?
Here is how to use concurrency with quickson:
```golang
result := ""
c := make(chan string)
			go func() {
				c <- quickson.Marshal(MyInterface)
			}()
result <- c
```


## How can I test it?
Follow the installation guide run the following code:
```golang
package main

import(
	"fmt"
	"github.com/ithirzty/quickson"
)

type testInterface struct {
  TestField  string           //"This is a test."
  TestPassed map[string]bool  //"My test": true
}

func main() {
	testVar := testInterface{"This is a test.", map[string]bool{"My test": true}}
	fmt.Printf("This is our struct converted in JSON: %v", quickson.Marshal(testVar))
	//should output {"TestField":"This is a test.","TestPassed":{"My test":true}}
}
```
