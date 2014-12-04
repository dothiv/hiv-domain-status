package hivdomainstatus

import (
	"database/sql"
	"fmt"
	"net/http"
	"log"
)

func Serve(c *Config) (err error) {
	// Open DB
	db, err := sql.Open("postgres", c.DSN())
	if err != nil {
		return
	}

	log.Println(fmt.Sprintf("Starting server on localhost:%d ...", c.Server.Port))

	domainCntrl := new(DomainController)
	domainCntrl.repo = NewDomainRepository(db)
	http.HandleFunc("/domain", domainCntrl.ListingHandler)

	entryPointCntrl := new(EntryPointController)
	http.HandleFunc("/", entryPointCntrl.EntryPointHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", c.Server.Port), nil))
	
	return
}
