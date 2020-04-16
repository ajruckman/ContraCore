package dbschema

import (
	"net"
	"time"
)

type LeaseDetailsByIPHostname struct {
	Time     time.Time
	IP       net.IP
	MAC      net.HardwareAddr
	Hostname *string
	Vendor   *string
}
