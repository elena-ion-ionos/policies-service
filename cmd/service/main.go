package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/policies-service/internal/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mainCmd = &cobra.Command{}

func main() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	mainCmd.AddCommand(cmd.WebserverUser())

	if err := mainCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v", err)
		os.Exit(1)
	}
}
