# Docker Compose - Sprint 1: Fundamentos e Estrutura

## docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: postgres_container
    environment:
      POSTGRES_DB: mundo_invest
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

## Instruções de Uso

### Iniciar PostgreSQL
```bash
docker-compose up -d
```

### Parar e remover containers
```bash
docker-compose down
```

### Ver logs do PostgreSQL
```bash
docker-compose logs postgres
```

### Acessar banco via psql
```bash
docker exec -it postgres_container psql -U postgres -d mundo_invest
```

### Verificar status dos containers
```bash
docker-compose ps
```

### Reiniciar o PostgreSQL
```bash
docker-compose restart postgres
```

## Configurações

### Variáveis de Ambiente
- `POSTGRES_DB`: Nome do banco de dados (mundo_invest)
- `POSTGRES_USER`: Usuário padrão (postgres)
- `POSTGRES_PASSWORD`: Senha do usuário (postgres)

### Portas
- `5432`: Porta padrão do PostgreSQL mapeada para o host

### Volumes
- `postgres_data`: Volume persistente para dados do PostgreSQL
- `./migrations`: Diretório de migrations montado para inicialização automática

### Healthcheck
- Verifica se o PostgreSQL está pronto a cada 10 segundos
- Timeout de 5 segundos por verificação
- 5 tentativas antes de marcar como unhealthy
