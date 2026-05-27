package dominio

import (
	"fmt"

	"github.com/MundoInvest/backend/internal/shared/config"
)

// PrioridadeCalculator define a interface para cálculo de nível de prioridade
type PrioridadeCalculator interface {
	CalcularNivelPrioridade(valorPatrimonio float64) string
}

// prioridadeCalculator implementa a interface PrioridadeCalculator
type prioridadeCalculator struct {
	config *config.Config
}

// NovoPrioridadeCalculator cria uma nova instância de PrioridadeCalculator
func NovoPrioridadeCalculator(cfg *config.Config) PrioridadeCalculator {
	return &prioridadeCalculator{
		config: cfg,
	}
}

// CalcularNivelPrioridade calcula o nível de prioridade baseado no valor do patrimônio
// Regras de negócio:
// - valor_patrimonio >= limite configurado → nivel_prioridade_alta
// - valor_patrimonio < limite configurado → nivel_prioridade_normal
func (p *prioridadeCalculator) CalcularNivelPrioridade(valorPatrimonio float64) string {
	if valorPatrimonio >= p.config.Business.LimitePrioridadeAlta {
		return PrioridadeAlta
	}

	return PrioridadeNormal
}

// CalcularNivelPrioridadeComDetalhes calcula o nível de prioridade e retorna detalhes
func (p *prioridadeCalculator) CalcularNivelPrioridadeComDetalhes(valorPatrimonio float64) (string, string) {
	if valorPatrimonio >= p.config.Business.LimitePrioridadeAlta {
		return PrioridadeAlta, fmt.Sprintf("Patrimônio %.2f atinge limite de prioridade alta (>= %.2f)", valorPatrimonio, p.config.Business.LimitePrioridadeAlta)
	}

	return PrioridadeNormal, fmt.Sprintf("Patrimônio %.2f abaixo do limite de prioridade alta (< %.2f)", valorPatrimonio, p.config.Business.LimitePrioridadeAlta)
}

// ValidarLimitePrioridadeAlta retorna o limite para prioridade alta
func (p *prioridadeCalculator) ValidarLimitePrioridadeAlta() float64 {
	return p.config.Business.LimitePrioridadeAlta
}

// EhPrioridadeAlta verifica se um determinado patrimônio é considerado prioridade alta
func (p *prioridadeCalculator) EhPrioridadeAlta(valorPatrimonio float64) bool {
	return valorPatrimonio >= p.config.Business.LimitePrioridadeAlta
}

// EhPrioridadeNormal verifica se um determinado patrimônio é considerado prioridade normal
func (p *prioridadeCalculator) EhPrioridadeNormal(valorPatrimonio float64) bool {
	return valorPatrimonio < p.config.Business.LimitePrioridadeAlta
}
