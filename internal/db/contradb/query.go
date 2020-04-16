package contradb

import (
	"context"
	"net"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"

	contradbschema "github.com/ajruckman/ContraCore/internal/db/contradb/dbschema"
	"github.com/ajruckman/ContraCore/internal/system"
)

func GetLeaseDetails() (res []contradbschema.LeaseDetailsByIP, err error) {
	if !system.ContraDBOnline.Load() {
		return nil, &ErrContraDBOffline{}
	}

	var rows *sqlx.Rows

	rows, err = xdb.Queryx(`SELECT time, op, ip, mac, hostname, vendor FROM lease_details;`)
	if err != nil {
		return nil, errOfflineOrOriginal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var n = _leaseDetails{}
		err = rows.StructScan(&n)
		if err != nil {
			return nil, err
		}

		res = append(res, contradbschema.LeaseDetailsByIP{
			Time:     n.Time,
			IP:       n.IP.IPNet.IP,
			MAC:      n.MAC.Addr,
			Hostname: n.Hostname,
			Vendor:   n.Vendor,
		})
	}

	return
}

type _leaseDetails struct {
	Time     time.Time      `db:"time"`
	Op       string         `db:"op"`
	IP       pgtype.Inet    `db:"ip"`
	MAC      pgtype.Macaddr `db:"mac"`
	Hostname *string        `db:"hostname"`
	Vendor   *string        `db:"vendor"`
}

func GetOUI() (res []contradbschema.OUI, err error) {
	if !system.ContraDBOnline.Load() {
		return nil, &ErrContraDBOffline{}
	}

	err = xdb.Select(&res, `SELECT * FROM oui;`)
	return res, errOfflineOrOriginal(err)
}

func GetConfig() (res contradbschema.Config, err error) {
	if !system.ContraDBOnline.Load() {
		return res, &ErrContraDBOffline{}
	}

	err = xdb.Get(&res, `SELECT * FROM config ORDER BY id DESC LIMIT 1`)
	return res, errOfflineOrOriginal(err)
}

func GetBlacklistRules() (res []contradbschema.Blacklist, err error) {
	if !system.ContraDBOnline.Load() {
		return nil, &ErrContraDBOffline{}
	}

	err = xdb.Select(&res, `SELECT * FROM blacklist;`)
	return res, errOfflineOrOriginal(err)
}

func GetWhitelistRules() (res []contradbschema.Whitelist, err error) {
	if !system.ContraDBOnline.Load() {
		return nil, &ErrContraDBOffline{}
	}

	var rows *sqlx.Rows

	rows, err = xdb.Queryx(`SELECT * FROM whitelist;`)
	if err != nil {
		return nil, errOfflineOrOriginal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var n = _whitelist{}
		err = rows.StructScan(&n)
		if err != nil {
			return nil, errOfflineOrOriginal(err)
		}

		var expires *time.Time
		if n.Expires != nil {
			expires = &n.Expires.Time
		}

		var ips *[]net.IP
		if n.IPs != nil {
			ips = &[]net.IP{}
			for _, ip := range n.IPs.Elements {
				*ips = append(*ips, ip.IPNet.IP)
			}
		}

		var subnets *[]net.IPNet
		if n.Subnets != nil {
			subnets = &[]net.IPNet{}
			for _, subnet := range n.Subnets.Elements {
				*subnets = append(*subnets, *subnet.IPNet)
			}
		}

		var macs *[]net.HardwareAddr
		if n.MACs != nil {
			macs = &[]net.HardwareAddr{}
			for _, mac := range n.MACs.Elements {
				*macs = append(*macs, mac.Addr)
			}
		}

		var vendors *[]string
		if n.Vendors != nil {
			vendors = &[]string{}
			for _, vendor := range n.Vendors.Elements {
				*vendors = append(*vendors, vendor.String)
			}
		}

		var hostnames *[]string
		if n.Hostnames != nil {
			hostnames = &[]string{}
			for _, hostname := range n.Hostnames.Elements {
				*hostnames = append(*hostnames, hostname.String)
			}
		}

		res = append(res, contradbschema.Whitelist{
			ID:        n.ID,
			Pattern:   n.Pattern,
			Expires:   expires,
			Creator:   n.Creator,
			IPs:       ips,
			Subnets:   subnets,
			MACs:      macs,
			Vendors:   vendors,
			Hostnames: hostnames,
		})
	}

	return
}

type _whitelist struct {
	ID        int                  `db:"id"`
	Pattern   string               `db:"pattern"`
	Expires   *pgtype.Timestamp    `db:"expires"`
	Creator   *string              `db:"creator"`
	IPs       *pgtype.InetArray    `db:"ips"`
	Subnets   *pgtype.CIDRArray    `db:"subnets"`
	MACs      *pgtype.MacaddrArray `db:"macs"`
	Vendors   *pgtype.TextArray    `db:"vendors"`
	Hostnames *pgtype.TextArray    `db:"hostnames"`
}

func Exec(query string, args ...interface{}) (cmd pgconn.CommandTag, err error) {
	if !system.ContraDBOnline.Load() {
		return cmd, &ErrContraDBOffline{}
	}

	cmd, err = pdb.Exec(context.Background(), query, args...)
	return cmd, errOfflineOrOriginal(err)
}

func CopyFrom(tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (numRows int64, err error) {
	if !system.ContraDBOnline.Load() {
		return 0, &ErrContraDBOffline{}
	}

	numRows, err = pdb.CopyFrom(context.Background(), tableName, columnNames, rowSrc)
	return numRows, errOfflineOrOriginal(err)
}

func Select(dest interface{}, query string, args ...interface{}) (err error) {
	if !system.ContraDBOnline.Load() {
		return &ErrContraDBOffline{}
	}

	err = xdb.Select(dest, query, args...)
	return errOfflineOrOriginal(err)
}

func insertDefaultConfig() (err error) {
	_, err = xdb.Exec(`INSERT INTO config (search_domains) VALUES(default);`)
	return errOfflineOrOriginal(err)
}
