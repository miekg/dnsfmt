%%%
title = "dnsfmt 1"
area = "System Administration"
workgroup = "DNS"
%%%

# NAME

dnsfmt - format DNS zone files

# SYNOPSIS

**dnsfmt** [**FILE**]...

# DESCRIPTION

**Dnsfmt** formats zone file from **FILE**. If no file is given, it reads from standard input.

The zone is formatted according to the following rules:

* ordering of the zone is left as-is
* all whitespace is removed
* unnecessary origins from names are stripped
* a new comments gets an empty line before it
* a new ownername gets an empty line before it
* repeated ownernames are suppressed
* TTLs are _all_ converted to human readable form (on minute accuracy) when they are larger than 600
* long records (DNSKEYs, RRSIGs) are wrapped and placed in braces
* names with only one, but equal, type are grouped together without newlines
* the SOA serial comment gets a written out timestamp

No semantic checks are done, this is purely text manipulation with some basic zone file syntax
understanding.

# OPTIONS

`-o` **ORIGIN**
: begin parsing with origin set to **ORIGIN**

`-i`
: increase the serial, for epoch serial the current time is used, for date+sequence serial it is
  just increased by one, defaults to true

# EXAMPLE

    % cat <<'EOF' | ./dnsfmt
    $TTL 6H
    $ORIGIN example.org.
    @       IN      SOA     ns miek.miek.nl. 1282630067  4H 1H 7D 7200
                    IN      NS  ns
    example.org.            IN      NS  ns-ext.nlnetlabs.nl.
    EOF

Returns:

    $TTL 6H
    $ORIGIN example.org.
    @               IN   SOA        ns miek.miek.nl. (
                                       1712997354   ; serial  Sat, 13 Apr 2024 08:35:54 UTC
                                       4H           ; refresh
                                       1H           ; retry
                                       7D           ; expire
                                       2H           ; minimum
                                       )
                    IN   NS         ns
                    IN   NS         ns-ext.nlnetlabs.nl.


# AUTHOR

Miek Gieben <miek@miek.nl>.
