package dominio

import (
	"fmt"

	"github.com/MundoInvest/backend/internal/shared/config"
)

// CalculadoraPrioridade define a interface para cálculo de nível de prioridade
type CalculadoraPrioridade interface {
	CalcularNivelPrioridade(valorPatrimonio float64) string
}

// calculadoraPrioridade implementa a interface CalculadoraPrioridade
type calculadoraPrioridade struct {
	config *config.Config
}

// NovaCalculadoraPrioridade cria uma nova instância de CalculadoraPrioridade
func NovaCalculadoraPrioridade(cfg *config.Config) CalculadoraPrioridade {
	return &calculadoraPrioridade{
		config: cfg,
	}
}

// CalcularNivelPrioridade calcula o nível de prioridade baseado no valor do patrimônio
// Regras de negócio:
// - valorPatrimonio >= limite configurado → nivelPrioridadeAlta
// - valorPatrimonio < limite configurado → nivelPrioridadeNormal
func (p *calculadoraPrioridade) CalcularNivelPrioridade(valorPatrimonio float64) string {
	if valorPatrimonio >= p.config.Business.LimitePrioridadeAlta {
		return PrioridadeAlta
	}

	return PrioridadeNormal
}

// CalcularNivelPrioridadeComDetalhes calcula o nível de prioridade e retorna detalhes
func (p *calculadoraPrioridade) CalcularNivelPrioridadeComDetalhes(valorPatrimonio float64) (string, string) {
	if valorPatrimonio >= p.config.Business.LimitePrioridadeAlta {
		return PrioridadeAlta, fmt.Sprintf("Patrimônio %.2f atinge limite de prioridade alta (>= %.2f)", valorPatrimonio, p.config.Business.LimitePrioridadeAlta)
	}

	return PrioridadeNormal, fmt.Sprintf("Patrimônio %.2f abaixo do limite de prioridade alta (< %.2f)", valorPatrimonio, p.config.Business.LimitePrioridadeAlta)
}

// ValidarLimitePrioridadeAlta retorna o limite para prioridade alta
func (p *calculadoraPrioridade) ValidarLimitePrioridadeAlta() float64 {
	return p.config.Business.LimitePrioridadeAlta
}

// EhPrioridadeAlta verifica se um determinado patrimônio é considerado prioridade alta
func (p *calculadoraPrioridade) EhPrioridadeAlta(valorPatrimonio float64) bool {
	return valorPatrimonio >= p.config.Business.LimitePrioridadeAlta
}

// EhPrioridadeNormal verifica se um determinado patrimônio é considerado prioridade normal
func (p *calculadoraPrioridade) EhPrioridadeNormal(valorPatrimonio float64) bool {
	return valorPatrimonio < p.config.Business.LimitePrioridadeAlta
}
