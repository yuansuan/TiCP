package jsonmap

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

// Map 表示一个能自动执行持久化的字典数据结构
type Map struct {
	f  *os.File
	mu sync.RWMutex
	m  map[string]interface{}
}

// Get 从字典中获取一个键为k的值
func (m *Map) Get(k string) (v interface{}) {
	m.mu.RLock()
	v = m.m[k]
	m.mu.RUnlock()
	return
}

// Set 向字典中设置一个键值对
func (m *Map) Set(k string, v interface{}) {
	m.mu.Lock()
	m.m[k] = v
	m.saveLocked()
	m.mu.Unlock()
}

// Del 从字典中删除一个键值对
func (m *Map) Del(k string) {
	m.mu.Lock()
	delete(m.m, k)
	m.saveLocked()
	m.mu.Unlock()
}

// Close 关闭字典并拒绝后续的所有写入操作
func (m *Map) Close() (err error) {
	m.mu.Lock()
	err = m.f.Close()
	m.f = nil
	m.m = nil
	m.mu.Unlock()
	return
}

// Visit 遍历字典中的所有键值对
func (m *Map) Visit(visit func(k string, v interface{}) (exit bool)) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.m {
		if visit(k, v) {
			return
		}
	}
}

// saveLocked 将字典持久化保存到文件中
func (m *Map) saveLocked() {
	if m.f != nil {
		bs, err := json.Marshal(m.m)
		if err != nil {
			panic(err)
		}

		nw, err := m.f.WriteAt(bs, 0)
		if err != nil {
			panic(err)
		}

		if err = m.f.Truncate(int64(nw)); err != nil {
			panic(err)
		}

		if err = m.f.Sync(); err != nil {
			panic(err)
		}
	}
}

// New 创建一个持久化字典
func New(filename string) (*Map, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return Load(f)
}

// Load 从文件中读取内容并加载到字典中
func Load(f *os.File) (*Map, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if fi.Size() != 0 {
		if err = json.NewDecoder(f).Decode(&m); err != nil {
			return nil, err
		}
	}

	return &Map{f: f, m: m}, nil
}
