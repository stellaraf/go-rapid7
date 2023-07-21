package rapid7_test

import (
	"os"
	"testing"

	"github.com/stellaraf/go-rapid7"
	"github.com/stellaraf/go-utils/environment"
	"github.com/stretchr/testify/assert"
)

const TEST_INVESTIGATION_ID string = "79c0e558-7dbe-49d3-b007-1ec93266214b"

type Env struct {
	Region string `env:"RAPID7_REGION"`
	APIKey string `env:"RAPID7_API_KEY"`
}

func LoadEnv() (env Env, err error) {
	ci := os.Getenv("CI")
	opts := &environment.EnvironmentOptions{
		DotEnv: ci == "",
	}
	err = environment.Load(&env, opts)
	return
}

func Test_NewIDR(t *testing.T) {
	t.Run("create client", func(t *testing.T) {
		env, err := LoadEnv()
		assert.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		assert.NoError(t, err)
		assert.NotNil(t, r7)
	})
}

func TestClient_Investigation(t *testing.T) {
	t.Run("get investigation", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		assert.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		assert.NoError(t, err)
		inv, err := r7.IDR.Investigation(TEST_INVESTIGATION_ID)
		assert.NoError(t, err)
		assert.Contains(t, inv.Title, "mlove")
	})
}

func TestClient_Investigations(t *testing.T) {
	t.Run("get investigations", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		assert.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		assert.NoError(t, err)
		invs, err := r7.IDR.Investigations()
		assert.NoError(t, err)
		assert.True(t, len(invs) != 0)
	})
	t.Run("get investigations with query", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		assert.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		assert.NoError(t, err)
		q := &rapid7.InvestigationsQuery{
			Statuses: []string{"CLOSED"},
		}
		invs, err := r7.IDR.Investigations(q)
		assert.NoError(t, err)
		assert.True(t, len(invs) != 0)
	})
}

func TestClient_InvestigationComments(t *testing.T) {
	t.Run("get investigation comments", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		assert.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		assert.NoError(t, err)
		inv, err := r7.IDR.Investigation(TEST_INVESTIGATION_ID)
		assert.NoError(t, err)
		comments, err := r7.IDR.InvestigationComments(inv)
		assert.NoError(t, err)
		assert.IsType(t, &rapid7.InvestigationComments{}, comments)
		assert.IsType(t, 0, len(comments.Data))
	})
}
