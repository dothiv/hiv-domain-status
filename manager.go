package hivdomainstatus

type Manager struct {
	repo DomainRepositoryInterface
}

func NewManager(repo DomainRepositoryInterface) (m *Manager) {
	m = new(Manager)
	m.repo = repo
	return
}

func (m *Manager) OnCheckDomainResult(r *DomainCheckResult) {
	domain, err := m.repo.FindByName(r.Domain)
	if err == nil {
		domain = new(Domain)
		domain.Name = r.Domain
	}
	domain.Valid = r.Valid
	m.repo.Persist(domain)
}