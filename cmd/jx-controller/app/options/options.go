package options

import (
	"flag"
	"time"
)

// ServerOption is the main context object for the controller manager.
type ServerOption struct {
	ChaosLevel             int
	ControllerConfigFile   string
	PrintVersion           bool
	GCInterval             time.Duration
	JsonLogFormat          bool
	EnableGangScheduling   bool
	EnableStreamScheduling bool
	NameSpace              string
	CreateCRD              bool
}

// NewServerOption creates a new CMServer with a default config.
func NewServerOption() *ServerOption {
	s := ServerOption{}
	return &s
}

// AddFlags adds flags for a specific CMServer to the specified FlagSet
func (s *ServerOption) AddFlags(fs *flag.FlagSet) {
	// chaos level will be removed once we have a formal tool to inject failures.
	fs.IntVar(&s.ChaosLevel, "chaos-level", -1, "DO NOT USE IN PRODUCTION - level of chaos injected into the TFJob created by the operator.")
	fs.BoolVar(&s.PrintVersion, "version", false, "Show version and quit")
	fs.DurationVar(&s.GCInterval, "gc-interval", 10*time.Minute, "GC interval")
	fs.StringVar(&s.ControllerConfigFile, "controller-config-file", "", "Path to file containing the controller config.")
	fs.StringVar(&s.NameSpace, "namespace", "default", "namespace for jx-controller,default as default.")
	fs.BoolVar(&s.JsonLogFormat, "json-log-format", true, "Set true to use json style log format. Set false to use plaintext style log format")
	fs.BoolVar(&s.EnableGangScheduling, "enable-gang-scheduling", false, "Set true to enable gang scheduling by kube-arbitrator.")
	fs.BoolVar(&s.EnableStreamScheduling, "enable-stream-scheduling", false, "Set true to enable steam scheduling by posedion.")
	fs.BoolVar(&s.CreateCRD, "create-crd", false, "set true to create crd resource for jx-controller")
}
