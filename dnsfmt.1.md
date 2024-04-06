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

**Dnsfmt** formats zone file. The zonefile must be piped to **dnsfmt**.
There are no options.

# EXAMPLE

    $TTL    6H
    $ORIGIN example.org.
    @       IN      SOA     ns miek.miek.nl. 1282630067  4H 1H 7D 4H
                    IN      NS  ns
    example.org.		IN	NS  ns-ext.nlnetlabs.nl.
