package hivdomainstatus

type DomainListModel struct {
	JsonLDTypedModel
	Items []*DomainModel `json:"items"`
	Total int                  `json:"total"`
}

type DomainModel struct {
	JsonLDTypedModel
	Id string
	Name         string
}