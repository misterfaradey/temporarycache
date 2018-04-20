package temporarycache

import (
	"strconv"
	"testing"
	"time"
)

var (
	correctSizesTest   = []int{0, 1, 100}
	incorrectSizesTest = []int{-1, -100}
)

func TestInitCacheCorrect(t *testing.T) {

	for _, size := range correctSizesTest {
		cache := InitMemcache(size)
		if len(cache.array) != size {
			t.Errorf("TestInitCacheCorrect len %v size %v", len(cache.cache), size)
		}
	}
}

func TestInitCacheIncorrect(t *testing.T) {

	for _, size := range incorrectSizesTest {
		cache := InitMemcache(size)
		if len(cache.array) != 0 {
			t.Errorf("TestInitCacheIncorrect len %v size %v", len(cache.cache), size)
		}
	}
}

func TestWriteGetCache(t *testing.T) {
	size := 10
	mem := InitMemcache(size)
	value := "something"

	for i := 0; i < size; i++ {
		value := "something" + strconv.Itoa(i)
		mem.Write(i, value)
	}

	for i := 0; i < size; i++ {
		m, ok := mem.Get(i)
		if !ok {
			t.Errorf("TestWriteCache Want_ok: true, Have ok: %v, key: %v", ok, i)
			return
		}
		if m.(string) != value+strconv.Itoa(i) {
			t.Errorf("TestWriteCache Want: %s, Have: %v, key: %v", value+strconv.Itoa(i), m, i)
		}
	}
}
func TestCleanCache(t *testing.T) {
	size := 100
	mem := InitMemcache(size)

	for i := 0; i < size; i++ {
		value := "something" + strconv.Itoa(i)
		mem.Write(i, value)
	}

	go mem.Cleaner(time.Millisecond, time.Millisecond)

	time.Sleep(time.Second * 5)

	if len(mem.cache) > 0 {
		t.Errorf("TestCleanCache Not Working. LenCache:%v", len(mem.cache))
	}
}
