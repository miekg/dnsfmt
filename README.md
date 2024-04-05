# Opiniated DNS zone file formatter

This tools (re)formats zone files keeping comments and $-pragmas intact.
It also sorts the zone file on:

* label count, and then
* alphabetically within a set a of names with the same labelcount

Furhter more:

* origin is removed from each name
* per label-count block the....

DNSSEC
