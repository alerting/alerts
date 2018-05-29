package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/RichardKnop/machinery/v1/config"
)

func TestAmqpRedis(t *testing.T) {
	amqpURL := os.Getenv("AMQP_URL")
	redisURL := os.Getenv("REDIS_URL")
	if amqpURL == "" || redisURL == "" {
		return
	}

	// AMQP broker, Redis result backend
	server := testSetup(&config.Config{
		Broker:        amqpURL,
		DefaultQueue:  "test_queue",
		ResultBackend: fmt.Sprintf("redis://%v", redisURL),
		AMQP: &config.AMQPConfig{
			Exchange:      "test_exchange",
			ExchangeType:  "direct",
			BindingKey:    "test_task",
			PrefetchCount: 1,
		},
	})
	worker := server.NewWorker("test_worker", 0)
	go worker.Launch()
	testAll(server, t)
	worker.Quit()
}
