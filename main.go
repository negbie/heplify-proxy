package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var (
	lAddrTCP = flag.String("ht", "", "TCP listen address to terminate HEP TLS traffic")
	lAddrUDP = flag.String("hu", "", "UDP listen address to terminate HEP UDP traffic")
	rAddr    = flag.String("hs", "", "UDP HEP server addresses")
	fileLog  = flag.Bool("fl", false, "Log to file")
)

func panicIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *fileLog {
		f, err := os.OpenFile("heplify-proxy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		panicIfErr(err)
		defer f.Close()
		log.SetOutput(f)
	}

	if *lAddrUDP == "" && *lAddrTCP == "" {
		log.Printf("Please type a UDP listen address like ./heplify-proxy -hu 0.0.0.0:9060\n")
		log.Printf("or type a TCP listen address like ./heplify-proxy -ht 0.0.0.0:9061\n\n")
		os.Exit(1)
	}
	if *rAddr == "" {
		log.Printf("Please type a UDP HEP server address like ./heplify-proxy -hs 192.168.1.1:9060\n")
		log.Printf("or type multiple UDP HEP server addresses like ./heplify-proxy -hs 192.168.1.1:9060,192.168.1.2:9060\n\n")
		os.Exit(1)
	}

	if *lAddrTCP != "" {
		tmp := strings.Split(*rAddr, ",")
		if len(tmp) > 1 {
			log.Printf("Please use only one UDP HEP server address for TLS\n")
			os.Exit(1)
		}
		serveTLS()
	}

	if *lAddrUDP != "" {
		serveUDP()
	}
}
