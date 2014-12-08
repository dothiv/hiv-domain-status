package hivdomainstatus

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"testing"
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
	result.Domain = "example.hiv"
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
}
