package network

type Protocol byte

const (
	TCP Protocol = 1
	UDP          = iota << 1
)

func IsTCP(protocol Protocol) bool {
	if (protocol & TCP) == TCP {
		return true
	}

	return false
}

func IsUDP(protocol Protocol) bool {
	if (protocol & UDP) == UDP {
		return true
	}

	return false
}
