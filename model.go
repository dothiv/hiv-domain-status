package hivdomainstatus

type DomainListModel struct {
	JsonLDTypedModel
	Items []*DomainModel `json:"items"`
	Total int            `json:"total"`
}

type DomainModel struct {
	JsonLDTypedModel
	Check *JsonLDTypedModel `json:"check"`
	Id    string
	Name  string
	Valid bool
}

type DomainCheckModel struct {
	JsonLDTypedModel
	Id             string
	Domain         string
	DnsOK          bool
	Addresses      []string
	URL            string
	StatusCode     int
	ScriptPresent  bool
	IframePresent  bool
	IframeTarget   string
	IframeTargetOk bool
	Valid          bool
}
