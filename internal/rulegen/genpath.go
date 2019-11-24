package rulegen

import (
    "strings"
)

func GenPath(domain string) []string {
    dcs := strings.Split(domain, ".")

    for i := len(dcs)/2 - 1; i >= 0; i-- {
        opp := len(dcs) - 1 - i
        dcs[i], dcs[opp] = dcs[opp], dcs[i]
    }

    //if len(dcs) < 2 {
    //    return nil
    //}

    return dcs
}

