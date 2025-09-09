package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitViperFlags(cmd *cobra.Command, args []string) error {
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err != nil {
			return
		}
		if !f.Changed && viper.IsSet(f.Name) {
			err = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", viper.Get(f.Name)))
		}
	})
	return err
}
