# Quickson is a JSON marshaler for the Go lang.
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
* To convert a struct into json:
```golang
myConvertedJson := quickson.Marshal(MyInterface)
```
* To parse JSON into a struct
```golang
data := myStruct{}
quickson.Unmarshal(json, &data)
```
* To parse JSON into a map/slice/string/bool/int
```golang
data := quickson.Unmarshal(json, false)
```

## When not to use it?
* If you have really complexe interfaces because it might be buggy with some conversions.

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


## I just want to test it
Then install the package and run the following:
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
