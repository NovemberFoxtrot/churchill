package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

func render(w http.ResponseWriter, data interface{}, filenames ...string) {
	t := template.New("layout")
	t.Delims("//", "//")

	t, err := t.ParseFiles(filenames...)
	sir.CheckError(err)

	t.Execute(w, data)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	stvs := make([]stv, 0)

	for index, location := range roosevelt.Query(r.FormValue("query")) {
		stvs = append(stvs, stv{location, index})
	}

	render(w, stvs, "templates/layout.html", "templates/search.html")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tvs := make([]tv, 0)
	tvs = append(tvs, tv{"Index", roosevelt.IndexDataLen()})

	for i := 0; i < len(winston.TheDocuments); i++ {
		tvs = append(tvs, tv{winston.TheDocuments[i].Location, len(winston.TheDocuments[i].Grams)})
	}

	render(w, tvs, "templates/layout.html", "templates/index.html")
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)

	if err != nil {
		s = ""
		return
	}

	s = string(b)
	return
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"success": true, "message": "Hello!"})
	return
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
