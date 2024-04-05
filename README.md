# Opiniated DNS zone file formatter

This tools (re)formats zone files keeping comments and $-pragmas intact.
It does remove in RR comments, i.e '; Serial' and friends.

It strips unnecessary origins from names. Every new 'set' of names will get a newline.
If a name is repeated for a different type, the name is stripped.

TTLs are converted back to human form.

The order of names is kept.

DNSSEC

Formatting of keys

Long records...
