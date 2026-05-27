package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MundoInvest/backend/internal/gestao_clientes"
	"github.com/MundoInvest/backend/internal/processamento_eventos"
	"github.com/MundoInvest/backend/internal/shared/database"
	"github.com/MundoInvest/backend/internal/shared/logger"
)

func main() {
	// Inicializar logger
	logger.Init()

	// Configurar conexão com banco de dados
	db, err := database.NovaConexao()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Obter PIPE_ID das variáveis de ambiente
	pipeID := os.Getenv("PIPEFY_PIPE_ID")
	if pipeID == "" {
		pipeID = "DEFAULT_PIPE_ID" // Valor padrão para desenvolvimento
		log.Println("AVISO: PIPEFY_PIPE_ID não definido, usando valor padrão")
	}

	// Criar controller de clientes
	clienteController := gestao_clientes.NovoClienteControllerComDB(db, pipeID)

	// Criar controller de webhooks
	webhookController := processamento_eventos.NovoWebhookControllerComDB(db)

	// Configurar router HTTP
	mux := http.NewServeMux()

	// Registrar rotas
	clienteController.RegistrarRotas(mux)
	webhookController.RegistrarRotas(mux)

	// Configurar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Mundo Invest iniciado na porta %s", port)
	log.Printf("Endpoint POST /clientes disponível")
	log.Printf("Endpoint POST /webhooks/pipefy/card-updated disponível")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
