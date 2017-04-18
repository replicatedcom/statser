package main

import (
	"fmt"
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

func Ping(server string) {
	fmt.Println("Pinging ", server)
	graphiteKey := fmt.Sprintf("ping.%s", cleanKey(server))

	pinger := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", server)
	if err != nil {
		fmt.Println("Failed to resolve %s: %v", server, err)
		return
	}
	pinger.AddIPAddr(ra)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		sendGauge(graphiteKey, rtt.Nanoseconds()/1000)
	}
	pinger.OnIdle = func() {
		// no-op
	}
	pinger.RunLoop()
}
