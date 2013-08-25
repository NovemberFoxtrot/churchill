package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"winston"
)

type Index struct {
	Data map[string][]*winston.Winston
}

func (i *Index) update(w *winston.Winston) {
	for _, gram := range w.Grams {
		if i.Data[gram] == nil {
			i.Data[gram] = make([]*winston.Winston, 0)
		}

		index.Data[gram] = append(index.Data[gram], w)
	}
}

func add(website string) {
	var w winston.Winston
	w.Location = website
	w.FetchUrl(website)
	w.CalcGrams()
	winstons = append(winstons, w)

	index.update(&w)
}

func AddHandler(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	winston.CheckError(err)

	fmt.Println(r.Form)

	go add(r.Form["website"][0])

	http.Redirect(rw, r, "/", http.StatusFound)
}

type tv struct {
	Location string
	GramsLen int
}

type stv struct {
	Location string
	Score    int
}

func SearchHandler(rw http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/search.html")
	winston.CheckError(err)

	stvs := make([]stv, 0)

	for k, v := range index.Data {
		if k == r.FormValue("query") {
			for _, w := range v {
				stvs = append(stvs, stv{w.Location, 1})
			}
		}
	}

	fmt.Println(index.Data, stvs)

	t.Execute(rw, stvs)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	winston.CheckError(err)

	tvs := make([]tv, 0)
	tvs = append(tvs, tv{"Index", len(index.Data)})

	for i := 0; i < len(winstons); i++ {
		tvs = append(tvs, tv{winstons[i].Location, len(winstons[i].Grams)})
	}

	t.Execute(w, tvs)
}

var winstons []winston.Winston
var index Index

func init() {
	winstons = make([]winston.Winston, 0)
	index.Data = make(map[string][]*winston.Winston)
}

func main() {
	wd, err := os.Getwd()

	winston.CheckError(err)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/search", SearchHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+`/public`))))

	err = http.ListenAndServe(":9090", nil)
	winston.CheckError(err)
}
