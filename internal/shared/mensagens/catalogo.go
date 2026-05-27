package mensagens

import (
	"fmt"
	"sync"
)

// Mensagem representa uma mensagem simplificada e focada
type Mensagem struct {
	Texto string
}

// Catalogo contém todas as mensagens centralizadas
type Catalogo struct {
	mensagens map[Chave]Mensagem
	mu        sync.RWMutex // Para thread-safety em extensões futuras
}

// Singleton thread-safe
var (
	catalogoSingleton *Catalogo
	catalogoOnce      sync.Once
)

// ObterCatalogo retorna o singleton do catálogo
func ObterCatalogo() *Catalogo {
	catalogoOnce.Do(func() {
		catalogoSingleton = NovoCatalogo()
	})
	return catalogoSingleton
}

// NovoCatalogo cria uma nova instância do catálogo (para testes)
func NovoCatalogo() *Catalogo {
	catalogo := &Catalogo{
		mensagens: make(map[Chave]Mensagem, 50), // Pre-alocar capacidade
	}
	catalogo.inicializarMensagens()
	return catalogo
}

// inicializarMensagens popula o catálogo com todas as mensagens
func (c *Catalogo) inicializarMensagens() {
	// Mensagens de Sucesso
	c.mensagens[CliCriadoSucesso] = Mensagem{
		Texto: "Cliente criado com sucesso",
	}

	c.mensagens[WebhookProcessado] = Mensagem{
		Texto: "Webhook processado com sucesso",
	}

	// Mensagens de Validação
	c.mensagens[ValNomeObrigatorio] = Mensagem{
		Texto: "nome é obrigatório",
	}

	c.mensagens[ValEmailObrigatorio] = Mensagem{
		Texto: "email é obrigatório",
	}

	c.mensagens[ValTipoSolicitacaoObrigatorio] = Mensagem{
		Texto: "tipoSolicitacao é obrigatório",
	}

	c.mensagens[ValValorPatrimonioObrigatorio] = Mensagem{
		Texto: "valorPatrimonio é obrigatório",
	}

	c.mensagens[ValNomeMinimoCaracteres] = Mensagem{
		Texto: "nome deve ter pelo menos 3 caracteres",
	}

	c.mensagens[ValNomeMaximoCaracteres] = Mensagem{
		Texto: "nome deve ter no máximo %d caracteres",
	}

	c.mensagens[ValEmailInvalido] = Mensagem{
		Texto: "email inválido",
	}

	c.mensagens[ValEmailMaximoCaracteres] = Mensagem{
		Texto: "email deve ter no máximo %d caracteres",
	}

	c.mensagens[ValValorPatrimonioNegativo] = Mensagem{
		Texto: "valorPatrimonio deve ser positivo",
	}

	c.mensagens[ValValorPatrimonioZero] = Mensagem{
		Texto: "valorPatrimonio deve ser maior que zero",
	}

	c.mensagens[ValorPatrimonioMaximo] = Mensagem{
		Texto: "valorPatrimonio deve ser menor ou igual a %.2f",
	}

	c.mensagens[ValTipoSolicitacaoMaximo] = Mensagem{
		Texto: "tipoSolicitacao deve ter no máximo %d caracteres",
	}

	// Validação Webhook
	c.mensagens[ValIdEventoObrigatorio] = Mensagem{
		Texto: "identificadorEvento é obrigatório",
	}

	c.mensagens[ValIdCardObrigatorio] = Mensagem{
		Texto: "identificadorCard é obrigatório",
	}

	c.mensagens[ValEmailClienteObrigatorio] = Mensagem{
		Texto: "emailCliente é obrigatório",
	}

	c.mensagens[ValDataEventoObrigatorio] = Mensagem{
		Texto: "dataEvento é obrigatório",
	}

	c.mensagens[ValEmailClienteInvalido] = Mensagem{
		Texto: "emailCliente inválido",
	}

	c.mensagens[ValDataEventoInvalido] = Mensagem{
		Texto: "dataEvento inválido",
	}

	// Erros de Serviço
	c.mensagens[ErrSalvarCliente] = Mensagem{
		Texto: "erro ao salvar cliente",
	}

	c.mensagens[ErrAtualizarCliente] = Mensagem{
		Texto: "erro ao atualizar cliente",
	}

	c.mensagens[ErrVerificarIdempotencia] = Mensagem{
		Texto: "erro ao verificar idempotência",
	}

	c.mensagens[ErrProcessarDataEvento] = Mensagem{
		Texto: "erro ao processar data do evento",
	}

	c.mensagens[ErrSalvarEvento] = Mensagem{
		Texto: "erro ao salvar evento",
	}

	// Erros de Integração
	c.mensagens[ErrCriarCardPipefy] = Mensagem{
		Texto: "erro ao criar card no Pipefy",
	}

	c.mensagens[ErrAtualizarCardPipefy] = Mensagem{
		Texto: "erro ao atualizar card no Pipefy",
	}

	c.mensagens[ErrEstruturarMutacaoCriar] = Mensagem{
		Texto: "erro ao estruturar mutation createCard",
	}

	c.mensagens[ErrEstruturarMutacaoAtualizar] = Mensagem{
		Texto: "erro ao estruturar mutation updateCard",
	}

	c.mensagens[ErrExecutarMutacaoCriar] = Mensagem{
		Texto: "erro ao executar mutation createCard",
	}

	c.mensagens[ErrExecutarMutacaoAtualizar] = Mensagem{
		Texto: "erro ao executar mutation updateCard",
	}

	// Erros HTTP
	c.mensagens[ErrMetodoNaoPermitido] = Mensagem{
		Texto: "Método não permitido",
	}

	c.mensagens[ErrProcessarJSON] = Mensagem{
		Texto: "Erro ao processar JSON",
	}

	c.mensagens[ErrGerarResposta] = Mensagem{
		Texto: "Erro ao gerar resposta",
	}

	c.mensagens[ErrClienteNaoEncontrado] = Mensagem{
		Texto: "Cliente não encontrado",
	}

	c.mensagens[ErrProcessarWebhook] = Mensagem{
		Texto: "Erro ao processar webhook",
	}

	// Erros de Banco de Dados
	c.mensagens[ErrAbrirConexao] = Mensagem{
		Texto: "erro ao abrir conexão com banco de dados",
	}

	c.mensagens[ErrTestarConexao] = Mensagem{
		Texto: "erro ao testar conexão com banco de dados",
	}

	c.mensagens[ErrIniciarTransacao] = Mensagem{
		Texto: "erro ao iniciar transação",
	}

	c.mensagens[ErrConfirmarTransacao] = Mensagem{
		Texto: "erro ao confirmar transação",
	}

	c.mensagens[ErrReverterTransacao] = Mensagem{
		Texto: "erro ao reverter transação",
	}

	c.mensagens[ErrReverterTransacaoComErro] = Mensagem{
		Texto: "erro ao fazer rollback: %w, erro original: %v",
	}

	// Erros de Tipos de Erro
	c.mensagens[ErrValidacaoCampo] = Mensagem{
		Texto: "erro de validação no campo '%s': %s",
	}

	c.mensagens[ErrValidacao] = Mensagem{
		Texto: "erro de validação: %s",
	}

	c.mensagens[ErrRepositorio] = Mensagem{
		Texto: "erro no repositório durante operação '%s': %v",
	}

	c.mensagens[ErrServico] = Mensagem{
		Texto: "erro no serviço '%s': %s",
	}

	c.mensagens[ErrServicoComErro] = Mensagem{
		Texto: "erro no serviço '%s': %s: %v",
	}

	c.mensagens[ErrIntegracao] = Mensagem{
		Texto: "erro de integração com sistema '%s': %s",
	}

	c.mensagens[ErrIntegracaoComErro] = Mensagem{
		Texto: "erro de integração com sistema '%s': %s: %v",
	}

	c.mensagens[ErrTimeout] = Mensagem{
		Texto: "timeout na operação '%s': %s",
	}

	c.mensagens[ErrNaoEncontrado] = Mensagem{
		Texto: "recurso '%s' não encontrado",
	}

	c.mensagens[ErrNaoEncontradoComID] = Mensagem{
		Texto: "recurso '%s' com ID '%s' não encontrado",
	}

	c.mensagens[ErrEstrategiaValidacao] = Mensagem{
		Texto: "erro na estratégia %s: %w",
	}

	c.mensagens[ErrObjetoInvalido] = Mensagem{
		Texto: "objeto inválido para validação de cliente",
	}

	c.mensagens[ErrEsperadoErro] = Mensagem{
		Texto: "Esperado erro, mas não houve erro",
	}

	c.mensagens[ErrViolacaoUnicidade] = Mensagem{
		Texto: "Esperado erro de violação de unicidade, mas não houve erro",
	}

	c.mensagens[ErrDespacharEvento] = Mensagem{
		Texto: "Erro ao despachar evento: %v",
	}

	c.mensagens[ErrProcessarEvento] = Mensagem{
		Texto: "Erro ao processar evento",
	}
}

// Texto retorna o texto de uma mensagem pela chave
func (c *Catalogo) Texto(chave Chave) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if msg, exists := c.mensagens[chave]; exists {
		return msg.Texto
	}
	return fmt.Sprintf("mensagem não encontrada: %s", chave)
}

// TextoFormatado retorna o texto de uma mensagem com parâmetros substituídos
func (c *Catalogo) TextoFormatado(chave Chave, args ...interface{}) string {
	texto := c.Texto(chave)
	return fmt.Sprintf(texto, args...)
}
