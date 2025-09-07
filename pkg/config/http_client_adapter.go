package config

import "net/http"

type HttpClientAdapter struct {
	*http.Client
}
