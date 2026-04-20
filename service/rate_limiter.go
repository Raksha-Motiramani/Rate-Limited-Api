package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	redisclient "rate-limited-api/redis"

	"github.com/redis/go-redis/v9"
)

const LIMIT = 5
const WINDOW = 60

func AllowRequest(userID string) (bool, int64) {

	ctx := context.Background()

	key := fmt.Sprintf("rate_limit:%s", userID)

	now := time.Now().Unix()

	client := redisclient.Client

	windowStart := now - WINDOW

	client.ZRemRangeByScore(
		ctx,
		key,
		"-inf",
		strconv.FormatInt(windowStart, 10),
	)

	count, _ := client.ZCard(ctx, key).Result()

	if count >= LIMIT {

		oldest, _ := client.ZRangeWithScores(ctx, key, 0, 0).Result()

		oldestTimestamp := int64(oldest[0].Score)

		retryAfter := WINDOW - (now - oldestTimestamp)

		if retryAfter < 1 {
			retryAfter = 1
		}

		return false, retryAfter
	}

	client.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d-%d", now, time.Now().UnixNano()),
	})

	client.Expire(ctx, key, time.Minute)

	return true, 0
}
