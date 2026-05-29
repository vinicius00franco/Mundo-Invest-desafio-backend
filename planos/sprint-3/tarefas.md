# Tarefas - Sprint 3: API de Webhook

## 1. Validação de Dados (Feature Folder: processamento_eventos)

- [ ] Implementar internal/processamento_eventos/validator.go com validador de payload de webhook
- [ ] Validar campos obrigatórios (identificador_evento, identificador_card, cliente_email, data_evento)
- [ ] Validar formato de cliente_email
- [ ] Criar mensagens de erro em Português

## 2. Camada de Serviço (Feature Folder: processamento_eventos)

- [ ] Implementar internal/processamento_eventos/service.go com WebhookService.processarWebhook()
- [ ] Implementar verificação de idempotencia por identificador_evento
- [ ] Implementar busca de cliente por cliente_email
- [ ] Implementar tratamento de cliente não encontrado

## 3. Cálculo de Prioridade (Feature Folder: dominio)

- [ ] Implementar internal/dominio/prioridade_calculator.go com PrioridadeCalculator
- [ ] Implementar regra: valor_patrimonio >= 200.000 → nivel_prioridade_alta
- [ ] Implementar regra: valor_patrimonio < 200.000 → nivel_prioridade_normal
- [ ] Implementar regra: valor_patrimonio == 200.000 → nivel_prioridade_alta

## 4. Atualização de Cliente (Feature Folder: processamento_eventos)

- [ ] Implementar atualização de status_cliente para "Processado"
- [ ] Implementar salvamento de nivel_prioridade calculado
- [ ] Implementar persistência de evento com foi_processado

## 5. Integração Pipefy (Feature Folder: integracao_pipefy)

- [ ] Pesquisar documentação oficial do Pipefy para updateCard
- [ ] Implementar internal/integracao_pipefy/mutations.go com updateCard
- [ ] Estruturar mutation updateCard conforme especificação
- [ ] Adicionar comentários com fonte da especificação
- [ ] Simular envio (sem requisição real)

## 6. Camada de Apresentação (Feature Folder: processamento_eventos)

- [ ] Implementar internal/processamento_eventos/controller.go com WebhookController
- [ ] Criar endpoint POST /webhooks/pipefy/card-updated
- [ ] Implementar parsing de JSON
- [ ] Implementar respostas HTTP (200, 400, 404)
