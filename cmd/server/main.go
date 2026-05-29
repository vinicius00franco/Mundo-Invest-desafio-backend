package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MundoInvest/backend/internal/server"
	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Não foi possível carregar arquivo .env: %v", err)
	}

	// Inicializar logger
	logger.Init()

	// Carregar configuração
	cfg := config.Load()

	// Criar servidor com dependency injection
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Erro ao criar servidor: %v", err)
	}

	// Configurar graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	log.Printf("Servidor Mundo Invest iniciado na porta %s", cfg.HTTP.Port)
	log.Printf("Endpoint POST /clientes disponível")
	log.Printf("Endpoint POST /webhooks/pipefy/card-updated disponível")

	// Aguardar sinal de shutdown
	<-done
	log.Println("Recebido sinal de shutdown, encerrando servidor...")

	// Encerrar servidor de forma graciosa
	if err := srv.Shutdown(cfg.HTTP.IdleTimeout); err != nil {
		log.Printf("Erro ao encerrar servidor: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
