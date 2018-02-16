package main

import (
	"log"
	"net"
	"strings"
)

func serveUDP() {

	sourceAddr, err := net.ResolveUDPAddr("udp", *lAddrUDP)
	if err != nil {
		log.Fatal("Could not resolve source address:", *lAddrUDP)
		return
	}

	var targetAddr []*net.UDPAddr
	for _, v := range strings.Split(*rAddr, ",") {
		addr, err := net.ResolveUDPAddr("udp", v)
		if err != nil {
			log.Fatal("Could not resolve target address:", v)
			return
		}
		targetAddr = append(targetAddr, addr)
	}

	sourceConn, err := net.ListenUDP("udp", sourceAddr)
	if err != nil {
		log.Fatal("Could not listen on address:", sourceAddr.String())
		return
	}

	defer sourceConn.Close()
	log.Printf(">> Starting heplify-proxy, listen at: %v", sourceAddr.String())
	var targetConn []*net.UDPConn
	for _, v := range targetAddr {
		conn, err := net.DialUDP("udp", nil, v)
		if err != nil {
			log.Fatal("Could not connect to target address:", v.String())
			return
		}

		defer conn.Close()
		log.Printf(">> Starting heplify-proxy, send to: %v", v.String())
		targetConn = append(targetConn, conn)
	}

	for {
		b := make([]byte, 8192)
		n, addr, err := sourceConn.ReadFromUDP(b)

		if err != nil {
			log.Printf("Could not receive a packet from: %v\n", addr.String())
			continue
		}

		for _, v := range targetConn {
			if _, err := v.Write(b[:n]); err != nil {
				log.Printf("Could not forward packet to: %v\n", v.RemoteAddr().String())
			}
		}
	}
}
