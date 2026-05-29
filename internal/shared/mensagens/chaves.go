package mensagens

// Chaves de mensagens - nomes curtos mas descritivos
const (
	// Sucesso
	CliCriadoSucesso  Chave = "cliente.criado.sucesso"
	WebhookProcessado Chave = "webhook.processado.sucesso"

	// Validação - Cliente
	ValNomeObrigatorio            Chave = "validacao.nome.obrigatorio"
	ValEmailObrigatorio           Chave = "validacao.email.obrigatorio"
	ValTipoSolicitacaoObrigatorio Chave = "validacao.tipo_solicitacao.obrigatorio"
	ValValorPatrimonioObrigatorio Chave = "validacao.valor_patrimonio.obrigatorio"
	ValNomeMinimoCaracteres       Chave = "validacao.nome.minimo_caracteres"
	ValNomeMaximoCaracteres       Chave = "validacao.nome.maximo_caracteres"
	ValEmailInvalido              Chave = "validacao.email.invalido"
	ValEmailMaximoCaracteres      Chave = "validacao.email.maximo_caracteres"
	ValValorPatrimonioNegativo    Chave = "validacao.valor_patrimonio.negativo"
	ValValorPatrimonioZero        Chave = "validacao.valor_patrimonio.zero"
	ValorPatrimonioMaximo         Chave = "validacao.valor_patrimonio.maximo"
	ValTipoSolicitacaoMaximo      Chave = "validacao.tipo_solicitacao.maximo"

	// Validação - Webhook
	ValIdEventoObrigatorio     Chave = "validacao.id_evento.obrigatorio"
	ValIdCardObrigatorio       Chave = "validacao.id_card.obrigatorio"
	ValEmailClienteObrigatorio Chave = "validacao.email_cliente.obrigatorio"
	ValDataEventoObrigatorio   Chave = "validacao.data_evento.obrigatorio"
	ValEmailClienteInvalido    Chave = "validacao.email_cliente.invalido"
	ValDataEventoInvalido      Chave = "validacao.data_evento.invalido"

	// Erros - Serviço
	ErrSalvarCliente         Chave = "erro.cliente.salvar"
	ErrAtualizarCliente      Chave = "erro.cliente.atualizar"
	ErrVerificarIdempotencia Chave = "erro.webhook.verificar_idempotencia"
	ErrProcessarDataEvento   Chave = "erro.webhook.processar_data_evento"
	ErrSalvarEvento          Chave = "erro.webhook.salvar_evento"

	// Erros - Integração
	ErrCriarCardPipefy            Chave = "erro.pipefy.criar_card"
	ErrAtualizarCardPipefy        Chave = "erro.pipefy.atualizar_card"
	ErrEstruturarMutacaoCriar     Chave = "erro.pipefy.estruturar_mutacao_criar"
	ErrEstruturarMutacaoAtualizar Chave = "erro.pipefy.estruturar_mutacao_atualizar"
	ErrExecutarMutacaoCriar       Chave = "erro.pipefy.executar_mutacao_criar"
	ErrExecutarMutacaoAtualizar   Chave = "erro.pipefy.executar_mutacao_atualizar"

	// Erros - HTTP
	ErrMetodoNaoPermitido   Chave = "erro.http.metodo_nao_permitido"
	ErrProcessarJSON        Chave = "erro.http.processar_json"
	ErrGerarResposta        Chave = "erro.http.gerar_resposta"
	ErrClienteNaoEncontrado Chave = "erro.http.cliente_nao_encontrado"
	ErrProcessarWebhook     Chave = "erro.http.processar_webhook"

	// Erros - Banco de Dados
	ErrAbrirConexao             Chave = "erro.banco_dados.abrir_conexao"
	ErrTestarConexao            Chave = "erro.banco_dados.testar_conexao"
	ErrIniciarTransacao         Chave = "erro.banco_dados.iniciar_transacao"
	ErrConfirmarTransacao       Chave = "erro.banco_dados.confirmar_transacao"
	ErrReverterTransacao        Chave = "erro.banco_dados.reverter_transacao"
	ErrReverterTransacaoComErro Chave = "erro.banco_dados.reverter_transacao_com_erro"

	// Erros - Tipos de Erro
	ErrValidacaoCampo      Chave = "erro.validacao.campo"
	ErrValidacao           Chave = "erro.validacao"
	ErrRepositorio         Chave = "erro.repositorio"
	ErrServico             Chave = "erro.servico"
	ErrServicoComErro      Chave = "erro.servico_com_erro"
	ErrIntegracao          Chave = "erro.integracao"
	ErrIntegracaoComErro   Chave = "erro.integracao_com_erro"
	ErrTimeout             Chave = "erro.timeout"
	ErrNaoEncontrado       Chave = "erro.nao_encontrado"
	ErrNaoEncontradoComID  Chave = "erro.nao_encontrado_com_id"
	ErrEstrategiaValidacao Chave = "erro.estrategia_validacao"
	ErrObjetoInvalido      Chave = "erro.objeto_invalido"
	ErrEsperadoErro        Chave = "erro.esperado_erro"
	ErrViolacaoUnicidade   Chave = "erro.violacao_unicidade"
	ErrDespacharEvento     Chave = "erro.despachar_evento"
	ErrProcessarEvento     Chave = "erro.processar_evento"
)
