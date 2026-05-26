# Cenários de Validação - Criação de Cliente

## Cenário: Criar cliente com e-mail inválido

**Descrição**: Tentar criar cliente com e-mail em formato inválido.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes com e-mail inválido:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "email-invalido",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica e-mail inválido
- Cliente não é salvo no banco

---

## Cenário: Criar cliente com patrimônio negativo

**Descrição**: Tentar criar cliente com valor de patrimônio negativo.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes com valor_patrimonio negativo:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": -1000
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica valor inválido
- Cliente não é salvo no banco

---

## Cenário: Criar cliente com patrimônio zero

**Descrição**: Tentar criar cliente com valor de patrimônio igual a zero.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes com valor_patrimonio zero:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 0
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica valor inválido
- Cliente não é salvo no banco

---

## Cenário: Criar cliente com nome vazio

**Descrição**: Tentar criar cliente com nome em branco.

**Pré-condições**:
- Endpoint POST /clientes está disponível

**Passos**:
1. Enviar requisição POST para /clientes com cliente_nome vazio:
   ```json
   {
     "cliente_nome": "",
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 400 Bad Request
- Mensagem de erro indica campo inválido
- Cliente não é salvo no banco
