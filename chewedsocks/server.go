package chewedsocks

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucas-clemente/quic-go"
)

//socks package based on https://github.com/brimstone/rsocks/ && https://github.com/kost/revsocks

// listen for chewedsocks agents connecting via quic.
func StartQUICListener(quicaddr, agentpassword, socksaddr, cert, certkey string) error {
	var err error
	portinc := 0

	// load tls cert
	cer, err := tls.LoadX509KeyPair(cert, certkey)
	if err != nil {

		return err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	listener, err := quic.ListenAddr(quicaddr, config, nil)
	if err != nil {
		log.Printf("Error listening on %s: %v", quicaddr, err)
		return err
	}

	for {
		conn, err := listener.Accept(context.Background())
		conn.RemoteAddr()
		agentstr := conn.RemoteAddr().String()
		log.Printf("[%s] Got a connection from %v: ", agentstr, conn.RemoteAddr())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Errors accepting!")
		}

		timeout, err := time.ParseDuration("3s")
		if err != nil {
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), timeout)

		session, err := conn.AcceptStream(ctx)
		if err != nil {
			return err
		}

		reader := bufio.NewReader(session)

		//read only 64 bytes with timeout=1-3 sec. So we haven't delay with browsers
		statusb := make([]byte, 64)
		_, _ = io.ReadFull(reader, statusb)

		if strings.Contains(string(statusb), agentpassword) {
			//magic bytes received.
			//disable socket read timeouts
			log.Printf("[%s] Got Client from %s", agentstr, conn.RemoteAddr())
			session.SetReadDeadline(time.Now().Add(100 * time.Hour))
			listenstr := strings.Split(socksaddr, ":")
			portnum, err := strconv.Atoi(listenstr[1])
			if err != nil {
				return err
			}

			go ListenForClients(agentstr, listenstr[0], portnum+portinc, session)

		} else {

			//Check for Agent Password failed, close session

			session.Close()
		}

	}

}

// Catches local clients and connects to QUIC stream. Agentstr is address for
func ListenForClients(agentstr string, listen string, port int, session quic.Stream) error {
	var ln net.Listener
	var address string
	var err error
	portinc := port

	//Start Socks listener over TCP

	for {
		address = fmt.Sprintf("%s:%d", listen, portinc)
		log.Printf("[%s] Waiting for clients on %s", agentstr, address)
		ln, err = net.Listen("tcp", address)
		if err != nil {
			log.Printf("[%s] Error listening on %s: %v", agentstr, address, err)
			portinc = portinc + 1
		} else {
			break
		}
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[%s] Error accepting on %s: %v", agentstr, address, err)
			return err
		}
		if session == nil {
			log.Printf("[%s] Session on %s is nil", agentstr, address)
			conn.Close()
			continue
		}
		log.Printf("[%s] Got client. Opening stream for %s", agentstr, conn.RemoteAddr())

		// connect both of conn and stream

		go func() {
			log.Printf("[%s] Starting to copy conn to stream for %s", agentstr, conn.RemoteAddr())
			io.Copy(conn, session)
			conn.Close()
			log.Printf("[%s] Done copying conn to stream for %s", agentstr, conn.RemoteAddr())
		}()
		go func() {
			log.Printf("[%s] Starting to copy stream to conn for %s", agentstr, conn.RemoteAddr())
			io.Copy(session, conn)
			session.Close()
			log.Printf("[%s] Done copying stream to conn for %s", agentstr, conn.RemoteAddr())
		}()
	}
}
