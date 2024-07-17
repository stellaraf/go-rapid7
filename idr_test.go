package rapid7_test

import (
	"testing"

	"github.com/stellaraf/go-rapid7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_IDR_URL(t *testing.T) {
	r7 := rapid7.IDR{}
	t.Run("paths 1", func(t *testing.T) {
		path := r7.URL("/thing")
		assert.Equal(t, "/idr/thing", path)
	})
	t.Run("paths 2", func(t *testing.T) {
		path := r7.URL("/thing/one/thing/two")
		assert.Equal(t, "/idr/thing/one/thing/two", path)
	})
}

func TestClient_IDR_Investigation(t *testing.T) {
	t.Run("get investigation", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		inv, err := r7.IDR.Investigation(env.InvestigationRRN)
		require.NoError(t, err)
		assert.Contains(t, inv.Title, env.InvestigationContains)
	})
}

func TestClient_IDR_Investigations(t *testing.T) {
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
			Statuses: []rapid7.InvestigationStatus{rapid7.CLOSED},
		}
		invs, err := r7.IDR.Investigations(q)
		require.NoError(t, err)
		assert.True(t, len(invs) != 0)
	})
	t.Run("get investigations with multiple statuses", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		q := &rapid7.InvestigationsQuery{
			Statuses: []rapid7.InvestigationStatus{rapid7.CLOSED, rapid7.WAITING},
		}
		invs, err := r7.IDR.Investigations(q)
		require.NoError(t, err)
		assert.True(t, len(invs) != 0)
	})

	t.Run("get investigations all", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		res, err := r7.IDR.InvestigationsResponse(&rapid7.InvestigationsQuery{Size: 2})
		assert.Len(t, res.Data, 2)
		require.NoError(t, err)
		invs, err := r7.IDR.InvestigationsAll()
		require.NoError(t, err)
		assert.Equal(t, res.Metadata.TotalData, int64(len(invs)))
	})
}

func TestClient_IDR_InvestigationComments(t *testing.T) {
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

func TestClient_IDR_UpdateInvestigation(t *testing.T) {
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

func TestClient_IDR_Assets(t *testing.T) {
	t.Run("assets", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		assets, err := r7.IDR.Assets()
		require.NoError(t, err)
		assert.True(t, len(assets) > 0, "no assets")
	})
}

func TestClient_AssetCount(t *testing.T) {
	env, err := LoadEnv()
	require.NoError(t, err)
	gql, err := rapid7.NewGraphQLClient(env.Region, env.APIKey)
	require.NoError(t, err)
	res, err := gql.AssetCount(env.OrgID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, int(res.Organization.Assets.TotalCount), 10)
}
