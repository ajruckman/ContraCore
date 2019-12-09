package schema

import (
    "net"
)

type OUI struct {
    MAC    net.HardwareAddr `db:"mac"`
    Vendor string           `db:"vendor"`
}
