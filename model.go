package hivdomainstatus

import "time"

type DomainListModel struct {
	JsonLDTypedModel
	Items []*DomainModel `json:"items"`
	Total int            `json:"total"`
}

type DomainCheckListModel struct {
	JsonLDTypedModel
	Items []*DomainCheckModel `json:"items"`
	Total int                 `json:"total"`
}

type DomainCheckModel struct {
	JsonLDTypedModel
	Id             string     `json:"-"`
	Domain         string     `json:"domain"`
	DnsOK          bool       `json:"dnsOk"`
	Addresses      []string   `json:"addresses"`
	URL            string     `json:"url"`
	StatusCode     int        `json:"statusCode"`
	ScriptPresent  bool       `json:"scriptPresent"`
	IframePresent  bool       `json:"iframePresent"`
	IframeTarget   string     `json:"iframeTarget"`
	IframeTargetOk bool       `json:"iframeTargetOk"`
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
