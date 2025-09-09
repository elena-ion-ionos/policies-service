package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func checkRegisterUser(t *testing.T, worker *RegisterUser, cmd *cobra.Command) {
	// Check Worker-specific flag
	assert.NotNil(t, cmd.Flag("sample-int-option"))

	// Check default value
	assert.Equal(t, 10, worker.SampleIntOption)

	checkService(t, &worker.Service, cmd)
}

func TestWorker_AddFlags(t *testing.T) {
	cmd := &cobra.Command{}
	worker := &RegisterUser{}
	worker.AddFlags(cmd)
	checkRegisterUser(t, worker, cmd)
}
