package idm

import (
	"fmt"
	"strings"
)

func main() {
	s := "a = b"
	stmt, err := NewParser(strings.NewReader(s)).Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", stmt)
}
