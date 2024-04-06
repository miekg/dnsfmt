%%%
title = "dnsfmt 1"
area = "System Administration"
workgroup = "DNS"
%%%

# NAME

dnsfmt - format DNS zone files

# SYNOPSIS

**dnsfmt**

# DESCRIPTION

**Dnsfmt** formats zone file. The zonefile must be piped to **dnsfmt**. There are no options.

The zone is formatted according to the following rules:

* ordering of the zone is left as-is
* unnecessary origins from names are stripped
* repeated ownernames are suppressed
* TTLs are _all_ converted to human readable form
* long records (DNSKEYs, RRSIGs) are wrapped and placed in braces

No semantic checks are done, this is purely text manipulation with some basic zone file syntax
understanding.

# EXAMPLE

    % cat <<'EOF' | ./dnsfmt
    $TTL 6H
    $ORIGIN example.org.
    @       IN      SOA     ns miek.miek.nl. 1282630067  4H 1H 7D 4H
                    IN      NS  ns
    example.org.            IN      NS  ns-ext.nlnetlabs.nl.
    EOF

Returns:

    $TTL 6H
    $ORIGIN example.org.
    @                 IN   SOA        ns miek.miek.nl. (
                                         1282630067   ; serial
                                         4H           ; refresh
                                         1H           ; retry
                                         7D           ; expire
                                         4H           ; minimum
                                      )
                      IN   NS         ns
                      IN   NS         ns-ext.nlnetlabs.nl.

# AUTHOR

Miek Gieben <miek@miek.nl>.
