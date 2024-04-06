# Opiniated DNS zone file formatter

This tools (re)formats zone files keeping comments and $-pragmas intact. It does remove in-RR
comments, i.e '; Serial' and friends, although for SOA records these get added back.

Builds up-on: https://github.com/bwesterb/go-zonefile which is butchered and vendored in ./zonefile.
(Only needed half of the functionality and comments weren't fleshed out.)

**dnsfmt** is a filter. See dnsfmt.1.md for more information.

Needs some tests. Maybe some more newlines in the correct places. Also strip origin from more
ownernames (right hand side: CNAME, NS, etc.).

Random ideas: sort record type per name, so this is consitent throughout the file. Group similar
ownernames.
