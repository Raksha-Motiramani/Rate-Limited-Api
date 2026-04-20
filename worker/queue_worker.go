package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"rate-limited-api/model"
	redisclient "rate-limited-api/redis"
	"rate-limited-api/service"
)

func StartWorker() {
	go func() {
		for {
			result, err := redisclient.Client.BRPop(
				context.Background(),
				0,
				"request_queue",
			).Result()

			if err != nil {
				continue
			}

			var req model.ApiRequest

			json.Unmarshal([]byte(result[1]), &req)

			allowed, _ := service.AllowRequest(req.UserId)

			if allowed {
				service.IncrementStats(req.UserId)
				queueCountKey := fmt.Sprintf("queue_count:%s", req.UserId)
				redisclient.Client.Decr(context.Background(), queueCountKey)
				log.Println("Processed queued request")
			} else {
				service.EnqueueRequest(req)
			}
		}
	}()
}
