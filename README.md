# Opiniated DNS zone file formatter

Works, but WIP.

This tools (re)formats zone files keeping comments and $-pragmas intact. It does remove in-RR
comments, i.e '; Serial' and friends, although for SOA records these get added back.

**dnsfmt** is a filter.

* ordering of the zone is left as-is
* unnecessary origins from names are stripped
* repeated ownernames are suppressed
* TTLs are _all_ converted to human readable form
* long records (DNSKEYs, RRSIGs) are wrapped and placed in braces

No semantic checks are done, this is purely text manipulation with some basic zone file syntax
understanding.
