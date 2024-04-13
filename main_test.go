package main

import (
	"bytes"
	"testing"
)

func TestFormat(t *testing.T) {
	const mess = `$TTL    6H
$ORIGIN example.org.
@       IN      SOA     ns miek.miek.nl. 1282630067  4H 1H 7D 7200
                IN      NS  ns
example.org.		IN	NS  ns.example.org.
`
	out := &bytes.Buffer{}
	Reformat([]byte(mess), nil, out)
	if out.String() != `$TTL 6H
$ORIGIN example.org.
@               IN   SOA        ns miek.miek.nl. (
                                   1282630067   ; serial  Tue, 24 Aug 2010 06:07:47 UTC
                                   4H           ; refresh
                                   1H           ; retry
                                   7D           ; expire
                                   2H           ; minimum
                                   )
                IN   NS         ns
                IN   NS         ns
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

func TestFormatKeepTogether(t *testing.T) {
	const mess = `$ORIGIN miek.nl.
@       IN      SOA     linode.miek.nl. miek.miek.nl. (
			     1282630063 ; Serial
                             4H         ; Refresh
                             1H         ; Retry
                             7D         ; Expire
                             4H )       ; Negative Cache TTL
                IN      NS      linode.atoom.net.

                IN      MX      10 aspmx3.googlemail.com.

                IN      A       127.0.0.1

a               IN      A       127.0.0.1
                IN      AAAA    1::53

mmark           IN      CNAME   a

bot             IN      CNAME   a

www             IN      CNAME   a
go.dns          IN      TXT     "Hello DNS developer!"
x               IN      CNAME   a

nlgids          IN      CNAME   a
`
	out := &bytes.Buffer{}
	Reformat([]byte(mess), nil, out)
	if out.String() != `$ORIGIN miek.nl.
@                    IN   SOA        linode miek (
                                        1282630063   ; serial  Tue, 24 Aug 2010 06:07:43 UTC
                                        4H           ; refresh
                                        1H           ; retry
                                        7D           ; expire
                                        4H           ; minimum
                                        )
                     IN   NS         linode.atoom.net.
                     IN   MX         10 aspmx3.googlemail.com.
                     IN   A          127.0.0.1

a                    IN   A          127.0.0.1
                     IN   AAAA       1::53

mmark                IN   CNAME      a
bot                  IN   CNAME      a
www                  IN   CNAME      a

go.dns               IN   TXT        "Hello DNS developer!"

x                    IN   CNAME      a
nlgids               IN   CNAME      a
` {
		t.Fatalf("failed to properly reformat\n%s\n", out.String())
	}
}
