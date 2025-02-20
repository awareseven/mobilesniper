# Port for SDC discovery over UDP multicasting.
UDP_MULTICAST_PORT = 3702

allow-udp-multicast:
	sudo nft add rule ip qubes custom-input udp dport $(UDP_MULTICAST_PORT) ct state new,established,related counter accept