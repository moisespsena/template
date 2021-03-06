package cache

import (
	"fmt"
	"sync"

	oscommon "github.com/moisespsena-go/os-common"

	"github.com/moisespsena-go/assetfs/assetfsapi"
	"github.com/moisespsena/template/text/template"
)

type ExecutorCache struct {
	Enable bool
	data   sync.Map
}

func NewCache() *ExecutorCache {
	return &ExecutorCache{}
}

var Cache = NewCache()

func (ec *ExecutorCache) Load(name string) *template.Executor {
	v, ok := ec.data.Load(name)
	if !ok {
		return nil
	}
	return v.(*template.Executor)
}

func (ec *ExecutorCache) LoadOrStore(name string, loader func(name string) (*template.Executor, error)) (*template.Executor, error) {
	if ec.Enable {
		v, ok := ec.data.Load(name)
		if !ok {
			v, err := loader(name)
			if err != nil {
				return nil, err
			}
			if v == nil {
				return nil, fmt.Errorf("nil value")
			}
			ec.data.Store(name, v)
			return v, nil
		}
		return v.(*template.Executor), nil
	}
	return loader(name)
}

func (ec *ExecutorCache) LoadOrStoreInfo(info assetfsapi.FileInfo, loader func(info assetfsapi.FileInfo) (*template.Executor, error)) (*template.Executor, error) {
	if ec.Enable {
		v, ok := ec.data.Load(info)
		if !ok {
			v, err := loader(info)
			if err != nil {
				return nil, err
			}
			if v == nil {
				return nil, fmt.Errorf("nil value")
			}
			ec.data.Store(info.RealPath(), v)
			return v, nil
		}
		return v.(*template.Executor), nil
	}
	return loader(info)
}

func (ec *ExecutorCache) LoadOrStoreNames(name string, loader func(name string) (*template.Executor, error), names ...string) (*template.Executor, error) {
	names = append([]string{name}, names...)
	for _, name := range names {
		v, ok := ec.data.Load(name)
		if ok && v != nil {
			return v.(*template.Executor), nil
		}

		t, err := loader(name)

		if err != nil {
			if oscommon.IsNotFound(err) {
				continue
			}
			return nil, err
		}

		if t != nil {
			if ec.Enable {
				ec.data.Store(name, t)
			}
			return t, nil
		}
	}
	return nil, oscommon.ErrNotFound(name)
}

func (ec *ExecutorCache) LoadOrStoreInfos(info assetfsapi.FileInfo, loader func(info assetfsapi.FileInfo) (*template.Executor, error), infos ...assetfsapi.FileInfo) (*template.Executor, error) {
	infos = append([]assetfsapi.FileInfo{info}, infos...)
	for _, info := range infos {
		v, ok := ec.data.Load(info.RealPath())
		if ok && v != nil {
			return v.(*template.Executor), nil
		}

		t, err := loader(info)

		if err != nil {
			return nil, err
		}

		if t != nil {
			if ec.Enable {
				ec.data.Store(info.RealPath(), t)
			}
			return t, nil
		}
	}
	return nil, nil
}
