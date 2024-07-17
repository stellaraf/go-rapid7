package rapid7_test

import (
	"os"
	"testing"

	"github.com/stellaraf/go-rapid7"
	"github.com/stellaraf/go-utils/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Env struct {
	Region                string `env:"RAPID7_REGION"`
	APIKey                string `env:"RAPID7_API_KEY"`
	InvestigationRRN      string `env:"RAPID7_INVESTIGATION_RRN"`
	OrgID                 string `env:"RAPID7_ORG_ID"`
	InvestigationContains string `env:"RAPID7_INVESTIGATION_CONTAINS"`
}

func LoadEnv() (env Env, err error) {
	ci := os.Getenv("CI")
	opts := &environment.EnvironmentOptions{
		DotEnv: ci == "",
	}
	err = environment.Load(&env, opts)
	return
}

func Test_New(t *testing.T) {
	t.Run("create client", func(t *testing.T) {
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		assert.NotNil(t, r7)
	})
}
