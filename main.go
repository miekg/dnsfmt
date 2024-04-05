package main

import (
	"bytes"
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

	longestname := 0
	// TODO take ORIGIN in consideration, and also strip off the origin.
	// If name _is_ origin, make it '@'
	origin := []byte{}
	for _, e := range zf.Entries() {
		if e.IsComment {
			continue
		}
		if e.IsControl {
			if bytes.Equal(e.Command(), []byte("$ORIGIN")) {
				origin = e.Values()[0]
			}
			continue
		}
		if len(origin) > 0 && bytes.HasSuffix(e.Domain(), origin) {
			// remove origin plus dot.
			l := len(e.Domain())

			if l == len(origin) {
				e.SetDomain([]byte("@"))
			} else {
				e.SetDomain(e.Domain()[:l-len(origin)-1])
			}
		}

		if l := len(e.Domain()); l > longestname {
			longestname = l
		}
	}
	longestname += 4

	prevname := []byte{}
	for _, e := range zf.Entries() {
		if e.IsComment {
			for _, c := range e.Comments() {
				fmt.Printf("%s\n", c)
			}
			continue
		}
		if e.IsControl {
			fmt.Printf("%s %s\n", e.Command(), bytes.Join(e.Values(), []byte(" ")))
			continue
		}

		if !bytes.Equal(prevname, e.Domain()) {
			fmt.Printf("%-*s", longestname, e.Domain())
		} else {
			fmt.Printf("%-*s", longestname, "")
		}

		if ttl := e.TTL(); ttl != nil {
			fmt.Printf("%5d", *ttl)
		} else {
			fmt.Printf("%5s", " ")
		}

		fmt.Printf("%5s", e.Class())
		fmt.Printf("   %-8s", e.Type())

		if bytes.Equal(e.Type(), []byte("TXT")) {
			values := e.Values()
			fmt.Printf("  ")
			for _, v := range values {
				fmt.Printf(" %q", v)
			}
			fmt.Println()
		} else {
			fmt.Printf("   %s\n", bytes.Join(e.Values(), []byte(" ")))
		}
		prevname = e.Domain()
	}
}
