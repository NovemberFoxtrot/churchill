package main

import (
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
