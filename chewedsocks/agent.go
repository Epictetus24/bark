package chewedsocks

import (
	"crypto/tls"
	"log"

	"github.com/armon/go-socks5"
	"github.com/lucas-clemente/quic-go"
)

// Connect to SOCKS server
func AgentConnect(address string, agentpassword string, verifytls bool) error {
	var conn quic.Connection
	server, err := socks5.New(&socks5.Config{})
	if err != nil {
		return err
	}

	conf := &tls.Config{
		InsecureSkipVerify: verifytls,
	}

	conn, err = quic.DialAddr(address, conf, nil)
	if err != nil {
		return err
	}

	//time.Sleep(time.Second * 1)
	session, err := newConn(conn)
	session.Write([]byte(agentpassword))

	for {
		log.Println("Acceping stream")
		if err != nil {
			return err
		}
		log.Println("Passing off to socks5")
		go func() {
			err = server.ServeConn(session)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
