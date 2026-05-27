package dominio

import (
	"fmt"
	"testing"

	"github.com/MundoInvest/backend/internal/shared/config"
)

func setupTestConfig() *config.Config {
	return &config.Config{
		Business: config.BusinessConfig{
			LimitePrioridadeAlta: 200000.00,
		},
	}
}

func TestCalculadoraPrioridade_CalcularNivelPrioridade(t *testing.T) {
	cfg := setupTestConfig()
	calculadora := NovaCalculadoraPrioridade(cfg)

	tests := []struct {
		name       string
		patrimonio float64
		esperado   string
	}{
		{"Patrimônio alto (acima do limite)", 250000.00, PrioridadeAlta},
		{"Patrimônio alto (no limite)", 200000.00, PrioridadeAlta},
		{"Patrimônio normal (abaixo do limite)", 150000.00, PrioridadeNormal},
		{"Patrimônio normal (zero)", 0.00, PrioridadeNormal},
		{"Patrimônio normal (pequeno)", 1000.00, PrioridadeNormal},
		{"Patrimônio muito alto", 1000000.00, PrioridadeAlta},
		{"Patrimônio limite - 1 centavo", 199999.99, PrioridadeNormal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := calculadora.CalcularNivelPrioridade(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("CalcularNivelPrioridade(%.2f) = %s, esperado %s", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

func TestCalculadoraPrioridade_CalcularNivelPrioridadeComDetalhes(t *testing.T) {
	cfg := setupTestConfig()
	calculadora := NovaCalculadoraPrioridade(cfg)

	// Type assert to concrete type to test implementation details
	concreteCalc, ok := calculadora.(*calculadoraPrioridade)
	if !ok {
		t.Fatal("Não foi possível fazer type assert para calculadoraPrioridade")
	}

	resultado, detalhes := concreteCalc.CalcularNivelPrioridadeComDetalhes(250000.00)
	if resultado != PrioridadeAlta {
		t.Errorf("Resultado esperado %s, obtido %s", PrioridadeAlta, resultado)
	}
	if detalhes == "" {
		t.Error("Detalhes não deveriam ser vazios")
	}
}

func TestCalculadoraPrioridade_EhPrioridadeAlta(t *testing.T) {
	cfg := setupTestConfig()
	calculadora := NovaCalculadoraPrioridade(cfg)

	tests := []struct {
		patrimonio float64
		esperado   bool
	}{
		{250000.00, true},
		{200000.00, true},
		{150000.00, false},
		{0.00, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Patrimônio %.2f", tt.patrimonio), func(t *testing.T) {
			// Type assert to concrete type to test implementation details
			concreteCalc, ok := calculadora.(*calculadoraPrioridade)
			if !ok {
				t.Fatal("Não foi possível fazer type assert para calculadoraPrioridade")
			}
			resultado := concreteCalc.EhPrioridadeAlta(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("EhPrioridadeAlta(%.2f) = %v, esperado %v", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

func TestCalculadoraPrioridade_EhPrioridadeNormal(t *testing.T) {
	cfg := setupTestConfig()
	calculadora := NovaCalculadoraPrioridade(cfg)

	tests := []struct {
		patrimonio float64
		esperado   bool
	}{
		{250000.00, false},
		{200000.00, false},
		{150000.00, true},
		{0.00, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Patrimônio %.2f", tt.patrimonio), func(t *testing.T) {
			// Type assert to concrete type to test implementation details
			concreteCalc, ok := calculadora.(*calculadoraPrioridade)
			if !ok {
				t.Fatal("Não foi possível fazer type assert para calculadoraPrioridade")
			}
			resultado := concreteCalc.EhPrioridadeNormal(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("EhPrioridadeNormal(%.2f) = %v, esperado %v", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

func TestCalculadoraPrioridade_ValidarLimitePrioridadeAlta(t *testing.T) {
	cfg := setupTestConfig()
	calculadora := NovaCalculadoraPrioridade(cfg)

	// Type assert to concrete type to test implementation details
	concreteCalc, ok := calculadora.(*calculadoraPrioridade)
	if !ok {
		t.Fatal("Não foi possível fazer type assert para calculadoraPrioridade")
	}

	limite := concreteCalc.ValidarLimitePrioridadeAlta()
	if limite != 200000.00 {
		t.Errorf("Limite esperado 200000.00, obtido %.2f", limite)
	}
}

func TestCalculadoraPrioridade_ConfigCustomizado(t *testing.T) {
	customCfg := &config.Config{
		Business: config.BusinessConfig{
			LimitePrioridadeAlta: 150000.00,
		},
	}
	calculadora := NovaCalculadoraPrioridade(customCfg)

	// Testar com limite customizado
	resultado := calculadora.CalcularNivelPrioridade(140000.00)
	if resultado != PrioridadeNormal {
		t.Errorf("Com limite customizado de 150000, patrimonio 140000 deveria ser normal, obtido %s", resultado)
	}

	resultado = calculadora.CalcularNivelPrioridade(160000.00)
	if resultado != PrioridadeAlta {
		t.Errorf("Com limite customizado de 150000, patrimonio 160000 deveria ser alta, obtido %s", resultado)
	}
}

func TestConstantes(t *testing.T) {
	// Testar se as constantes estão definidas corretamente
	if PrioridadeAlta != "prioridade_alta" {
		t.Errorf("PrioridadeAlta esperado 'prioridade_alta', obtido '%s'", PrioridadeAlta)
	}

	if PrioridadeNormal != "prioridade_normal" {
		t.Errorf("PrioridadeNormal esperado 'prioridade_normal', obtido '%s'", PrioridadeNormal)
	}
}
