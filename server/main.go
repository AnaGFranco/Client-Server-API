package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	serverPort = ":8080"
	apiTimeout = 200 * time.Millisecond
	dbTimeout  = 10 * time.Millisecond
)

type Quote struct {
	gorm.Model
	Bid string `json:"bid"`
}

func main() {
	db := initDB()

	http.HandleFunc("/quote", handleQuote(db))

	log.Printf("Servidor iniciado na porta %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

// Inicializando o banco de dados SQLite com GORM
func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("quotes.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	// Migrando o esquema do banco de dados
	err = db.AutoMigrate(&Quote{})
	if err != nil {
		log.Fatalf("Erro ao migrar o esquema do banco de dados: %v", err)
	}

	return db
}

func handleQuote(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Criando um contexto para chamada da API de cotação do dólar (timeout: 200ms)
		ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
		defer cancel()

		apiResp, err := getQuoteFromAPI(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao obter cotação da API: %v", err), http.StatusInternalServerError)
			return
		}

		// Criando um contexto para inserção no banco de dados (timeout: 10ms)
		dbCtx, dbCancel := context.WithTimeout(context.Background(), dbTimeout)
		defer dbCancel()

		if err := saveQuote(dbCtx, db, apiResp.USDBRL.Bid); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao salvar cotação no banco de dados: %v", err), http.StatusInternalServerError)
			return
		}

		// Retorna a cotação para o cliente
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apiResp.USDBRL)
	}
}

type APIResponse struct {
	USDBRL Quote `json:"USDBRL"`
}

// Chama a API de cotação do dólar para obter a cotação
func getQuoteFromAPI(ctx context.Context) (*APIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar requisição para a API: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Erro ao chamar a API: %v", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("Erro ao decodificar resposta da API: %v", err)
	}

	return &apiResp, nil
}

// Salva a cotação no banco de dados usando GORM
func saveQuote(ctx context.Context, db *gorm.DB, bid string) error {
	quote := Quote{
		Bid: bid,
	}
	result := db.WithContext(ctx).Create(&quote)
	if result.Error != nil {
		return fmt.Errorf("Erro ao salvar cotação no banco de dados: %v", result.Error)
	}
	return nil
}
