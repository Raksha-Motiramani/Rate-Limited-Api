package service

import (
	"context"
	"encoding/json"
	"fmt"

	"rate-limited-api/model"
	redisclient "rate-limited-api/redis"
)

const queueKey = "request_queue"

func EnqueueRequest(req model.ApiRequest) {

	ctx := context.Background()

	data, _ := json.Marshal(req)

	redisclient.Client.LPush(
		ctx,
		queueKey,
		data,
	)

	queueCountKey := fmt.Sprintf("queue_count:%s", req.UserId)

	redisclient.Client.Incr(ctx, queueCountKey)
}
