package handler

import (
	"encoding/json"
	"net/http"

	"rate-limited-api/model"
	"rate-limited-api/service"
)

func RequestHandler(w http.ResponseWriter, r *http.Request) {

	var req model.ApiRequest

	json.NewDecoder(r.Body).Decode(&req)

	allowed, _ := service.AllowRequest(req.UserId)

	if !allowed {
		service.EnqueueRequest(req)
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	service.IncrementStats(req.UserId)

	w.WriteHeader(http.StatusOK)
}
