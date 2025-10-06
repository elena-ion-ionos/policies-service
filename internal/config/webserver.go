package config

import (
	"github.com/spf13/cobra"
)

type Webserver struct {
	Port    int
	Service Service
}

func (w *Webserver) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&w.Port, "port", 8080, "Port to start the webserver API on.")
	w.Service.AddFlags(cmd)
}
