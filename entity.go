package hivdomainstatus

type Entity interface {
}

type Domain struct {
	Entity
	Id           int64
	Name         string
}
