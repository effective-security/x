package configloader

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/effective-security/x/fileutil/resolve"
	"github.com/effective-security/x/netutil"
	"github.com/effective-security/xlog"
	"github.com/oleiade/reflections"
	yamlcfg "go.uber.org/config"
	"gopkg.in/yaml.v3"
)

var logger = xlog.NewPackageLogger("github.com/effective-security/x", "configloader")

// Factory is used to create Configuration instance
type Factory struct {
	nodeInfo    netutil.NodeInfo
	envPrefix   string
	environment string
	overrideCfg string
	searchDirs  []string
	user        *string

	secrets SecretProvider
}

// NewFactory returns new configuration factory
func NewFactory(nodeInfo netutil.NodeInfo, searchDirs []string, envPrefix string) (*Factory, error) {
	var err error
	if nodeInfo == nil {
		nodeInfo, err = netutil.NewNodeInfo(nil)
		if err != nil {
			return nil, err
		}
	}

	return &Factory{
		searchDirs: searchDirs,
		nodeInfo:   nodeInfo,
		envPrefix:  envPrefix,
	}, nil
}

// WithSecretProvider allows to specify secret provider
func (f *Factory) WithSecretProvider(p SecretProvider) *Factory {
	f.secrets = p
	return f
}

// WithOverride allows to specify additional override config file
func (f *Factory) WithOverride(file string) *Factory {
	f.overrideCfg = file
	return f
}

// WithEnvironment allows to override environment in Configuration
func (f *Factory) WithEnvironment(environment string) *Factory {
	f.environment = environment
	return f
}

// GetAbsFilename returns absolute path for the file
// from the relative path to projFolder
func GetAbsFilename(file, projFolder string) (string, error) {
	if !filepath.IsAbs(projFolder) {
		wd, err := os.Getwd() // package dir
		if err != nil {
			return "", errors.WithMessage(err, "unable to determine current directory")
		}

		projFolder, err = filepath.Abs(filepath.Join(wd, projFolder))
		if err != nil {
			return "", errors.WithMessagef(err, "unable to determine project directory: %q", projFolder)
		}
	}

	return filepath.Join(projFolder, file), nil
}

// Load will load the configuration from the named config file,
// apply any overrides, and resolve relative directory locations.
func (f *Factory) Load(configFile string, config interface{}) (absConfigFile string, err error) {
	return f.LoadForHostName(configFile, "", config)
}

// LoadForHostName will load the configuration from the named config file for specified host name,
// apply any overrides, and resolve relative directory locations.
func (f *Factory) LoadForHostName(configFile, hostnameOverride string, config interface{}) (absConfigFile string, err error) {
	logger.KV(xlog.TRACE, "cfg", configFile, "hostname", hostnameOverride)

	configFile, baseDir, err := f.ResolveConfigFile(configFile)
	if err != nil {
		return "", err
	}

	logger.KV(xlog.DEBUG, "cfg", configFile, "baseDir", baseDir)

	err = f.load(configFile, hostnameOverride, baseDir, config)
	if err != nil {
		return "", err
	}

	environment := f.environment
	if environment != "" {
		// ignore error as Environment may not exist in the config
		_ = reflections.SetField(config, "Environment", environment)
	} else if value, err := reflections.GetField(config, "Environment"); err == nil {
		environment = value.(string)
	}

	variables := f.getVariableValues(environment)

	envName := f.envPrefix + "CONFIG_DIR"
	if variables[envName] == "" {
		variables[envName] = baseDir
		os.Setenv(envName, baseDir)
	}

	expander := &Expander{
		Variables:      variables,
		SecretProvider: f.secrets,
	}
	err = expander.ExpandAll(config)
	if err != nil {
		return configFile, err
	}

	return configFile, nil
}

// Hostmap provides overrides info
type Hostmap struct {
	// Override is a map of host name to file location
	Override map[string]string
}

// Load will attempt to load the configuration from the supplied filename.
// Overrides defined in the config file will be applied based on the hostname
// the hostname used is dervied from [in order]
//  1. the hostnameOverride parameter if not ""
//  2. the value of the Environment variable in envKeyName, if not ""
//  3. the OS supplied hostname
func (f *Factory) load(configFilename, hostnameOverride, baseDir string, config interface{}) error {
	var err error
	ops := []yamlcfg.YAMLOption{yamlcfg.File(configFilename)}

	// load hostmap schema
	if hmapraw, err := os.ReadFile(configFilename + ".hostmap"); err == nil {
		var hmap Hostmap
		err = yaml.Unmarshal(hmapraw, &hmap)
		if err != nil {
			return errors.Wrapf(err, "failed to load hostmap file")
		}

		hn := hostnameOverride
		if hn == "" {
			if f.envPrefix != "" {
				hn = os.Getenv(f.envPrefix + "HOSTNAME")
			}
			if hn == "" {
				hn, err = os.Hostname()
				if err != nil {
					logger.KV(xlog.ERROR, "reason", "hostname", "err", err)
				}
			}
		}

		if hn != "" && hmap.Override[hn] != "" {
			override := hmap.Override[hn]
			override, err = resolve.File(override, baseDir)
			if err != nil {
				return errors.WithMessagef(err, "failed to resolve file")
			}
			logger.KV(xlog.TRACE, "hostname", hn, "override", override)
			ops = append(ops, yamlcfg.File(override))
		}
	}

	if len(f.overrideCfg) > 0 {
		overrideCfg, _, err := f.ResolveConfigFile(f.overrideCfg)
		if err != nil {
			return err
		}
		logger.KV(xlog.TRACE, "override", overrideCfg)
		ops = append(ops, yamlcfg.File(overrideCfg))
	}

	provider, err := yamlcfg.NewYAML(ops...)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}

	err = provider.Get(yamlcfg.Root).Populate(config)
	if err != nil {
		return errors.Wrap(err, "failed to parse configuration")
	}

	return nil
}

func (f *Factory) getVariableValues(environment string) map[string]string {
	ret := map[string]string{
		"HOSTNAME":              f.nodeInfo.HostName(),
		"NODENAME":              f.nodeInfo.NodeName(),
		"LOCALIP":               f.nodeInfo.LocalIP(),
		"USER":                  f.userName(),
		"NORMALIZED_USER":       f.normalizedUserName(),
		"ENVIRONMENT":           environment,
		"ENVIRONMENT_UPPERCASE": strings.ToUpper(environment),
	}

	if len(f.envPrefix) > 0 {
		for k, v := range ret {
			ret[f.envPrefix+k] = v
		}

		for _, x := range os.Environ() {
			kvp := strings.SplitN(x, "=", 2)

			env, val := kvp[0], kvp[1]
			if strings.HasPrefix(env, f.envPrefix) {
				if _, ok := ret[env]; !ok {
					logger.KV(xlog.DEBUG, "set", env)
					ret[env] = val
				}
			}
		}
	}

	return ret
}

// ResolveConfigFile returns absolute path for the config file
func (f *Factory) ResolveConfigFile(configFile string) (absConfigFile, baseDir string, err error) {
	if configFile == "" {
		panic("config file not provided!")
		//configFile = ConfigFileName
	}

	if filepath.IsAbs(configFile) {
		// for absolute, use the folder containing the config file
		baseDir = filepath.Dir(configFile)
		absConfigFile = configFile
		return
	}

	for _, absDir := range f.searchDirs {
		absConfigFile, err = resolve.File(configFile, absDir)
		if err == nil && absConfigFile != "" {
			baseDir = absDir
			logger.KV(xlog.DEBUG, "resolved", absConfigFile)
			return
		}
	}

	err = errors.Errorf("file %q not found in [%s]", configFile, strings.Join(f.searchDirs, ","))
	return
}

func (f *Factory) userName() string {
	if f.user == nil {
		userName := userName()
		f.user = &userName
	}
	return *f.user
}

func (f *Factory) normalizedUserName() string {
	username := f.userName()
	return strings.ReplaceAll(username, ".", "")
}

func userName() string {
	u, err := user.Current()
	if err != nil {
		logger.Panicf("unable to determine current user: %v", err)
	}
	return u.Username
}
