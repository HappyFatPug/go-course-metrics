package agent

import (
	"bytes"
	"compress/gzip"
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
	c       *http.Client
	scheme  string
	host    string
	port    uint
	logger  *slog.Logger
	useGzip bool
}

func NewClient(config *config.AgentClientConfig, logger *slog.Logger) *Client {
	return &Client{
		scheme: config.Scheme,
		host:   config.NetAddr.Host,
		port:   uint(config.NetAddr.Port),
		c: &http.Client{
			Timeout: 1 * time.Second,
		},
		logger:  logger,
		useGzip: config.UseGzip,
	}
}

func (client *Client) postMetrics(m *models.Metrics) (*models.Metrics, bool) {
	endpoint := fmt.Sprintf("%s://%s:%d/update/", client.scheme, client.host, client.port)

	body, err := json.Marshal(m)
	if err != nil {
		client.logger.Error("JSON encode error", slog.String("msg", err.Error()))
		return nil, false
	}

	var req *http.Request

	if client.useGzip {
		var gBody bytes.Buffer
		gw := gzip.NewWriter(&gBody)

		_, err = gw.Write(body)
		if err != nil {
			client.logger.Error("gzip encode error", slog.String("msg", err.Error()))
			return nil, false
		}
		if err = gw.Close(); err != nil {
			client.logger.Error("gzip encode error", slog.String("msg", err.Error()))
			return nil, false
		}

		req, err = http.NewRequest(http.MethodPost, endpoint, &gBody)
		if err != nil {
			client.logger.Error("build http request error", slog.String("msg", err.Error()))
			return nil, false
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")

	} else {
		req, err = http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			client.logger.Error("build http request error", slog.String("msg", err.Error()))
			return nil, false
		}
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.c.Do(req)
	if err != nil {
		client.logger.Error("http error", slog.String("msg", err.Error()))
		return nil, false
	}
	if resp.StatusCode != http.StatusOK {
		client.logger.Warn("Failed posting metrics", slog.String("URL", endpoint), slog.Int("code", resp.StatusCode))
		return nil, false
	}
	client.logger.Debug("Metrics sent", slog.String("ID", m.ID), slog.Int("code", resp.StatusCode))

	defer resp.Body.Close()
	var respBody []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			client.logger.Error("gzip encode error", slog.String("msg", err.Error()))
			return nil, false
		}
		respBody, err = io.ReadAll(gr)
		if err != nil {
			client.logger.Error("gzip encode error", slog.String("msg", err.Error()))
			return nil, false
		}
	} else {
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			client.logger.Error("failed response read", slog.String("msg", err.Error()))
			return nil, false
		}
	}

	var respM models.Metrics
	err = json.Unmarshal(respBody, &respM)
	if err != nil {
		client.logger.Error("JSON decode error", slog.String("msg", err.Error()))
		return nil, false
	}
	return &respM, true
}
