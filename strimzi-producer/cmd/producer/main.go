package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"strimzi-producer/internal/platform/kafka"
	log "strimzi-producer/internal/platform/logger"
	"strimzi-producer/internal/producer"
	"strimzi-producer/internal/rest/api"
	"strings"
	"syscall"
	"time"
)

func main() {
	app := &cli.App{
		Name:        "Strimzi Producer",
		Usage:       "Testing Strimzi Cluster",
		Description: "A simple producer to test our Strimzi cluster inside Kubernetes",
		Version:     "0.1.0",
		Authors: []*cli.Author{
			{Name: "Sina Nourian", Email: "sina.nourian@gmail.com"},
		},
		Flags:  flags,
		Action: start,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func start(c *cli.Context) error {
	// Init logger (Zap)
	err := log.InitLogger(log.LogProfile(c.String("log")))
	if err != nil {
		return err
	}

	// Init Kafka Producer
	servers := strings.Split(c.String("bootstrap_servers"), ",")
	writer := kafka.Producer(servers, c.String("topic"))
	defer writer.Close()
	client := producer.NewProducerClient(c.Duration("sleep_time_ms"), c.Uint("buffer_size"), writer)

	// Init Health Check api
	http.HandleFunc("/health", api.HealthCheck)
	srv := &http.Server{
		Addr:         c.String("serve"),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	// Listen and serve
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Fatalf("Error in listening and serving requests: %s", err.Error())
		}
	}()

	intChan := interruptChan() //interrupt channel for graceful shutdown

	// Start the producer
	go func() {
		client.Produce()
		intChan <- os.Interrupt //if the producer stops, exit the app
	}()

	gracefulShutdown(srv, intChan)
	return nil //return without error
}

func interruptChan() chan os.Signal {
	intChan := make(chan os.Signal)
	// SIGKILL (os.Kill) can't be caught, no need to add it
	signal.Notify(intChan, os.Interrupt, syscall.SIGTERM)
	return intChan
}

func gracefulShutdown(srv *http.Server, intChan chan os.Signal) {
	// Block until we receive our signal.
	<-intChan
	close(intChan)
	log.Logger.Infof("Shutting down server...")
	/*	the context is used to inform the server it has 5 seconds to finish
		the request it is currently handling */
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger.Warn("Server forced to shutdown: %s", err.Error())
	}
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "serve",
		Aliases:     []string{"s", "srv", "a", "addr"},
		Usage:       "bind address for http server",
		EnvVars:     []string{"SP_SERVE_ADDR"},
		Value:       ":8083",
		DefaultText: ":8083",
	},
	&cli.StringFlag{
		Name:        "bootstrap_servers",
		Aliases:     []string{"bs"},
		Usage:       "Kafka bootstrap servers",
		EnvVars:     []string{"SP_BOOTSTRAP_SERVERS"},
		Value:       "localhost:9092",
		DefaultText: "localhost:9092",
	},
	&cli.StringFlag{
		Name:        "topic",
		Aliases:     []string{"t"},
		Usage:       "Kafka topic to produce messages",
		EnvVars:     []string{"SP_KAFKA_TOPIC"},
		Value:       "my-topic",
		DefaultText: "my-topic",
	},
	&cli.DurationFlag{
		Name:        "sleep_time_ms",
		Aliases:     []string{"st"},
		Usage:       "How much to wait between producing messages in millisecond",
		EnvVars:     []string{"SP_SLEEP_TIME_MS"},
		Value:       1000,
		DefaultText: "1000",
	},
	&cli.UintFlag{
		Name:        "buffer_size",
		Aliases:     []string{"b"},
		Usage:       "How much to wait between producing messages in millisecond",
		EnvVars:     []string{"SP_BUFFER_SIZE"},
		Value:       5,
		DefaultText: "5",
	},
	&cli.StringFlag{
		Name:        "log",
		Aliases:     []string{"l"},
		Usage:       "logger preset (development, production)",
		EnvVars:     []string{"SP_LOGGER_PRESET"},
		Value:       "development",
		DefaultText: "development (showing debug logs)",
	},
}
