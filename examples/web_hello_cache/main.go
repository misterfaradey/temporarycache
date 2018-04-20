package main

import (
	"net/http"
	"time"

	"github.com/misterfaradey/temporarycache"
)

const (
	size         = 1000
	timeGap      = time.Millisecond
	liveDuration = time.Second
)

func helloHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	key := r.FormValue("key")

	value, ok := mem.Get(key)
	if ok {
		w.Write(value.([]byte))
		return
	}

	out := []byte("Hello World " + time.Now().String())
	w.Write(out)
	mem.Write(key, out)
}

var mem temporarycache.Mem

func main() {

	if size > 0 {
		mem = temporarycache.InitMemCache(size)
		go mem.Cleaner(timeGap, liveDuration)
	} else {
		mem = temporarycache.InitMemCache(0)
	}

	http.HandleFunc("/", helloHandler)

	http.ListenAndServe(":8080", nil)

}
