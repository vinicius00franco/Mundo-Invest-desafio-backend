package gestao_clientes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
)

// Cliente representa a entidade Cliente do contexto de gestão de clientes
type Cliente struct {
	IdentificadorInterno int64     // gcl_cli_int
	IdentificadorExterno string    // gcl_cli_ext
	Nome                 string    // gcl_cli_nom
	Email                string    // gcl_cli_ema
	ValorPatrimonio      float64   // gcl_cli_pat
	TipoSolicitacao      string    // gcl_cli_tso
	Status               string    // gcl_cli_stc
	NivelPrioridade      string    // gcl_cli_npr
	DataCriacao          time.Time // gcl_cli_dcr
	DataAtualizacao      time.Time // gcl_cli_dat
}

// ClienteRepository define a interface para operações de persistência de clientes
type ClienteRepository interface {
	Salvar(ctx context.Context, cliente Cliente) (*Cliente, error)
	BuscarPorIdentificadorInterno(ctx context.Context, identificadorInterno int64) (*Cliente, error)
	BuscarPorEmail(ctx context.Context, email string) (*Cliente, error)
	BuscarPorIdentificadorExterno(ctx context.Context, identificadorExterno string) (*Cliente, error)
	Atualizar(ctx context.Context, cliente Cliente) error
}

// clienteRepository implementa a interface ClienteRepository
type clienteRepository struct {
	banco  *sql.DB
	config *config.Config
}

// NovoClienteRepository cria uma nova instância de ClienteRepository
func NovoClienteRepository(banco *sql.DB, cfg *config.Config) ClienteRepository {
	return &clienteRepository{
		banco:  banco,
		config: cfg,
	}
}

// Salvar insere um novo cliente no banco de dados
func (r *clienteRepository) Salvar(ctx context.Context, cliente Cliente) (*Cliente, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		INSERT INTO gestao_clientes.cliente (
			gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, gcl_cli_pat, 
			gcl_cli_tso, gcl_cli_stc, gcl_cli_npr
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING gcl_cli_int, gcl_cli_dcr
	`

	var identificadorInterno int64
	var dataCriacao time.Time

	err := r.banco.QueryRowContext(
		ctx,
		query,
		cliente.IdentificadorExterno,
		cliente.Nome,
		cliente.Email,
		cliente.ValorPatrimonio,
		cliente.TipoSolicitacao,
		cliente.Status,
		cliente.NivelPrioridade,
	).Scan(&identificadorInterno, &dataCriacao)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("salvar cliente", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("salvar cliente", err)
	}

	cliente.IdentificadorInterno = identificadorInterno
	cliente.DataCriacao = dataCriacao

	return &cliente, nil
}

// BuscarPorIdentificadorInterno busca um cliente pelo identificador interno
func (r *clienteRepository) BuscarPorIdentificadorInterno(ctx context.Context, identificadorInterno int64) (*Cliente, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		SELECT 
			gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, 
			gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr, 
			gcl_cli_dcr, COALESCE(gcl_cli_dat, gcl_cli_dcr)
		FROM gestao_clientes.cliente
		WHERE gcl_cli_int = $1
	`

	var cliente Cliente

	err := r.banco.QueryRowContext(ctx, query, identificadorInterno).Scan(
		&cliente.IdentificadorInterno,
		&cliente.IdentificadorExterno,
		&cliente.Nome,
		&cliente.Email,
		&cliente.ValorPatrimonio,
		&cliente.TipoSolicitacao,
		&cliente.Status,
		&cliente.NivelPrioridade,
		&cliente.DataCriacao,
		&cliente.DataAtualizacao,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("cliente", fmt.Sprintf("%d", identificadorInterno))
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("buscar cliente por identificador interno", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("buscar cliente por identificador interno", err)
	}

	return &cliente, nil
}

// BuscarPorEmail busca um cliente pelo email
func (r *clienteRepository) BuscarPorEmail(ctx context.Context, email string) (*Cliente, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		SELECT 
			gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, 
			gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr, 
			gcl_cli_dcr, COALESCE(gcl_cli_dat, gcl_cli_dcr)
		FROM gestao_clientes.cliente
		WHERE gcl_cli_ema = $1
	`

	var cliente Cliente

	err := r.banco.QueryRowContext(ctx, query, email).Scan(
		&cliente.IdentificadorInterno,
		&cliente.IdentificadorExterno,
		&cliente.Nome,
		&cliente.Email,
		&cliente.ValorPatrimonio,
		&cliente.TipoSolicitacao,
		&cliente.Status,
		&cliente.NivelPrioridade,
		&cliente.DataCriacao,
		&cliente.DataAtualizacao,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("cliente", email)
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("buscar cliente por email", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("buscar cliente por email", err)
	}

	return &cliente, nil
}

// BuscarPorIdentificadorExterno busca um cliente pelo identificador externo (card_id)
func (r *clienteRepository) BuscarPorIdentificadorExterno(ctx context.Context, identificadorExterno string) (*Cliente, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		SELECT 
			gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, 
			gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr, 
			gcl_cli_dcr, COALESCE(gcl_cli_dat, gcl_cli_dcr)
		FROM gestao_clientes.cliente
		WHERE gcl_cli_ext = $1
	`

	var cliente Cliente

	err := r.banco.QueryRowContext(ctx, query, identificadorExterno).Scan(
		&cliente.IdentificadorInterno,
		&cliente.IdentificadorExterno,
		&cliente.Nome,
		&cliente.Email,
		&cliente.ValorPatrimonio,
		&cliente.TipoSolicitacao,
		&cliente.Status,
		&cliente.NivelPrioridade,
		&cliente.DataCriacao,
		&cliente.DataAtualizacao,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("cliente", identificadorExterno)
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("buscar cliente por identificador externo", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("buscar cliente por identificador externo", err)
	}

	return &cliente, nil
}

// Atualizar atualiza um cliente existente no banco de dados
func (r *clienteRepository) Atualizar(ctx context.Context, cliente Cliente) error {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		UPDATE gestao_clientes.cliente
		SET 
			gcl_cli_ext = $2,
			gcl_cli_nom = $3,
			gcl_cli_ema = $4,
			gcl_cli_pat = $5,
			gcl_cli_tso = $6,
			gcl_cli_stc = $7,
			gcl_cli_npr = $8,
			gcl_cli_dat = CURRENT_TIMESTAMP
		WHERE gcl_cli_int = $1
	`

	resultado, err := r.banco.ExecContext(
		ctx,
		query,
		cliente.IdentificadorInterno,
		cliente.IdentificadorExterno,
		cliente.Nome,
		cliente.Email,
		cliente.ValorPatrimonio,
		cliente.TipoSolicitacao,
		cliente.Status,
		cliente.NivelPrioridade,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return errors.NewTimeoutError("atualizar cliente", "timeout ao executar operação de banco de dados")
		}
		return errors.NewRepositoryError("atualizar cliente", err)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return errors.NewRepositoryError("verificar linhas afetadas", err)
	}

	if linhasAfetadas == 0 {
		return errors.NewNotFoundError("cliente", fmt.Sprintf("%d", cliente.IdentificadorInterno))
	}

	return nil
}
