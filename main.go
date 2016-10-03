package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
	geoip2 "github.com/oschwald/geoip2-golang"
)

const dom = "whoami.fluffcomputing.com."

func handleResponse(w dns.ResponseWriter, r *dns.Msg) {

	var (
		rr dns.RR
		a  net.IP
	)

	m := new(dns.Msg)
	m.SetReply(r)

	if ip, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		a = ip.IP
	}

	rr = &dns.A{
		Hdr: dns.RR_Header{Name: dom, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   a.To4(),
	}

	// Getting the system ready:
	// Open DB
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	city, country := GeoIP(a.To4(), db)

	t := &dns.TXT{
		Hdr: dns.RR_Header{Name: dom, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
		Txt: []string{city + ", " + country},
	}

	switch r.Question[0].Qtype {
	case dns.TypeTXT:
		m.Answer = append(m.Answer, t)
		m.Extra = append(m.Extra, rr)
	default:
		m.Answer = append(m.Answer, rr)
	}
	w.WriteMsg(m)
}

func main() {
	// Start Server
	go func() {
		srv := &dns.Server{Addr: ":" + "53", Net: "udp"}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	dns.HandleFunc("whoami.fluffcomputing.com.", handleResponse)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)

}

// GeoIP Resolution
func GeoIP(IPAddr net.IP, db *geoip2.Reader) (city string, country string) {

	record, err := db.City(IPAddr)
	if err != nil {
		log.Fatal(err)
	}

	return record.City.Names["en"], record.Country.Names["en"]
}
