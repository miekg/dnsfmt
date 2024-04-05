# Opiniated DNS zone file formatter

This tools (re)formats zone files keeping comments and $-pragmas intact.
It also sorts the zone file on:

* label count, and then
* and type, A, AAAA first, then the rest.

Furhter more:

* origin is removed from each name
* per label-count block the....

DNSSEC

Formatting of keys

Comment inside RRs are discarded (the SOA ;Serial thing).
