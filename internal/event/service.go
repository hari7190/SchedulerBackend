package event

import "errors"

type Store interface {
	Create(e EvenCreationRequest) error
	GetById(id string) (Event, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

// func (s *Service) GetById(id string) (Event, error) {

// }

func (s *Service) CreateEvent(eReq EvenCreationRequest) (Event, error) {
	if eReq.Title == "" {
		return Event{}, errors.New("Title is required")
	}

	e := Event{
		Title: eReq.Title,
	}

	if err := s.store.Create(eReq); err != nil {
		return Event{}, err
	}

	return e, nil
}
