$TTL    6H
$ORIGIN atoom.net.
@       IN      SOA     linode.atoom.net. miek.miek.nl. (
                       1282630067       ; Serial
                             4H         ; Refresh
                             1H         ; Retry
                             7D         ; Expire
                             4H )       ; Negative Cache TTL
                IN      NS  linode
		IN	NS  ns-ext.nlnetlabs.nl.
		IN	NS  omval.tednet.nl.
                IN      A       45.138.52.215
                IN      AAAA    2a10:3781:2dc2:3::53

linode          IN      A       45.138.52.215
                IN      AAAA    2a10:3781:2dc2:3::53

nuc             IN      A       45.138.52.215
                IN      AAAA    2a10:3781:2dc2:3::53

www         	IN  	CNAME   nuc
g         IN  	CNAME   nuc

; bramk
; nog een comment
google.com.		1394	IN	TXT	"google-site-verification=wD8N7i1JTNTkezJ49swvWW48f8_9xveREV4oB-0Hf5o"
lafhart         IN      A       178.79.160.171
voordeur        IN      A       77.249.87.46
