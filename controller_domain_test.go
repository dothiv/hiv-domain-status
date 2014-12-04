package hivdomainstatus

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"regexp"
	"fmt"

	"github.com/stretchr/testify/assert"
)

// Test for registrations

type DomainList struct {
	Total int
	Items []struct {
		Name         string
	}
}

func SetupDomainTest(t *testing.T) (cntrl *DomainController) {
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())

	cntrl = new(DomainController)
	cntrl.repo = NewDomainRepository(db)
	db.Exec("TRUNCATE domain RESTART IDENTITY")

	data := []string{"example.hiv", "acme.hiv"}
	for _,name := range data {
		d := new(Domain)
		d.Name = name
		cntrl.repo.Persist(d)
	}
	return
}

func TestThatItListsDomains(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.ListingHandler(w, r, nil)
    }))
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

	assert.Equal(`<` + ts.URL + `/domain?offsetKey=2>; rel="next"`, res.Header.Get("Link"))
}

func TestThatItReturnsNextUrlAfterEndOfDomains(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.ListingHandler(w, r, nil)
    }))
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

	assert.Equal(`<` + ts.URL + `/domain?offsetKey=2>; rel="next"`, res.Header.Get("Link"))
}

func TestThatItAddsNewDomain(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.ListingHandler(w, r, nil)	
    }))
    defer ts.Close()

	var data = []byte(`{"name":"example.hiv"}`)
	res, err := http.Post(ts.URL + "/domain", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(http.StatusCreated, res.StatusCode)
}

func TestThatItFetchesNewDomain(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.ItemHandler(w, r, regexp.MustCompile("^/domain/([0-9]+)$").FindStringSubmatch(r.URL.String()))	
    }))
    defer ts.Close()

    d := new(Domain)
    d.Name = "example.hiv"
    cntrl.repo.Persist(d)

	// Fetch the new domain
	fetchRes, fetchErr := http.Get(fmt.Sprintf("%s/domain/%d", ts.URL, d.Id))
	if fetchErr != nil {
		t.Fatal(fetchErr)
	}
	assert.Equal(http.StatusOK, fetchRes.StatusCode)

	b, readErr := ioutil.ReadAll(fetchRes.Body)
	fetchRes.Body.Close()
	if readErr != nil {
		t.Fatal(readErr)
	}
	assert.Equal("application/json", fetchRes.Header.Get("Content-Type"))

	var m DomainModel
	unmarshalErr := json.Unmarshal(b, &m)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	assert.Equal("example.hiv", m.Name)
	assert.Equal(fmt.Sprintf("%s/domain/%d", ts.URL, d.Id), m.JsonLDId)
	assert.Equal("http://jsonld.click4life.hiv/Domain", m.JsonLDContext)
}

func TestThatItDeletesDomain(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.ItemHandler(w, r, regexp.MustCompile("^/domain/([0-9]+)$").FindStringSubmatch(r.URL.String()))	
    }))
    defer ts.Close()

    d := new(Domain)
    d.Name = "example.hiv"
    cntrl.repo.Persist(d)

	// Delete
	deleteReq, deleteReqErr := http.NewRequest("DELETE", fmt.Sprintf("%s/domain/%d", ts.URL, d.Id), nil)
	if deleteReqErr != nil {
		t.Fatal(deleteReqErr)
	}
	client := &http.Client{}
	deleteRes, deleteErr := client.Do(deleteReq)
	if deleteErr != nil {
		t.Fatal(deleteErr)
	}
	assert.Equal(http.StatusNoContent, deleteRes.StatusCode)
	all, _ := cntrl.repo.FindAll()
	assert.Equal(2, len(all))
}