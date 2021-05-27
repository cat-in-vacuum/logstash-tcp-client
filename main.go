package main

import (
	"io"
	"test/logstash"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogstashHook struct {
	w io.Writer
}

func (l LogstashHook) Run(e *zerolog.Event, level zerolog.Level, message string) {

	output := log.Output(l.w)
	output.Debug().Msg("555")

}

func main() {
	ls := logstash.New("0.0.0.0", 5046, 1)
	ls.Connect()
	l := zerolog.New(ls)

	l.Debug().Msg("1234565")

}
