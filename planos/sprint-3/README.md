# Sprint 3: API de Webhook (2 dias)

## Objetivo

Implementar endpoint POST /webhooks/pipefy/card-updated seguindo linguagem ubíqua

## Visão Geral

Esta sprint implementa o processamento de webhooks do Pipefy, que permite que o sistema receba atualizações sobre cards. Inclui:

- Validação de dados de webhook
- Implementação de idempotência
- Cálculo automático de prioridade
- Atualização de clientes
- Integração com Pipefy para updateCard

## Contexto

Esta sprint é crítica pois implementa o fluxo de processamento de eventos do Pipefy, permitindo que o sistema reaja a mudanças nos cards de forma idempotente e confiável. O cálculo de prioridade é uma regra de negócio importante que impacta o tratamento dos clientes.

## Arquivos Relacionados

- [tarefas.md](./tarefas.md) - Lista detalhada de tarefas
- [criterios-aceite.md](./criterios-aceite.md) - Critérios de aceite da sprint
