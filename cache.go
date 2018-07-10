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

//InitMemCache инициализация кэша. head- указывает на ячейку, в которую будет положен следующий элемент. tail- самый старый элемент
func InitMemCache(size int) Mem {
	if size <= 0 {
		return Mem{cache: make(cacheMap), end: 0, array: make([]timestampStruct, 0), head: 0, tail: 0, emptyTimestamp: timestampStruct{}}
	}
	return Mem{cache: make(cacheMap), end: size - 1, array: make([]timestampStruct, size), head: 0, tail: 0, emptyTimestamp: timestampStruct{}}
}

//Write записать элемент
func (m *Mem) Write(key interface{}, value interface{}) {
	if len(m.array) == 0 {
		return
	}

	//предотвращает повторную перезапись теми же данными
	m.RLock()
	val, ok := m.cache[key]
	m.RUnlock()
	if ok && val == value {
		return
	}

	//перезапись самого старого элемента в кэше
	m.Lock()
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
		m.Unlock()
		return
	}
	m.head++
	m.Unlock()
}

//Get получить элемент
func (m *Mem) Get(key interface{}) (value interface{}, ok bool) {
	m.RLock()
	value, ok = m.cache[key]
	m.RUnlock()
	return
}

//GetAll получить все элементы из кэша. Если ok, то элементы в кэше есть
func (m *Mem) GetAll() (mas []interface{}, ok bool) {
	m.RLock()

	switch {

	case m.head == m.tail:

		if m.array[m.head] != m.emptyTimestamp {
			mas = make([]interface{}, len(m.array))

			j := 0
			for _, i := range m.array[m.tail:] {
				value, _ := m.cache[i.key]
				mas[j] = value
				j++
			}
			for _, i := range m.array[:m.head] {
				value, _ := m.cache[i.key]
				mas[j] = value
				j++
			}

		} else {
			m.RUnlock()
			return nil, false
		}

	case m.head > m.tail:

		mas = make([]interface{}, m.head-m.tail)
		for j, i := range m.array[m.tail:m.head] {
			value, _ := m.cache[i.key]
			mas[j] = value
		}

	default:

		mas = make([]interface{}, m.tail-m.end+m.head)
		j := 0
		for _, i := range m.array[m.tail:] {
			value, _ := m.cache[i.key]
			mas[j] = value
			j++
		}
		for _, i := range m.array[:m.head] {
			value, _ := m.cache[i.key]
			mas[j] = value
			j++
		}
	}
	m.RUnlock()

	return mas, true
}

func (m *Mem) deleteOld(liveDuration time.Duration) {
	m.Lock()

	for {
		ar := m.array[m.tail]
		if ar == m.emptyTimestamp {
			m.Unlock()
			return
		}
		if time.Now().Sub(ar.time) < liveDuration {
			m.Unlock()
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
		m.deleteOld(liveDuration)
		time.Sleep(timeGap)
	}
}
