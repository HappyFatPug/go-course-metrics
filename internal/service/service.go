package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
	models "github.com/HappyFatPug/go-course-metrics/internal/model"
)

var gaugeBuf = make(map[string]float64)
var counterBuf = make(map[string]int64)

func UpdateMetric(mType, mName, mValue string) error {
	var err error

	switch mType {
	case models.Gauge:
		err = updateGauge(mName, mValue)
	case models.Counter:
		err = updateCounter(mName, mValue)
	default:
		err = errors.New("invalid metric type")
	}

	if err != nil {
		return err
	}

	return nil
}

func updateGauge(name string, strValue string) error {
	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return err
	}
	gaugeBuf[name] = value
	logger.Logger.Debug(fmt.Sprintf("Gauge value %s updated, new value: %.2f", name, gaugeBuf[name]))
	logger.Logger.Debug(fmt.Sprintf("Received gauge value: %.2f", value))
	return nil
}

func updateCounter(name, strValue string) error {
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return err
	}
	counterBuf[name] += value
	logger.Logger.Debug(fmt.Sprintf("Counter value %s updated, new value: %d", name, counterBuf[name]))
	logger.Logger.Debug(fmt.Sprintf("Received counter value: %d", value))
	return nil
}
