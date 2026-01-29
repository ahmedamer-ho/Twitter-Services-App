package services
import(
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/domain"
) 

type TimelineService struct {
	redis  *redis.Client
	tweets domain.TweetRepository
}
func (s *TimelineService) OnUserFollowed(follower, followee string) error {
	ctx := context.Background()

	// follower → following
	s.redis.SAdd(ctx, "following:"+follower, followee)

	// followee → followers
	s.redis.SAdd(ctx, "followers:"+followee, follower)

	return nil
}

func (s *TimelineService) shouldFanOutOnWrite(userID string) bool {
	count, _ := s.redis.SCard(context.Background(), "followers:"+userID).Result()
	return count < 100_000
}
func (s *TimelineService) FanOutTweet(tweet domain.Tweet) error {
	ctx := context.Background()

	followers, _ := s.redis.SMembers(ctx, "followers:"+tweet.AuthorID).Result()

	for _, followerID := range followers {
		s.redis.ZAdd(ctx,
			"timeline:"+followerID,
			&redis.Z{
				Score:  float64(tweet.CreatedAt.Unix()),
				Member: tweet.ID,
			},
		)
	}

	return nil
}
func (s *TimelineService) GetTimeline(userID string) ([]domain.Tweet, error) {
	ctx := context.Background()

	tweetIDs, _ := s.redis.ZRevRange(ctx, "timeline:"+userID, 0, 50).Result()

	if len(tweetIDs) == 0 {
		return s.buildTimelineOnRead(userID)
	}

	return s.tweets.GetByIDs(ctx,tweetIDs)
}

func (s *TimelineService) buildTimelineOnRead(userID string) ([]domain.Tweet, error) {
	following, _ := s.redis.SMembers(context.Background(), "following:"+userID).Result()

	return s.tweets.GetLatestByAuthors(context.Background(),following)
}

func (s *TimelineService) OnTweetCreated(tweet domain.Tweet) error {
	if s.shouldFanOutOnWrite(tweet.AuthorID) {
		return s.FanOutTweet(tweet)
	}

	// celebrity path
	return s.StoreForFanOutOnRead(tweet)
}

func (s *TimelineService) StoreForFanOutOnRead(tweet domain.Tweet) error {
	ctx := context.Background()

	// Store the tweet in the author's outbox for fan-out on read
	s.redis.ZAdd(ctx,
		"outbox:"+tweet.AuthorID,
		&redis.Z{
			Score:  float64(tweet.CreatedAt.Unix()),
			Member: tweet.ID,
		},
	)
	return nil
}