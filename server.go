package hivdomainstatus

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type route struct {
	re      *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request, []string)
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) AddRoute(re string, handler func(http.ResponseWriter, *http.Request, []string)) {
	r := &route{regexp.MustCompile(re), handler}
	h.routes = append(h.routes, r)
}

func (h *RegexpHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		matches := route.re.FindStringSubmatch(r.URL.Path)
		if matches != nil {
			route.handler(rw, r, matches)
			break
		}
	}
}

func Serve(c *Config) (err error) {
	// Open DB
	db, err := sql.Open("postgres", c.DSN())
	if err != nil {
		return
	}

	log.Println(fmt.Sprintf("Starting server on localhost:%d ...", c.Server.Port))

	domainCntrl := new(DomainController)
	domainCntrl.domainRepo = NewDomainRepository(db)
	domainCntrl.domainCheckRepo = NewDomainCheckRepository(db)
	entryPointCntrl := new(EntryPointController)

	reHandler := new(RegexpHandler)
	reHandler.AddRoute("^/domain/([0-9]+)$", domainCntrl.ItemHandler)
	reHandler.AddRoute("^/domain$", domainCntrl.ListingHandler)
	reHandler.AddRoute("^/$", entryPointCntrl.EntryPointHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", c.Server.Port), reHandler))

	return
}
