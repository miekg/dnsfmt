# Opiniated DNS zone file formatter

This tools (re)formats zone files keeping comments and $-pragmas intact. It does remove in-RR
comments, i.e '; serial' and friends, although for SOA records these get added back.

Builds up-on: https://github.com/bwesterb/go-zonefile which is butchered and vendored in ./zonefile.
(Only needed half of the functionality and comments weren't fleshed out.)

**dnsfmt** is a filter. See dnsfmt.1.md for more information.

See this [screencast](https://asciinema.org/a/E5B2d7lfDV0X17wMkL5ouoybD).
