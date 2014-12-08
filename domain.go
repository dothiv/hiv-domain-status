package hivdomainstatus

import "time"

type Domain struct {
	EntityInterface
	Id           int64
	Name         string
	Valid        bool
	Created      *time.Time
}
