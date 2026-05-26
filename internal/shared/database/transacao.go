package database

import (
	"context"
	"database/sql"
	"fmt"
)

// Transacao representa uma transação de banco de dados
type Transacao struct {
	tx *sql.Tx
}

// NovaTransacao cria uma nova transação
func NovaTransacao(banco *sql.DB, nivelIsolamento sql.IsolationLevel) (*Transacao, error) {
	ctx := context.Background()
	tx, err := banco.BeginTx(ctx, &sql.TxOptions{Isolation: nivelIsolamento})
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	return &Transacao{tx: tx}, nil
}

// Commit confirma a transação
func (t *Transacao) Commit() error {
	if err := t.tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit da transação: %w", err)
	}
	return nil
}

// Rollback desfaz a transação
func (t *Transacao) Rollback() error {
	if err := t.tx.Rollback(); err != nil {
		return fmt.Errorf("erro ao fazer rollback da transação: %w", err)
	}
	return nil
}

// Banco retorna a conexão de banco de dados da transação
func (t *Transacao) Banco() *sql.DB {
	return t.tx
}

// ExecutarEmTransacao executa uma função dentro de uma transação com rollback automático em caso de erro
func ExecutarEmTransacao(banco *sql.DB, nivelIsolamento sql.IsolationLevel, fn func(*Transacao) error) error {
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
			return fmt.Errorf("erro ao fazer rollback: %w, erro original: %v", rbErr, err)
		}
		return err
	}

	// Se tudo deu certo, faz commit
	if err := transacao.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}

	return nil
}
