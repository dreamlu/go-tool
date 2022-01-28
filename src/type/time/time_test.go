package time

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestTime(t *testing.T) {

	ti := time.Now()
	t.Log(ti)
	te := CTime(time.Now()).String()
	t.Log(te)
	var tt CTime
	t.Log(tt.MarshalJSON())
	t.Log(tt.String())
	_ = tt.UnmarshalJSON([]byte(te))
	t.Log(tt)
}

// test time Marshal
func TestWebTime(t *testing.T) {
	http.HandleFunc("/", sayhelloName)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func TestCSTime(t *testing.T) {
	ti := time.Now()
	t.Log(ti)
	te := CSTime(time.Now()).String()
	t.Log(te)
	var tt CSTime
	t.Log(tt.MarshalJSON())
	t.Log(tt.String())
	_ = tt.UnmarshalJSON([]byte(te))
	t.Log(tt)
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w,  CTime(time.Now())) // 这个写入到 w 的是输出到客户端的
	b, err := json.Marshal(CSTime(time.Now()))
	log.Print(err)
	w.Write(b)
}
