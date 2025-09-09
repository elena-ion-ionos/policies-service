package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitViperFlags_SetsFlagFromViper(t *testing.T) {
	viper.Reset()

	cmd := &cobra.Command{}
	var foo string
	cmd.Flags().StringVar(&foo, "foo", "default", "foo flag")

	viper.Set("foo", "bar")

	err := InitViperFlags(cmd, []string{})
	assert.NoError(t, err)
	assert.Equal(t, "bar", foo)
}

func TestInitViperFlags_DoesNotOverrideChangedFlag(t *testing.T) {
	viper.Reset()

	cmd := &cobra.Command{}
	var foo string
	cmd.Flags().StringVar(&foo, "foo", "default", "foo flag")

	// Simulate user setting the flag
	assert.NoError(t, cmd.Flags().Set("foo", "user"))
	viper.Set("foo", "bar")

	err := InitViperFlags(cmd, []string{})
	assert.NoError(t, err)
	assert.Equal(t, "user", foo)
}

func TestInitViperFlags_NoViperValueKeepsDefault(t *testing.T) {
	viper.Reset()
	cmd := &cobra.Command{}
	var foo string
	cmd.Flags().StringVar(&foo, "foo", "default", "foo flag")

	// viper does not have "foo" set
	err := InitViperFlags(cmd, []string{})
	assert.NoError(t, err)
	assert.Equal(t, "default", foo)
}
