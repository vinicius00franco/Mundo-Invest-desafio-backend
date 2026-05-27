package gestao_clientes

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/MundoInvest/backend/internal/shared/config"
	_ "github.com/lib/pq"
)

// TestRepositoryIntegracao testa o repository com banco de dados real
func TestRepositoryIntegracao(t *testing.T) {
	// Este teste requer banco de dados real
	// Executar com: go test ./internal/gestao_clientes/... -v -tags=integration

	t.Skip("Teste de integração requer banco de dados real - executar com -tags=integration")
}

// TestRepositoryIntegracao_Salvar testa salvar cliente no banco
func TestRepositoryIntegracao_Salvar(t *testing.T) {
	// Configurar conexão de teste
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=mundo_invest sslmode=disable")
	if err != nil {
		t.Skipf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	repository := NovoClienteRepository(db, cfg)

	// Criar cliente de teste
	cliente := Cliente{
		IdentificadorExterno: "card_test_123",
		Nome:                 "Cliente Teste Integração",
		Email:                "teste.integracao@example.com",
		ValorPatrimonio:      100000.00,
		TipoSolicitacao:      "teste",
		Status:               "Aguardando Análise",
		NivelPrioridade:      "",
		DataCriacao:          time.Now(),
		DataAtualizacao:      time.Now(),
	}

	// Salvar cliente
	salvo, err := repository.Salvar(context.Background(), cliente)
	if err != nil {
		t.Fatalf("Erro ao salvar cliente: %v", err)
	}

	// Verificar se foi salvo
	if salvo.IdentificadorInterno == 0 {
		t.Error("Identificador interno não foi gerado")
	}

	// Limpar
	db.Exec("DELETE FROM gestao_clientes.cliente WHERE gcl_cli_int = $1", salvo.IdentificadorInterno)
}

// TestRepositoryIntegracao_BuscarPorEmail testa busca por email
func TestRepositoryIntegracao_BuscarPorEmail(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=mundo_invest sslmode=disable")
	if err != nil {
		t.Skipf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	repository := NovoClienteRepository(db, cfg)

	// Criar cliente de teste
	cliente := Cliente{
		IdentificadorExterno: "card_test_456",
		Nome:                 "Cliente Teste Busca",
		Email:                "teste.busca@example.com",
		ValorPatrimonio:      50000.00,
		TipoSolicitacao:      "teste",
		Status:               "Aguardando Análise",
		NivelPrioridade:      "",
		DataCriacao:          time.Now(),
		DataAtualizacao:      time.Now(),
	}

	salvo, err := repository.Salvar(context.Background(), cliente)
	if err != nil {
		t.Fatalf("Erro ao salvar cliente: %v", err)
	}

	// Buscar por email
	encontrado, err := repository.BuscarPorEmail(context.Background(), "teste.busca@example.com")
	if err != nil {
		t.Fatalf("Erro ao buscar por email: %v", err)
	}

	if encontrado.Email != "teste.busca@example.com" {
		t.Errorf("Email esperado 'teste.busca@example.com', obtido '%s'", encontrado.Email)
	}

	// Limpar
	db.Exec("DELETE FROM gestao_clientes.cliente WHERE gcl_cli_int = $1", salvo.IdentificadorInterno)
}

// TestRepositoryIntegracao_Atualizar testa atualização de cliente
func TestRepositoryIntegracao_Atualizar(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5434 user=postgres password=postgres dbname=mundo_invest sslmode=disable")
	if err != nil {
		t.Skipf("Não foi possível conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	cfg := config.Load()
	repository := NovoClienteRepository(db, cfg)

	// Criar cliente de teste
	cliente := Cliente{
		IdentificadorExterno: "card_test_789",
		Nome:                 "Cliente Teste Atualização",
		Email:                "teste.atualizacao@example.com",
		ValorPatrimonio:      75000.00,
		TipoSolicitacao:      "teste",
		Status:               "Aguardando Análise",
		NivelPrioridade:      "",
		DataCriacao:          time.Now(),
		DataAtualizacao:      time.Now(),
	}

	salvo, err := repository.Salvar(context.Background(), cliente)
	if err != nil {
		t.Fatalf("Erro ao salvar cliente: %v", err)
	}

	// Atualizar cliente
	salvo.Nome = "Cliente Atualizado"
	salvo.Status = "Processado"
	salvo.NivelPrioridade = "prioridade_normal"

	err = repository.Atualizar(context.Background(), *salvo)
	if err != nil {
		t.Fatalf("Erro ao atualizar cliente: %v", err)
	}

	// Verificar atualização
	atualizado, err := repository.BuscarPorIdentificadorInterno(context.Background(), salvo.IdentificadorInterno)
	if err != nil {
		t.Fatalf("Erro ao buscar cliente atualizado: %v", err)
	}

	if atualizado.Nome != "Cliente Atualizado" {
		t.Errorf("Nome esperado 'Cliente Atualizado', obtido '%s'", atualizado.Nome)
	}

	if atualizado.Status != "Processado" {
		t.Errorf("Status esperado 'Processado', obtido '%s'", atualizado.Status)
	}

	// Limpar
	db.Exec("DELETE FROM gestao_clientes.cliente WHERE gcl_cli_int = $1", salvo.IdentificadorInterno)
}
