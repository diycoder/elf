package nacos

import (
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/diycoder/elf/config"
	"github.com/diycoder/elf/config/reader"
	"github.com/diycoder/elf/config/source"
	"github.com/diycoder/elf/plugin/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var cfg config.Config

type configSource struct {
	confClient config_client.IConfigClient
	opts       source.Options
	watch      []*Watch
}

func newSource(opts *Options) source.Source {
	s := &configSource{
		opts: source.Options{},
	}

	if err := sourceConfiguration(s, opts); err != nil {
		log.Errorf("nacos source config err:%v", err)
	}
	return s
}

func load(opts *Options) error {
	var err error
	cfg, err = config.NewConfig()
	if err != nil {
		return err
	}
	if err := cfg.Load(newSource(opts)); err != nil {
		return err
	}
	return nil
}

func sourceConfiguration(s *configSource, opts *Options) error {
	s.opts = source.NewOptions()
	clientConfig := constant.ClientConfig{
		CacheDir:            "./cache/nacos",
		LogDir:              "./log/nacos",
		NotLoadCacheAtStart: true,
	}
	serverConfigs := make([]constant.ServerConfig, 0)

	clientConfig.NamespaceId = opts.NamespaceId
	clientConfig.SecretKey = opts.SecretKey
	clientConfig.AccessKey = opts.AccessKey
	addrs := strings.Split(opts.Address, ",")
	watches, err := getWatchData(opts.WatchConfig)
	if err != nil {
		return err
	}
	s.watch = watches

	for _, addr := range addrs {
		// check we have a port
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}

		p, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			return err
		}

		serverConfigs = append(serverConfigs, constant.ServerConfig{
			// Scheme:      "go.micro",
			IpAddr: host,
			Port:   p,
		})
	}

	ic, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		log.Errorf("nacos init client err:%v", err)
		return err
	}

	s.confClient = ic
	return nil
}

type metadata struct {
	Watch []*Watch `yaml:"watch"`
}

func getWatchData(inputPath string) ([]*Watch, error) {
	// 获取绝对地址
	absFilename, err := filepath.Abs(inputPath)
	if err != nil {
		log.Errorf("file abs err:%v", err)
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(absFilename)
	if err != nil {
		log.Errorf("read store config err:%v", err)
		return nil, err
	}
	metadata := new(metadata)
	err = yaml.Unmarshal(yamlFile, metadata)
	if err != nil {
		log.Errorf("unmarshal store config err:%v", err)
		return nil, err
	}
	return metadata.Watch, nil
}

func (n *configSource) Read() (*source.ChangeSet, error) {
	snapMap := make(map[string]map[string]interface{})
	for _, val := range n.watch {
		str, err := n.confClient.GetConfig(vo.ConfigParam{
			DataId: val.DataId,
			Group:  val.Group,
		})
		if err != nil {
			log.Errorf("nacos get config group:%v, dataId:%v, err:%v", val.Group, val.DataId, err)
			return nil, err
		}
		value, ok := snapMap[val.Group]
		if !ok {
			values := make(map[string]interface{})
			values[val.DataId] = str
			snapMap[val.Group] = values
		} else {
			value[val.DataId] = map[string]interface{}{
				val.DataId: str,
			}
			snapMap[val.Group] = value
		}
	}

	encode, err := n.opts.Encoder.Encode(&snapMap)
	if err != nil {
		log.Errorf("nacos encode data:%+v, err:%v", snapMap, err)
		return nil, err
	}
	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    n.opts.Encoder.String(),
		Source:    n.String(),
		Data:      encode,
	}
	cs.Checksum = cs.Sum()
	return cs, nil
}

func (n *configSource) Write(set *source.ChangeSet) error {
	return nil
}

func (n *configSource) Watch() (source.Watcher, error) {
	return newConfigWatcher(n.confClient, n.opts.Encoder, n.String(), n.watch)
}

func (n *configSource) String() string {
	return "nacos"
}

// Return config as raw json
func Bytes() []byte {
	return cfg.Bytes()
}

// Return config as a map
func Map() map[string]interface{} {
	return cfg.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return cfg.Scan(v)
}

// Force a source changeset sync
func Sync() error {
	return cfg.Sync()
}

// Get a value from the config
func Get(path ...string) reader.Value {
	return cfg.Get(path...)
}
