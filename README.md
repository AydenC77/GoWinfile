# Go Win file 
A Golang lib to open files using windows explore

## Example

```go
package main

import (
	"fmt"
	"path"

	"github.com/AydenC77/GoWinfile"
)

func main() {
	path, err GoWinfile.PickTextFile("Text Files (*.txt)", "*.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected path %s\n", path)
}
```
