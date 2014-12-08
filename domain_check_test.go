package hivdomainstatus

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

// Test for the domain check equals

func TestThatDomainCheckEquals(t *testing.T) {
	assert := assert.New(t)

	c1 := new(DomainCheck)
	c2 := new(DomainCheck)

	assert.True(c1.Equals(c2))

	c1.Domain = "example.hiv"
	c1.URL = "http://example.hiv"
	c1.StatusCode = 200
	c1.ScriptPresent = true
	c1.IframeTarget = "http://example.com/"
	c1.IframeTargetOk = true
	c1.Valid = true
	c1.DnsOK = true
	c1.Addresses = []string{"1.2.3.4"}

	c2.Domain = "example.hiv"
	c2.URL = "http://example.hiv"
	c2.StatusCode = 200
	c2.ScriptPresent = true
	c2.IframeTarget = "http://example.com/"
	c2.IframeTargetOk = true
	c2.Valid = true
	c2.DnsOK = true
	c2.Addresses = []string{"1.2.3.4"}

	assert.True(c1.Equals(c2))
	c2.Domain = "example2.hiv"
	assert.False(c1.Equals(c2))
	c1.Domain = c2.Domain
	assert.True(c1.Equals(c2))

	c2.URL = "http://example2.hiv"
	assert.False(c1.Equals(c2))
	c1.URL = c2.URL
	assert.True(c1.Equals(c2))

	c2.StatusCode = 404
	assert.False(c1.Equals(c2))
	c1.StatusCode = c2.StatusCode
	assert.True(c1.Equals(c2))

	c2.ScriptPresent = false
	assert.False(c1.Equals(c2))
	c1.ScriptPresent = c2.ScriptPresent
	assert.True(c1.Equals(c2))

	c2.IframeTarget = "http://example2.com"
	assert.False(c1.Equals(c2))
	c1.IframeTarget = c2.IframeTarget
	assert.True(c1.Equals(c2))

	c2.IframeTargetOk = false
	assert.False(c1.Equals(c2))
	c1.IframeTargetOk = c2.IframeTargetOk
	assert.True(c1.Equals(c2))

	c2.Valid = false
	assert.False(c1.Equals(c2))
	c1.Valid = c2.Valid
	assert.True(c1.Equals(c2))

	c2.DnsOK = false
	assert.False(c1.Equals(c2))
	c1.DnsOK = c2.DnsOK
	assert.True(c1.Equals(c2))

	c2.Addresses = []string{"::1"}
	assert.False(c1.Equals(c2))
	c1.Addresses = c2.Addresses
	assert.True(c1.Equals(c2))
}
