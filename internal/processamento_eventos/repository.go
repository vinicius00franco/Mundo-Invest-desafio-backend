package processamento_eventos

import (
	"database/sql"
	"fmt"
	"time"
)

// Evento representa a entidade Evento do contexto de processamento de eventos
type Evento struct {
	IdentificadorInterno int64     // pev_eve_int
	IdentificadorEvento  string    // pev_eve_ide
	IdentificadorCard    string    // pev_eve_idc
	EmailCliente         string    // pev_eve_ema
	TimestampEvento      time.Time // pev_eve_tms
	FoiProcessado        bool      // pev_eve_fpr
	DataCriacao          time.Time // pev_eve_dcr
	DataAtualizacao      time.Time // pev_eve_dat
}

// EventoRepository define a interface para operações de persistência de eventos
type EventoRepository interface {
	Salvar(evento Evento) (*Evento, error)
	BuscarPorIdentificadorEvento(identificadorEvento string) (*Evento, error)
	VerificarFoiProcessado(identificadorEvento string) (bool, error)
	BuscarPorIdentificadorCard(identificadorCard string) ([]Evento, error)
}

// eventoRepository implementa a interface EventoRepository
type eventoRepository struct {
	banco *sql.DB
}

// NovoEventoRepository cria uma nova instância de EventoRepository
func NovoEventoRepository(banco *sql.DB) EventoRepository {
	return &eventoRepository{banco: banco}
}

// Salvar insere um novo evento no banco de dados
func (r *eventoRepository) Salvar(evento Evento) (*Evento, error) {
	query := `
		INSERT INTO processamento_eventos.evento (
			pev_eve_ide, pev_eve_idc, pev_eve_ema, pev_eve_tms, pev_eve_fpr
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING pev_eve_int, pev_eve_dcr
	`

	var identificadorInterno int64
	var dataCriacao time.Time

	err := r.banco.QueryRow(
		query,
		evento.IdentificadorEvento,
		evento.IdentificadorCard,
		evento.EmailCliente,
		evento.TimestampEvento,
		evento.FoiProcessado,
	).Scan(&identificadorInterno, &dataCriacao)

	if err != nil {
		return nil, fmt.Errorf("erro ao salvar evento: %w", err)
	}

	evento.IdentificadorInterno = identificadorInterno
	evento.DataCriacao = dataCriacao

	return &evento, nil
}

// BuscarPorIdentificadorEvento busca um evento pelo identificador de evento
func (r *eventoRepository) BuscarPorIdentificadorEvento(identificadorEvento string) (*Evento, error) {
	query := `
		SELECT 
			pev_eve_int, pev_eve_ide, pev_eve_idc, pev_eve_ema, 
			pev_eve_tms, pev_eve_fpr, pev_eve_dcr, pev_eve_dat
		FROM processamento_eventos.evento
		WHERE pev_eve_ide = $1
	`

	var evento Evento

	err := r.banco.QueryRow(query, identificadorEvento).Scan(
		&evento.IdentificadorInterno,
		&evento.IdentificadorEvento,
		&evento.IdentificadorCard,
		&evento.EmailCliente,
		&evento.TimestampEvento,
		&evento.FoiProcessado,
		&evento.DataCriacao,
		&evento.DataAtualizacao,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("evento não encontrado com identificador %s", identificadorEvento)
		}
		return nil, fmt.Errorf("erro ao buscar evento por identificador: %w", err)
	}

	return &evento, nil
}

// VerificarFoiProcessado verifica se um evento já foi processado (idempotência)
func (r *eventoRepository) VerificarFoiProcessado(identificadorEvento string) (bool, error) {
	query := `
		SELECT pev_eve_fpr
		FROM processamento_eventos.evento
		WHERE pev_eve_ide = $1
	`

	var foiProcessado bool

	err := r.banco.QueryRow(query, identificadorEvento).Scan(&foiProcessado)
	if err != nil {
		if err == sql.ErrNoRows {
			// Evento não existe, então não foi processado
			return false, nil
		}
		return false, fmt.Errorf("erro ao verificar se evento foi processado: %w", err)
	}

	return foiProcessado, nil
}

// BuscarPorIdentificadorCard busca todos os eventos relacionados a um card
func (r *eventoRepository) BuscarPorIdentificadorCard(identificadorCard string) ([]Evento, error) {
	query := `
		SELECT 
			pev_eve_int, pev_eve_ide, pev_eve_idc, pev_eve_ema, 
			pev_eve_tms, pev_eve_fpr, pev_eve_dcr, pev_eve_dat
		FROM processamento_eventos.evento
		WHERE pev_eve_idc = $1
		ORDER BY pev_eve_dcr DESC
	`

	linhas, err := r.banco.Query(query, identificadorCard)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar eventos por identificador de card: %w", err)
	}
	defer linhas.Close()

	var eventos []Evento

	for linhas.Next() {
		var evento Evento
		err := linhas.Scan(
			&evento.IdentificadorInterno,
			&evento.IdentificadorEvento,
			&evento.IdentificadorCard,
			&evento.EmailCliente,
			&evento.TimestampEvento,
			&evento.FoiProcessado,
			&evento.DataCriacao,
			&evento.DataAtualizacao,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear evento: %w", err)
		}
		eventos = append(eventos, evento)
	}

	if err = linhas.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre eventos: %w", err)
	}

	return eventos, nil
}
