package usecase

import (
	_interface "ai-agent/interface"
)

type (
	DomainUsecase struct {
		cr _interface.DomainRepo
	}
)

func NewDomainUsecase(dr _interface.DomainRepo) *DomainUsecase {
	return &DomainUsecase{dr}
}

func (b *DomainUsecase) GetInstructions(name string) (string, error) {
	return b.cr.GetDomainInstructions(name)
}
