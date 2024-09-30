package nacos

import (
	"time"

	"github.com/diycoder/elf/config/encoder"
	"github.com/diycoder/elf/config/source"
	"github.com/diycoder/elf/plugin/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type watcher struct {
	configClient config_client.IConfigClient
	e            encoder.Encoder
	name         string
	watch        []*Watch
	ch           chan *source.ChangeSet
	exit         chan bool
}

func newConfigWatcher(cc config_client.IConfigClient, e encoder.Encoder, name string, watch []*Watch) (source.Watcher, error) {
	w := &watcher{
		e:            e,
		name:         name,
		configClient: cc,
		watch:        watch,
		ch:           make(chan *source.ChangeSet),
		exit:         make(chan bool),
	}

	for _, val := range watch {
		err := w.configClient.ListenConfig(vo.ConfigParam{
			DataId:   val.DataId,
			Group:    val.Group,
			OnChange: w.callback,
		})
		if err != nil {
			log.Errorf("nacos listen config group:%v, dataId:%v, err:%v", val.Group, val.DataId, err)
			return nil, err
		}
	}
	return w, nil
}

func (w *watcher) callback(namespace, group, dataId, data string) {
	var snapMap map[string]map[string]interface{}
	snap := cfg.Get().Bytes()
	err := w.e.Decode(snap, &snapMap)
	if err != nil {
		log.Errorf("nacos callback get snap:%v, err:%v", string(snap), err)
		return
	}
	value, ok := snapMap[group]
	if !ok {
		snapMap[group] = map[string]interface{}{
			dataId: data,
		}
	} else {
		value[dataId] = data
		snapMap[group] = value
	}

	encode, err := w.e.Encode(&snapMap)
	if err != nil {
		log.Errorf("nacos callback encode data:%+v, err:%v", snapMap, err)
		return
	}

	log.Infof("nacos callback namespaceId:%v, group:%v, dataId:%v, data:%v, result:%v",
		namespace, group, dataId, data, string(encode))

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    w.e.String(),
		Source:    w.name,
		Data:      encode,
	}
	cs.Checksum = cs.Sum()
	w.ch <- cs
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.ch:
		return cs, nil
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
		close(w.exit)
	}

	return nil
}
