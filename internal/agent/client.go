package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/config"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type Client struct {
	c      *http.Client
	scheme string
	host   string
	port   uint
	logger *slog.Logger
}

func NewClient(config *config.AgentClientConfig, logger *slog.Logger) *Client {
	return &Client{
		scheme: config.Scheme,
		host:   config.NetAddr.Host,
		port:   uint(config.NetAddr.Port),
		c: &http.Client{
			Timeout: 1 * time.Second,
		},
		logger: logger,
	}
}

func (client *Client) postMetrics(m *models.Metrics) error {
	endpoint := fmt.Sprintf("%s://%s:%d/update/", client.scheme, client.host, client.port)

	body, err := json.Marshal(m)
	if err != nil {
		return err
	}

	response, err := client.c.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	defer response.Body.Close()
	io.Copy(io.Discard, response.Body)
	if response.StatusCode != http.StatusOK {
		client.logger.Warn("Failed posting metrics", slog.String("URL", endpoint), slog.Int("code", response.StatusCode))
	} else {
		client.logger.Debug("Metrics sent", slog.String("ID", m.ID), slog.Int("code", response.StatusCode))
	}

	return nil
}
