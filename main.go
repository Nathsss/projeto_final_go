
package main

import (
    "fmt"
    "net/http"

	"gorm.io/driver/sqlite"
  	"gorm.io/gorm"
)

type Dados struct {
  gorm.Model
  Nome string
  Preco uint
}

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func main() {
	 
	db, err := gorm.Open(sqlite.Open("inventario.db"), &gorm.Config{})
	if err != nil {
		panic(err) 
	}


	if err := db.AutoMigrate(&Dados{}); err != nil {
		panic("Erro no automigrate")
	}

    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)

    http.ListenAndServe(":8090", nil)
}