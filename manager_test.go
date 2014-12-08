package hivdomainstatus

import (
	"database/sql"
	"testing"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func SetupManagerTest(t *testing.T) (domainRepo *DomainRepository, domainCheckRepo *DomainCheckRepository) {
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	db.Exec("TRUNCATE domain RESTART IDENTITY")
	db.Exec("TRUNCATE domain_check RESTART IDENTITY")
	domainRepo = NewDomainRepository(db)
	domainCheckRepo = NewDomainCheckRepository(db)
	return
}

func TestThatItStoresResultForNewDomain(t *testing.T) {
	assert := assert.New(t)
	domainRepo, domainCheckRepo := SetupManagerTest(t)

	// New domain
	r := new(DomainCheckResult)
	r.Domain = "example.hiv"
	r.URL, _ = url.Parse("http://example.hiv")
	m := NewManager(domainRepo, domainCheckRepo)
	err := m.OnCheckDomainResult(r)
	assert.Nil(err)

	// Verify domain
	domains, findErr := domainRepo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(domains))

	d := domains[0]

	assert.Equal(1, d.Id)
	assert.Equal("example.hiv", d.Name)
	assert.False(d.Valid)

	// Verify domain check result
	results, resultsErr := domainCheckRepo.FindAll()
	assert.Nil(resultsErr)
	assert.Equal(1, len(results))

	res := results[0]
	assert.Equal(1, res.Id)
	assert.Equal("example.hiv", res.Domain)
	assert.Equal("http://example.hiv", res.URL)
	assert.False(res.Valid)
}

func TestThatItStoresResultForExistingDomain(t *testing.T) {
	assert := assert.New(t)
	domainRepo, domainCheckRepo := SetupManagerTest(t)

	d := new(Domain)
	d.Name = "example.hiv"
	d.Valid = false
	domainRepo.Persist(d)

	res := new(DomainCheck)
	res.Domain = "example.hiv"
	res.Valid = false
	domainCheckRepo.Persist(res)

	// Existing domain
	r := new(DomainCheckResult)
	r.Domain = "example.hiv"
	r.URL, _ = url.Parse("http://example.hiv")		
	r.Valid = true
	m := NewManager(domainRepo, domainCheckRepo)
	err := m.OnCheckDomainResult(r)
	assert.Nil(err)

	// Verify domain
	domains, findErr := domainRepo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(domains))

	d2 := domains[0]

	assert.Equal(1, d2.Id)
	assert.Equal("example.hiv", d2.Name)
	assert.True(d2.Valid)

	// Verify domain check result
	results, resultsErr := domainCheckRepo.FindAll()
	assert.Nil(resultsErr)
	assert.Equal(2, len(results))

	res1 := results[0]
	assert.Equal(1, res1.Id)
	assert.Equal("example.hiv", res1.Domain)
	assert.False(res1.Valid)

	res2 := results[1]
	assert.Equal(2, res2.Id)
	assert.Equal("example.hiv", res2.Domain)
	assert.True(res2.Valid)
}