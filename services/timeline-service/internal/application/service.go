package application

import (
	"context"
	"encoding/json"

	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/domain"
	"github.com/yourusername/twitter-services-app/shared/nats"
)

type TimelineRepository interface {
	SaveTweet(ctx context.Context, userID string, tweet domain.TimelineTweet) error
	GetTimeline(ctx context.Context, userID string) (*domain.Timeline, error)
}

type TimelineService struct {
	repo TimelineRepository
}

func NewTimelineService(repo TimelineRepository) *TimelineService {
	return &TimelineService{repo: repo}
}

func (s *TimelineService) AddTweetToTimeline(ctx context.Context, tweet domain.TimelineTweet) error {
	// In a real app, we'd fetch followers of tweet.AuthorID and fan out.
	// For this demonstration, we'll just store it in a "global" timeline
	// for a user named "everyone" and the author themselves.

	err := s.repo.SaveTweet(ctx, tweet.AuthorID, tweet)
	if err != nil {
		return err
	}

	return s.repo.SaveTweet(ctx, "global", tweet)
}

func (s *TimelineService) GetTimeline(ctx context.Context, userID string) (*domain.Timeline, error) {
	return s.repo.GetTimeline(ctx, userID)
}

func (s *TimelineService) HandleTweetEvent(ctx context.Context, event nats.Event) error {
	var tweet domain.TimelineTweet

	bytes, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &tweet); err != nil {
		return err
	}

	return s.AddTweetToTimeline(ctx, tweet)
}
