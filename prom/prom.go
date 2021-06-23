package prom

import (
	"fmt"
	"time"

	"github.com/janse180/modem-scraper/scrape"
	"go.uber.org/zap"
)

func Publish(logger *zap.Logger, modemInformation scrape.ModemInformation) error {

	start := time.Now()

	logger.Debug("publishing prometheus metrics:",
		zap.String("op", "prometheus.Publish"),
	)

	modemInformation.ConnectionStatus.UpdateGauge()

	elapsed := time.Since(start)
	logger.Debug(fmt.Sprintf("finished exporting prometheus metrics, took %s", elapsed),
		zap.String("op", "prometheus.Publish"),
	)

	return nil
}
