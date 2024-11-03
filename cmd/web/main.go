package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"completed"`
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/", rootHandler)
	m.HandleFunc("/todos/", todosHandler)

	log.Println("Listening to :8000")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Panicln(err.Error())
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	serveTemplate("templates/index.html", nil, w)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {

	pageData := TodoPageData{
		PageTitle: "Todo's",
		Todos:     fetchTodo(),
	}

	serveTemplate("templates/todo-list.html", pageData, w)

}

func serveTemplate(name string, v any, w io.Writer) {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte("Unable to parse the template"))
		return
	}

	if err := tmpl.Execute(w, v); err != nil {
		log.Println(err.Error())
		w.Write([]byte("Unable to parse the template"))
		return

	}
}

func fetchTodo() []Todo {
	var todos []Todo
	response, err := http.Get("https://jsonplaceholder.typicode.com/todos/")
	if err != nil {
		log.Println(err.Error())
		return todos
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&todos); err != nil {
		log.Println(err.Error())
		return todos
	}

	return todos
}
