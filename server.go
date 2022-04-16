package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

type Server struct {
	host     string
	port     int
	rTimeout time.Duration
	wTimeout time.Duration
}

func (s *Server) Addr() string {
	return net.JoinHostPort(s.host, strconv.Itoa(s.port))
}

func Ips() (map[string]string, error) {

	ips := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {

			var ip net.IP
			switch v := v.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}

			ips[byName.Name] = ip.String()
		}
	}
	return ips, nil
}

func (s *Server) Run() {
	Handler := NewHandler()

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", Handler.DoTCP)

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", Handler.DoUDP)

	mInterfaces, err := Ips()
	if err == nil {
		fmt.Println(mInterfaces)
	}

	fmt.Println(s.host)
	// fmt.Println(s.Addr())

	addrOpt := s.host
	flagMatch := false

	for _, v := range mInterfaces {
		if v == addrOpt {
			flagMatch = true
		}
	}
	if !flagMatch {
		addrOpt = ""
		for _, v := range mInterfaces {
			addrOpt = v
			break
		}
	}
	if addrOpt == "" {
		panic("网络出错 地址出错")
	}
	DNSHOST := os.Getenv("DNSHOST")
	if DNSHOST != "" {
		addrOpt = DNSHOST
		fmt.Println("使用配置的NDS : ", DNSHOST)
	}
	addrJoin := net.JoinHostPort(addrOpt, strconv.Itoa(s.port))

	tcpServer := &dns.Server{Addr: addrJoin,
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout}

	udpServer := &dns.Server{Addr: addrJoin,
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout}

	go s.start(udpServer)
	go s.start(tcpServer)

}

func (s *Server) start(ds *dns.Server) {

	logger.Info("Start %s listener on %s", ds.Net, s.Addr())
	err := ds.ListenAndServe()
	if err != nil {
		logger.Error("Start %s listener on %s failed:%s", ds.Net, s.Addr(), err.Error())
	}

}
