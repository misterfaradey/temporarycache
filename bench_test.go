package temporarycache

import (
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	globmem = InitMemCache(size)
	value := "something"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globmem.Write(i, value)
	}
}

func BenchmarkGet(b *testing.B) {

	for i := 0; i < b.N; i++ {
		t, ok := globmem.Get(i & 0xFFFFF)
		_, _ = t, ok
	}
}

func BenchmarkDeleteOld(b *testing.B) {
	mem := InitMemCache(size)
	value := "something"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mem.Write(i, value)
		b.StartTimer()

		mem.deleteOld(0)
	}

}
