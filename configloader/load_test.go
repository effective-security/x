package configloader_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/effective-security/x/configloader"
	"github.com/effective-security/x/guid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadConfigWithSchema_plain(t *testing.T) {
	c, err := configloader.ResolveValue("test_data")
	require.NoError(t, err)
	assert.Equal(t, "test_data", c)
}

func Test_LoadConfigWithSchema_file(t *testing.T) {
	c, err := configloader.ResolveValue("file://./load.go")
	require.NoError(t, err)
	require.NotEmpty(t, c)
	assert.Contains(t, c, "package configloader")
}

func Test_SaveConfigWithSchema_file(t *testing.T) {
	tmpDir := t.TempDir()
	file := path.Join(tmpDir, guid.MustCreate())
	cfg := "file://" + file

	err := os.WriteFile(file, []byte("test"), os.ModePerm)
	require.NoError(t, err)

	c, err := configloader.ResolveValue(cfg)
	require.NoError(t, err)
	assert.Equal(t, "test", c)

	t.Setenv("PORTO_TEST", "test")
	c, err = configloader.ResolveValue("env://PORTO_TEST")
	require.NoError(t, err)
	assert.Equal(t, "test", c)
}

func Test_ConfigWithSchema_Secret(t *testing.T) {
	cfg := "secret://key1"

	configloader.SecretProviderInstance = nil
	_, err := configloader.ResolveValue(cfg)
	assert.EqualError(t, err, "secret loader not provided: unable to expand: secret://key1")

	configloader.SecretProviderInstance = &mockSecret{
		secrets: map[string]string{
			"key1": "value1",
		},
	}

	val, err := configloader.ResolveValue(cfg)
	require.NoError(t, err)
	assert.Equal(t, "value1", val)

	_, err = configloader.ResolveValue("secret://key2")
	assert.EqualError(t, err, "unable to load secret: key2: secret not found: key2")
}

type mockSecret struct {
	secrets map[string]string
}

func (s *mockSecret) GetSecret(name string) (string, error) {
	tokens := strings.Split(name, "/")
	sec := s.secrets[tokens[0]]
	if sec != "" {
		return sec, nil

	}
	return "", errors.Errorf("secret not found: %s", name)
}

type config struct {
	Service     string
	Region      string
	Cluster     string
	Environment string
}

func Test_Unmarshal(t *testing.T) {
	tmp := t.TempDir()

	var v config
	err := configloader.Unmarshal("testdata/test_config.yaml", &v)
	require.NoError(t, err)

	assert.Equal(t, "porto-pod", v.Service)
	assert.Equal(t, "local", v.Region)
	assert.Equal(t, "env://NODENAME", v.Cluster)
	assert.Equal(t, "test", v.Environment)

	fn := path.Join(tmp, "test_config.yaml")
	err = configloader.Marshal(fn, &v)
	require.NoError(t, err)

	var v2 config
	err = configloader.Unmarshal(fn, &v2)
	require.NoError(t, err)
	assert.Equal(t, v, v2)

	err = configloader.Unmarshal("testdata/test_config.json", &v)
	require.NoError(t, err)

	assert.Equal(t, "porto-pod", v.Service)
	assert.Equal(t, "local", v.Region)
	assert.Equal(t, "${NODENAME}", v.Cluster)
	assert.Equal(t, "test", v.Environment)

	fn = path.Join(tmp, "test_config.json")
	err = configloader.Marshal(fn, &v)
	require.NoError(t, err)
	encoded, err := os.ReadFile(fn)
	require.NoError(t, err)
	assert.Equal(t,
		`{
  "Service": "porto-pod",
  "Region": "local",
  "Cluster": "${NODENAME}",
  "Environment": "test"
}`,
		string(encoded))

	err = configloader.Unmarshal(fn, &v2)
	require.NoError(t, err)
	assert.Equal(t, v, v2)
}

func Test_UnmarshalAndExpand(t *testing.T) {
	configloader.SecretProviderInstance = nil
	t.Setenv("NODENAME", "")

	v := new(config)
	err := configloader.UnmarshalAndExpand("testdata/test_config.yaml", v)
	assert.EqualError(t, err, "environment variable not set: NODENAME")

	configloader.SecretProviderInstance = &mockSecret{
		secrets: map[string]string{
			"secret1": "api-key1",
			"secret2": "api-key2",
		},
	}
	t.Setenv("NODENAME", "cluster1")

	v = new(config)
	err = configloader.UnmarshalAndExpand("testdata/test_config.yaml", v)
	require.NoError(t, err)
	assert.Equal(t, "porto-pod", v.Service)
	assert.Equal(t, "local", v.Region)
	assert.Equal(t, "cluster1", v.Cluster)
	assert.Equal(t, "test", v.Environment)
}
