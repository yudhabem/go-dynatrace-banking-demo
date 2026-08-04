package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const defaultAPIURL = "http://localhost:8080"

type config struct {
	apiURL         string
	workers        int
	interval       time.Duration
	bootstrapUsers int
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

type loginResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AccountNumber string `json:"accountNumber"`
	} `json:"data"`
	Message string `json:"message"`
}

func main() {
	cfg := loadConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	api := &client{
		baseURL: strings.TrimRight(cfg.apiURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	if err := api.bootstrap(ctx, cfg.bootstrapUsers); err != nil {
		log.Printf("bootstrap failed: %v; workers will keep retrying", err)
	}

	log.Printf("load generator started: workers=%d interval=%s api=%s", cfg.workers, cfg.interval, cfg.apiURL)

	var workers sync.WaitGroup
	for workerID := 1; workerID <= cfg.workers; workerID++ {
		workers.Add(1)
		go func(id int) {
			defer workers.Done()
			api.runWorker(ctx, id, cfg.interval)
		}(workerID)
	}

	<-ctx.Done()
	log.Println("load generator stopping")
	workers.Wait()
}

func loadConfig() config {
	return config{
		apiURL:         envString("API_BASE_URL", defaultAPIURL),
		workers:        envInt("WORKER_COUNT", 100),
		interval:       envDuration("REQUEST_INTERVAL", 200*time.Millisecond),
		bootstrapUsers: envInt("BOOTSTRAP_USERS", 10),
	}
}

func (c *client) bootstrap(ctx context.Context, users int) error {
	for i := 0; i < users; i++ {
		for {
			if err := c.request(ctx, http.MethodPost, "/users/random", nil, nil); err == nil {
				break
			} else if ctx.Err() != nil {
				return ctx.Err()
			} else {
				log.Printf("bootstrap user=%d unavailable: %v; retrying", i+1, err)
				if !wait(ctx, time.Second) {
					return ctx.Err()
				}
			}
		}
	}

	return nil
}

func (c *client) runWorker(ctx context.Context, workerID int, interval time.Duration) {
	for iteration := 1; ctx.Err() == nil; iteration++ {
		if err := c.runFlow(ctx); err != nil && iteration%20 == 1 {
			log.Printf("worker=%d flow failed: %v", workerID, err)
		}

		if !wait(ctx, interval) {
			return
		}
	}
}

func (c *client) runFlow(ctx context.Context) error {
	from, err := c.randomLogin(ctx)
	if err != nil {
		return fmt.Errorf("login source account: %w", err)
	}

	if err := c.request(ctx, http.MethodGet, "/accounts/"+from, nil, nil); err != nil {
		return fmt.Errorf("inquiry: %w", err)
	}

	to, err := c.randomLogin(ctx)
	if err != nil {
		return fmt.Errorf("login destination account: %w", err)
	}

	transfer := map[string]any{
		"fromAccount": from,
		"toAccount":   to,
		"amount":      rand.Float64()*900 + 100,
	}
	if err := c.request(ctx, http.MethodPost, "/transfer", transfer, nil); err != nil {
		return fmt.Errorf("transfer: %w", err)
	}

	payment := map[string]any{
		"accountNumber": from,
		"merchant":      randomMerchant(),
		"amount":        rand.Float64()*400 + 50,
	}
	if err := c.request(ctx, http.MethodPost, "/payment", payment, nil); err != nil {
		return fmt.Errorf("payment: %w", err)
	}

	if err := c.request(ctx, http.MethodGet, "/transactions", nil, nil); err != nil {
		return fmt.Errorf("history: %w", err)
	}

	return nil
}

func (c *client) randomLogin(ctx context.Context) (string, error) {
	var response loginResponse
	if err := c.request(ctx, http.MethodPost, "/login/random", nil, &response); err != nil {
		return "", err
	}
	if !response.Success || response.Data.AccountNumber == "" {
		return "", fmt.Errorf("unexpected login response: %s", response.Message)
	}

	return response.Data.AccountNumber, nil
}

func (c *client) request(ctx context.Context, method, path string, body any, output any) error {
	var payload *bytes.Reader
	if body == nil {
		payload = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected status %d", response.StatusCode)
	}

	if output != nil {
		return json.NewDecoder(response.Body).Decode(output)
	}

	return nil
}

func randomMerchant() string {
	merchants := []string{"PLN", "TELKOMSEL", "TOKOPEDIA", "SHOPEE", "GOJEK"}
	return merchants[rand.IntN(len(merchants))]
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(key))
	if err != nil || value < 0 {
		return fallback
	}
	return value
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
