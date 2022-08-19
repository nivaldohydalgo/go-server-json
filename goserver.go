package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Mensagem struct {
	Id    int
	Texto string
	Data  string
	Hora  string
	Autor string
}

const fileJson string = "./mensagens.json"

func main() {
	// define as rotas
	r := mux.NewRouter()
	r.HandleFunc("/", handleGet).Methods("GET")
	r.HandleFunc("/", handlePost).Methods("POST")
	// levanta o servidor
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Printf("****** Servidor executando na porta 8080 ******")
	srv.ListenAndServe()
}

/*====== GET das mensagens gravadas ======*/
func handleGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Consultando as Mensagens gravadas...")
	// Carrega a lista de mensagens gravadas
	var mensagemList = carregaJson()
	// Retorna a lista de mensagens
	retornaLista(w, mensagemList)
}

/*====== POST de nova mensagem e retorna lista atualizada ======*/
func handlePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incluindo nova Mensagem...")
	// Carrega a lista de mensagens gravadas
	var mensagemList = carregaJson()
	// Verifica Id da última mensagem
	var ultimoPost = mensagemList[len(mensagemList)-1].Id
	decoder := json.NewDecoder(r.Body)
	var m Mensagem
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}
	// Gera nova mensagem e atualiza lista
	var mensagem Mensagem = Mensagem{ultimoPost + 1, m.Texto, dataAtual(), horaAtual(), m.Autor}
	mensagemList = append(mensagemList, mensagem)
	// Grava o JSon atualizado no disco
	gravaJson(mensagemList)
	// Retorna a lista de mensagens
	retornaLista(w, mensagemList)
}

/*====== Carrega a lista de mensagens no Response ======*/
func retornaLista(w http.ResponseWriter, mensagemList []Mensagem) {
	// Faz o Parse da lista de mensagens para formato JSon
	mensagemJson, err := json.Marshal(mensagemList)
	if err != nil {
		fmt.Fprintf(w, "Erro ao tratar a resposta!")
	}
	// Gera o Response da requisição
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(mensagemJson)
}

/*====== Recupera o JSon gravado no disco (se não existir gera um inicial) ======*/
func carregaJson() []Mensagem {
	_, err := os.Stat(fileJson)
	if os.IsNotExist(err) {
		// Se o arquivo "mensagens.json" não existir, grava no disco com uma mensagem inicial
		var mensagens = []Mensagem{}
		var mensagem Mensagem = Mensagem{1, "Benvindo ao ZapZap!!!", dataAtual(), horaAtual(), "ZapZap"}
		mensagens = append(mensagens, mensagem)
		gravaJson(mensagens)
		return mensagens
	} else {
		// Se o arquivo "mensagens.json" já existir, recupera as mensagens gravadas
		return recuperaJson()
	}

}

/*====== Grava o arquivo JSon no disco ======*/
func gravaJson(mensagens []Mensagem) {
	f, _ := os.Create(fileJson)
	defer f.Close()
	json, _ := json.MarshalIndent(mensagens, "", "\t")
	f.Write(json)
}

/*====== Le o arquivo JSon no disco ======*/
func recuperaJson() []Mensagem {
	arquivo, err := ioutil.ReadFile(fileJson)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var mensagens []Mensagem
	err = json.Unmarshal(arquivo, &mensagens)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return mensagens
}

/*====== Retorna data atual formatada ======*/
func dataAtual() string {
	currentTime := time.Now()
	return currentTime.Format("02/01/2006")
	// Formatação de Data:   02 é dia com dois digitos
	//                       01 é mês com dois digitos
	//                       2006 é ano com 4 digitos
}

/*====== Retorna hora atual formatada ======*/
func horaAtual() string {
	currentTime := time.Now()
	return currentTime.Format("15:04:05 ")
	// Formatação de Hora:   15 é hora com dois digitos
	//                       04 é minuto com dois digitos
	//                       05 é segundo com dois digitos
}
