package options

import (
	"flag"
)

// ServerOption is the main context object for the controller manager.
type ServerOption struct {
	Master        string
	Kubeconfig    string
	JsonLogFormat bool
	CreateCRD     bool
	// Debugging is reserved options
	Debugging *DebuggingOptions
}

// DebuggingOptions holds the Debugging options.
type DebuggingOptions struct {
	EnableProfiling           bool
	EnableContentionProfiling bool
}

// NewServerOption creates a new CMServer with a default config.
func NewServerOption() *ServerOption {
	s := ServerOption{}
	return &s
}

// AddFlags adds flags for a specific CMServer to the specified FlagSet
func (s *ServerOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	fs.BoolVar(&s.JsonLogFormat, "json-log-format", true, "Set true to use json style log format. Set false to use plaintext style log format")
	fs.BoolVar(&s.CreateCRD, "create-crd", false, "set true to create crd resource for jx-controller")
}
