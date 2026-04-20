package service

import (
	"context"
	"fmt"

	redisclient "rate-limited-api/redis"
)

type UserStats struct {
	TotalRequests      int64 `json:"total_requests"`
	LastMinuteRequests int64 `json:"last_minute_requests"`
	QueuedRequests     int64 `json:"queued_requests"`
}

func GetUserStats(userID string) UserStats {

	ctx := context.Background()

	totalKey := fmt.Sprintf("stats:%s", userID)

	rateLimitKey := fmt.Sprintf("rate_limit:%s", userID)

	totalRequests, _ := redisclient.Client.Get(ctx, totalKey).Int64()

	lastMinuteRequests, _ := redisclient.Client.ZCard(ctx, rateLimitKey).Result()

	queuedRequests, _ := redisclient.Client.LLen(
		ctx,
		"request_queue",
	).Result()

	return UserStats{
		TotalRequests:      totalRequests,
		LastMinuteRequests: lastMinuteRequests,
		QueuedRequests:     queuedRequests,
	}
}

func IncrementStats(userID string) {

	key := fmt.Sprintf("stats:%s", userID)

	redisclient.Client.Incr(
		context.Background(),
		key,
	)
}
