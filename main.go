package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miekg/dnsfmt/zonefile"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("dnsfmt: %s", err)
	}

	zf, perr := zonefile.Load(data)
	if perr != nil {
		log.Fatalf("dnsfmt: error on line %d: %s", perr.LineNo, perr)
	}

	longestname := 0
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
		// remove origin from other
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
	longestname += 4 // Extra indent

	prevname := []byte{}
	prevttl := 0
	prevcom := false
	firstname := true
	for _, e := range zf.Entries() {
		if e.IsComment {
			if !prevcom {
				fmt.Println()
			}
			for _, c := range e.Comments() {
				fmt.Printf("%s\n", c)
			}
			prevcom = true
			continue
		}
		if e.IsControl {
			fmt.Printf("%s %s\n", e.Command(), bytes.Join(e.Values(), []byte(" ")))
			prevcom = false
			continue
		}

		if !bytes.Equal(prevname, e.Domain()) {
			// keep comments near, don't add a newline when previous line was comment.
			// first record doesn't need a newline
			if len(e.Domain()) > 0 && !prevcom && !firstname {
				fmt.Println()
			}
			fmt.Printf("%-*s", longestname, e.Domain())
		} else {
			fmt.Printf("%-*s", longestname, "")
		}

		prevcom = false
		firstname = false

		if ttl := e.TTL(); ttl != nil && *ttl != prevttl {
			prevttl = *ttl
			fmt.Printf("%10s", TimeToHuman(ttl))
		} else {
			fmt.Printf("%10s", " ")
		}

		fmt.Printf("%5s", e.Class())
		fmt.Printf("   %-8s", e.Type())

		// Specicial handling for certain RR types
		switch {
		case bytes.Equal(e.Type(), []byte("TXT")):
			values := e.Values()
			fmt.Printf(Space3)
			for _, v := range values {
				fmt.Printf(" %q", v)
			}
			fmt.Println()
		case bytes.Equal(e.Type(), []byte("SOA")):
			values := e.Values()
			fmt.Printf("%s%s (\n", Space3, bytes.Join(values[:2], []byte(" ")))
			for i, v := range values[2:] {
				fmt.Printf("%-*s%s%-13s%s\n", longestname+Indent, " ", Space3, v, soacomment[i])
			}
			fmt.Printf("%-*s)\n", longestname+Indent, " ")

		case bytes.Equal(e.Type(), []byte("DNSKEY")):
			values := e.Values()
			fmt.Printf("%s%s (\n", Space3, bytes.Join(values[:3], []byte(" ")))

			all := bytes.Join(values[3:], nil)
			pieces := Split(all, 55)
			for _, p := range pieces {
				fmt.Printf("%-*s%s%-13s\n", longestname+Indent, " ", Space3, p)
			}
			fmt.Printf("%-*s)\n", longestname+Indent, " ")

		case bytes.Equal(e.Type(), []byte("RRSIG")):
			values := e.Values()
			fmt.Printf("%s%s (\n", Space3, bytes.Join(values[:8], []byte(" ")))

			all := bytes.Join(values[8:], nil)
			pieces := Split(all, 55)
			for _, p := range pieces {
				fmt.Printf("%-*s%s%-13s\n", longestname+Indent, " ", Space3, p)
			}
			fmt.Printf("%-*s)\n", longestname+Indent, " ")

		default:
			fmt.Printf("%s%s\n", Space3, bytes.Join(e.Values(), []byte(" ")))
		}
		if len(e.Domain()) > 0 {
			prevname = e.Domain()
		}
	}
}

const (
	Space3 = "   "
	Indent = 29
)

var soacomment = []string{"; serial", "; refresh", "; retry", "; expire", "; minimum"}

func Split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}
