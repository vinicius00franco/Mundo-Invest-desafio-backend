package processamento_eventos

import (
	"context"
	"database/sql"
	"time"

	"github.com/MundoInvest/backend/internal/shared/config"
	"github.com/MundoInvest/backend/internal/shared/errors"
)

// Evento representa a entidade Evento do contexto de processamento de eventos
type Evento struct {
	IdentificadorInterno int64      // pev_eve_int
	IdentificadorEvento  string     // pev_eve_ide
	IdentificadorCard    string     // pev_eve_idc
	EmailCliente         string     // pev_eve_ema
	TimestampEvento      time.Time  // pev_eve_tms
	FoiProcessado        bool       // pev_eve_fpr
	DataCriacao          time.Time  // pev_eve_dcr
	DataAtualizacao      *time.Time // pev_eve_dat (nullable)
}

// EventoRepository define a interface para operações de persistência de eventos
type EventoRepository interface {
	Salvar(ctx context.Context, evento Evento) (*Evento, error)
	BuscarPorIdentificadorEvento(ctx context.Context, identificadorEvento string) (*Evento, error)
	VerificarFoiProcessado(ctx context.Context, identificadorEvento string) (bool, error)
	BuscarPorIdentificadorCard(ctx context.Context, identificadorCard string) ([]Evento, error)
}

// eventoRepository implementa a interface EventoRepository
type eventoRepository struct {
	banco  *sql.DB
	config *config.Config
}

// NovoEventoRepository cria uma nova instância de EventoRepository
func NovoEventoRepository(banco *sql.DB, cfg *config.Config) EventoRepository {
	return &eventoRepository{
		banco:  banco,
		config: cfg,
	}
}

// Salvar insere um novo evento no banco de dados
func (r *eventoRepository) Salvar(ctx context.Context, evento Evento) (*Evento, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		INSERT INTO processamento_eventos.evento (
			pev_eve_ide, pev_eve_idc, pev_eve_ema, pev_eve_tms, pev_eve_fpr
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING pev_eve_int, pev_eve_dcr
	`

	var identificadorInterno int64
	var dataCriacao time.Time

	err := r.banco.QueryRowContext(
		ctx,
		query,
		evento.IdentificadorEvento,
		evento.IdentificadorCard,
		evento.EmailCliente,
		evento.TimestampEvento,
		evento.FoiProcessado,
	).Scan(&identificadorInterno, &dataCriacao)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("salvar evento", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("salvar evento", err)
	}

	evento.IdentificadorInterno = identificadorInterno
	evento.DataCriacao = dataCriacao

	return &evento, nil
}

// BuscarPorIdentificadorEvento busca um evento pelo identificador de evento
func (r *eventoRepository) BuscarPorIdentificadorEvento(ctx context.Context, identificadorEvento string) (*Evento, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		SELECT 
			pev_eve_int, pev_eve_ide, pev_eve_idc, pev_eve_ema, 
			pev_eve_tms, pev_eve_fpr, pev_eve_dcr, pev_eve_dat
		FROM processamento_eventos.evento
		WHERE pev_eve_ide = $1
	`

	var evento Evento
	var dataAtualizacao sql.NullTime

	err := r.banco.QueryRowContext(ctx, query, identificadorEvento).Scan(
		&evento.IdentificadorInterno,
		&evento.IdentificadorEvento,
		&evento.IdentificadorCard,
		&evento.EmailCliente,
		&evento.TimestampEvento,
		&evento.FoiProcessado,
		&evento.DataCriacao,
		&dataAtualizacao,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("evento", identificadorEvento)
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("buscar evento por identificador", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("buscar evento por identificador", err)
	}

	if dataAtualizacao.Valid {
		evento.DataAtualizacao = &dataAtualizacao.Time
	}

	return &evento, nil
}

// VerificarFoiProcessado verifica se um evento já foi processado (idempotência)
func (r *eventoRepository) VerificarFoiProcessado(ctx context.Context, identificadorEvento string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		SELECT pev_eve_fpr
		FROM processamento_eventos.evento
		WHERE pev_eve_ide = $1
	`

	var foiProcessado bool

	err := r.banco.QueryRowContext(ctx, query, identificadorEvento).Scan(&foiProcessado)
	if err != nil {
		if err == sql.ErrNoRows {
			// Evento não existe, então não foi processado
			return false, nil
		}
		if ctx.Err() == context.DeadlineExceeded {
			return false, errors.NewTimeoutError("verificar se evento foi processado", "timeout ao executar operação de banco de dados")
		}
		return false, errors.NewRepositoryError("verificar se evento foi processado", err)
	}

	return foiProcessado, nil
}

// BuscarPorIdentificadorCard busca todos os eventos relacionados a um card
func (r *eventoRepository) BuscarPorIdentificadorCard(ctx context.Context, identificadorCard string) ([]Evento, error) {
	ctx, cancel := context.WithTimeout(ctx, r.config.Timeouts.Database)
	defer cancel()

	query := `
		SELECT 
			pev_eve_int, pev_eve_ide, pev_eve_idc, pev_eve_ema, 
			pev_eve_tms, pev_eve_fpr, pev_eve_dcr, pev_eve_dat
		FROM processamento_eventos.evento
		WHERE pev_eve_idc = $1
		ORDER BY pev_eve_dcr DESC
	`

	linhas, err := r.banco.QueryContext(ctx, query, identificadorCard)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("buscar eventos por identificador de card", "timeout ao executar operação de banco de dados")
		}
		return nil, errors.NewRepositoryError("buscar eventos por identificador de card", err)
	}
	defer linhas.Close()

	var eventos []Evento

	for linhas.Next() {
		var evento Evento
		var dataAtualizacao sql.NullTime
		err := linhas.Scan(
			&evento.IdentificadorInterno,
			&evento.IdentificadorEvento,
			&evento.IdentificadorCard,
			&evento.EmailCliente,
			&evento.TimestampEvento,
			&evento.FoiProcessado,
			&evento.DataCriacao,
			&dataAtualizacao,
		)
		if err != nil {
			return nil, errors.NewRepositoryError("escanear evento", err)
		}

		if dataAtualizacao.Valid {
			evento.DataAtualizacao = &dataAtualizacao.Time
		}

		eventos = append(eventos, evento)
	}

	if err = linhas.Err(); err != nil {
		return nil, errors.NewRepositoryError("iterar sobre eventos", err)
	}

	return eventos, nil
}
