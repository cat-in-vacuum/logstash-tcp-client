package logstash

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	hostname   string
	port       int
	connection *net.TCPConn
	timeout    int
}

func New(hostname string, port int, timeout int) *Client {
	l := Client{}
	l.hostname = hostname
	l.port = port
	l.connection = nil
	l.timeout = timeout
	return &l
}

func (l *Client) Connect() error {
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

func (l *Client) setTimeouts() error {
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

func (l *Client) Write(p []byte) (int, error) {
	if l.connection == nil {
		if err := l.Connect(); err != nil {
			return 0, fmt.Errorf("connect to logstash errored: %s", err)
		}
	}

	n, err := l.connection.Write(p)
	if err != nil {

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			l.connection.Close()
			if err := l.Connect(); err != nil {
				return n, fmt.Errorf("reconnect to logstash errored: %s", err)
			}

		} else {
			l.connection.Close()
			l.connection = nil
			return n, err
		}
	} else {
		err = l.setTimeouts()
		if err != nil {
			return 0, fmt.Errorf("set timeouts errored: %s", err)
		}
		return n, nil
	}

	return 0, fmt.Errorf("write to logstash errored: %s", err)
}
