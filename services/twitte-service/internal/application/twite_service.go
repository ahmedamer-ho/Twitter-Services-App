package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Twitter-Services-App/twite-service/internal/domain"
)

type TwiteService struct {
	repo domain.TwiteRepository
}

func NewTwiteService(repo domain.TwiteRepository) *TwiteService {
	return &TwiteService{repo: repo}
}
func (s *TwiteService) CreateTwite(ctx context.Context, authorID, content string) error {
	t := domain.Twite{
			ID:        uuid.NewString(),
			AuthorID:  authorID,
			Content:   content,
			CreatedAt: time.Now().UTC(),
	}

	return s.repo.Insert(ctx, t)
}