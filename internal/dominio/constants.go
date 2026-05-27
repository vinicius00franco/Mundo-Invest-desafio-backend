package dominio

// Constantes de negócio para cálculo de prioridade e status
const (
	// Níveis de prioridade
	PrioridadeAlta   = "prioridade_alta"
	PrioridadeNormal = "prioridade_normal"

	// Status de cliente
	StatusAguardandoAnalise = "Aguardando Análise"
	StatusProcessado        = "Processado"
	StatusEmAnalise         = "Em Análise"
	StatusAprovado          = "Aprovado"
	StatusReprovado         = "Reprovado"
)
