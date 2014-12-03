package hivdomainstatus

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"testing"
	assert "github.com/stretchr/testify/assert"
)

// Test for the domain repository

func TestThatItPersists(t *testing.T) {
	assert := assert.New(t)

	c := NewDefaultConfig()
	configErr := gcfg.ReadFileInto(c, "config.ini")
	if configErr != nil {
		t.Fatal(configErr)
	}
	db, _ := sql.Open("postgres", c.DSN())
	
	// Persist
	domain := new(Domain)
	domain.Name = "example.hiv"
	repo := NewDomainRepository(db)
	_, persistErr := repo.Persist(domain)
	assert.Nil(persistErr)
	
	// Verify import
	domains, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(1, len(domains))

	assert.Equal("example.hiv", domains[0].Name)
}
