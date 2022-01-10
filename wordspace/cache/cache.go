package cache

import (
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/manbeardo/wordle_solver/wordspace/params"
)

type lockStruct struct {
	lock     *sync.Mutex
	sizeLock *sync.RWMutex
}

var lockMapLock = &sync.Mutex{}
var params2locks = map[params.Params]lockStruct{}

var wordSpaceCache *lru.ARCCache
var wordSpaceSizeCache = &sync.Map{}

func init() {
	var err error
	wordSpaceCache, err = lru.NewARC(10000)
	if err != nil {
		panic(err)
	}
}

func getLocks(p params.Params) lockStruct {
	lockMapLock.Lock()
	defer lockMapLock.Unlock()

	locks, ok := params2locks[p]
	if ok {
		return locks
	}
	params2locks[p] = lockStruct{
		lock:     &sync.Mutex{},
		sizeLock: &sync.RWMutex{},
	}
	return params2locks[p]
}

func Lock(p params.Params) {
	locks := getLocks(p)
	locks.lock.Lock()
}

func Unlock(p params.Params) {
	locks := getLocks(p)
	locks.lock.Unlock()
}

func Get(p params.Params) (interface{}, bool) {
	return wordSpaceCache.Get(p)
}

func Set(p params.Params, wordSpace interface{}) {
	wordSpaceCache.Add(p, wordSpace)
}

func LockSize(p params.Params) {
	locks := getLocks(p)
	locks.sizeLock.Lock()
}

func UnlockSize(p params.Params) {
	locks := getLocks(p)
	locks.sizeLock.Unlock()
}

func GetSize(p params.Params) (int, bool) {
	size, ok := wordSpaceSizeCache.Load(p)
	if !ok {
		return -1, false
	}
	return size.(int), ok
}

func SetSize(p params.Params, size int) {
	wordSpaceSizeCache.Store(p, size)
}
