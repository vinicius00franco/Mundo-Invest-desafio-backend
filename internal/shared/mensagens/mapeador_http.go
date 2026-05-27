package mensagens

import "net/http"

// MapeadorHTTP separa responsabilidade de mapeamento de status HTTP
type MapeadorHTTP struct {
	statusPorTipo map[TipoMensagem]int
}

// NovoMapeadorHTTP cria um novo mapeador HTTP
func NovoMapeadorHTTP() *MapeadorHTTP {
	return &MapeadorHTTP{
		statusPorTipo: map[TipoMensagem]int{
			TipoSucesso:    http.StatusCreated, // ou http.StatusOK dependendo do contexto
			TipoValidacao:  http.StatusBadRequest,
			TipoServico:    http.StatusInternalServerError,
			TipoIntegracao: http.StatusBadGateway,
			TipoHTTP:       http.StatusBadRequest,
			TipoBancoDados: http.StatusInternalServerError,
		},
	}
}

// StatusPara retorna o status HTTP para um tipo de mensagem
func (m *MapeadorHTTP) StatusPara(tipo TipoMensagem) int {
	if status, exists := m.statusPorTipo[tipo]; exists {
		return status
	}
	return http.StatusInternalServerError // Fallback seguro
}
