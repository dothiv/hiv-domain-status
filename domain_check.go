package hivdomainstatus

import "time"

type DomainCheck struct {
	EntityInterface
	Id           int64
	Domain         string
	URL         string
	StatusCode int
	ScriptPresent        bool
	IframeTarget         string
	IframeTargetOk        bool
	Valid        bool
	Created      *time.Time
}