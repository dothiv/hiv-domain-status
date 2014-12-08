package hivdomainstatus

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupManagerTest(t *testing.T) (repo *DomainRepository) {
	c, configErr := NewConfig()
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	db.Exec("TRUNCATE domain RESTART IDENTITY")
	repo = NewDomainRepository(db)
	return
}

func TestThatItStoresResultForNewDomain(t *testing.T) {
	assert := assert.New(t)
	repo := SetupManagerTest(t)

	// New domain
	r := new(DomainCheckResult)
	r.Domain = "example.hiv"
	m := NewManager(repo)
	m.OnCheckDomainResult(r)

	// Verify domain
	domains, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(domains))

	assert.Equal(1, domains[0].Id)
	assert.Equal("example.hiv", domains[0].Name)
	assert.False(domains[0].Valid)
}

func TestThatItStoresResultForExistingDomain(t *testing.T) {
	assert := assert.New(t)
	repo := SetupManagerTest(t)

	d := new(Domain)
	d.Name = "example.hiv"
	d.Valid = false
	repo.Persist(d)

	// Existing domain
	r := new(DomainCheckResult)
	r.Domain = "example.hiv"
	r.Valid = true
	m := NewManager(repo)
	m.OnCheckDomainResult(r)

	// Verify domain
	domains, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(domains))

	assert.Equal(1, domains[0].Id)
	assert.Equal("example.hiv", domains[0].Name)
	assert.True(domains[0].Valid)
}