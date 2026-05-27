package mensagens

import (
	"sync"
	"testing"
)

// BenchmarkSingletonVsNew compara performance do singleton vs criar nova instância
func BenchmarkSingletonVsNew(b *testing.B) {
	b.Run("Singleton", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ObterCatalogo()
		}
	})

	b.Run("NovaInstancia", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NovoCatalogo()
		}
	})
}

// BenchmarkTextoAcesso benchmark de acesso a mensagens
func BenchmarkTextoAcesso(b *testing.B) {
	catalogo := ObterCatalogo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = catalogo.Texto(CliCriadoSucesso)
	}
}

// BenchmarkTextoFormatado benchmark de acesso formatado
func BenchmarkTextoFormatado(b *testing.B) {
	catalogo := ObterCatalogo()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = catalogo.TextoFormatado(ErrClienteNaoEncontrado, "Teste")
	}
}

// BenchmarkConcorrencia benchmark de acesso concorrente
func BenchmarkConcorrencia(b *testing.B) {
	catalogo := ObterCatalogo()
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = catalogo.Texto(CliCriadoSucesso)
		}()
	}
	wg.Wait()
}

// BenchmarkMapeadorHTTP benchmark do mapeador HTTP
func BenchmarkMapeadorHTTP(b *testing.B) {
	mapeador := NovoMapeadorHTTP()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapeador.StatusPara(TipoValidacao)
	}
}
