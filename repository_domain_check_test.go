package hivdomainstatus

import (
	"database/sql"
	"testing"

	"code.google.com/p/gcfg"
	assert "github.com/stretchr/testify/assert"
)

// Test for the domain check result repository

func TestThatItPersistsDomainCheck(t *testing.T) {
	assert := assert.New(t)

	c := NewDefaultConfig()
	configErr := gcfg.ReadFileInto(c, "config.ini")
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	db.Exec("TRUNCATE domain_check RESTART IDENTITY")

	// Persist
	result := new(DomainCheck)
	result.DnsOK = true
	result.Addresses = []string{"127.0.0.1", "::1"}
	result.Domain = "example.hiv"
	result.URL = "http://example.hiv"
	result.StatusCode = 200
	result.ScriptPresent = true
	result.IframeTarget = "http://example.com/"
	result.IframeTargetOk = true
	result.Valid = true
	repo := NewDomainCheckRepository(db)
	persistErr := repo.Persist(result)
	assert.Nil(persistErr)

	// Verify
	results, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(results))

	r := results[0]

	assert.Equal(1, r.Id)
	assert.Equal("example.hiv", r.Domain)
	assert.True(r.DnsOK)
	assert.Equal("127.0.0.1", r.Addresses[0])
	assert.Equal("::1", r.Addresses[1])
	assert.Equal("http://example.hiv", r.URL)
	assert.Equal(200, r.StatusCode)
	assert.True(r.ScriptPresent)
	assert.Equal("http://example.com/", r.IframeTarget)
	assert.True(r.IframeTargetOk)
	assert.True(r.Valid)

	// Verify By Domain
	resultsByName, findNameErr := repo.FindByDomain("example.hiv")
	assert.Nil(findNameErr)
	assert.Equal(1, len(resultsByName))

	r2 := results[0]

	assert.Equal(1, r2.Id)
	assert.Equal("example.hiv", r2.Domain)

	// Verify Latest By Domain
	r3, findLatestByNameErr := repo.FindLatestByDomain("example.hiv")
	assert.Nil(findLatestByNameErr)
	assert.Equal(1, r3.Id)
	assert.Equal("example.hiv", r3.Domain)
}
