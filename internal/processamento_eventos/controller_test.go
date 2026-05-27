package processamento_eventos

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MundoInvest/backend/internal/dominio"
	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
)

// MockWebhookService é um mock do WebhookService para testes
type MockWebhookService struct {
	processarWebhookFunc func(ctx context.Context, requisicao RequisicaoWebhook) error
	erro                 error
	clienteExistente     bool
}

func (m *MockWebhookService) ProcessarWebhook(ctx context.Context, requisicao RequisicaoWebhook) error {
	if m.processarWebhookFunc != nil {
		return m.processarWebhookFunc(ctx, requisicao)
	}
	if m.erro != nil {
		return m.erro
	}
	if !m.clienteExistente {
		return errors.NewNotFoundError("cliente", requisicao.EmailCliente)
	}
	return nil
}

func TestProcessarWebhookHandler_PrioridadeAlta(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{
		clienteExistente: true,
	}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_prioridade_alta",
		IdentificadorCard:   "card_123",
		EmailCliente:        "cliente.alto@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
	}

	if response["mensagem"] != "Webhook processado com sucesso" {
		t.Errorf("Mensagem esperada 'Webhook processado com sucesso', obtido '%v'", response["mensagem"])
	}
}

func TestProcessarWebhookHandler_PrioridadeNormal(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{
		clienteExistente: true,
	}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_prioridade_normal",
		IdentificadorCard:   "card_456",
		EmailCliente:        "cliente.normal@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
	}

	if response["mensagem"] != "Webhook processado com sucesso" {
		t.Errorf("Mensagem esperada 'Webhook processado com sucesso', obtido '%v'", response["mensagem"])
	}
}

func TestProcessarWebhookHandler_Idempotencia(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{
		clienteExistente: true,
	}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_idempotente",
		IdentificadorCard:   "card_789",
		EmailCliente:        "cliente.teste@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)

	// Primeira requisição
	req1 := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	// Act - Primeira requisição
	controller.ProcessarWebhookHandler(w1, req1)

	// Assert - Primeira requisição
	if w1.Code != http.StatusOK {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusOK, w1.Code)
	}

	// Segunda requisição com mesmo identificador
	req2 := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	// Act - Segunda requisição
	controller.ProcessarWebhookHandler(w2, req2)

	// Assert - Segunda requisição (idempotência)
	if w2.Code != http.StatusOK {
		t.Errorf("Status code esperado %d (idempotência), obtido %d", http.StatusOK, w2.Code)
	}
}

func TestProcessarWebhookHandler_ClienteNaoEncontrado(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{
		processarWebhookFunc: func(ctx context.Context, requisicao RequisicaoWebhook) error {
			return errors.NewNotFoundError("cliente", requisicao.EmailCliente)
		},
	}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_inexistente",
		IdentificadorCard:   "card_999",
		EmailCliente:        "nao.existe@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusNotFound, w.Code)
	}

	if !containsString(w.Body.String(), "Cliente não encontrado") {
		t.Errorf("Erro esperado contendo 'Cliente não encontrado', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_CampoObrigatorioEventoId(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorCard: "card_123",
		EmailCliente:      "cliente@example.com",
		DataEvento:        "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "identificadorEvento") {
		t.Errorf("Erro esperado contendo 'identificadorEvento', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_CampoObrigatorioCardId(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		EmailCliente:        "cliente@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "identificadorCard") {
		t.Errorf("Erro esperado contendo 'identificadorCard', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_CampoObrigatorioClienteEmail(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "card_123",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "emailCliente") {
		t.Errorf("Erro esperado contendo 'emailCliente', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_CampoObrigatorioTimestamp(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "card_123",
		EmailCliente:        "cliente@example.com",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "dataEvento") {
		t.Errorf("Erro esperado contendo 'dataEvento', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_EventIdVazio(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "",
		IdentificadorCard:   "card_123",
		EmailCliente:        "cliente@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "identificadorEvento") {
		t.Errorf("Erro esperado contendo 'identificadorEvento', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_CardIdVazio(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "",
		EmailCliente:        "cliente@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "identificadorCard") {
		t.Errorf("Erro esperado contendo 'identificadorCard', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_ClienteEmailVazio(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "card_123",
		EmailCliente:        "",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "emailCliente") {
		t.Errorf("Erro esperado contendo 'emailCliente', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_EmailInvalido(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "card_123",
		EmailCliente:        "email-invalido",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "email") {
		t.Errorf("Erro esperado contendo 'email', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_TimestampInvalido(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "card_123",
		EmailCliente:        "cliente@example.com",
		DataEvento:          "timestamp-invalido",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}

	if !containsString(w.Body.String(), "dataEvento") {
		t.Errorf("Erro esperado contendo 'dataEvento', obtido: %s", w.Body.String())
	}
}

func TestProcessarWebhookHandler_ErroBancoDados(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{
		erro: sql.ErrConnDone,
	}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_123",
		IdentificadorCard:   "card_123",
		EmailCliente:        "cliente@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusInternalServerError, w.Code)
	}
}

func TestProcessarWebhookHandler_MetodoNaoPermitido(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	req := httptest.NewRequest(http.MethodGet, "/webhooks/pipefy/card-updated", nil)
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestProcessarWebhookHandler_JSONInvalido(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{}
	controller := NovoEventoController(mockService)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBufferString("json-invalido"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusBadRequest, w.Code)
	}
}

// Teste de cálculo de prioridade no limite
func TestProcessarWebhookHandler_PrioridadeNoLimite(t *testing.T) {
	// Arrange
	mockService := &MockWebhookService{
		clienteExistente: true,
	}
	controller := NovoEventoController(mockService)

	requestBody := RequisicaoWebhook{
		IdentificadorEvento: "evt_limite",
		IdentificadorCard:   "card_limite",
		EmailCliente:        "cliente.limite@example.com",
		DataEvento:          "2026-05-27T10:00:00Z",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/pipefy/card-updated", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	controller.ProcessarWebhookHandler(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Status code esperado %d, obtido %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Erro ao decodificar resposta: %v", err)
	}

	if response["mensagem"] != "Webhook processado com sucesso" {
		t.Errorf("Mensagem esperada 'Webhook processado com sucesso', obtido '%v'", response["mensagem"])
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
func TestProcessarWebhookHandler_IntegracaoBancoDados(t *testing.T) {
	// Este teste requer banco de dados real e deve ser executado com tag integration
	// Exemplo: go test ./internal/processamento_eventos/... -v -tags=integration

	// Arranque: Iniciar banco de dados de teste
	// db := setupTestDB(t)
	// defer teardownTestDB(t, db)

	// controller := NovoEventoControllerComDB(db)

	// ... implementar teste de integração
	t.Skip("Teste de integração requer banco de dados real - executar com -tags=integration")
}

// Teste unitário do PrioridadeCalculator
func TestPrioridadeCalculator_CalcularNivelPrioridade(t *testing.T) {
	cfg := config.Load()
	calculator := dominio.NovaCalculadoraPrioridade(cfg)

	tests := []struct {
		name       string
		patrimonio float64
		esperado   string
	}{
		{"Patrimônio alto (acima do limite)", 250000.00, "prioridade_alta"},
		{"Patrimônio alto (no limite)", 200000.00, "prioridade_alta"},
		{"Patrimônio normal (abaixo do limite)", 150000.00, "prioridade_normal"},
		{"Patrimônio normal (zero)", 0.00, "prioridade_normal"},
		{"Patrimônio normal (pequeno)", 1000.00, "prioridade_normal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := calculator.CalcularNivelPrioridade(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("CalcularNivelPrioridade(%.2f) = %s, esperado %s", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

// Teste unitário do validador de webhook
func TestValidarRequisicaoWebhook(t *testing.T) {
	tests := []struct {
		name    string
		request RequisicaoWebhook
		erro    bool
	}{
		{"Payload válido", RequisicaoWebhook{
			IdentificadorEvento: "evt_123",
			IdentificadorCard:   "card_123",
			EmailCliente:        "cliente@example.com",
			DataEvento:          time.Now().Format(time.RFC3339),
		}, false},
		{"Evento ID vazio", RequisicaoWebhook{
			IdentificadorEvento: "",
			IdentificadorCard:   "card_123",
			EmailCliente:        "cliente@example.com",
			DataEvento:          time.Now().Format(time.RFC3339),
		}, true},
		{"Email inválido", RequisicaoWebhook{
			IdentificadorEvento: "evt_123",
			IdentificadorCard:   "card_123",
			EmailCliente:        "email-invalido",
			DataEvento:          time.Now().Format(time.RFC3339),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarRequisicaoWebhook(tt.request)
			if (err != nil) != tt.erro {
				t.Errorf("ValidarRequisicaoWebhook() erro = %v, esperado erro = %v", err, tt.erro)
			}
		})
	}
}
