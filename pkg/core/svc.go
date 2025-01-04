package core

type ResoureRepo interface {
}

type Service struct {
	resource ResoureRepo
}

func New(resource ResoureRepo) *Service {
	return &Service{
		resource: resource,
	}
}
