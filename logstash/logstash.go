package logstash

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Logstash struct {
	hostname   string
	port       int
	connection *net.TCPConn
	timeout    int
}

func New(hostname string, port int, timeout int) *Logstash {
	l := Logstash{}
	l.hostname = hostname
	l.port = port
	l.connection = nil
	l.timeout = timeout
	return &l
}

func (l *Logstash) Connect() error {
	var connection *net.TCPConn
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", l.hostname, l.port))
	if err != nil {
		return err
	}

	connection, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	err = connection.SetLinger(0)
	if err != nil {
		return err
	}
	err = connection.SetNoDelay(true)
	if err != nil {
		return err
	}
	err = connection.SetKeepAlive(true)
	if err != nil {
		return err
	}
	err = connection.SetKeepAlivePeriod(time.Duration(5) * time.Second)
	if err != nil {
		return err
	}

	l.connection = connection

	err = l.setTimeouts()
	if err != nil {
		return err
	}

	return nil
}

func (l *Logstash) setTimeouts() error {
	timeout := time.Now().Add(time.Duration(l.timeout) * time.Second)
	err := l.connection.SetDeadline(timeout)
	if err != nil {
		return err
	}
	err = l.connection.SetWriteDeadline(timeout)
	if err != nil {
		return err
	}
	err = l.connection.SetReadDeadline(timeout)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logstash) Write(p []byte) (int, error) {
	err := fmt.Errorf("tcp conn is nil")
	if l.connection != nil {
		n, err := l.connection.Write(p)
		if err != nil {

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				l.connection.Close()
				if err := l.Connect(); err != nil {
					log.Printf("connection error: %s\n", err)
					return n, err
				}

			} else {
				l.connection.Close()
				l.connection = nil
				return n, err
			}
		} else {
			err = l.setTimeouts()
			return n, nil
		}
	}
	return 0, err
}
