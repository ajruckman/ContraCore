package dbschema

import (
	"net"
	"time"
)

type Whitelist struct {
	ID        int
	Pattern   string
	Expires   *time.Time
	IPs       *[]net.IP
	Subnets   *[]net.IPNet
	MACs      *[]net.HardwareAddr
	Vendors   *[]string
	Hostnames *[]string
}
