package hivdomainstatus

import (
	"database/sql"
)

type Manager struct {
	repo DomainRepositoryInterface
}

func NewManager(repo DomainRepositoryInterface) (m *Manager) {
	m = new(Manager)
	m.repo = repo
	return
}

func (m *Manager) OnCheckDomainResult(r *DomainCheckResult) (err error) {
	domain, err := m.repo.FindByName(r.Domain)
	if err == sql.ErrNoRows {
		domain = new(Domain)
		domain.Name = r.Domain
	} else if err != nil {
		return
	}
	domain.Valid = r.Valid
	m.repo.Persist(domain)
	return
}