package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"winston"
)

func add(website string) {
	var w winston.Winston
	w.Location = website
	w.FetchUrl(website)
	w.CalcGrams()
	fmt.Println(len(w.Text), len(w.Grams), len(w.Freq))
	winstons = append(winstons, w)
}

func AddHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	err := r.ParseForm()
	winston.CheckError(err)

	fmt.Println(r.Form)

	go add(r.Form["website"][0])
}

type tv struct {
	Location string
	GramsLen int
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.template", "templates/index.template")
	winston.CheckError(err)

	tvs := make([]tv, 0)

	for i := 0; i < len(winstons); i++ {
		tvs = append(tvs, tv{winstons[i].Location, len(winstons[i].Grams)})
	}

	t.Execute(w, tvs)
}

var winstons []winston.Winston

func init() {
	winstons = make([]winston.Winston, 0)
}

func main() {
	wd, err := os.Getwd()

	winston.CheckError(err)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/add", AddHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(wd+`/public`))))

	err = http.ListenAndServe(":9090", nil)
	winston.CheckError(err)
}
