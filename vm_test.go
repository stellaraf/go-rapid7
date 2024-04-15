package rapid7_test

import (
	"testing"

	"github.com/stellaraf/go-rapid7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_VM_URL(t *testing.T) {
	r7 := rapid7.VM{}
	t.Run("paths 1", func(t *testing.T) {
		path := r7.URL("/thing")
		assert.Equal(t, "/vm/thing", path)
	})
	t.Run("paths 2", func(t *testing.T) {
		path := r7.URL("/thing/one/thing/two")
		assert.Equal(t, "/vm/thing/one/thing/two", path)
	})
}

func TestClient_VM_Assets(t *testing.T) {
	t.Run("assets", func(t *testing.T) {
		t.Parallel()
		env, err := LoadEnv()
		require.NoError(t, err)
		r7, err := rapid7.New(env.Region, env.APIKey)
		require.NoError(t, err)
		assets, err := r7.VM.Assets()
		require.NoError(t, err)
		assert.True(t, len(assets) > 0, "no assets")
		assert.IsType(t, "", assets[0].HostName)
		assert.True(t, len(assets[0].HostName) > 0)
	})
}
