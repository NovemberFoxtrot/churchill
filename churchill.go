package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nimitz"
	"os"
	"roosevelt"
	"sir"
	"winston"
)

type tv struct {
	Location string
	GramsLen int
}

type stv struct {
	Result roosevelt.QueryResult
	Score  int
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	go roosevelt.Add(r.FormValue("website"))
	http.Redirect(w, r, "/", http.StatusFound)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	stvs := make([]stv, 0)

	for index, location := range roosevelt.Query(r.FormValue("query")) {
		stvs = append(stvs, stv{location, index})
	}

	nimitz.ThePool.Pools["search"].Execute(w, stvs)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tvs := make([]tv, 0)
	tvs = append(tvs, tv{"Index", roosevelt.IndexDataLen()})

	for i := 0; i < len(winston.TheDocuments); i++ {
		tvs = append(tvs, tv{winston.TheDocuments[i].Location, len(winston.TheDocuments[i].Grams)})
	}

	nimitz.ThePool.Pools["index"].Execute(w, tvs)
}

type Response map[string]interface{}

func (r Response) String() string {
	b, err := json.Marshal(r)

	sir.CheckError(err)

	return string(b)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true, "message": "Hello!"})
	return
}

func init() {
	nimitz.ThePool.Fill("index", "templates/layout.html", "templates/index.html")
	nimitz.ThePool.Fill("search", "templates/layout.html", "templates/search.html")
}

func main() {
	wd, err := os.Getwd()
	sir.CheckError(err)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/test", TestHandler)
	http.HandleFunc("/search", SearchHandler)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+`/public`))))

	err = http.ListenAndServe(":9090", nil)
	sir.CheckError(err)
}
