package main

import (
	"GOLEARN/pkg/models"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl.html",
		"./ui/html/base.layout.tmpl.html",
		"./ui/html/footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {

		app.errorLog.Println(err.Error())
		app.ServerError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.ServerError(w, err)
	}

	//snippets, err := app.snippets.Latest()
	//if err != nil {
	//	app.ServerError(w, err)
	//	return
	//}
	//
	//for _, snippet := range snippets {
	//	fmt.Fprintf(w, "%v\n", snippet)
	//}

	w.Write([]byte("Skbidi dop dop dop yes yes yes"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.ServerError(w, err)
		}
		return
	}

	data := &templateData{Snippet: snippet}

	files := []string{
		"./ui/html/show.page.tmpl.html",
		"./ui/html/base.layout.tmpl.html",
		"./ui/html/footer.partial.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.ServerError(w, err)
	}
	fmt.Fprintf(w, "Displaying a specific snippet %v", snippet)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, 405)
		return
	}

	title := "BEBRA"
	content := "SKIBIDI DOP DOP DOP"
	expires := "7"
	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
