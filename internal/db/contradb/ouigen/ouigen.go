package ouigen

import (
    "bufio"
    "fmt"
    "net/http"
    "regexp"
    "strings"

    . "github.com/ajruckman/xlib"
    "github.com/jackc/pgx/v4"

    "github.com/ajruckman/ContraCore/internal/db/contradb"
)

var matchOUI = regexp.MustCompile(`^([A-z0-9]{2}-[A-z0-9]{2}-[A-z0-9]{2})\s+\(hex\)\s+(.*)$`)

func GenOUI() {
    _, err := contradb.Exec(`TRUNCATE TABLE oui;`)
    Err(err)

    var res [][]interface{}

    resp, err := http.Get(`https://linuxnet.ca/ieee/oui.txt`)
    Err(err)
    s := bufio.NewScanner(resp.Body)

    for s.Scan() {
        t := s.Text()

        if matchOUI.MatchString(t) {
            m := matchOUI.FindStringSubmatch(t)
            mac := strings.ToLower(strings.Replace(m[1], "-", ":", -1))
            vendor := m[2]

            fmt.Println(mac, "->", vendor)

            res = append(res, []interface{}{mac, vendor})
        }
    }

    _, err = contradb.CopyFrom(pgx.Identifier{"oui"}, []string{"mac", "vendor"}, pgx.CopyFromRows(res))
    Err(err)
}
