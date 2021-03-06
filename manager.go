package hivdomainstatus

import (
	"database/sql"
	"log"
)

type Manager struct {
	domainRepo      DomainRepositoryInterface
	domainCheckRepo DomainCheckRepositoryInterface
}

func NewManager(domainRepo DomainRepositoryInterface, domainCheckRepo DomainCheckRepositoryInterface) (m *Manager) {
	m = new(Manager)
	m.domainRepo = domainRepo
	m.domainCheckRepo = domainCheckRepo
	return
}

func (m *Manager) OnCheckDomainResult(r *DomainCheckResult) (err error) {
	domain, err := m.domainRepo.FindByName(r.Domain)
	if err == sql.ErrNoRows {
		domain = new(Domain)
		domain.Name = r.Domain
		err = nil
	} else if err != nil {
		return
	}
	domain.Valid = r.Valid
	m.domainRepo.Persist(domain)

	result := new(DomainCheck)
	result.Domain = r.Domain
	result.DnsOK = r.DnsOk
	result.Addresses = r.Addresses
	result.URL = r.URL.String()
	result.StatusCode = r.StatusCode
	result.ScriptPresent = r.ScriptPresent
	result.IframePresent = r.IframePresent
	result.IframeTarget = r.IframeTarget
	result.IframeTargetOk = r.IframeTargetOk
	result.Valid = r.Valid
	lastResult, resultErr := m.domainCheckRepo.FindLatestByDomain(domain.Name)
	if resultErr == sql.ErrNoRows {
		m.domainCheckRepo.Persist(result)
		err = nil
	} else if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if !lastResult.Equals(result) {
		err = m.domainCheckRepo.Persist(result)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	}
	return
}
