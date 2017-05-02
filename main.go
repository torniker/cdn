package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleAll)
	http.ListenAndServe(":8000", mux)
}

func handleAll(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("mysql", "root@/cdn")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	stmtOut, err := db.Prepare("SELECT redirect FROM routes WHERE path = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(req.URL.Path)
	if err != nil {
		panic(err.Error())
	}
	if rows.Next() {
		var url string
		rows.Scan(&url)
		http.Redirect(w, req, url, 301)
		return
	}
	http.NotFound(w, req)
	return
}
