# Cenários de Sucesso - Criação de Cliente

## Cenário: Criar cliente com dados válidos

**Descrição**: Criar um novo cliente com todos os campos obrigatórios preenchidos corretamente.

**Pré-condições**:
- Endpoint POST /clientes está disponível
- Banco de dados está acessível

**Passos**:
1. Enviar requisição POST para /clientes com payload:
   ```json
   {
     "cliente_nome": "João Silva",
     "cliente_email": "joao.silva@example.com",
     "tipo_solicitacao": "Atualização cadastral",
     "valor_patrimonio": 250000
   }
   ```

**Resultados Esperados**:
- HTTP 201 Created
- Cliente salvo no banco de dados
- Status do cliente é "Aguardando Análise"
- Card_id é gerado
- Mutation GraphQL createCard é estruturada no código
- Response contém dados do cliente criado

**Dados de Verificação**:
- ID do cliente gerado (UUID)
- Nome: "João Silva"
- Email: "joao.silva@example.com"
- Tipo de solicitação: "Atualização cadastral"
- Valor patrimônio: 250000
- Status: "Aguardando Análise"
- Card_id presente
