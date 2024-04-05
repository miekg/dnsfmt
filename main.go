package main

import (
	"fmt"
	"os"

	"github.com/miekg/dnsfmt/zonefile"
)

// Increments the serial of a zonefile
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<path to zonefile>")
		os.Exit(1)
	}

	// Load zonefile
	data, ioerr := os.ReadFile(os.Args[1])
	if ioerr != nil {
		fmt.Println(os.Args[1], ioerr)
		os.Exit(2)
	}

	zf, perr := zonefile.Load(data)
	if perr != nil {
		fmt.Println(os.Args[1], perr.LineNo(), perr)
		os.Exit(3)
	}

	for _, e := range zf.Entries() {
		fmt.Printf("%v\n", e)
	}
}
