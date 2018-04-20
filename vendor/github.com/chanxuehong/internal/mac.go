package internal

import (
	"bytes"
	"net"

	"github.com/chanxuehong/rand"
)

var MAC [6]byte = getMAC() // One MAC of this machine; Particular case, it is a random bytes.

var zeroMAC [8]byte

func getMAC() (mac [6]byte) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return genMAC()
	}

	// Gets a MAC from interfaces of this machine,
	// the MAC of up state interface is preferred.
	found := false // Says it has found a MAC
	for _, itf := range interfaces {
		if itf.Flags&net.FlagLoopback == net.FlagLoopback ||
			itf.Flags&net.FlagPointToPoint == net.FlagPointToPoint {
			continue
		}

		switch len(itf.HardwareAddr) {
		case 6: // MAC-48, EUI-48
			if bytes.Equal(itf.HardwareAddr, zeroMAC[:6]) {
				continue
			}
			if itf.Flags&net.FlagUp == 0 {
				if !found {
					copy(mac[:], itf.HardwareAddr)
					found = true
				}
				continue
			}
			copy(mac[:], itf.HardwareAddr)
			return
		case 8: // EUI-64
			if bytes.Equal(itf.HardwareAddr, zeroMAC[:]) {
				continue
			}
			if itf.Flags&net.FlagUp == 0 {
				if !found {
					copy(mac[:3], itf.HardwareAddr)
					copy(mac[3:], itf.HardwareAddr[5:])
					found = true
				}
				continue
			}
			copy(mac[:3], itf.HardwareAddr)
			copy(mac[3:], itf.HardwareAddr[5:])
			return
		}
	}
	if found {
		return
	}

	return genMAC()
}

// generates a random MAC.
func genMAC() (mac [6]byte) {
	rand.Read(mac[:])
	mac[0] |= 0x01 // multicast
	return
}
