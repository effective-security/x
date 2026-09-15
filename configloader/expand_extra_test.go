package configloader_test

import (
	"testing"

	"github.com/effective-security/x/configloader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type envAlias map[string]string

type expandTarget struct {
	Name    string
	Env     envAlias
	Any     map[string]any
	Nested  map[string]expandNested
	Ptrs    map[string]*expandNested
	Missing string
}

type expandNested struct {
	Value string
}

func TestExpand_MapsAndAliases(t *testing.T) {
	t.Setenv("EXPAND_NAME", "from-env")

	e := &configloader.Expander{
		Variables: map[string]string{
			"NAME": "alice",
		},
	}

	cfg := &expandTarget{
		Name: "${NAME}",
		Env: envAlias{
			"user": "${NAME}",
		},
		Any: map[string]any{
			"greet": "hi ${NAME}",
		},
		Nested: map[string]expandNested{
			"one": {Value: "${NAME}"},
		},
		Ptrs: map[string]*expandNested{
			"two": {Value: "env://EXPAND_NAME"},
		},
		Missing: "${NOT_SET_XYZ}",
	}

	err := e.ExpandAll(cfg)
	require.NoError(t, err)
	assert.Equal(t, "alice", cfg.Name)
	assert.Equal(t, "alice", cfg.Env["user"])
	assert.Equal(t, "hi alice", cfg.Any["greet"])
	assert.Equal(t, "alice", cfg.Nested["one"].Value)
	require.NotNil(t, cfg.Ptrs["two"])
	assert.Equal(t, "from-env", cfg.Ptrs["two"].Value)
	assert.Equal(t, "", cfg.Missing)
}

func TestExpand_SecretInterpolationError(t *testing.T) {
	e := &configloader.Expander{}
	_, err := e.Expand("token=${secret://api}")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "secret loader not provided")

	e.SecretProvider = &mockSecret{secrets: map[string]string{}}
	_, err = e.Expand("token=${secret://missing}")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unable to load secret: missing")
}

func TestExpand_EnvironmentNonStringIgnored(t *testing.T) {
	t.Parallel()
	type cfg struct {
		Environment int
		Name        string
	}
	c := &cfg{Environment: 7, Name: "ok"}
	err := configloader.ExpandAll(c)
	require.NoError(t, err)
	assert.Equal(t, 7, c.Environment)
	assert.Equal(t, "ok", c.Name)
}
