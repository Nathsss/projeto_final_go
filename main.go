
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

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

func index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "templates/index.html")
}

func salvar(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Método incorreto", http.StatusMethodNotAllowed)
		return
	}

	nome := req.FormValue("nome")
	precoStr := req.FormValue("preco")
	
	preco, _ := strconv.ParseUint(precoStr, 10, 64)

	novoDado := Dados{
		Nome:  nome,
		Preco: uint(preco),
	}

	if err := db.Create(&novoDado).Error; err != nil {
		http.Error(w, "Erro ao salvar", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Salvo com sucesso! <a href='/'>Voltar</a>")
}

func deletar(w http.ResponseWriter, req *http.Request) {
	idStr := strings.TrimPrefix(req.URL.Path, "/deletar/")
	
	if err := db.Delete(&Dados{}, idStr).Error; err != nil {
		http.Error(w, "Erro ao deletar", http.StatusInternalServerError)
		return
	}
	
	fmt.Fprintf(w, "Deletado! <a href='/'>Voltar</a>")
}

func editar(w http.ResponseWriter, req *http.Request) {
	idStr := strings.TrimPrefix(req.URL.Path, "/editar/")
	
	var produto Dados
	result := db.First(&produto, idStr)
	
	if result.Error != nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/update.html")
	if err != nil {
		http.Error(w, "Erro no template", http.StatusInternalServerError)
		return
	}
	
	tmpl.Execute(w, produto)
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("inventario.db"), &gorm.Config{})
	if err != nil {
		panic("falha ao conectar no banco de dados")
	}
	
	if err := db.AutoMigrate(&Dados{}); err != nil {
		panic("Erro no automigrate")
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/salvar", salvar)
	http.HandleFunc("/deletar/", deletar)
	http.HandleFunc("/editar/", editar)

	fmt.Println("Servidor rodando na porta :8090")
	http.ListenAndServe(":8090", nil)
}