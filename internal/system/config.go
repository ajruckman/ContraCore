package system

// Config values as loaded from Corefiles or from the ContraDB database.
var (
	ContraDBURL  string // The connection string to ContraDB.
	ContraLogURL string // The connection string to ContraLog.

	RuleSources   []string // URLs from which to generate blacklist rules.
	SearchDomains []string // DNS query suffixes.
	DomainNeeded  bool     // If true, queries without dots or domain parts will never be forwarded to upstream servers.

	// Fake values with which to respond to blocked DNS queries.
	SpoofedA       string
	SpoofedAAAA    string
	SpoofedCNAME   string
	SpoofedDefault string
)
