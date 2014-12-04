package hivdomainstatus

type Entity interface {
}

type Domain struct {
	Entity
	Id           int
	Name         string
}
