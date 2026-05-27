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
	if valorPatrimonio >= LimitePrioridadeAlta {
		return PrioridadeAlta
	}

	return PrioridadeNormal
}

// CalcularNivelPrioridadeComDetalhes calcula o nível de prioridade e retorna detalhes
func (p *prioridadeCalculator) CalcularNivelPrioridadeComDetalhes(valorPatrimonio float64) (string, string) {
	if valorPatrimonio >= LimitePrioridadeAlta {
		return PrioridadeAlta, fmt.Sprintf("Patrimônio %.2f atinge limite de prioridade alta (>= %.2f)", valorPatrimonio, LimitePrioridadeAlta)
	}

	return PrioridadeNormal, fmt.Sprintf("Patrimônio %.2f abaixo do limite de prioridade alta (< %.2f)", valorPatrimonio, LimitePrioridadeAlta)
}

// ValidarLimitePrioridadeAlta retorna o limite para prioridade alta
func (p *prioridadeCalculator) ValidarLimitePrioridadeAlta() float64 {
	return LimitePrioridadeAlta
}

// EhPrioridadeAlta verifica se um determinado patrimônio é considerado prioridade alta
func (p *prioridadeCalculator) EhPrioridadeAlta(valorPatrimonio float64) bool {
	return valorPatrimonio >= LimitePrioridadeAlta
}

// EhPrioridadeNormal verifica se um determinado patrimônio é considerado prioridade normal
func (p *prioridadeCalculator) EhPrioridadeNormal(valorPatrimonio float64) bool {
	return valorPatrimonio < LimitePrioridadeAlta
}
