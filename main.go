package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miekg/dns"
	"github.com/miekg/dnsfmt/zonefile"
)

var (
	flagOrigin = flag.String("o", "", "set the origin, otherwise taken from $ORIGIN or the owner name of the SOA record.")
	flagInc    = flag.Bool("i", true, "increase the serial")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("dnsfmt: %s", err)
		}
		Reformat(data, []byte(*flagOrigin), os.Stdout)
		return
	}

	for _, a := range flag.Args() {
		data, err := os.ReadFile(a)
		if err != nil {
			log.Fatalf("dnsfmt: %s", err)
		}
		Reformat(data, []byte(*flagOrigin), os.Stdout)
	}
}

func Reformat(data, origin []byte, w io.Writer) error {
	origin = zonefile.Fqdn(origin)

	zf, perr := zonefile.Load(data)
	if perr != nil {
		log.Fatalf("dnsfmt: error on line %d: %s", perr.LineNo, perr)
	}

	// 2 loops: finding and striping the  origin and some admin, and then actually reformatting.

	single := map[string]int{}
	longestname := 0
	prevname := []byte{}
	for _, e := range zf.Entries() {
		if e.IsComment {
			continue
		}
		if e.IsControl {
			if bytes.Equal(e.Command(), []byte("$ORIGIN")) {
				origin = zonefile.Fqdn(e.Values()[0])
			}
			continue
		}

		e.SetDomain(StripOrigin(origin, e.Domain()))

		// count number of types per name, as we want to group singletons.
		if !bytes.Equal(prevname, e.Domain()) && len(prevname) > 0 {
			if len(e.Domain()) > 0 {
				single[string(e.Domain())] += 1
			} else {
				single[string(prevname)] += 1
			}
		}

		// Strip origin from selected records.
		values := e.Values()
		switch e.RRType() {
		case dns.TypeSOA:
			if len(values) < 3 {
				return fmt.Errorf("malformed SOA RR: %v", values)
			}
			if len(origin) == 0 { // $ORIGIN not set take from SOA
				origin = zonefile.Fqdn(e.Domain())
			}

			e.SetValue(0, StripOrigin(origin, values[0]))
			e.SetValue(1, StripOrigin(origin, values[1]))

		case dns.TypeSRV:
			if len(values) < 4 {
				return fmt.Errorf("malformed SRV RR: %v", values)
			}
			e.SetValue(3, StripOrigin(origin, values[3]))

		case dns.TypeRRSIG:
			if len(values) < 8 {
				return fmt.Errorf("malformed RRSIG RR: %v", values)
			}
			e.SetValue(7, StripOrigin(origin, values[7]))

		case dns.TypeMX:
			if len(values) < 2 {
				return fmt.Errorf("malformed MX RR: %v", values)
			}
			e.SetValue(1, StripOrigin(origin, values[1]))

		case dns.TypePTR:
			fallthrough
		case dns.TypeNS:
			fallthrough
		case dns.TypeCNAME:
			fallthrough
		case dns.TypeNSEC:
			if len(values) < 1 {
				return fmt.Errorf("malformed RR: %v", values)
			}
			e.SetValue(0, StripOrigin(origin, values[0]))
		}

		if l := len(e.Domain()); l > longestname {
			longestname = l
		}
		if len(e.Domain()) > 0 {
			prevname = e.Domain()
		}
	}
	longestname += 2 // extra indent (we already take the origin into account)

	prevname = []byte{}
	prevtype := []byte{}
	prevttl := 0
	prevcom := false
	firstname := true
	for _, e := range zf.Entries() {
		if e.IsComment {
			if !prevcom && !firstname {
				fmt.Fprintln(w)
			}
			for _, c := range e.Comments() {
				fmt.Fprintf(w, "%s\n", c)
			}
			prevcom = true
			prevname = []byte{}
			prevtype = []byte{}
			continue
		}
		if e.IsControl {
			fmt.Fprintf(w, "%s %s\n", e.Command(), bytes.Join(e.Values(), []byte(" ")))
			prevcom = false
			prevname = []byte{}
			prevtype = []byte{}
			continue
		}

		if !bytes.Equal(prevname, e.Domain()) {
			// keep comments near, don't add a newline when previous line was comment.
			// first record doesn't need a newline
			if len(e.Domain()) > 0 && !prevcom && !firstname {
				v, _ := single[string(prevname)]
				// names /w multiple types get a newline
				if v > 1 {
					fmt.Fprintln(w)
				}
				// single type names together, except when types differ
				if v == 1 && !bytes.Equal(prevtype, e.Type()) {
					fmt.Fprintln(w)
				}
			}
			fmt.Fprintf(w, "%-*s", longestname, e.Domain())
		} else {
			fmt.Fprintf(w, "%-*s", longestname, "")
		}

		prevcom = false
		firstname = false

		if ttl := e.TTL(); ttl != nil && *ttl != prevttl {
			prevttl = *ttl
			fmt.Fprintf(w, "%10s", TimeToHuman(ttl))
		} else {
			fmt.Fprintf(w, "%10s", " ")
		}

		if len(e.Class()) > 0 {
			fmt.Fprintf(w, "%5s", e.Class())
		} else {
			fmt.Fprintf(w, "%5s", "IN")
		}
		fmt.Fprintf(w, "   %-8s", e.Type())

		// Specicial handling for certain RR types
		values := e.Values()
		switch e.RRType() {
		case dns.TypeTXT:
			fmt.Fprintf(w, Space3)
			space := ""
			// TODO: insert new lines when multiple blocks and longer then certain....
			for _, v := range values {
				fmt.Fprintf(w, "%s%q", space, v)
				space = " "
			}
			fmt.Fprintln(w)

		case dns.TypeCAA:
			fmt.Fprintf(w, Space3)
			space := ""
			for i, v := range values {
				if i < 2 {
					fmt.Fprintf(w, "%s%s", space, v)
				} else {
					fmt.Fprintf(w, "%s%q", space, v)
				}
				space = " "
			}
			fmt.Fprintln(w)

		case dns.TypeSOA:
			fmt.Fprintf(w, "%s%s (\n", Space3, bytes.Join(values[:2], []byte(" ")))
			for i, v := range values[2:] {
				if i == 0 {
					if *flagInc {
						v = Increase(v)
					}
					humandate := SerialToHuman(v)
					fmt.Fprintf(w, "%-*s%s%-13s%s%s\n", longestname+Indent, " ", Space3, v, soacomment[i], humandate)
				} else {
					fmt.Fprintf(w, "%-*s%s%-13s%s\n", longestname+Indent, " ", Space3, bytes.ToUpper(TimeToHumanByte(v)), soacomment[i])
				}
			}
			closeBrace(w, longestname)

		case dns.TypeTLSA:
			fallthrough
		case dns.TypeCDS, dns.TypeDS:
			fallthrough
		case dns.TypeCDNSKEY:
			fallthrough
		case dns.TypeDNSKEY:
			if len(values) < 4 {
				return fmt.Errorf("malformed RR: %v", values)
			}
			all := bytes.Join(values[3:], nil)
			pieces := Split(all, 55)
			if len(pieces) == 1 {
				fmt.Fprintf(w, "%s%s\n", Space3, bytes.Join(e.Values(), []byte(" ")))
				break
			}

			fmt.Fprintf(w, "%s%s (\n", Space3, bytes.Join(values[:3], []byte(" ")))
			for _, p := range pieces {
				fmt.Fprintf(w, "%-*s%s%-13s\n", longestname+Indent, " ", Space3, p)
			}
			closeBrace(w, longestname)

		case dns.TypeRRSIG:
			fmt.Fprintf(w, "%s%s (\n", Space3, bytes.Join(values[:8], []byte(" ")))
			all := bytes.Join(values[8:], nil)
			pieces := Split(all, 55)
			for _, p := range pieces {
				fmt.Fprintf(w, "%-*s%s%-13s\n", longestname+Indent, " ", Space3, p)
			}
			closeBrace(w, longestname)

		default:
			fmt.Fprintf(w, "%s%s\n", Space3, bytes.Join(values, []byte(" ")))
		}

		if len(e.Domain()) > 0 {
			prevname = e.Domain()
		}
		prevtype = e.Type()
	}
	return nil
}

const (
	Space3 = "   "
	Indent = 29
)

var soacomment = []string{"; serial", "; refresh", "; retry", "; expire", "; minimum"}

func closeBrace(w io.Writer, longestname int) {
	fmt.Fprintf(w, "%-*s)\n", longestname+Indent+3, " ")
}

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

func StripOrigin(origin, name []byte) []byte {
	if len(origin) > 0 && bytes.HasSuffix(name, origin) {
		// remove origin plus dot.
		l := len(name)
		if l == len(origin) {
			return []byte("@")
		} else {
			return name[:l-len(origin)-1]
		}
	}
	return name
}
