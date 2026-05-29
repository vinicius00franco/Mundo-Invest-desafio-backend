package mensagens

// Chave type-safe para evitar "magic strings"
type Chave string

// ProvedorMensagens interface para testabilidade
type ProvedorMensagens interface {
	Texto(chave Chave) string
	TextoFormatado(chave Chave, args ...interface{}) string
}

// TipoMensagem define o tipo da mensagem no contexto do domínio
type TipoMensagem string

const (
	TipoSucesso    TipoMensagem = "sucesso"
	TipoValidacao  TipoMensagem = "validacao"
	TipoServico    TipoMensagem = "servico"
	TipoIntegracao TipoMensagem = "integracao"
	TipoHTTP       TipoMensagem = "http"
	TipoBancoDados TipoMensagem = "banco_dados"
)

// ContextoMensagem define o contexto de negócio da mensagem
type ContextoMensagem string

const (
	ContextoGestaoClientes      ContextoMensagem = "gestao_clientes"
	ContextoProcessamentoEventos ContextoMensagem = "processamento_eventos"
	ContextoIntegracaoPipefy    ContextoMensagem = "integracao_pipefy"
	ContextoBancoDados          ContextoMensagem = "banco_dados"
	ContextoGenerico            ContextoMensagem = "generico"
)
