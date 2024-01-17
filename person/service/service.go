package service

import (
	"go.uber.org/zap"
	"test-task/person/entity"
	"test-task/person/info"
	"test-task/person/repository"
)

type service struct {
	repo   repository.Repository
	logger *zap.SugaredLogger // zap is probably the most loved logging library for Go
}

type Service interface {
	AddPerson(req *entity.PersonReq) error
	GetPeople(nation string, page, pageSize int) ([]*entity.Person, error)
	DeletePerson(id int) error
	GetPersonById(id int) (entity.Person, error)
	UpdatePerson(id int, updated entity.Person) error
}

func NewService(repo repository.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) AddPerson(req *entity.PersonReq) error {
	personInfo, err := info.GetAllAdditionalInfo(req.Name)
	if err != nil {
		s.logger.Infof("Error in getting additional info: %s", err)
	} else {
		s.logger.Debugf("Additional info: %+v", personInfo)
	}

	person := entity.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         personInfo.Age,
		Gender:      personInfo.Gender,
		Nationality: personInfo.Nationality,
	}

	if err := s.repo.AddPerson(&person); err != nil {
		s.logger.Fatalf("Error adding person: %s", err)
	} else {
		s.logger.Debugf("Person added successfully: %+v", person)
	}

	return nil
}

func (s *service) GetPeople(nation string, page, pageSize int) ([]*entity.Person, error) {
	s.logger.Infof("Fetching people with filters: nation=%s, page=%d, pageSize=%d\n", nation, page, pageSize)

	people, err := s.repo.GetPeople(nation, page, pageSize)
	if err != nil {
		s.logger.Infof("Error fetching people: %v\n", err)
		return nil, err
	}

	s.logger.Debugf("Fetched people: %+v\n", people)
	return people, nil
}

func (s *service) DeletePerson(id int) error {
	s.logger.Infof("Deleting person...")
	if err := s.repo.DeleteById(id); err != nil {
		s.logger.Errorf("can not delete the person: %s", err)
		return err
	}
	s.logger.Infof("Delete success")
	return nil
}

func (s *service) GetPersonById(id int) (entity.Person, error) {
	s.logger.Infof("Getting person by id...")
	person, err := s.repo.GetPersonById(id)
	if err != nil {
		s.logger.Errorf("can not getting person by id or this person by this id does not exists: %s", err)
		return entity.Person{}, err
	}

	return person, nil
}

func (s *service) UpdatePerson(id int, updated entity.Person) error {
	s.logger.Infof("Updating person...")
	if err := s.repo.UpdatePerson(id, updated); err != nil {
		s.logger.Errorf("error in updating person: %s", err)
		return err
	}

	return nil
}
