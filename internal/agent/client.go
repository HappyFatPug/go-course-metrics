package agent

import (
	"errors"
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
	var urlPath, endpoint string

	switch m.MType {
	case models.Counter:
		urlPath = fmt.Sprintf("/update/counter/%s/%s", m.ID, m.ValueStr())
	case models.Gauge:
		urlPath = fmt.Sprintf("/update/gauge/%s/%s", m.ID, m.ValueStr())
	default:
		return errors.New("wrong type")
	}

	endpoint = fmt.Sprintf("%s://%s:%d%s", client.scheme, client.host, client.port, urlPath)

	client.logger.Debug("Calling POST", slog.String("URL", endpoint))
	response, err := client.c.Post(endpoint, "text/plain", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	io.Copy(io.Discard, response.Body)
	client.logger.Debug("Metrics sent", slog.String("ID", m.ID), slog.Int("code", response.StatusCode))

	return nil
}
