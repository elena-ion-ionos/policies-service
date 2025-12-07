package config

import (
	"github.com/spf13/cobra"
)

type Webserver struct {
	Port       int
	ServerHost string
	Service    Service
}

func (w *Webserver) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&w.Port, "port", 8081, "Port to start the webserver API on.")
	cmd.Flags().StringVar(&w.ServerHost, "server-host", "localhost", "Port to start the webserver API on.")
	w.Service.AddFlags(cmd)
}
