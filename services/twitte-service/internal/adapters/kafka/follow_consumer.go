package kafka

import "encoding/json"

type FollowEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}

func (c *Consumer) HandleFollowEvent(msg []byte) error {
	var event FollowEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	return c.timelineService.OnUserFollowed(
		event.FollowerID,
		event.FolloweeID,
	)
}
