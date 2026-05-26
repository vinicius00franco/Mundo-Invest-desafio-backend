# Cenários de Falha - Criação de Cliente

## Cenário: Criar cliente sem campos obrigatórios

**Descrição**: Tentar criar cliente sem fornecer todos os campos obrigatórios.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes sem o campo cliente_nome:
   ```json
   {
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Cliente não é salvo no banco

---

## Cenário: Criar cliente sem tipo de solicitação

**Descrição**: Tentar criar cliente sem fornecer o tipo de solicitação.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes sem o campo tipo_solicitacao:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "joao.silva@example.com",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Cliente não é salvo no banco

---

## Cenário: Criar cliente sem valor de patrimônio

**Descrição**: Tentar criar cliente sem fornecer o valor de patrimônio.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes sem o campo valor_patrimonio:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral"
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Cliente não é salvo no banco

---

## Cenário: Criar cliente sem e-mail

**Descrição**: Tentar criar cliente sem fornecer o e-mail.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes sem o campo cliente_email:
   ```json
   {
     "cliente_nome": "João Silva",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo obrigatório faltando
- Cliente não é salvo no banco

---

## Cenário: Criar cliente com erro de banco de dados

**Descrição**: Tentar criar cliente quando banco de dados está indisponível.

**Pré-condições**:
- Endpoint POST /clientes está disponível
- Banco de dados está indisponível ou em erro

**Passos**:
1. Enviar requisição POST para /clientes com payload válido:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 500 Internal Server Error
- Mensagem de erro indica problema no banco de dados
- Cliente não é salvo no banco

**Dados de Verificação**:
- Response contém erro "ERRO_BANCO_DADOS"
- Sistema permanece em estado consistente
