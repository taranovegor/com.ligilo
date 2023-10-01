package repository

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"gorm.io/gorm"
)

type linkRepository struct {
	domain.LinkRepository
	orm *gorm.DB
}

func NewLinkRepository(
	orm *gorm.DB,
) domain.LinkRepository {
	return &linkRepository{
		orm: orm,
	}
}

func NewRepository() domain.LinkRepository {
	return &linkRepository{}
}

func (repo linkRepository) Store(link *domain.Link) error {
	return repo.orm.Create(link).Error
}

func (repo linkRepository) GetByToken(token string) (domain.Link, error) {
	var link domain.Link
	if err := repo.orm.First(&link, "token = ?", token).Error; err != nil {
		return link, err
	}

	return link, nil
}

func (repo linkRepository) Update(link *domain.Link) error {
	return repo.orm.Updates(link).Error
}

func (repo linkRepository) DeleteByToken(token string) error {
	return repo.orm.Delete(&domain.Link{}, "token = ?", token).Error
}

func (repo linkRepository) Paginate(p kontrakto.Paginator) ([]domain.Link, error) {
	var links []domain.Link
	err := repo.orm.Order("created_at DESC").Limit(p.Limit).Offset(p.Offset).Find(&links).Error

	return links, err
}
