package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"bytes"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/config"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(cfg *config.Config) *Client {

	return &Client{
		baseURL: cfg.BaseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) LoadUsers() ([]model.User, error) {

	url := fmt.Sprintf("%s/users", c.baseURL)

	resp, err := c.http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result model.UserResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) post(path string, body any) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", c.baseURL, path)

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.http.Do(req)
}

func (c *Client) Login() error {

	resp, err := c.post("/login/random", struct{}{})

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Inquiry(account string) error {

	url := fmt.Sprintf(
		"%s/accounts/%s",
		c.baseURL,
		account,
	)

	resp, err := c.http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inquiry failed")
	}

	return nil
}

func (c *Client) Transfer(from, to string, amount float64) error {

	req := model.TransferRequest{
		FromAccount: from,
		ToAccount:   to,
		Amount:      amount,
	}

	resp, err := c.post("/transfer", req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("transfer failed")
	}

	return nil
}

func (c *Client) Payment(account string, merchant string, amount float64) error {

	req := model.PaymentRequest{
		AccountNumber: account,
		Merchant:      merchant,
		Amount:        amount,
	}

	resp, err := c.post("/payment", req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment failed")
	}

	return nil
}

func (c *Client) History() error {

	url := fmt.Sprintf("%s/transactions", c.baseURL)

	resp, err := c.http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("history failed")
	}

	return nil
}


