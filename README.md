# Opiniated DNS zone file formatter

This tools (re)formats zone files keeping comments and $-pragmas intact. It does remove in-RR
comments, i.e '; serial' and friends, although for SOA records these get added back.

Builds up-on: https://github.com/bwesterb/go-zonefile which is butchered and vendored in ./zonefile.
(Only needed half of the functionality and comments weren't fleshed out.)

**dnsfmt** is a filter. See dnsfmt.1.md for more information.

[![asciicast](https://asciinema.org/a/E5B2d7lfDV0X17wMkL5ouoybD.svg)](https://asciinema.org/a/E5B2d7lfDV0X17wMkL5ouoybD)

## Why not miekg/dns?

Pondered this, and yes, it has a better parser, but then re-arranging the []dns.RR and pretty
printing would have been (IMO) more work. Also miekg/dns does not have an option to leave
$-directives as-is.
