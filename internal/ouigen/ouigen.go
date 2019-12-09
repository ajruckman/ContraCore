package ouigen

import (
    "bufio"
    "fmt"
    "net"
    "net/http"
    "regexp"

    . "github.com/ajruckman/xlib"
    "github.com/jackc/pgx"

    "github.com/ajruckman/ContraCore/internal/db"
)

var matchOUI = regexp.MustCompile(`^([A-z0-9]{2}-[A-z0-9]{2}-[A-z0-9]{2})\s+\(hex\)\s+(.*)$`)

func GenOUI() {
    _, err := db.PDB.Exec(`TRUNCATE TABLE oui;`)
    Err(err)

    var res [][]interface{}

    resp, err := http.Get(`https://linuxnet.ca/ieee/oui.txt`)
    Err(err)
    s := bufio.NewScanner(resp.Body)

    for s.Scan() {
        t := s.Text()

        if matchOUI.MatchString(t) {
            m := matchOUI.FindStringSubmatch(t)
            mac, err := net.ParseMAC(m[1] + "-00-00-00")
            Err(err)
            vendor := m[2]

            fmt.Println(mac, "->", vendor)

            res = append(res, []interface{}{mac, vendor})
        }
    }

    _, err = db.PDB.CopyFrom(pgx.Identifier{"oui"}, []string{"mac", "vendor"}, pgx.CopyFromRows(res))
    Err(err)
}
