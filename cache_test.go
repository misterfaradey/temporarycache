package temporarycache

import (
	"strconv"
	"testing"
	"time"
)

const size = 5000000

var (
	correctSizesTest   = []int{1, 100}
	incorrectSizesTest = []int{0, -1, -100}
	globmem            Mem
)

func init() {
	globmem = InitMemCache(size)

	for i := 0; i < size; i++ {
		value := "something" + strconv.Itoa(i)
		globmem.Write(i, value)
	}
}
func TestInitCacheCorrect(t *testing.T) {

	value := "something"
	for _, size := range correctSizesTest {
		mem := InitMemCache(size)
		if len(mem.array) != size {
			t.Errorf("TestInitCacheCorrect len %v size %v", len(mem.array), size)
		}
		for i := 0; i < size; i++ {
			value := "something" + strconv.Itoa(i)
			mem.Write(i, value)
		}

		for i := 0; i < size; i++ {
			m, ok := mem.Get(i)
			if !ok {
				t.Errorf("TestInitCacheCorrect Want_ok: true, Have ok: %v, key: %v", ok, i)
				return
			}
			if m.(string) != value+strconv.Itoa(i) {
				t.Errorf("TestInitCacheCorrect Want: %s, Have: %v, key: %v", value+strconv.Itoa(i), m, i)
			}
			value := "something+1"
			mem.Write(i+1, value)
		}

	}

}

func TestInitCacheIncorrect(t *testing.T) {

	for _, size := range incorrectSizesTest {
		mem := InitMemCache(size)
		if len(mem.array) != 0 {
			t.Errorf("TestInitCacheIncorrect len %v size %v", len(mem.array), size)
		}

		mem.Write(1, "something")

		for i := 0; i < size; i++ {
			_, ok := mem.Get(i)
			if ok {
				t.Errorf("TestInitCacheIncorrect Want_ok: false, Have ok: %v, key: %v", ok, i)
			}
		}
	}
}
func TestDeleteFiveMillionsValues(t *testing.T) {

	globmem.deleteOld(time.Microsecond)

	if len(globmem.cache) > 0 {
		t.Errorf("TestCleanCache Not Working. LenCache:%v", len(globmem.cache))
	}
}

func TestCleaner(t *testing.T) {

	go globmem.Cleaner(time.Microsecond, time.Microsecond)
	time.Sleep(time.Second)
}
