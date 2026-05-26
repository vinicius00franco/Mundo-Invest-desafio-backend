# TESTE TÉCNICO: Desenvolvedor Backend (Client Management & Pipefy Integration)

## 1. O Contexto e o Problema

Você passará a desenvolver o esqueleto de um sistema interno para o Mundo Invest. A aplicação deve gerenciar clientes e seus respectivos patrimônios investidos, além de mapear essas ações para o Pipefy (nossa ferramenta de controle de processos).

Em vez de se conectar aos servidores reais do Pipefy, sua aplicação fará a persistência em um banco de dados local, mas o formato das queries e mutations GraphQL no seu código deve seguir rigorosamente a especificação da documentação oficial do Pipefy.

## 2. Requisitos Funcionais e Fluxo da Aplicação

Sua API (desenvolvida em Python ou Golang) deve expor dois fluxos principais:

### 2.1 Fluxo 1: Criação de Cliente e Mapeamento de Card

**Endpoint:** `POST /clientes`

**Payload de Entrada esperado:**
```json
{
  "cliente_nome": "João Silva",
  "cliente_email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização cadastral",
  "valor_patrimonio": 250000
}
```

**O que a aplicação deve fazer:**

1. **Validação:** Validar se os campos obrigatórios estão presentes e se o e-mail é válido.
2. **Persistência Local:** Salvar o cliente em um banco de dados local (SQLite ou PostgreSQL via Docker) com o status inicial "Aguardando Análise".
3. **Mapeamento Pipefy (GraphQL):** Você deve pesquisar na documentação pública do Pipefy como é estruturada a mutation de criação de card (createCard). No seu código, na camada de client/serviço, você deve deixar estruturada essa string/payload exata da mutation (com as variáveis corretas de nome, e-mail e patrimônio), simulando o envio para o Pipefy.

### 2.2 Fluxo 2: Atualização de Card (Simulação de Webhook)

**Endpoint:** `POST /webhooks/pipefy/card-updated`

Esse endpoint simula o momento em que o time operacional mexe no card dentro do Pipefy e dispara um aviso de volta para o nosso sistema.

**Payload de Entrada esperado:**
```json
{
  "event_id": "evt_123",
  "card_id": "card_456",
  "cliente_email": "joao.silva@example.com",
  "timestamp": "2026-05-18T12:00:00Z"
}
```

**O que a aplicação deve fazer:**

1. **Idempotência:** Verificar pelo event_id se esse evento já foi processado para evitar duplicidade.
2. **Regra de Negócio:** Buscar o cliente no banco local usando o cliente_email.

   - Se valor_patrimonio for maior ou igual a 200.000, defina a prioridade como prioridade_alta.
   - Se valor_patrimonio for menor que 200.000, defina a prioridade como prioridade_normal.

3. **Mapeamento de Update (GraphQL):** Você deve pesquisar na documentação pública do Pipefy qual é a mutation exata para atualizar os campos de um card (updateCardField ou similar). Deixe estruturada no seu código a mutation real que enviaria ao Pipefy o novo status do cliente ("Processado") e a prioridade calculada.

4. **Atualizar Banco Local:** Mudar o status do cliente no banco local para "Processado" e salvar a prioridade definida.

## 3. Testes Obrigatórios

Implemente testes automatizados cobrindo:

1. Criação de cliente com payload válido e salvamento no banco.
2. Processamento do webhook aplicando a regra de prioridade correta com base no patrimônio.
3. Bloqueio de processamento caso o event_id do webhook seja duplicado.

## 4. O que deve conter no seu README.md

1. Instruções de execução local do projeto e dos testes.
2. Exemplos de requisição (curl) para os dois endpoints.
3. Visão de Produção (AWS): Explique textualmente como essa estrutura de banco de dados + processamento de webhook escalaria na AWS (usando serviços como Lambda, API Gateway e DynamoDB/RDS) - Opcional.

## 5. Diretrizes para o Vídeo de Defesa (Obrigatório)

Seu vídeo deve ter no máximo 7 minutos e demonstrar:

1. Uma breve apresentação do seu background técnico.
2. Explicação da estrutura de pastas e camadas que você escolheu para isolar as regras de negócio dos clientes.
3. O ponto alto: Mostrar no código onde você estruturou as duas mutations GraphQL e explicar brevemente como localizou a sintaxe delas na documentação do Pipefy.
4. Demonstração dos testes rodando no terminal e execução dos endpoints.

## 6. Tempo para o desenvolvimento

O teste deve ser desenvolvido e enviado até o dia 29 de Maio.

## 7. Como deve ser entregue

Preencher o form: https://docs.google.com/forms/d/e/1FAIpQLSfRhrxW4B2JGreNkrNCXlKJA9eCONGvL2JaRuQFelQz7PTkw/viewform?usp=publish-editor com as informações solicitadas.
