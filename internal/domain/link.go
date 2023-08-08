package domain

import (
	"github.com/google/uuid"
	kontrakto "github.com/taranovegor/com.kontrakto"
	"time"
)

type Link struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;<-:create"`
	Token     string    `gorm:"index;type:char(5)"`
	Location  string    `gorm:"type:text"`
	CreatedAt time.Time
}

type LinkRepository interface {
	Store(*Link) error
	GetByToken(string) (Link, error)
	Update(*Link) error
	DeleteByToken(string) error
	Paginate(kontrakto.Paginator) ([]Link, error)
}

func NewLink(token string, location string) *Link {
	return &Link{
		ID:        uuid.New(),
		Token:     token,
		Location:  location,
		CreatedAt: time.Now(),
	}
}
