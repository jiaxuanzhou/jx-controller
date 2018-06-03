package main

import (
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	"github.com/jiaxuanzhou/jx-controller/cmd/jx-controller/app/options"

	"flag"
)

func init() {
	// Add filename as one of the fields of the structured log message
	filenameHook := filename.NewHook()
	filenameHook.Field = "filename"
	log.AddHook(filenameHook)
}

func main() {
	s := options.NewServerOption()
	s.AddFlags(flag.CommandLine)

	flag.Parse()

	if s.JsonLogFormat {
		// Output logs in a json format so that it can be parsed by services like Stackdriver
		log.SetFormatter(&log.JSONFormatter{})
	}

	if err := app.Run(s); err != nil {
		log.Fatalf("%v\n", err)
	}

}

