package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
)

var (
	statsdClient      statsd.Statter
	httpClient        *http.Client
	curStatsdEndpoint = ""
	newStatsdEndpoint = ""
	endPointMutex     sync.Mutex
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient = &http.Client{Transport: tr}

	go monitorStatsdEndpoint()
	go Ping("8.8.8.8")
	go Ping("8.8.4.4")
	go Disk("/")
	// go generateGauge("myapp.mygauge1", 555)
	// go generateGauge("myapp.mygauge2", 9999)

	<-make(chan interface{})
}

func generateGauge(name string, max int) {
	delay := func() {}
	for {
		delay()
		delay = func() { <-time.After(2 * time.Second) }

		sendGauge(name, int64(rand.Intn(max)))
	}
}

func sendGauge(name string, val int64) {
	maybeReplaceStatsdClient()
	if statsdClient == nil {
		return
	}

	//fmt.Println("Sending gauge")
	statsdClient.Gauge(name, val, 1)
}

func maybeReplaceStatsdClient() {
	endPointMutex.Lock()
	defer endPointMutex.Unlock()

	if statsdClient != nil && curStatsdEndpoint == newStatsdEndpoint {
		return
	}

	if statsdClient != nil {
		statsdClient.Close()
		statsdClient = nil
	}

	curStatsdEndpoint = newStatsdEndpoint
	if curStatsdEndpoint == "" {
		fmt.Println("Cannot send stats. No statsd endpoint defined.")
		return
	}

	client, err := statsd.NewBufferedClient(curStatsdEndpoint, "myapp100", 0, 0)
	if err != nil {
		fmt.Println("Could not create statsd client", err)
		return
	}

	fmt.Println("Made new statsd client, pointing at", curStatsdEndpoint)
	statsdClient = client
}

func monitorStatsdEndpoint() {
	delay := func() {}
	for {
		delay()
		delay = func() { <-time.After(30 * time.Second) }

		fmt.Println("Checking statsd endpoint")
		apiEndpoint := os.Getenv("REPLICATED_INTEGRATIONAPI")
		if apiEndpoint == "" {
			fmt.Println("REPLICATED_INTEGRATIONAPI is not set")
			os.Exit(-1)
		}

		url := fmt.Sprintf("%s/console/v1/option?name=statsd.endpoint", apiEndpoint)
		fmt.Printf("Calling %s\n", url)
		resp, err := httpClient.Get(url)
		if err != nil {
			fmt.Printf("%s: ERROR: %v\n", url, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%v: BAD STATUS: %v\n", resp.StatusCode, err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Could not read responce body: %v\n", err)
			continue
		}

		fmt.Printf("Got new statsd endpoint: %s\n", body)
		newStatsdEndpoint = string(body)
	}
}

func round(v float64) int64 {
	return int64(math.Floor(v + 0.5))
}

func cleanKey(key string) string {
	return strings.Replace(key, ".", "_", -1)
}
