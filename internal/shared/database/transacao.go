package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MundoInvest/backend/internal/shared/mensagens"
)

// Transacao representa uma transação de banco de dados
type Transacao struct {
	tx *sql.Tx
}

// NovaTransacao cria uma nova transação
func NovaTransacao(banco *sql.DB, nivelIsolamento sql.IsolationLevel) (*Transacao, error) {
	catalogo := mensagens.ObterCatalogo()
	ctx := context.Background()
	tx, err := banco.BeginTx(ctx, &sql.TxOptions{Isolation: nivelIsolamento})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", catalogo.Texto(mensagens.ErrIniciarTransacao), err)
	}

	return &Transacao{tx: tx}, nil
}

// Commit confirma a transação
func (t *Transacao) Commit() error {
	catalogo := mensagens.ObterCatalogo()
	if err := t.tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", catalogo.Texto(mensagens.ErrConfirmarTransacao), err)
	}
	return nil
}

// Rollback desfaz a transação
func (t *Transacao) Rollback() error {
	catalogo := mensagens.ObterCatalogo()
	if err := t.tx.Rollback(); err != nil {
		return fmt.Errorf("%s: %w", catalogo.Texto(mensagens.ErrReverterTransacao), err)
	}
	return nil
}

// Banco retorna a conexão de banco de dados da transação
func (t *Transacao) Banco() *sql.Tx {
	return t.tx
}

// ExecutarEmTransacao executa uma função dentro de uma transação com rollback automático em caso de erro
func ExecutarEmTransacao(banco *sql.DB, nivelIsolamento sql.IsolationLevel, fn func(*Transacao) error) error {
	catalogo := mensagens.ObterCatalogo()
	transacao, err := NovaTransacao(banco, nivelIsolamento)
	if err != nil {
		return err
	}

	// Garante rollback se ocorrer erro
	defer func() {
		if p := recover(); p != nil {
			_ = transacao.Rollback()
			panic(p)
		}
	}()

	// Executa a função
	if err := fn(transacao); err != nil {
		// Se houver erro, faz rollback
		if rbErr := transacao.Rollback(); rbErr != nil {
			return fmt.Errorf(catalogo.Texto(mensagens.ErrReverterTransacaoComErro), rbErr, err)
		}
		return err
	}

	// Se tudo deu certo, faz commit
	if err := transacao.Commit(); err != nil {
		return fmt.Errorf("%s: %w", catalogo.Texto(mensagens.ErrConfirmarTransacao), err)
	}

	return nil
}
