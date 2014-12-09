package hivdomainstatus

import "fmt"

func transformCheckEntity(check *DomainCheck, route string) (m *DomainCheckModel) {
	m = new(DomainCheckModel)
	m.JsonLDContext = "http://jsonld.click4life.hiv/DomainCheck"
	m.JsonLDId = fmt.Sprintf(route, check.Id)
	m.Id = fmt.Sprintf("%d", check.Id)
	m.Domain = check.Domain
	m.DnsOK = check.DnsOK
	m.Addresses = check.Addresses
	m.URL = check.URL
	m.StatusCode = check.StatusCode
	m.ScriptPresent = check.ScriptPresent
	m.IframeTarget = check.IframeTarget
	m.IframeTargetOk = check.IframeTargetOk
	m.Valid = check.Valid
	m.Created = check.Created
	return
}

func transformEntity(e *Domain, route string) (m *DomainModel) {
	m = new(DomainModel)
	m.JsonLDContext = "http://jsonld.click4life.hiv/Domain"
	m.JsonLDId = fmt.Sprintf(route, e.Id)
	m.Id = fmt.Sprintf("%d", e.Id)
	m.Name = e.Name
	m.Created = e.Created
	return
}
