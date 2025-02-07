package scrape

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/influxdata/influxdb1-client" // this is important because of a bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// DownstreamBondedChannel holds all info from the
// "Downstream Bonded Channels" table on
// /cmconnectionstatus.html.
type DownstreamBondedChannel struct {
	ChannelID      int
	LockStatus     string
	Modulation     string
	FrequencyHz    int
	PowerdBmV      float64
	SNRdB          float64
	Corrected      int
	Uncorrectables int
}

var (
	DownstreamBondedChannelPowerGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "downstream_bonded_channel_powerdbmv",
		Help: "The downstream bonded channel power",
	}, []string{
		"ChannelID",
		"LockStatus",
		"Modulation",
		"FrequencyHz",
	})
	DownstreamBondedChannelSNRGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "downstream_bonded_channel_snrdb",
		Help: "The downstream bonded channel snr",
	}, []string{
		"ChannelID",
		"LockStatus",
		"Modulation",
		"FrequencyHz",
	})
	DownstreamBondedChannelCorrectedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "downstream_bonded_channel_error_corrected",
		Help: "The downstream bonded channel corrected errors",
	}, []string{
		"ChannelID",
		"LockStatus",
		"Modulation",
		"FrequencyHz",
	})
	DownstreamBondedChannelUncorrectedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "downstream_bonded_channel_error_uncorrected",
		Help: "The downstream bonded channel uncorrected errors",
	}, []string{
		"ChannelID",
		"LockStatus",
		"Modulation",
		"FrequencyHz",
	})
)

func (d DownstreamBondedChannel) UpdateGauge() error {

	DownstreamBondedChannelPowerGauge.WithLabelValues(
		strconv.Itoa(d.ChannelID),
		d.LockStatus,
		d.Modulation,
		strconv.Itoa(d.FrequencyHz)).Set(d.PowerdBmV)

	DownstreamBondedChannelSNRGauge.WithLabelValues(
		strconv.Itoa(d.ChannelID),
		d.LockStatus,
		d.Modulation,
		strconv.Itoa(d.FrequencyHz)).Set(d.SNRdB)

	DownstreamBondedChannelCorrectedGauge.WithLabelValues(
		strconv.Itoa(d.ChannelID),
		d.LockStatus,
		d.Modulation,
		strconv.Itoa(d.FrequencyHz)).Set(float64(d.Corrected))

	DownstreamBondedChannelUncorrectedGauge.WithLabelValues(
		strconv.Itoa(d.ChannelID),
		d.LockStatus,
		d.Modulation,
		strconv.Itoa(d.FrequencyHz)).Set(float64(d.Uncorrectables))

	return nil

}

// ToInfluxPoints converts DownstreamBondedChannel to "points"
func (d DownstreamBondedChannel) ToInfluxPoints() ([]*client.Point, error) {
	var points []*client.Point

	channelIDString := strconv.Itoa(d.ChannelID)
	tags := map[string]string{
		"channel_id": channelIDString,
	}
	fields := map[string]interface{}{
		"lock_status":    d.LockStatus,
		"modulation":     d.Modulation,
		"frequency_hz":   d.FrequencyHz,
		"power_dbmv":     d.PowerdBmV,
		"snr_db":         d.SNRdB,
		"corrected":      d.Corrected,
		"uncorrectables": d.Uncorrectables,
	}
	point, err := client.NewPoint("downstream_bonded_channel", tags, fields, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error generating points data for DownstreamBondedChannel: %s", err.Error())
	}

	points = append(points, point)

	return points, nil
}

const downstreamBondedChannelTableSelector = "#bg3 > div.container > div.content > center:nth-child(5) > table"

func scrapeDownstreamBondedChannels(doc *goquery.Document) []DownstreamBondedChannel {
	downstreamBondedChannelTable := doc.Find(downstreamBondedChannelTableSelector)
	downstreamBondedChannelTableTbody := downstreamBondedChannelTable.Children()
	downstreamBondedChannelTableTbodyRows := downstreamBondedChannelTableTbody.Children()

	downstreamBondedChannels := []DownstreamBondedChannel{}
	downstreamBondedChannelTableTbodyRows.Each(func(index int, row *goquery.Selection) {
		// Skip the "title" row as well as the "header" row.
		// These are both regular old <tr> rows on this page.
		if index > 1 {
			downstreamBondedChannels = append(downstreamBondedChannels, makeDownstreamBondedChannel(row))
		}
	})

	return downstreamBondedChannels
}

func makeDownstreamBondedChannel(selection *goquery.Selection) DownstreamBondedChannel {
	rowData := selection.Children()
	downstreamBondedChannel := DownstreamBondedChannel{
		ChannelID:      getIntRowData(rowData, 0),
		LockStatus:     rowData.Get(1).FirstChild.Data,
		Modulation:     rowData.Get(2).FirstChild.Data,
		FrequencyHz:    getIntRowData(rowData, 3),
		PowerdBmV:      getFloatRowData(rowData, 4),
		SNRdB:          getFloatRowData(rowData, 5),
		Corrected:      getIntRowData(rowData, 6),
		Uncorrectables: getIntRowData(rowData, 7),
	}

	return downstreamBondedChannel
}
