package howl

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"

	"github.com/lucas-clemente/quic-go"
)

//Streams are used only with QUIC, and allow for the sending and recieving byte streams

type Validator func([]byte) bool

// Defines the howl struct
type Howl struct {
	//quic bits.
	Conn   quic.Connection
	Stream quic.Stream //Bidirectional QUIC Stream

	Listener quic.Listener //Pure Quic Listener

	Validator Validator // Function for validating the the connection before accepting unidirectional comms.
	MsgDelim  string    //Delimeter to set end of message - "/n" by default.
}

// DefaultValidator accepts any stream which contains the string "woofwoof"
func DefaultValidator(b []byte) bool {
	if strings.Contains(string(b), "woofwoof") {
		return true

	} else {
		return false
	}

}

// NewHowlListener starts a listener for Howl
func (h *Howl) NewHowlListener(addr string, tlsconfig *tls.Config) error {
	var err error
	h.Listener, err = quic.ListenAddr(addr, tlsconfig, nil)
	if err != nil {
		return err
	}
	for {
		sess, err := h.Listener.Accept(context.Background())
		if err != nil {
			return err
		}

		msg, err := sess.ReceiveMessage()
		if err != nil {
			return err
		}

		valid := h.Validator(msg)

		if !valid {
			sess.CloseWithError(0, "Your connection had an issue while streaming data, connection has been closed.")
		}

		go func() {
			defer func() {
				_ = sess.CloseWithError(0, "bye")
				fmt.Errorf("close session: %s", sess.RemoteAddr().String())
			}()
			h.communicate(sess)
		}()
	}

}

// Connect to a quic server and start a stream.
func (h *Howl) Connect(addr string, verifytls bool, validator []byte) error {
	var err error
	tlsConf := &tls.Config{
		InsecureSkipVerify: verifytls,
		NextProtos:         []string{"howl"},
	}
	h.Conn, err = quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}

	err = h.Conn.SendMessage(validator)
	if err != nil {
		return err
	}

	h.Stream, err = h.Conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Bi-directional comms for Howl.
func (h *Howl) communicate(sess quic.Connection) error {
	var err error
	for {
		h.Stream, err = sess.AcceptStream(context.Background())
		if err != nil {
			return err
		}

	}
}

// Send data over the bi-directional stream
func (h *Howl) Send(message []byte) (len int, err error) {
	len, err = h.Stream.Write(message)
	if err != nil {
		return 0, err
	}
	return len, err

}

// Recieve data over the bi-directional stream
func (h *Howl) Recieve(msglen int) (len int, message []byte, err error) {

	if h.MsgDelim == "" {
		h.MsgDelim = "\n"
	}

	reader := bufio.NewReader(h.Stream)

	len, err = reader.Read(message)
	if err != nil {
		if err != io.EOF {
			return len, message, err
		}

	}

	return len, message, nil

}
