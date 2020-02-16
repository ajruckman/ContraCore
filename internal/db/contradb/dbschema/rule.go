package dbschema

import (
    "net"
    "time"
)

type Blacklist struct {
    ID      int    `db:"id"`
    Pattern string `db:"pattern"`
    Class   int    `db:"class"`
    Domain  string `db:"domain"`
    TLD     string `db:"tld"`
    SLD     string `db:"sld"`
}

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
