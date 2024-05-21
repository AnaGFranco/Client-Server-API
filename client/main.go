package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL = "http://localhost:8080/quote"
	timeout   = 300 * time.Millisecond
)

type Quote struct {
	Bid string `json:"bid"`
}

func main() {
	// Criando um contexto para a requisição ao servidor (timeout: 300ms)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	quote, err := getQuoteFromServer(ctx)
	if err != nil {
		log.Fatalf("Erro ao obter cotação do servidor: %v", err)
	}

	if err := saveQuoteToFile(quote); err != nil {
		log.Fatalf("Erro ao salvar cotação no arquivo: %v", err)
	}

	fmt.Println("Cotação salva em cotacao.txt")
}

// Faz a requisição ao servidor para obter a cotação
func getQuoteFromServer(ctx context.Context) (*Quote, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Resposta inesperada do servidor: %v", resp.Status)
	}

	var quote Quote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, fmt.Errorf("Erro ao decodificar resposta: %v", err)
	}

	return &quote, nil
}

// Salva a cotação em um arquivo
func saveQuoteToFile(quote *Quote) error {
	content := fmt.Sprintf("Dólar: %s", quote.Bid)
	if err := os.WriteFile("cotacao.txt", []byte(content), 0644); err != nil {
		return fmt.Errorf("Erro ao salvar arquivo: %v", err)
	}
	return nil
}
