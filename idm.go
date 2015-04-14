// package idm (it doesn't matter) is a toy implementation of an APL interpreter.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\t") // human lines start at tab. machine lines are without tab

	for scanner.Scan() {
		s := scanner.Text()
		stmt, err := NewParser(strings.NewReader(s)).Parse()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", stmt)
		fmt.Printf("\t")
	}
}
