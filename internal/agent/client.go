package agent

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

type Client struct {
	c      *http.Client
	scheme string
	host   string
	port   uint
}

func NewClient(scheme, host string, port uint) *Client {
	return &Client{
		scheme: scheme,
		host:   host,
		port:   port,
		c: &http.Client{
			Timeout: 1 * time.Second,
		},
	}
}

func (client *Client) postMetrics(m *models.Metrics) error {
	var urlPath, endpoint string

	switch m.MType {
	case models.Counter:
		urlPath = fmt.Sprintf("/update/counter/%s/%d", m.ID, *m.Delta)
	case models.Gauge:
		urlPath = fmt.Sprintf("/update/gauge/%s/%s", m.ID, m.ValueStr())
	default:
		return errors.New("wrong type")
	}

	endpoint = fmt.Sprintf("%s://%s:%d%s", client.scheme, client.host, client.port, urlPath)

	logger.Logger.Debug(fmt.Sprintf("Calling POST: %s", endpoint))
	response, err := client.c.Post(endpoint, "text/plain", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	io.Copy(io.Discard, response.Body)
	logger.Logger.Debug(fmt.Sprintf("Metrics %s sent, code received: %d", m.ID, response.StatusCode))

	return nil
}
