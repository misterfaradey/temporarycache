package temporarycache

import (
	"strconv"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	size := 4000000
	mem := InitMemcache(size)
	value := "something"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		mem.Write(i, value)
	}
}

func BenchmarkGet(b *testing.B) {
	size := 4000000
	mem := InitMemcache(size)

	for i := 0; i < size; i++ {
		value := "something" + strconv.Itoa(i)
		mem.Write(i, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t, ok := mem.Get(i)
		_, _ = t, ok
	}
}
