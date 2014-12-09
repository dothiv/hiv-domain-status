package hivdomainstatus

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupDomainCheckTest(t *testing.T) (cntrl *DomainCheckController) {
	assert := assert.New(t)
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())

	cntrl = new(DomainCheckController)
	domainRepo := NewDomainRepository(db)
	cntrl.domainCheckRepo = NewDomainCheckRepository(db)
	db.Exec("TRUNCATE domain RESTART IDENTITY")
	db.Exec("TRUNCATE domain_check RESTART IDENTITY")

	data := []string{"example.hiv", "acme.hiv"}
	for _, name := range data {
		d := new(Domain)
		d.Name = name
		assert.Nil(domainRepo.Persist(d))
		check := new(DomainCheck)
		check.Domain = name
		check.DnsOK = true
		check.Addresses = []string{"127.0.0.1", "::1"}
		check.URL = "http://example.hiv"
		check.StatusCode = 200
		check.ScriptPresent = true
		check.IframeTarget = "http://example.com/"
		check.IframeTargetOk = true
		check.Valid = true
		assert.Nil(cntrl.domainCheckRepo.Persist(check))
	}

	return
}

func TestThatItListsDomainChecks(t *testing.T) {
	assert := assert.New(t)

	cntrl := SetupDomainCheckTest(t)

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

	var l DomainCheckListModel
	unmarshalErr := json.Unmarshal(b, &l)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	assert.Equal(2, l.Total)

	assert.Equal("example.hiv", l.Items[0].Domain)

	assert.Equal("acme.hiv", l.Items[1].Domain)

	assert.Equal(`<`+ts.URL+`/check?offsetKey=2>; rel="next"`, res.Header.Get("Link"))

	assert.Equal("example.hiv", l.Items[0].Domain)
	assert.Equal(fmt.Sprintf("%s/check/1", ts.URL), l.Items[0].JsonLDId)
	assert.Equal("http://jsonld.click4life.hiv/DomainCheck", l.Items[0].JsonLDContext)
	assert.True(l.Items[0].DnsOK)
	assert.Equal("127.0.0.1", l.Items[0].Addresses[0])
	assert.Equal("::1", l.Items[0].Addresses[1])
	assert.Equal("http://example.hiv", l.Items[0].URL)
	assert.Equal(200, l.Items[0].StatusCode)
	assert.True(l.Items[0].ScriptPresent)
	assert.Equal("http://example.com/", l.Items[0].IframeTarget)
	assert.True(l.Items[0].IframeTargetOk)
	assert.True(l.Items[0].Valid)
}
