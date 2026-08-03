# Go Win file 
A Golang lib to open files using windows explore

## Example

```go
package main

import (
	"fmt"

	gowinfile "github.com/AydenC77/GoWinfile"
)

func main() {
	path, err := gowinfile.PickTextFile("Text Files (*.txt)", "*.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected path %s\n", path)
}
```