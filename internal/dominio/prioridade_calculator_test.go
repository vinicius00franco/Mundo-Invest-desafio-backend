package dominio

import (
	"fmt"
	"testing"
)

func TestPrioridadeCalculator_CalcularNivelPrioridade(t *testing.T) {
	calculator := NovoPrioridadeCalculator()

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
			resultado := calculator.CalcularNivelPrioridade(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("CalcularNivelPrioridade(%.2f) = %s, esperado %s", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

func TestPrioridadeCalculator_CalcularNivelPrioridadeComDetalhes(t *testing.T) {
	// Skip test for now as method was removed in refactor
	t.Skip("CalcularNivelPrioridadeComDetalhes não está disponível na versão atual")
}

func TestPrioridadeCalculator_EhPrioridadeAlta(t *testing.T) {
	calculator := NovoPrioridadeCalculator()

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
			concreteCalc, ok := calculator.(*prioridadeCalculator)
			if !ok {
				t.Fatal("Não foi possível fazer type assert para prioridadeCalculator")
			}
			resultado := concreteCalc.EhPrioridadeAlta(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("EhPrioridadeAlta(%.2f) = %v, esperado %v", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

func TestPrioridadeCalculator_EhPrioridadeNormal(t *testing.T) {
	calculator := NovoPrioridadeCalculator()

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
			concreteCalc, ok := calculator.(*prioridadeCalculator)
			if !ok {
				t.Fatal("Não foi possível fazer type assert para prioridadeCalculator")
			}
			resultado := concreteCalc.EhPrioridadeNormal(tt.patrimonio)
			if resultado != tt.esperado {
				t.Errorf("EhPrioridadeNormal(%.2f) = %v, esperado %v", tt.patrimonio, resultado, tt.esperado)
			}
		})
	}
}

func TestPrioridadeCalculator_ValidarLimitePrioridadeAlta(t *testing.T) {
	calculator := NovoPrioridadeCalculator()

	// Type assert to concrete type to test implementation details
	concreteCalc, ok := calculator.(*prioridadeCalculator)
	if !ok {
		t.Fatal("Não foi possível fazer type assert para prioridadeCalculator")
	}

	limite := concreteCalc.ValidarLimitePrioridadeAlta()
	if limite != 200000.00 {
		t.Errorf("Limite esperado 200000.00, obtido %.2f", limite)
	}
}

func TestConstantes(t *testing.T) {
	// Testar se as constantes estão definidas corretamente
	if LimitePrioridadeAlta != 200000.00 {
		t.Errorf("LimitePrioridadeAlta esperado 200000.00, obtido %.2f", LimitePrioridadeAlta)
	}

	if PrioridadeAlta != "prioridade_alta" {
		t.Errorf("PrioridadeAlta esperado 'prioridade_alta', obtido '%s'", PrioridadeAlta)
	}

	if PrioridadeNormal != "prioridade_normal" {
		t.Errorf("PrioridadeNormal esperado 'prioridade_normal', obtido '%s'", PrioridadeNormal)
	}
}
