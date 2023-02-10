package reporter

import "net/http"

type Config struct {
	httpClient *http.Client

	DSN         string `json:"dsn"`
	Debug       bool   `json:"debug"`
	Disabled    bool   `json:"disabled"`
	Environment string `json:"-"`
	Release     string `json:"-"`
}
