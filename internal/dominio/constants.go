package dominio

// Constantes de negócio para cálculo de prioridade
const (
	// LimitePrioridadeAlta define o limite de patrimônio para prioridade alta
	LimitePrioridadeAlta = 200000.00
	
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

// Constantes de configuração de banco de dados
const (
	// MaxOpenConns define o número máximo de conexões abertas
	MaxOpenConns = 25
	
	// MaxIdleConns define o número máximo de conexões ociosas
	MaxIdleConns = 5
	
	// ConnMaxLifetime define o tempo máximo de vida de uma conexão (em segundos)
	ConnMaxLifetime = 300 // 5 minutos
	
	// ConnMaxIdleTime define o tempo máximo que uma conexão pode ficar ociosa (em segundos)
	ConnMaxIdleTime = 60 // 1 minuto
)

// Constantes de timeout
const (
	// DefaultTimeout define o timeout padrão para operações (em segundos)
	DefaultTimeout = 30
	
	// DatabaseTimeout define o timeout para operações de banco de dados (em segundos)
	DatabaseTimeout = 10
	
	// ExternalAPITimeout define o timeout para chamadas de API externa (em segundos)
	ExternalAPITimeout = 15
)

// Constantes de HTTP
const (
	// DefaultPort define a porta padrão do servidor
	DefaultPort = "8080"
	
	// ReadTimeout define o timeout de leitura HTTP (em segundos)
	ReadTimeout = 15
	
	// WriteTimeout define o timeout de escrita HTTP (em segundos)
	WriteTimeout = 15
	
	// IdleTimeout define o timeout de ociosidade HTTP (em segundos)
	IdleTimeout = 60
)

// Constantes de validação
const (
	// MaxNomeLength define o tamanho máximo do nome
	MaxNomeLength = 255
	
	// MaxEmailLength define o tamanho máximo do email
	MaxEmailLength = 255
	
	// MaxTipoSolicitacaoLength define o tamanho máximo do tipo de solicitação
	MaxTipoSolicitacaoLength = 50
	
	// MinValorPatrimonio define o valor mínimo do patrimônio
	MinValorPatrimonio = 0.01
	
	// MaxValorPatrimonio define o valor máximo do patrimônio
	MaxValorPatrimonio = 999999999.99
)
