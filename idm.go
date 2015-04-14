// package idm as It Doesn't Matter) is a toy implementation of an APL interpreter.
package main // package github.com/santiaago/idm

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func init() {
	log.SetFlags(log.Ltime | log.Ldate | log.Lshortfile)
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\t") // human lines start at tab. machine lines are without tab
	for scanner.Scan() {
		s := scanner.Text()
		expr, err := NewParser(strings.NewReader(s)).Parse()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("\t")
			continue
		}
		fmt.Printf("%+v\n", (*expr).Evaluate())
		fmt.Printf("\t")
	}
}
