package sql_console

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

// TODO: s/panic(/log.Println( or similar
// TODO: should this be an http.Handler?

//go:embed sql_console.html
var templateFile []byte

type Handler struct {
	DB *sql.DB
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("console").Parse(string(templateFile))
	if err != nil {
		panic(err)
	}

	switch r.Method {
	case "GET":
		tmplData := map[string]interface{}{
			// "Tables": tables,
		}

		if err = tmpl.Execute(w, tmplData); err != nil {
			panic(err)
		}
	case "POST":
		if r.FormValue("query") == "" {
			log.Println(r.PostForm)
			http.Error(w, "Expected value for query", http.StatusBadRequest)
			return
		}
		rows, err := h.DB.Query(r.FormValue("query"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results := [][]interface{}{}
		columns, err := rows.Columns()
		for rows.Next() {
			row := make([]interface{}, len(columns))
			for i := range columns {
				row[i] = new(interface{})
			}
			rows.Scan(row...)
			results = append(results, row)
		}
		b, err := json.Marshal(struct {
			Columns []string
			Rows    interface{}
		}{
			Columns: columns,
			Rows:    results,
		})
		if err != nil {
			panic(err)
		}
		w.Write(b)
	default:
		http.Error(w, "Expected GET or POST", http.StatusMethodNotAllowed)
	}
}
