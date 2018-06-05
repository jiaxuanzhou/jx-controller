package main

import (
	"github.com/jiaxuanzhou/jx-controller/cmd/jx-controller/app"
	"github.com/jiaxuanzhou/jx-controller/cmd/jx-controller/app/options"
	"github.com/golang/glog"

	"flag"
)


func main() {
	s := options.NewServerOption()
	s.AddFlags(flag.CommandLine)

	flag.Parse()

	if err := app.Run(s); err != nil {
		glog.Fatalf("%v\n", err)
	}
}
