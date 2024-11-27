package main

import "sync"

type KV struct {
	mutex sync.RWMutex
	data  map[string][]byte
}

func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

func (kv *KV) Set(key, val []byte) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.data[string(key)] = []byte(val)
	return nil
}

func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	val, ok := kv.data[string(key)]
	return val, ok
}
