# ouiTooler
oui.txt lorder

# Simple

```go
package main

import (
  "fmt"
  "github.com/mrYuandq/ouiTool-txt"
)

var Oui *OuiDb

func main() {
	Oui = Oui.Open("../oui.txt")
	if nil == Oui {
		fmt.Println("open oui.txt failed")
		return
	}
  
	addr, err := Oui.VendorLookup("00254B")
	if nil != err {
		fmt.Println("err:", err.Error())
		return
	}
	fmt.Println(addr.Organization)
}
```
