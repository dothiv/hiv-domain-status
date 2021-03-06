package hivdomainstatus

import (
	"reflect"
	"time"
)

type EntityInterface interface {
}

type Domain struct {
	EntityInterface
	Id      int64
	Name    string
	Valid   bool
	Created *time.Time
}

type DomainCheck struct {
	EntityInterface
	Id             int64
	Domain         string
	DnsOK          bool
	AddressesJson  []byte
	Addresses      []string
	URL            string
	StatusCode     int
	ScriptPresent  bool
	IframePresent  bool
	IframeTarget   string
	IframeTargetOk bool
	Valid          bool
	Created        *time.Time
}

func (self *DomainCheck) Equals(other *DomainCheck) bool {
	if self.Domain != other.Domain {
		return false
	}
	if self.DnsOK != other.DnsOK {
		return false
	}
	if !reflect.DeepEqual(self.Addresses, other.Addresses) {
		return false
	}
	if self.URL != other.URL {
		return false
	}
	if self.StatusCode != other.StatusCode {
		return false
	}
	if self.ScriptPresent != other.ScriptPresent {
		return false
	}
	if self.IframePresent != other.IframePresent {
		return false
	}
	if self.IframeTarget != other.IframeTarget {
		return false
	}
	if self.IframeTargetOk != other.IframeTargetOk {
		return false
	}
	if self.Valid != other.Valid {
		return false
	}
	return true
}
