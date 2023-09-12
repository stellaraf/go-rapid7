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
	Region           string `env:"RAPID7_REGION"`
	APIKey           string `env:"RAPID7_API_KEY"`
	InvestigationRRN string `env:"RAPID7_INVESTIGATION_RRN"`
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
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		assert.NotNil(t, r7)
	})
}

func TestClient_Investigation(t *testing.T) {
	t.Run("get investigation", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		inv, err := r7.IDR.Investigation(env.InvestigationRRN)
		require.NoError(t, err)
		assert.Contains(t, inv.Title, "mlove")
	})
}

func TestClient_Investigations(t *testing.T) {
	t.Run("get investigations", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		invs, err := r7.IDR.Investigations()
		require.NoError(t, err)
		assert.True(t, len(invs) != 0)
	})
	t.Run("get investigations with query", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		q := &rapid7.InvestigationsQuery{
			Statuses: []string{"CLOSED"},
		}
		invs, err := r7.IDR.Investigations(q)
		require.NoError(t, err)
		assert.True(t, len(invs) != 0)
	})
}

func TestClient_InvestigationComments(t *testing.T) {
	t.Run("get investigation comments", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		inv, err := r7.IDR.Investigation(env.InvestigationRRN)
		require.NoError(t, err)
		comments, err := r7.IDR.InvestigationComments(inv)
		require.NoError(t, err)
		assert.IsType(t, &rapid7.InvestigationComments{}, comments)
		assert.IsType(t, 0, len(comments.Data))
	})
}

var statusMap = map[rapid7.InvestigationStatus]rapid7.InvestigationStatus{
	rapid7.CLOSED:        rapid7.OPEN,
	rapid7.OPEN:          rapid7.INVESTIGATING,
	rapid7.INVESTIGATING: rapid7.OPEN,
}

func TestClient_UpdateInvestigation(t *testing.T) {
	t.Run("update investigation", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		inv, err := r7.IDR.Investigation(env.InvestigationRRN)
		require.NoError(t, err)
		status := statusMap[inv.Status]
		update := &rapid7.InvestigationUpdateRequest{
			Status: status,
		}
		updated, err := r7.IDR.UpdateInvestigation(inv.RRN, update)
		require.NoError(t, err)
		assert.Equal(t, status, updated.Status)
		revert := &rapid7.InvestigationUpdateRequest{
			Status: inv.Status,
		}
		reverted, err := r7.IDR.UpdateInvestigation(inv.RRN, revert)
		require.NoError(t, err)
		assert.Equal(t, inv.Status, reverted.Status)
	})
}
