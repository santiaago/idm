// package idm (it doesn't matter) is a toy implementation of an APL interpreter.
package main

import (
	"fmt"
	"strings"
)

func main() {
	// for {
	// 	fmt.Printf("\t")
	// 	var s string
	// 	fmt.Scanf("%s", &s)
	// 	stmt, err := NewParser(strings.NewReader(s)).Parse()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Printf("%+v\n\t", stmt)
	// }
	s := "a = b"
	stmt, err := NewParser(strings.NewReader(s)).Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", stmt)
}
