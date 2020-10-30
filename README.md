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
## Why?
* It is up to 3x as fast as the native one (encoding/json).
* It is really easy to use: 
```golang
myConvertedJson := quickson.Marshal(MyInterface)
```
## When not to use it?
* If you have really complexe interfaces because it might be buggy with some conversions.

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
fmt.PrintF("This is our struct converted in JSON: %v", quickson.Marshal(testVar))

}
```
