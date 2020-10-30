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
Then just 
## Why?
* It is up to 3x as fast as the native one (encoding/json).
* It is really easy to use: 
```golang
myConvertedJson := quickson.Marshal(MyInterface)
```
## When not to use it?
* If you have really complexe interfaces because it might be buggy with some conversions.
