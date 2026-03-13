package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	Redis struct {
		Addr string `yaml:"addr"`
		DB   int    `yaml:"db"`
	} `yaml:"redis"`
	Channels struct {
		Target  string   `yaml:"target"`
		Sources []string `yaml:"sources"`
	} `yaml:"channels"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func run(ctx context.Context, cfg *Config) error {
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	log.Println("Connected to Redis at", cfg.Redis.Addr)

	pubsub := rdb.Subscribe(ctx, cfg.Channels.Sources...)
	defer pubsub.Close()

	log.Printf("Subscribed to %d source channel(s), forwarding to %q", len(cfg.Channels.Sources), cfg.Channels.Target)

	msgCh := pubsub.Channel()
	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				log.Println("Subscription channel closed")
				return nil
			}
			if err := rdb.Publish(ctx, cfg.Channels.Target, msg.Payload).Err(); err != nil {
				log.Printf("Failed to publish to %q: %v", cfg.Channels.Target, err)
			}
		case <-ctx.Done():
			log.Println("Shutting down")
			return nil
		}
	}
}

func main() {
	// Load .env for local development; ignore errors (env vars may already be set).
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config from %q: %v", configPath, err)
	}

	if cfg.Channels.Target == "" {
		log.Fatal("Config error: channels.target must not be empty")
	}
	if len(cfg.Channels.Sources) == 0 {
		log.Fatal("Config error: channels.sources must contain at least one channel")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg); err != nil {
		log.Fatalf("Service error: %v", err)
	}
}
