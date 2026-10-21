package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/HappyFatPug/go-course-metrics/internal/logger"
)

var gaugeBuf = make(map[string]float64)
var counterBuf = make(map[string]int64)

func UpdateMetric(mType, mName, mValue string) error {
	var err error

	switch mType {
	case "gauge":
		err = updateGauge(mName, mValue)
	case "counter":
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
	return nil
}

func updateCounter(name, strValue string) error {
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return err
	}
	counterBuf[name] += value
	logger.Logger.Debug(fmt.Sprintf("Counter value %s updated, new value: %d", name, counterBuf[name]))
	return nil
}
