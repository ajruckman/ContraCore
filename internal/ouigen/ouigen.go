package ouigen

import (
    "bufio"
    "net/http"
    "regexp"
    "strings"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db"
)

var matchOUI = regexp.MustCompile(`^([A-z0-9]{2}-[A-z0-9]{2}-[A-z0-9]{2})\s+\(hex\)\s+(.*)$`)

func GenOUI() {
    _, err := db.CDB.Exec(`TRUNCATE TABLE contralog.oui;`)
    Err(err)

    tx, err := db.CDB.Begin()
    Err(err)
    stmt, err := tx.Prepare(`INSERT INTO contralog.oui (mac, vendor) VALUES (?, ?);`)
    Err(err)

    resp, err := http.Get(`https://linuxnet.ca/ieee/oui.txt`)
    Err(err)
    s := bufio.NewScanner(resp.Body)

    for s.Scan() {
        t := s.Text()

        if matchOUI.MatchString(t) {
            m := matchOUI.FindStringSubmatch(t)

            mac := strings.ToLower(strings.Replace(m[1], "-", ":", -1))
            vendor := m[2]

            _, err = stmt.Exec(mac, vendor)
            Err(err)
        }
    }

    err = tx.Commit()
    Err(err)
}
