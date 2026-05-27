package gestao_clientes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockClienteService é um mock do ClienteService para testes
type MockClienteService struct {
	criarClienteFunc func(request CriarClienteRequest) (*Cliente, error)
	erro             error
}

func (m *MockClienteService) CriarCliente(request CriarClienteRequest) (*Cliente, error) {
	if m.criarClienteFunc != nil {
		return m.criarClienteFunc(request)
	}
	if m.erro != nil {
		return nil, m.erro
	}
	return &Cliente{
		IdentificadorInterno: 1,
		IdentificadorExterno: "card_123",
		Nome:                 request.Nome,
		Email:                request.Email,
		ValorPatrimonio:      request.ValorPatrimonio,
		TipoSolicitacao:      request.TipoSolicitacao,
		Status:               "Aguardando Análise",
	}, nil
}

func TestCriarClienteHandler_PayloadValido(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
	}

	if response["mensagem"] != "Cliente criado com sucesso" {
		t.Errorf("Mensagem esperada 'Cliente criado com sucesso', obtido '%v'", response["mensagem"])
	}
}

func TestCriarClienteHandler_CampoObrigatorioNome(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "nome") {
		t.Errorf("Erro esperado contendo 'nome', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_CampoObrigatorioEmail(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "email") {
		t.Errorf("Erro esperado contendo 'email', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_CampoObrigatorioTipoSolicitacao(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "joao.silva@example.com",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "tipo_solicitacao") {
		t.Errorf("Erro esperado contendo 'tipo_solicitacao', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_CampoObrigatorioValorPatrimonio(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "valor_patrimonio") {
		t.Errorf("Erro esperado contendo 'valor_patrimonio', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_EmailInvalido(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "email-invalido",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "email") {
		t.Errorf("Erro esperado contendo 'email', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_PatrimonioNegativo(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: -1000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "valor_patrimonio") {
		t.Errorf("Erro esperado contendo 'valor_patrimonio', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_PatrimonioZero(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 0.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "valor_patrimonio") {
		t.Errorf("Erro esperado contendo 'valor_patrimonio', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_NomeVazio(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "",
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "nome") {
		t.Errorf("Erro esperado contendo 'nome', obtido: %s", w.Body.String())
	}
}

func TestCriarClienteHandler_ErroBancoDados(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{
		erro: sql.ErrConnDone,
	}
	controller := NovoClienteController(mockService)

	requestBody := CriarClienteRequest{
		Nome:            "João Silva",
		Email:           "joao.silva@example.com",
		TipoSolicitacao: "abertura_conta",
		ValorPatrimonio: 150000.00,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCriarClienteHandler_MetodoNaoPermitido(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/clientes", nil)
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestCriarClienteHandler_JSONInvalido(t *testing.T) {
	// Arrange
	mockService := &MockClienteService{}
	controller := NovoClienteController(mockService)

	req := httptest.NewRequest(http.MethodPost, "/clientes", bytes.NewBufferString("json-invalido"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.CriarClienteHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Teste de integração com banco de dados real
func TestCriarClienteHandler_IntegracaoBancoDados(t *testing.T) {
	// Este teste requer banco de dados real e deve ser executado com tag integration
	// Exemplo: go test ./internal/gestao_clientes/... -v -tags=integration

	// Arranque: Iniciar banco de dados de teste
	// db := setupTestDB(t)
	// defer teardownTestDB(t, db)

	// controller := NovoClienteControllerComDB(db, "test_pipe_id")

	// ... implementar teste de integração
	t.Skip("Teste de integração requer banco de dados real - executar com -tags=integration")
}
