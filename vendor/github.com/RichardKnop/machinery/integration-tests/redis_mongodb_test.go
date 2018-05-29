package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/RichardKnop/machinery/v1/config"
)

func TestRedisMongodb(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	mongodbURL := os.Getenv("MONGODB_URL")
	if redisURL == "" || mongodbURL == "" {
		return
	}

	// Redis broker, MongoDB result backend
	server := testSetup(&config.Config{
		Broker:        fmt.Sprintf("redis://%v", redisURL),
		DefaultQueue:  "test_queue",
		ResultBackend: fmt.Sprintf("mongodb://%v", mongodbURL),
	})
	worker := server.NewWorker("test_worker", 0)
	go worker.Launch()
	testAll(server, t)
	worker.Quit()
}
