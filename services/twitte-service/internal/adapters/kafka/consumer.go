package kafka

import "github.com/yourusername/twitter-services-app/services/twitte-service/internal/core/services"

type Consumer struct {
	timelineService *services.TimelineService
}

func NewConsumer(ts *services.TimelineService) *Consumer {
	return &Consumer{
		timelineService: ts,
	}
}
