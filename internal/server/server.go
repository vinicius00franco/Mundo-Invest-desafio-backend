package server

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/gestao_clientes"
	"github.com/MundoInvest/backend/internal/integracao_pipefy"
	"github.com/MundoInvest/backend/internal/processamento_eventos"
	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/database"
)

// Server representa o servidor HTTP
type Server struct {
	httpServer *http.Server
	config     *config.Config
	db         *sql.DB
}

// NewServer cria uma nova instância do servidor com dependency injection
func NewServer(cfg *config.Config) (*Server, error) {
	// Configurar conexão com banco de dados
	db, err := database.NovaConexao(cfg)
	if err != nil {
		return nil, err
	}

	// Criar repositórios
	clienteRepository := gestao_clientes.NovoClienteRepository(db, cfg)
	eventoRepository := processamento_eventos.NovoEventoRepository(db, cfg)

	// Criar serviços de domínio
	prioridadeCalculator := dominio.NovoPrioridadeCalculator(cfg)
	pipefyClient := integracao_pipefy.NovoPipefyGraphQLClient(cfg.Pipefy.APIToken, cfg.Pipefy.APIURL)
	eventDispatcher := dominio.NewInMemoryEventDispatcher()

	// Criar serviços de aplicação
	clienteService := gestao_clientes.NovoClienteService(
		clienteRepository,
		integracao_pipefy.NewPipefyIntegrationService(pipefyClient, cfg),
		cfg.Pipefy.PipeID,
		eventDispatcher,
		cfg,
	)
	webhookService := processamento_eventos.NovoWebhookService(
		eventoRepository,
		clienteRepository,
		prioridadeCalculator,
		pipefyClient,
		cfg,
	)

	// Criar controllers
	clienteController := gestao_clientes.NovoClienteController(clienteService)
	eventoController := processamento_eventos.NovoEventoController(webhookService)

	// Configurar router HTTP
	mux := http.NewServeMux()

	// Registrar rotas
	clienteController.RegistrarRotas(mux)
	eventoController.RegistrarRotas(mux)

	// Configurar servidor HTTP
	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      mux,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
		db:         db,
	}, nil
}

// Start inicia o servidor HTTP
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown encerra o servidor HTTP de forma graciosa
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Fechar conexão com banco de dados
	if err := s.db.Close(); err != nil {
		return err
	}

	return s.httpServer.Shutdown(ctx)
}
