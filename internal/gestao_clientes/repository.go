package gestao_clientes

import (
	"database/sql"
	"fmt"
	"time"
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
	Salvar(cliente Cliente) (*Cliente, error)
	BuscarPorIdentificadorInterno(identificadorInterno int64) (*Cliente, error)
	BuscarPorEmail(email string) (*Cliente, error)
	BuscarPorIdentificadorExterno(identificadorExterno string) (*Cliente, error)
	Atualizar(cliente Cliente) error
}

// clienteRepository implementa a interface ClienteRepository
type clienteRepository struct {
	banco *sql.DB
}

// NovoClienteRepository cria uma nova instância de ClienteRepository
func NovoClienteRepository(banco *sql.DB) ClienteRepository {
	return &clienteRepository{banco: banco}
}

// Salvar insere um novo cliente no banco de dados
func (r *clienteRepository) Salvar(cliente Cliente) (*Cliente, error) {
	query := `
		INSERT INTO gestao_clientes.cliente (
			gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, gcl_cli_pat, 
			gcl_cli_tso, gcl_cli_stc, gcl_cli_npr
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING gcl_cli_int, gcl_cli_dcr
	`

	var identificadorInterno int64
	var dataCriacao time.Time

	err := r.banco.QueryRow(
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
		return nil, fmt.Errorf("erro ao salvar cliente: %w", err)
	}

	cliente.IdentificadorInterno = identificadorInterno
	cliente.DataCriacao = dataCriacao

	return &cliente, nil
}

// BuscarPorIdentificadorInterno busca um cliente pelo identificador interno
func (r *clienteRepository) BuscarPorIdentificadorInterno(identificadorInterno int64) (*Cliente, error) {
	query := `
		SELECT 
			gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, 
			gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr, 
			gcl_cli_dcr, COALESCE(gcl_cli_dat, gcl_cli_dcr)
		FROM gestao_clientes.cliente
		WHERE gcl_cli_int = $1
	`

	var cliente Cliente

	err := r.banco.QueryRow(query, identificadorInterno).Scan(
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
			return nil, fmt.Errorf("cliente não encontrado com identificador interno %d", identificadorInterno)
		}
		return nil, fmt.Errorf("erro ao buscar cliente por identificador interno: %w", err)
	}

	return &cliente, nil
}

// BuscarPorEmail busca um cliente pelo email
func (r *clienteRepository) BuscarPorEmail(email string) (*Cliente, error) {
	query := `
		SELECT 
			gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, 
			gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr, 
			gcl_cli_dcr, COALESCE(gcl_cli_dat, gcl_cli_dcr)
		FROM gestao_clientes.cliente
		WHERE gcl_cli_ema = $1
	`

	var cliente Cliente

	err := r.banco.QueryRow(query, email).Scan(
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
			return nil, fmt.Errorf("cliente não encontrado com email %s", email)
		}
		return nil, fmt.Errorf("erro ao buscar cliente por email: %w", err)
	}

	return &cliente, nil
}

// BuscarPorIdentificadorExterno busca um cliente pelo identificador externo (card_id)
func (r *clienteRepository) BuscarPorIdentificadorExterno(identificadorExterno string) (*Cliente, error) {
	query := `
		SELECT 
			gcl_cli_int, gcl_cli_ext, gcl_cli_nom, gcl_cli_ema, 
			gcl_cli_pat, gcl_cli_tso, gcl_cli_stc, gcl_cli_npr, 
			gcl_cli_dcr, COALESCE(gcl_cli_dat, gcl_cli_dcr)
		FROM gestao_clientes.cliente
		WHERE gcl_cli_ext = $1
	`

	var cliente Cliente

	err := r.banco.QueryRow(query, identificadorExterno).Scan(
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
			return nil, fmt.Errorf("cliente não encontrado com identificador externo %s", identificadorExterno)
		}
		return nil, fmt.Errorf("erro ao buscar cliente por identificador externo: %w", err)
	}

	return &cliente, nil
}

// Atualizar atualiza um cliente existente no banco de dados
func (r *clienteRepository) Atualizar(cliente Cliente) error {
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

	resultado, err := r.banco.Exec(
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
		return fmt.Errorf("erro ao atualizar cliente: %w", err)
	}

	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if linhasAfetadas == 0 {
		return fmt.Errorf("nenhum cliente encontrado para atualizar com identificador interno %d", cliente.IdentificadorInterno)
	}

	return nil
}
