package main

import (
	"test/logstash"

	"github.com/rs/zerolog"
)

func main() {
	ls := logstash.New("0.0.0.0", 5046, 1)
	ls.Connect()
	l := zerolog.New(ls)

	l.Debug().Msg("1234565")

}
