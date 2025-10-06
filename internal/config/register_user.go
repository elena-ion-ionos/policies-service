package config

import "github.com/spf13/cobra"

// Worker is the configuration for the worker.
type RegisterUser struct {
	Service

	SampleIntOption int
}

const defaultSampleIntOption int = 10

// AddFlags adds the flags for the worker.
func (w *RegisterUser) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&w.SampleIntOption, "sample-int-option", defaultSampleIntOption, "Sample int option.")

	w.Service.AddFlags(cmd)
}
