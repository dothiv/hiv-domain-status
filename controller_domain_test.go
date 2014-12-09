package hivdomainstatus

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupDomainTest(t *testing.T) (cntrl *DomainController) {
	assert := assert.New(t)
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())

	cntrl = new(DomainController)
	cntrl.domainRepo = NewDomainRepository(db)
	cntrl.domainCheckRepo = NewDomainCheckRepository(db)
	db.Exec("TRUNCATE domain RESTART IDENTITY")
	db.Exec("TRUNCATE domain_check RESTART IDENTITY")

	data := []string{"example.hiv", "acme.hiv"}
	for _, name := range data {
		d := new(Domain)
		d.Name = name
		cntrl.domainRepo.Persist(d)
	}

	check := new(DomainCheck)
	check.Domain = "example.hiv"
	check.DnsOK = true
	check.Addresses = []string{"127.0.0.1", "::1"}
	check.URL = "http://example.hiv"
	check.StatusCode = 200
	check.ScriptPresent = true
	check.IframeTarget = "http://example.com/"
	check.IframeTargetOk = true
	check.Valid = true
	assert.Nil(cntrl.domainCheckRepo.Persist(check))

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

	var l DomainListModel
	unmarshalErr := json.Unmarshal(b, &l)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	assert.Equal(2, l.Total)

	assert.Equal("example.hiv", l.Items[0].Name)
	assert.Equal("acme.hiv", l.Items[1].Name)

	assert.Equal(`<`+ts.URL+`/domain?offsetKey=2>; rel="next"`, res.Header.Get("Link"))

	// Contains latest domain check
	assert.Equal("example.hiv", l.Items[0].Check.Domain)
	assert.Equal(fmt.Sprintf("%s/check/1", ts.URL), l.Items[0].Check.JsonLDId)
	assert.Equal("http://jsonld.click4life.hiv/DomainCheck", l.Items[0].Check.JsonLDContext)
	assert.True(l.Items[0].Check.DnsOK)
	assert.Equal("127.0.0.1", l.Items[0].Check.Addresses[0])
	assert.Equal("::1", l.Items[0].Check.Addresses[1])
	assert.Equal("http://example.hiv", l.Items[0].Check.URL)
	assert.Equal(200, l.Items[0].Check.StatusCode)
	assert.True(l.Items[0].Check.ScriptPresent)
	assert.Equal("http://example.com/", l.Items[0].Check.IframeTarget)
	assert.True(l.Items[0].Check.IframeTargetOk)
	assert.True(l.Items[0].Check.Valid)

	assert.Nil(l.Items[1].Check)
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

	assert.Equal(`<`+ts.URL+`/domain?offsetKey=2>; rel="next"`, res.Header.Get("Link"))
}

func TestThatItAddsNewDomain(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainTest(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntrl.ListingHandler(w, r, nil)
	}))
	defer ts.Close()

	var data = []byte(`{"name":"test.hiv"}`)
	res, err := http.Post(ts.URL+"/domain", "application/json", bytes.NewBuffer(data))
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
	d.Name = "test.hiv"
	cntrl.domainRepo.Persist(d)

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
	assert.Equal("test.hiv", m.Name)
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
	d.Name = "test.hiv"
	cntrl.domainRepo.Persist(d)

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
	all, _ := cntrl.domainRepo.FindAll()
	assert.Equal(2, len(all))
}
