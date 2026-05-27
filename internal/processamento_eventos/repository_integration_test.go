package processamento_eventos

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/MundoInvest/backend/internal/shared/config"
	_ "github.com/lib/pq"
)

// TestEventoRepositoryIntegracao testa o repository de eventos com banco real
func TestEventoRepositoryIntegracao(t *testing.T) {
	// Este teste requer banco de dados real
	// Executar com: go test ./internal/processamento_eventos/... -v -tags=integration

	t.Skip("Teste de integração requer banco de dados real - executar com -tags=integration")
}

// TestEventoRepositoryIntegracao_Salvar testa salvar evento no banco
func TestEventoRepositoryIntegracao_Salvar(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=mundo_invest sslmode=disable")
	if err != nil {
		t.Skipf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	repository := NovoEventoRepository(db, cfg)

	evento := Evento{
		IdentificadorEvento: "evt_test_123",
		IdentificadorCard:   "card_test_123",
		EmailCliente:        "teste.evento@example.com",
		TimestampEvento:     time.Now(),
		FoiProcessado:       true,
		DataCriacao:         time.Now(),
		DataAtualizacao:     nil, // Field is nullable
	}

	salvo, err := repository.Salvar(context.Background(), evento)
	if err != nil {
		t.Fatalf("Erro ao salvar evento: %v", err)
	}

	if salvo.IdentificadorInterno == 0 {
		t.Error("Identificador interno não foi gerado")
	}

	// Limpar
	db.Exec("DELETE FROM processamento_eventos.evento WHERE pev_eve_int = $1", salvo.IdentificadorInterno)
}

// TestEventoRepositoryIntegracao_VerificarFoiProcessado testa verificação de processamento
func TestEventoRepositoryIntegracao_VerificarFoiProcessado(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=mundo_invest sslmode=disable")
	if err != nil {
		t.Skipf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	repository := NovoEventoRepository(db, cfg)

	// Testar evento não existente
	foiProcessado, err := repository.VerificarFoiProcessado(context.Background(), "evt_inexistente")
	if err != nil {
		t.Fatalf("Erro ao verificar evento inexistente: %v", err)
	}

	if foiProcessado {
		t.Error("Evento inexistente não deve estar processado")
	}

	// Criar evento processado
	evento := Evento{
		IdentificadorEvento: "evt_test_456",
		IdentificadorCard:   "card_test_456",
		EmailCliente:        "teste.processado@example.com",
		TimestampEvento:     time.Now(),
		FoiProcessado:       true,
		DataCriacao:         time.Now(),
		DataAtualizacao:     nil, // Field is nullable
	}

	salvo, err := repository.Salvar(context.Background(), evento)
	if err != nil {
		t.Fatalf("Erro ao salvar evento: %v", err)
	}

	// Testar evento processado
	foiProcessado, err = repository.VerificarFoiProcessado(context.Background(), "evt_test_456")
	if err != nil {
		t.Fatalf("Erro ao verificar evento processado: %v", err)
	}

	if !foiProcessado {
		t.Error("Evento deve estar processado")
	}

	// Limpar
	db.Exec("DELETE FROM processamento_eventos.evento WHERE pev_eve_int = $1", salvo.IdentificadorInterno)
}

// TestEventoRepositoryIntegracao_BuscarPorIdentificadorCard testa busca por card
func TestEventoRepositoryIntegracao_BuscarPorIdentificadorCard(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=mundo_invest sslmode=disable")
	if err != nil {
		t.Skipf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	repository := NovoEventoRepository(db, cfg)

	// Limpar eventos de teste anteriores
	cardID := "card_test_789"
	db.Exec("DELETE FROM processamento_eventos.evento WHERE pev_eve_idc = $1", cardID)

	// Criar eventos de teste
	for i := 0; i < 3; i++ {
		evento := Evento{
			IdentificadorEvento: fmt.Sprintf("evt_test_%d", i),
			IdentificadorCard:   cardID,
			EmailCliente:        "teste.card@example.com",
			TimestampEvento:     time.Now(),
			FoiProcessado:       true,
			DataCriacao:         time.Now(),
			DataAtualizacao:     nil, // Field is nullable
		}

		_, err := repository.Salvar(context.Background(), evento)
		if err != nil {
			t.Fatalf("Erro ao salvar evento: %v", err)
		}
	}

	// Buscar eventos por card
	eventos, err := repository.BuscarPorIdentificadorCard(context.Background(), cardID)
	if err != nil {
		t.Fatalf("Erro ao buscar eventos por card: %v", err)
	}

	if len(eventos) != 3 {
		t.Errorf("Esperado 3 eventos, obtido %d", len(eventos))
	}
}
