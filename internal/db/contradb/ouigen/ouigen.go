// Package ouigen loads MAC OUIs into ContraDB so they can be used to look up
// device vendors by their MAC address.
package ouigen

import (
	"bufio"
	"net/http"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v4"

	"github.com/ajruckman/ContraCore/internal/db/contradb"
	"github.com/ajruckman/ContraCore/internal/functions"
)

// Matches a line of the form: 00-D0-EF   (hex)		IGT
var matchOUI = regexp.MustCompile(`^([A-z0-9]{2}-[A-z0-9]{2}-[A-z0-9]{2})\s+\(hex\)\s+(.*)$`)

const ouiSource = `https://linuxnet.ca/ieee/oui.txt`

// Loads MAC OUIs into ContraDB.
func GenOUI(callback functions.ProgressCallback) {
	callback("Getting OUI list from: " + ouiSource)
	resp, err := http.Get(ouiSource)
	if err != nil {
		callback(err.Error())
		return
	}
	s := bufio.NewScanner(resp.Body)

	_, err = contradb.Exec(`TRUNCATE TABLE oui;`)
	if err != nil {
		callback("Failed to truncate OUI table: " + err.Error())
		return
	}

	var res [][]interface{}

	for s.Scan() {
		t := s.Text()

		if matchOUI.MatchString(t) {
			m := matchOUI.FindStringSubmatch(t)
			mac := strings.ToLower(strings.Replace(m[1], "-", ":", -1))
			vendor := m[2]

			if callback(mac + " -> " + vendor) {
				return
			}

			res = append(res, []interface{}{mac, vendor})
		}
	}

	_, err = contradb.CopyFrom(pgx.Identifier{"oui"}, []string{"mac", "vendor"}, pgx.CopyFromRows(res))
	if err != nil {
		callback(err.Error())
		return
	}
}
