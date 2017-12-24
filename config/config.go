package asticonfig

import (
	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

// New builds a new configuration based on a ptr to the global configuration, the path to the optional toml local
// configuration and a ptr to the flag configuration
func New(global interface{}, localPath string, flag interface{}) (_ interface{}, err error) {
	// Local config
	if localPath != "" {
		if _, err = toml.DecodeFile(localPath, global); err != nil {
			err = errors.Wrapf(err, "toml decoding %s failed", localPath)
			return
		}
	}

	// Merge configs
	if err = mergo.Merge(flag, global); err != nil {
		err = errors.Wrap(err, "merging configs failed")
		return
	}
	return flag, nil
}
