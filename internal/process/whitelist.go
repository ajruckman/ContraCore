package process

import (
	"github.com/ajruckman/ContraCore/internal/cache"
	"github.com/ajruckman/ContraCore/internal/system"
)

func whitelist(q *queryContext) (ret bool) {
	if cache.WhitelistCache.Check(q._domain, q._client, q.mac, q.hostname, q.vendor) {
		q.action = ActionWhitelisted
		system.Console.Infof("Whitelisting query %d", q.r.Id)
		return true
	}

	return false
}
