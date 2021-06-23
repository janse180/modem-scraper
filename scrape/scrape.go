package scrape

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/janse180/modem-scraper/config"
	"go.uber.org/zap"
)

// Scrape scrapes data from the modem.
func Scrape(logger *zap.Logger, conf config.Configuration) (*ModemInformation, error) {

	token, err := getToken(logger, conf)
	if err != nil {
		return nil, err
	}

	doc, err := getDocumentFromURL(logger, conf.Modem.Url+"/cmconnectionstatus.html", conf, token)
	if err != nil {
		return nil, err
	}
	connectionStatus := scrapeConnectionStatus(doc)

	doc, err = getDocumentFromURL(logger, conf.Modem.Url+"/cmswinfo.html", conf, token)
	if err != nil {
		return nil, err
	}
	softwareInformation := scrapeSoftwareInformation(logger, doc)

	doc, err = getDocumentFromURL(logger, conf.Modem.Url+"/cmeventlog.html", conf, token)
	if err != nil {
		return nil, err
	}
	eventLog := scrapeEventLogs(logger, doc)

	// Logout to let the modem reclaim resources, per https://github.com/mdonoughe/modem_status
	getDocumentFromURL(logger, conf.Modem.Url+"/logout.html", conf, token)

	modemInformation := ModemInformation{
		ConnectionStatus:    *connectionStatus,
		SoftwareInformation: *softwareInformation,
		EventLog:            eventLog,
	}

	return &modemInformation, nil
}

func getDocumentFromURL(logger *zap.Logger, address string, conf config.Configuration, token string) (*goquery.Document, error) {
	logger.Debug(fmt.Sprintf("grabbing %s", address),
		zap.String("op", "scrape.getDocumentFromURL"),
	)

	start := time.Now()

	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "credential",
		Value:  token,
		Path:   "/",
		Domain: "",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(address)
	jar.SetCookies(u, cookies)

	// The modem has an ancient cert loaded and there is no option to replace it
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(conf.Modem.Username, conf.Modem.Password)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	logger.Debug(fmt.Sprintf("got %s, took %s", address, elapsed),
		zap.String("op", "scrape.getDocumentFromURL"),
	)

	// logger.Debug(fmt.Sprintf("%s", doc.Text()),
	// 	zap.String("op", "scrape.getDocumentFromURL"),
	// )

	return doc, nil
}

func getToken(logger *zap.Logger, conf config.Configuration) (string, error) {

	logger.Info(fmt.Sprintf("Attempting to renew token"),
		zap.String("op", "scrape.getToken"),
	)
	// The modem has an ancient cert loaded and there is no option to replace it
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	authString := conf.Modem.Username + ":" + conf.Modem.Password
	basicAuthString := base64.StdEncoding.EncodeToString([]byte(authString))

	req, err := http.NewRequest("GET", conf.Modem.Url+"/cmconnectionstatus.html?"+basicAuthString, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(conf.Modem.Username, conf.Modem.Password)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	token := string(bodyBytes)
	if len(token) != 31 {
		return "", fmt.Errorf("did not retrieve auth token successfully")
	}

	logger.Info(fmt.Sprintf("Got new token: %s", token),
		zap.String("op", "scrape.getToken"),
	)

	return token, nil
}
