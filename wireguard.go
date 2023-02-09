package waygate

import (
	"fmt"
)

type WireGuardConfig struct {
	Address    string           `json:"address"`
	PrivateKey string           `json:"private_key"`
	ListenPort int              `json:"listen_port"`
	Peers      []*WireGuardPeer `json:"peers"`
}

type WireGuardPeer struct {
	PublicKey           string   `json:"public_key"`
	AllowedIps          []string `json:"AllowedIps"`
	Endpoint            string   `json:"endpoint"`
	PersistentKeepalive int      `json:"persistent_keepalive"`
}

func (c *WireGuardConfig) String() string {
	s := "[Interface]\n"

	if c.Address != "" {
		s += fmt.Sprintf("Address = %s\n", c.Address)
	}

	s += fmt.Sprintf("PrivateKey = %s\n", c.PrivateKey)

	if c.ListenPort != 0 {
		s += fmt.Sprintf("ListenPort = %d\n", c.ListenPort)
	}

	s += "\n"

	for _, peer := range c.Peers {
		s += "[Peer]\n"
		s += fmt.Sprintf("PublicKey = %s\n", peer.PublicKey)

		s += "AllowedIPs = "
		for i, ip := range peer.AllowedIps {
			s += fmt.Sprintf("%s", ip)
			if i != len(peer.AllowedIps)-1 {
				s += ","
			}
		}
		s += "\n"

		if peer.Endpoint != "" {
			s += fmt.Sprintf("Endpoint = %s\n", peer.Endpoint)
		}

		if peer.PersistentKeepalive != 0 {
			s += fmt.Sprintf("PersistentKeepalive = %d\n", peer.PersistentKeepalive)
		}

		s += "\n"
	}

	return s
}
