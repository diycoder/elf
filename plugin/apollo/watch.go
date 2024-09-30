package apollo

import (
	"time"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/diycoder/elf/config/encoder"
	"github.com/diycoder/elf/config/source"
	"github.com/diycoder/elf/plugin/log"
)

type watcher struct {
	e    encoder.Encoder
	name string
	ch   chan *source.ChangeSet
	exit chan bool
}

func (w *watcher) OnNewestChange(event *storage.FullChangeEvent) {
	// log.Infof(fmt.Sprintf("OnNewestChange namespace:%v, changes:%v", event.Namespace, event.Changes))
}

func (w *watcher) OnChange(changeEvent *storage.ChangeEvent) {
	var snapMap map[string]map[string]interface{}
	snap := apolloConfig.Get().Bytes()
	err := w.e.Decode(snap, &snapMap)
	if err != nil {
		log.Errorf("OnChange get snap:%v, err:%v", string(snap), err)
		return
	}

	for k, v := range changeEvent.Changes {
		switch v.ChangeType {
		case storage.ADDED, storage.MODIFIED:
			snapMap[changeEvent.Namespace][k] = v.NewValue
		case storage.DELETED:
			delete(snapMap[changeEvent.Namespace], k)
		default:
			// 异常case 不作处理
			log.Warn("apollo Onchange ChangeType error, changeEvent:%v", changeEvent)
		}
	}

	b, err := w.e.Encode(snapMap)
	if err != nil {
		log.Errorf("apollo Onchange encode data:%v, err:%v", snapMap, err.Error())
		return
	}
	log.WithField("changeEvent.Changes", changeEvent.Changes).Infof("apollo OnChange")

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    w.e.String(),
		Source:    w.name,
		Data:      b,
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
		// TODO apollo stop
	}
	return nil
}

func newWatcher(name string, e encoder.Encoder) (*watcher, error) {
	return &watcher{
		e:    e,
		name: name,
		exit: make(chan bool),
		ch:   make(chan *source.ChangeSet),
	}, nil
}
