package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type TempUrlMap struct {
	MailGRPC string
	Mail     string
	Log      string
	Auth     string
	Broker   string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	http.HandleFunc("/hidden", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This route is hidden and should not be exposed outside of the k8s kluster"))
	})

	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TempUrlMap{
		Broker: os.Getenv("BROKER_ENDPOINT"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
