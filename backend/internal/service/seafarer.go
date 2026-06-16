package service

import (
	"errors"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
)

type SeafarerService struct {
	repo *repository.SeafarerRepo
	log  *logrus.Logger
}

func NewSeafarerService(repo *repository.SeafarerRepo, log *logrus.Logger) *SeafarerService {
	return &SeafarerService{repo: repo, log: log}
}

func (s *SeafarerService) List(page, pageSize int, keyword string) ([]model.Seafarer, int64, error) {
	s.log.Debugf("SeafarerService.List page=%d pageSize=%d keyword=%s", page, pageSize, keyword)
	return s.repo.List(page, pageSize, keyword)
}

func (s *SeafarerService) GetByID(id int64) (*model.Seafarer, error) {
	s.log.Debugf("SeafarerService.GetByID id=%d", id)
	return s.repo.GetByID(id)
}

func (s *SeafarerService) Create(sf *model.Seafarer) error {
	s.log.Debugf("SeafarerService.Create name=%s", sf.Name)
	if sf.Name == "" {
		return errors.New("船员姓名不能为空")
	}
	return s.repo.Create(sf)
}

func (s *SeafarerService) Update(sf *model.Seafarer) error {
	s.log.Debugf("SeafarerService.Update id=%d", sf.ID)
	if sf.ID <= 0 {
		return errors.New("无效的船员ID")
	}
	return s.repo.Update(sf)
}

func (s *SeafarerService) Delete(id int64) error {
	s.log.Debugf("SeafarerService.Delete id=%d", id)
	sf, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if sf.Status == model.SeafarerStatusOnboard {
		return errors.New("在船船员无法删除")
	}
	return s.repo.Delete(id)
}
