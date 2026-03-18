package _interface

type (
	DomainRepo interface {
		GetDomainInstructions(name string) (string, error)
	}

	DomainUsecase interface {
		GetInstructions(n string) (string, error)
	}
)
