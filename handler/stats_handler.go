package handler

import (
	"context"
	"encoding/json"
	"net/http"
	redisclient "rate-limited-api/redis"
	"strconv"
	"strings"
	"time"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	keys, err := redisclient.Client.Keys(ctx, "stats:*").Result()
	if err != nil {
		http.Error(w, "failed to fetch stats", http.StatusInternalServerError)
		return
	}

	stats := make(map[string]map[string]int64)

	now := time.Now().Unix()
	windowStart := now - 60

	for _, key := range keys {

		userID := strings.TrimPrefix(key, "stats:")

		totalRequests, _ := redisclient.Client.Get(ctx, key).Int64()

		rateLimitKey := "rate_limit:" + userID

		lastMinuteRequests, _ := redisclient.Client.ZCount(
			ctx,
			rateLimitKey,
			strconv.FormatInt(windowStart, 10),
			strconv.FormatInt(now, 10),
		).Result()

		queuedRequests, _ := redisclient.Client.LLen(
			ctx,
			"request_queue",
		).Result()

		stats[userID] = map[string]int64{
			"total_requests":       totalRequests,
			"last_minute_requests": lastMinuteRequests,
			"queued_requests":      queuedRequests,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
