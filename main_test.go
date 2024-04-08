package main

import (
	"bytes"
	"testing"
)

func TestFormat(t *testing.T) {
	const mess = `$TTL    6H
$ORIGIN example.org.
@       IN      SOA     ns miek.miek.nl. 1282630067  4H 1H 7D 4H
                IN      NS  ns
example.org.		IN	NS  ns.example.org
`
	out := &bytes.Buffer{}
	Reformat([]byte(mess), nil, out)
	if out.String() != `$TTL 6H
$ORIGIN example.org.
@                 IN   SOA        ns miek.miek.nl. (
                                     1282630067   ; serial
                                     4H           ; refresh
                                     1H           ; retry
                                     7D           ; expire
                                     4H           ; minimum
                                     )
                  IN   NS         ns
                  IN   NS         ns.example.org
` {
		t.Fatalf("failed to properly reformat\n%s\n", out.String())
	}
}

func TestFormatCommentStart(t *testing.T) {
	const mess = `; example.nl,v 1.00 2015/03/19 14:31:47 root Exp
$ORIGIN example.nl.
`
	out := &bytes.Buffer{}
	Reformat([]byte(mess), nil, out)
	if out.String() != `; example.nl,v 1.00 2015/03/19 14:31:47 root Exp
$ORIGIN example.nl.
` {
		t.Fatalf("failed to properly reformat\n%s\n", out.String())
	}
}
