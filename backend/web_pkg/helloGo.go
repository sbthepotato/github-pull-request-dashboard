package web_pkg

import (
	"net/http"
	"time"
)

func HelloGo(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "text")
	w.Write([]byte("Hello, from the golang backend " + time.Now().String()))
}
