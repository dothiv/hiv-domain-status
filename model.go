package hivdomainstatus

type Model interface {
}

type Domain struct {
	Model
	Id           int
	Name         string
}
