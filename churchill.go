package main

import (
	"html/template"
	"net/http"
	"os"
	"sir"
	"winston"
)

type tv struct {
	Location string
	GramsLen int
}

type stv struct {
	Location string
	Score    int
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	go winston.Add(r.FormValue("website"))
	http.Redirect(w, r, "/", http.StatusFound)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/search.html")
	sir.CheckError(err)

	stvs := make([]stv, 0)

	for index, location := range winston.Query(r.FormValue("query")) {
		stvs = append(stvs, stv{location, index})
	}

	t.Execute(w, stvs)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/index.html")
	sir.CheckError(err)

	tvs := make([]tv, 0)
	tvs = append(tvs, tv{"Index", winston.IndexDataLen()})
	/*
		for i := 0; i < len(winston.Documents); i++ {
			tvs = append(tvs, tv{winstons[i].Location, len(winstons[i].Grams)})
		}
	*/
	t.Execute(w, tvs)
}

func main() {
	wd, err := os.Getwd()

	sir.CheckError(err)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/search", SearchHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+`/public`))))

	err = http.ListenAndServe(":9090", nil)
	sir.CheckError(err)
}
