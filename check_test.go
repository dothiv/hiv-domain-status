package hivdomainstatus

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupCheckTest(t *testing.T) (c *Config) {
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	db.Exec("TRUNCATE domain RESTART IDENTITY")
	repo := NewDomainRepository(db)

	data := []string{"click4life.hiv"}
	for _, name := range data {
		d := new(Domain)
		d.Name = name
		repo.Persist(d)
	}
	return
}

func TestThatItChecksDomain(t *testing.T) {
	assert := assert.New(t)
	SetupCheckTest(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<script src="` + CLICKCOUNTER_SCRIPT + `">`))
		w.Write([]byte(`<iframe id="clickcounter-target-iframe" src="//` + r.Host + `">`))
	}))
	defer ts.Close()

	lookupHost = func(domain string) (addresses []string, err error) {
		return []string{"1.2.3.4"}, nil
	}

	n := 0

	isValidDomain := func(domain string) bool {
		n += 1
		if n == 1 {
			return true
		}
		if n == 2 {
			return true
		}
		if n == 3 {
			return true
		}
		if n == 5 {
			return true
		}
		return false
	}

	testUrl, _ := url.Parse(ts.URL)
	testChecker := NewDomainCheckResult(testUrl.Host, isValidDomain)
	testChecker.URL = testUrl

	testChecker.SaveBody = false
	err := testChecker.Check()
	assert.Nil(err)
	assert.Equal(http.StatusOK, testChecker.StatusCode)
	assert.True(testChecker.ScriptPresent)
	assert.True(testChecker.IframePresent)
	assert.Equal(testChecker.IframeTarget, ts.URL)
	assert.True(testChecker.IframeTargetOk)
	assert.True(testChecker.DnsOk)
	assert.Equal("1.2.3.4", testChecker.Addresses[0])
	assert.True(testChecker.Valid)
}

func TestThatItDetectsHivDomain(t *testing.T) {
	assert := assert.New(t)
	assert.True(isHivDomain("hanseventures.hiv"))
}
