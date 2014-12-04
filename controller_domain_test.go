package hivdomainstatus

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for registrations

type DomainList struct {
	Total int
	Items []struct {
		Name         string
	}
}

func SetupDomainTest(t *testing.T) (ts *httptest.Server) {
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())

	cntrl := new(DomainController)
	cntrl.repo = NewDomainRepository(db)
	db.Exec("TRUNCATE domain RESTART IDENTITY")

	data := []string{"example.hiv", "acme.hiv"}
	for _,name := range data {
		d := new(Domain)
		d.Name = name
		cntrl.repo.Persist(d)
	}

	ts = httptest.NewServer(http.HandlerFunc(cntrl.ListingHandler))
	return
}

func TestThatItListsDomains(t *testing.T) {
	assert := assert.New(t)

	ts := SetupDomainTest(t)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("application/json", res.Header.Get("Content-Type"))

	var l DomainList
	unmarshalErr := json.Unmarshal(b, &l)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	assert.Equal(2, l.Total)

	assert.Equal("example.hiv", l.Items[0].Name)
	assert.Equal("acme.hiv", l.Items[1].Name)

	assert.Equal(`</domain?offsetKey=2>; rel="next"`, res.Header.Get("Link"))
}

func TestThatItReturnsNextUrlAfterEndOfDomains(t *testing.T) {
	assert := assert.New(t)

	ts := SetupDomainTest(t)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/domain?offsetKey=2")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(`</domain?offsetKey=2>; rel="next"`, res.Header.Get("Link"))

}
