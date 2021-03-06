package hivdomainstatus

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

// Test for the domain repository

func TestThatItPersistsADomain(t *testing.T) {
	assert := assert.New(t)

	c := NewDefaultConfig()
	configErr := gcfg.ReadFileInto(c, "config.ini")
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	db.Exec("TRUNCATE domain RESTART IDENTITY")
	
	// Persist
	domain := new(Domain)
	domain.Name = "example.hiv"
	repo := NewDomainRepository(db)
	persistErr := repo.Persist(domain)
	assert.Nil(persistErr)
	
	// Verify
	domains, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(domains))

	d := domains[0]

	assert.Equal(1, d.Id)
	assert.Equal("example.hiv", d.Name)
}
