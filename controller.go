package hivdomainstatus

import (
	"fmt"
	"net/http"
)

func getHttpHost(r *http.Request) string {
	proto := "https"
	if r.TLS == nil {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s", proto, r.Host)
}
