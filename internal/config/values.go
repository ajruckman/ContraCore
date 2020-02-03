package config

var (
    ContraDBURL  string
    ContraLogURL string

    RuleSources    []string
    SearchDomains  []string
    DomainNeeded   bool
    SpoofedA       string
    SpoofedAAAA    string
    SpoofedCNAME   string
    SpoofedDefault string
)
