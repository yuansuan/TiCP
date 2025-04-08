package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

func cleanData(data interface{}) interface{} {
	typ := reflect.ValueOf(data).Elem().Type()
	return reflect.New(typ).Interface()
}

// RegisterRemoteConfigDefault : Use default etcd
func (mw *Middleware) RegisterRemoteConfigDefault(ctx context.Context, mut *sync.Mutex, forceFrist bool, path string, data interface{}, f func(interface{})) error {
	return mw.RegisterRemoteConfig(ctx, mw.DefaultEtcd(), forceFrist, mut, path, data, f)
}

func writeData(mut *sync.Mutex, path string, value []byte, data interface{}, f func(interface{})) error {
	data = cleanData(data)

	logger := logging.Default()
	d := yaml.NewDecoder(bytes.NewReader(value))
	mut.Lock()
	// logger.Infof("decode %p", data)
	// logger.Infof("value %v", string(value))
	err := d.Decode(data)
	mut.Unlock()
	if err != nil {
		logger.Warnf("err in read remote config %v, %v", path, err)
	}
	if err == nil {
		if f != nil {
			f(data)
		}
		er := ioutil.WriteFile(filepath.Join(config.ConfigDir, path), value, os.ModePerm)
		if er != nil {
			logger.Warnf("err in read remote config %v, %v", path, er)
		}
	}
	return err
}

// RegisterRemoteConfigV2 RegisterRemoteConfig
// f will be called when config change, it's your duty to make sure data used in f should be synced
// f may be called multiple times
func (mw *Middleware) RegisterRemoteConfig(ctx context.Context, client *clientv3.Client, forceFirst bool, mut *sync.Mutex, path string, data interface{}, f func(interface{})) error {
	if client == nil {
		panic("client should not be nil")
	}
	logger := logging.Default()

	// read local file first
	previous, err := ioutil.ReadFile(filepath.Join(config.ConfigDir, path))
	if err == nil {
		d := yaml.NewDecoder(bytes.NewReader(previous))
		if env.Env.LogLevel < env.LevelWarn {
			d.SetStrict(true)
		}
		mut.Lock()
		util.ChkErr(d.Decode(data))
		mut.Unlock()
		if f != nil {
			f(data)
		}
	} else {
		err = nil
	}

	ch := make(chan error, 1)
	go func() {
		err := os.MkdirAll(filepath.Dir(filepath.Join(config.ConfigDir, path)), os.ModePerm)
		if err != nil {
			logger.Warnf("err in read remote config %v, %v", path, err)
			err = nil
		}
		val, err := client.Get(ctx, "/config/"+path)
		logger.Infof("fetch remote %v", "/config/"+path)
		if err != nil || len(val.Kvs) == 0 {
			logger.Warnf("err in read remote config %v, %v", path, err)
			if forceFirst {
				ch <- fmt.Errorf("err in read remote config %v, %v", path, err)
			}
			err = nil
		} else {
			err = writeData(mut, path, val.Kvs[0].Value, data, f)
			if err == io.EOF {
				err = nil
			}
			if forceFirst {
				ch <- err
			}
			err = nil
		}
		if f == nil {
			return
		}

		rev := val.Header.Revision

		if mw.etcd.polling {
			for {
				select {
				case <-time.After(time.Second * 20):
					val, err := client.Get(ctx, "/config/"+path)
					if err != nil || len(val.Kvs) == 0 {
						logger.Warnf("err in read remote config %v, %v", path, err)
						err = nil
					} else {
						if rev != val.Header.Revision {
							_ = writeData(mut, path, val.Kvs[0].Value, data, f)
							rev = val.Header.Revision
						}
					}
				case <-ctx.Done():
					break
				}
			}
		} else {
			for {
				ch := client.Watch(ctx, "/config/"+path, clientv3.WithRev(rev+1))
				select {
				case v := <-ch:
					if v.Err() != nil {
						logger.Warnf("err in watch remote config %v, %v", path, v.Err())
						val, err := client.Get(ctx, "/config/"+path)
						if err != nil || len(val.Kvs) == 0 {
							logger.Warnf("err in read remote config %v, %v", path, err)
							err = nil
						} else {
							rev = val.Header.Revision
							err = writeData(mut, path, val.Kvs[0].Value, data, f)
							err = nil
						}
						continue
					}
					for _, evt := range v.Events {
						if evt == nil {
							continue
						}
						logger.Infof("%v, %v, %v, %v", string(evt.Kv.Key), evt.Type, evt.PrevKv, evt.Kv.ModRevision)
						rev = evt.Kv.ModRevision
						_ = writeData(mut, path, evt.Kv.Value, data, f)
					}
				case <-ctx.Done():
					break
				}
			}
		}
	}()
	if forceFirst {
		err = <-ch
		close(ch)
	}
	return err
}
