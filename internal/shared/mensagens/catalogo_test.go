package mensagens

import (
	"testing"
)

func TestNovoCatalogo(t *testing.T) {
	catalogo := NovoCatalogo()
	if catalogo == nil {
		t.Fatal("NovoCatalogo() retornou nil")
	}

	if catalogo.mensagens == nil {
		t.Fatal("mensagens map não foi inicializado")
	}
}

func TestObterCatalogoSingleton(t *testing.T) {
	catalogo1 := ObterCatalogo()
	catalogo2 := ObterCatalogo()

	if catalogo1 != catalogo2 {
		t.Error("ObterCatalogo() não retornou o mesmo singleton")
	}
}

func TestTextoMensagemExistente(t *testing.T) {
	catalogo := NovoCatalogo()

	texto := catalogo.Texto(CliCriadoSucesso)
	if texto != "Cliente criado com sucesso" {
		t.Errorf("Texto() = %s, queria 'Cliente criado com sucesso'", texto)
	}

	texto = catalogo.Texto(ValNomeObrigatorio)
	if texto != "nome é obrigatório" {
		t.Errorf("Texto() = %s, queria 'nome é obrigatório'", texto)
	}
}

func TestTextoMensagemInexistente(t *testing.T) {
	catalogo := NovoCatalogo()

	chaveInexistente := Chave("chave.inexistente")
	texto := catalogo.Texto(chaveInexistente)

	esperado := "mensagem não encontrada: chave.inexistente"
	if texto != esperado {
		t.Errorf("Texto() = %s, queria %s", texto, esperado)
	}
}

func TestTextoFormatado(t *testing.T) {
	catalogo := NovoCatalogo()

	// Testar sem parâmetros (deve funcionar como Texto normal)
	texto := catalogo.TextoFormatado(CliCriadoSucesso)
	esperado := "Cliente criado com sucesso"
	if texto != esperado {
		t.Errorf("TextoFormatado() = %s, queria %s", texto, esperado)
	}

	// Testar que TextoFormatado funciona como Texto quando não há parâmetros
	texto2 := catalogo.TextoFormatado(ValNomeObrigatorio)
	esperado2 := "nome é obrigatório"
	if texto2 != esperado2 {
		t.Errorf("TextoFormatado() = %s, queria %s", texto2, esperado2)
	}
}

func TestTodasMensagensDefinidas(t *testing.T) {
	catalogo := NovoCatalogo()

	// Verificar se todas as chaves estão definidas
	chaves := []Chave{
		CliCriadoSucesso,
		WebhookProcessado,
		ValNomeObrigatorio,
		ValEmailObrigatorio,
		ValTipoSolicitacaoObrigatorio,
		ValValorPatrimonioObrigatorio,
		ValNomeMinimoCaracteres,
		ValEmailInvalido,
		ValValorPatrimonioNegativo,
		ValValorPatrimonioZero,
		ValIdEventoObrigatorio,
		ValIdCardObrigatorio,
		ValEmailClienteObrigatorio,
		ValDataEventoObrigatorio,
		ValEmailClienteInvalido,
		ValDataEventoInvalido,
		ErrSalvarCliente,
		ErrAtualizarCliente,
		ErrVerificarIdempotencia,
		ErrProcessarDataEvento,
		ErrSalvarEvento,
		ErrCriarCardPipefy,
		ErrAtualizarCardPipefy,
		ErrEstruturarMutacaoCriar,
		ErrEstruturarMutacaoAtualizar,
		ErrExecutarMutacaoCriar,
		ErrExecutarMutacaoAtualizar,
		ErrMetodoNaoPermitido,
		ErrProcessarJSON,
		ErrGerarResposta,
		ErrClienteNaoEncontrado,
		ErrProcessarWebhook,
		ErrAbrirConexao,
		ErrTestarConexao,
		ErrIniciarTransacao,
		ErrConfirmarTransacao,
		ErrReverterTransacao,
	}

	for _, chave := range chaves {
		texto := catalogo.Texto(chave)
		if texto == "" || texto == "mensagem não encontrada" {
			t.Errorf("Chave %s não está definida no catálogo", chave)
		}
	}
}

func TestProvedorMensagensInterface(t *testing.T) {
	catalogo := NovoCatalogo()

	// Verificar se Catalogo implementa a interface
	var provedor ProvedorMensagens = catalogo

	if provedor == nil {
		t.Fatal("Catalogo não implementa ProvedorMensagens")
	}

	// Testar métodos da interface
	texto := provedor.Texto(CliCriadoSucesso)
	if texto != "Cliente criado com sucesso" {
		t.Errorf("Provedor.Texto() = %s, queria 'Cliente criado com sucesso'", texto)
	}
}
