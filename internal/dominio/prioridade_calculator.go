package dominio

import (
	"fmt"
)

// PrioridadeCalculator define a interface para cálculo de nível de prioridade
type PrioridadeCalculator interface {
	CalcularNivelPrioridade(valorPatrimonio float64) string
}

// prioridadeCalculator implementa a interface PrioridadeCalculator
type prioridadeCalculator struct{}

// NovoPrioridadeCalculator cria uma nova instância de PrioridadeCalculator
func NovoPrioridadeCalculator() PrioridadeCalculator {
	return &prioridadeCalculator{}
}

// CalcularNivelPrioridade calcula o nível de prioridade baseado no valor do patrimônio
// Regras de negócio:
// - valor_patrimonio >= 200.000 → nivel_prioridade_alta
// - valor_patrimonio < 200.000 → nivel_prioridade_normal
// - valor_patrimonio == 200.000 → nivel_prioridade_alta
func (p *prioridadeCalculator) CalcularNivelPrioridade(valorPatrimonio float64) string {
	const LIMITE_PRIORIDADE_ALTA = 200000.00

	if valorPatrimonio >= LIMITE_PRIORIDADE_ALTA {
		return "prioridade_alta"
	}

	return "prioridade_normal"
}

// CalcularNivelPrioridadeComDetalhes calcula o nível de prioridade e retorna detalhes
func (p *prioridadeCalculator) CalcularNivelPrioridadeComDetalhes(valorPatrimonio float64) (string, string) {
	const LIMITE_PRIORIDADE_ALTA = 200000.00

	if valorPatrimonio >= LIMITE_PRIORIDADE_ALTA {
		return "prioridade_alta", fmt.Sprintf("Patrimônio %.2f atinge limite de prioridade alta (>= %.2f)", valorPatrimonio, LIMITE_PRIORIDADE_ALTA)
	}

	return "prioridade_normal", fmt.Sprintf("Patrimônio %.2f abaixo do limite de prioridade alta (< %.2f)", valorPatrimonio, LIMITE_PRIORIDADE_ALTA)
}

// ValidarLimitePrioridadeAlta retorna o limite para prioridade alta
func (p *prioridadeCalculator) ValidarLimitePrioridadeAlta() float64 {
	return 200000.00
}

// EhPrioridadeAlta verifica se um determinado patrimônio é considerado prioridade alta
func (p *prioridadeCalculator) EhPrioridadeAlta(valorPatrimonio float64) bool {
	return valorPatrimonio >= 200000.00
}

// EhPrioridadeNormal verifica se um determinado patrimônio é considerado prioridade normal
func (p *prioridadeCalculator) EhPrioridadeNormal(valorPatrimonio float64) bool {
	return valorPatrimonio < 200000.00
}
