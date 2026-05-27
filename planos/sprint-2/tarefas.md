# Tarefas - Sprint 2: API de Criação de Cliente

## 1. Validação de Dados (Feature Folder: gestao_clientes)

- [ ] Implementar internal/gestao_clientes/validator.go com validador de payload de cliente
- [ ] Validar campos obrigatórios (cliente_nome, cliente_email, tipo_solicitacao, valor_patrimonio)
- [ ] Validar formato de cliente_email
- [ ] Validar valor_patrimonio positivo
- [ ] Criar mensagens de erro em Português

## 2. Camada de Serviço (Feature Folder: gestao_clientes)

- [ ] Implementar internal/gestao_clientes/service.go com ClienteService.criarCliente()
- [ ] Implementar lógica de status_cliente inicial "Aguardando Análise"
- [ ] Implementar geração de identificador_card simulado
- [ ] Implementar regras de negócio de criação

## 3. Camada de Apresentação (Feature Folder: gestao_clientes)

- [ ] Implementar internal/gestao_clientes/controller.go com ClienteController
- [ ] Criar endpoint POST /clientes
- [ ] Implementar parsing de JSON
- [ ] Implementar respostas HTTP (201, 400)

## 4. Integração Pipefy (Feature Folder: integracao_pipefy)

- [ ] Pesquisar documentação oficial do Pipefy para createCard
- [ ] Implementar internal/integracao_pipefy/client.go com PipefyGraphQLClient
- [ ] Implementar internal/integracao_pipefy/mutations.go com createCard
- [ ] Estruturar mutation createCard conforme especificação
- [ ] Adicionar comentários com fonte da especificação
- [ ] Simular envio (sem requisição real)
