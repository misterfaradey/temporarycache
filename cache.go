package temporarycache

import (
	"sync"
	"time"
)

type cacheMap map[interface{}]interface{}

//Mem временный кэш на основе бинарного дерева и массива
type Mem struct {
	//key==timestampStruct.key
	cache cacheMap
	end   int
	array []timestampStruct
	head  int
	tail  int
	def   time.Time
	sync.RWMutex
	emptyTimestamp timestampStruct
}

type timestampStruct struct {
	time time.Time
	key  interface{}
}

//InitMemcache инициализация кэша
func InitMemcache(size int) Mem {
	if size <= 0 {
		return Mem{cache: make(cacheMap), end: 0, array: make([]timestampStruct, 0), head: 0, tail: 0, emptyTimestamp: timestampStruct{}}
	}
	return Mem{cache: make(cacheMap), end: size - 1, array: make([]timestampStruct, size), head: 0, tail: 0, emptyTimestamp: timestampStruct{}}
}

//Write записать элемент
func (m *Mem) Write(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()
	if m.end == 0 {
		return
	}

	//предотвращает повторную перезапись теми же данными
	if _, ok := m.cache[key]; ok {
		return
	}

	//перезапись самого старого элемента в кэше
	if m.tail == m.head {
		if m.array[m.tail] != m.emptyTimestamp {
			delete(m.cache, m.array[m.tail].key)
			if m.tail < m.end {
				m.tail++
			} else {
				m.tail = 0
			}
		}
	}
	m.array[m.head] = timestampStruct{key: key, time: time.Now()}
	m.cache[key] = value

	if m.head == m.end {
		m.head = 0
		return
	}
	m.head++
}

//Get получить элемент
func (m *Mem) Get(key interface{}) (value interface{}, ok bool) {
	m.RLock()
	defer m.RUnlock()

	value, ok = m.cache[key]
	return
}

func (m *Mem) del(liveDuration time.Duration) {
	m.Lock()
	defer m.Unlock()

	for {
		ar := m.array[m.tail]
		if ar == m.emptyTimestamp {
			return
		}
		if time.Now().Sub(ar.time) < liveDuration {
			return
		}
		delete(m.cache, ar.key)
		m.array[m.tail] = m.emptyTimestamp
		if m.tail == m.end {
			m.tail = 0
			continue
		}
		m.tail++
	}
}

//Cleaner : сборщик мусора. timeGap делать разумно маленьким, удаление старых элементов блокирует кэш
func (m *Mem) Cleaner(timeGap, liveDuration time.Duration) {
	for {
		m.del(liveDuration)
		time.Sleep(timeGap)
	}
}
