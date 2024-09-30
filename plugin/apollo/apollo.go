package apollo

import (
	"fmt"
	"strings"
	"time"

	apo "github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	cfg "github.com/diycoder/elf/config"
	"github.com/diycoder/elf/config/reader"
	"github.com/diycoder/elf/config/source"
	"github.com/diycoder/elf/plugin/log"
	"github.com/diycoder/elf/utils/convert"
)

var apolloConfig cfg.Config

type apolloSource struct {
	client        apo.Client
	namespaceName string
	opts          source.Options
}

func (a *apolloSource) String() string {
	return "apollo"
}

func (a *apolloSource) Read() (*source.ChangeSet, error) {
	data := map[string]interface{}{}
	split := strings.Split(a.namespaceName, ",")
	for _, namespace := range split {
		c := a.client.GetConfig(namespace).GetCache()
		values := map[string]interface{}{}
		c.Range(func(key interface{}, value interface{}) bool {
			values[convert.ToString(key)] = value
			return true
		})
		data[namespace] = values
	}

	b, err := a.opts.Encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    a.opts.Encoder.String(),
		Source:    a.String(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (a *apolloSource) Watch() (source.Watcher, error) {
	watcher, err := newWatcher(a.String(), a.opts.Encoder)
	a.client.AddChangeListener(watcher)
	return watcher, err
}

func (a *apolloSource) Write(cs *source.ChangeSet) error {
	return nil
}

func load(opts *Options) error {
	var err error
	apolloConfig, err = cfg.NewConfig()
	if err != nil {
		return err
	}
	if err := apolloConfig.Load(newSource(opts)); err != nil {
		return err
	}
	return nil
}

func newSource(opts *Options) source.Source {
	options := source.NewOptions()
	readyConfig := &config.AppConfig{
		IP:               opts.Address,
		AppID:            opts.AppID,
		Cluster:          opts.Cluster,
		IsBackupConfig:   opts.Backup,
		NamespaceName:    opts.Namespace,
		BackupConfigPath: opts.BackupPath,
	}
	client, err := apo.StartWithConfig(func() (*config.AppConfig, error) {
		return readyConfig, nil
	})
	if err != nil {
		log.Errorf("apollo init err:%v", err)
		return nil
	}
	return &apolloSource{
		client:        client,
		opts:          options,
		namespaceName: opts.Namespace,
	}
}

// Return config as raw json
func Bytes() []byte {
	return apolloConfig.Bytes()
}

// Return config as a map
func Map() map[string]interface{} {
	return apolloConfig.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return apolloConfig.Scan(v)
}

// Force a source changeset sync
func Sync() error {
	return apolloConfig.Sync()
}

// Get a value from the config
func Get(path ...string) reader.Value {
	return apolloConfig.Get(path...)
}
