package main

import (
	"log"
	"net/http"

	"rate-limited-api/handler"
	"rate-limited-api/worker"
)

func main() {

	worker.StartWorker()

	http.HandleFunc("/request", handler.RequestHandler)

	http.HandleFunc("/stats", handler.StatsHandler)

	log.Println("Server started at :8080")

	http.ListenAndServe(":8080", nil)
}
