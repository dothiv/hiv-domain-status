package hivdomainstatus

import "time"

type DomainListModel struct {
	JsonLDTypedModel
	Items []*DomainModel `json:"items"`
	Total int            `json:"total"`
}

type DomainCheckModel struct {
	JsonLDTypedModel
	Id             string     `json:"-"`
	Domain         string     `json:"domain"`
	DnsOK          bool       `json:"dns_ok"`
	Addresses      []string   `json:"addresses"`
	URL            string     `json:"url"`
	StatusCode     int        `json:"status_code"`
	ScriptPresent  bool       `json:"script_present"`
	IframePresent  bool       `json:"iframe_present"`
	IframeTarget   string     `json:"iframe_target"`
	IframeTargetOk bool       `json:"iframe_target_ok"`
	Valid          bool       `json:"valid"`
	Created        *time.Time `json:"created"`
}

type DomainModel struct {
	JsonLDTypedModel
	Id      string            `json:"-"`
	Name    string            `json:"name"`
	Valid   bool              `json:"valid"`
	Check   *DomainCheckModel `json:"check"`
	Created *time.Time        `json:"created"`
}
