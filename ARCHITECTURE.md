# Arquitetura

## Decisão central

**Monolito modular sobre um único Postgres.** Um binário, um deploy, um
banco. Cada módulo é um contexto de domínio (bounded context) com fronteira
explícita, para que as partes críticas possam virar serviço depois sem
reescrever o resto.

Começamos monolito porque o projeto tem 1 time, 1 VM e carga de uso
familiar. Microserviço agora compraria rede, deploy e consistência
distribuída sem nenhum problema real para resolver.

---

## Estrutura

```
cmd/api/                  entrypoint: config, pool, servidor, shutdown
  main.go
  router.go               único lugar onde módulos são plugados
internal/
  config/                 leitura e validação das variáveis de ambiente
  platform/               infraestrutura compartilhada, sem domínio
    postgres/             pool, DBTX, query builder, helpers de erro
    httpx/                JSON, respostas de erro, middlewares
  health/                 /healthz e /readyz
  identity/               usuários (módulo de referência)
    model.go              tipos e erros do domínio
    service.go            regras de negócio, API pública do módulo
    repository.go         SQL, funções não exportadas
    handler.go            HTTP: decodifica, chama service, mapeia erro
    routes.go             registra as rotas do módulo
migrations/               SQL versionado (golang-migrate)
```

Novos módulos seguem o mesmo formato de `identity`: `family`, `media`,
`stream`, `room`.

---

## Regras de fronteira

Estas regras são o que mantém o monolito extraível:

1. **Um módulo só fala com outro pelas funções exportadas do `service.go`.**
   Nunca importa o repository nem consulta a tabela do outro.
2. **Cada tabela tem um dono.** `users` é de `identity`. Outro módulo
   guarda `user_id` como referência e pergunta a `identity` quando precisa.
3. **Repository é não exportado.** O compilador impede o atalho da regra 1.
4. **`platform` não importa nenhum módulo.** A dependência só desce:
   `cmd` → módulos → `platform`.
5. **Foreign key entre módulos é permitida enquanto for um banco só**, mas
   anote no SQL. Na extração, ela vira validação na aplicação.

---

## Estilo de código

- **Funções em vez de structs com estado.** Dependências entram por
  parâmetro (`db postgres.DBTX`), handlers são closures sobre elas.
- **`postgres.DBTX`** é satisfeito por `*pgxpool.Pool` e `pgx.Tx`. A mesma
  função de repository roda dentro ou fora de transação. Para transação,
  use `pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error { ... })`.
- **Lógica pura separada de I/O.** Validação e normalização (ex.:
  `normalizeNewUser`) não tocam banco e são testadas sem container.
- **Erros de domínio são sentinelas** (`ErrNotFound`, `ErrEmailTaken`). O
  repository traduz erro do Postgres para eles, o handler traduz para HTTP.
- Interfaces pequenas, declaradas por quem consome (ex.: `health.Pinger`).

---

## Decisões técnicas

| Tema | Escolha | Por quê |
|---|---|---|
| HTTP | `net/http` (Go 1.22+) | Rotas com método e `{id}` nativos, sem framework |
| Banco | Postgres 18 | Relacional cobre usuários, famílias, mídias, histórico |
| Driver | `jackc/pgx/v5` + `pgxpool` | Driver mais rápido e mantido, pool próprio |
| Queries | `Masterminds/squirrel` | Filtros opcionais sem concatenar string |
| IDs | `uuidv7()` do Postgres 18 | Ordenável por tempo, bom para índice B-tree |
| Migrations | `golang-migrate` via Docker | SQL puro, sem instalar CLI |
| Logs | `log/slog` | Estruturado e da stdlib |
| Config | variáveis de ambiente | `.env` só em desenvolvimento |

Decisões que ficam para quando o módulo chegar:

- **Storage:** Backblaze B2 via `aws-sdk-go-v2` (API S3).
- **Auth:** JWT (`golang-jwt/jwt`) e `bcrypt` (`golang.org/x/crypto`).
- **Tempo real:** WebSocket (`coder/websocket` ou `gorilla/websocket`).
- **Transcodificação:** FFmpeg chamado como processo externo.

### Servidor HTTP

O `http.Server` não define `WriteTimeout` de propósito. Ele cortaria
streaming de vídeo e WebSocket longos. Rotas comuns que precisarem de prazo
usam `http.ResponseController.SetWriteDeadline`.

O shutdown é gracioso: `SIGINT`/`SIGTERM` param de aceitar conexões e
esperam até 10s pelas requisições em andamento.

### Health

- `GET /healthz` é liveness: só diz que o processo responde. Não toca o
  banco, para uma queda do Postgres não reiniciar um processo saudável.
- `GET /readyz` é readiness: faz ping no banco e devolve 503 se falhar.

---

## Caminho de evolução

Extrair um módulo vale a pena quando ele tem **perfil de recurso ou de
deploy diferente do resto**. Tamanho de código não é critério.

| Módulo | Sinal para extrair | Como |
|---|---|---|
| `media` (transcode) | FFmpeg disputando CPU com a API | Worker separado lendo uma fila; começa como tabela de jobs no Postgres |
| `room` (watch party) | Muitas conexões WebSocket simultâneas | Serviço próprio; Elixir/Phoenix é candidato pelo modelo de atores |
| `stream` | Banda de saída saturando a VM | Servir segmentos HLS direto do B2/CDN com URL assinada |

Passos de uma extração, possíveis porque as regras de fronteira valem
desde o início:

1. As chamadas ao `service.go` do módulo viram cliente HTTP/gRPC com a
   mesma assinatura.
2. As tabelas do módulo migram para um schema e depois um banco próprio.
3. FKs entre módulos viram validação na aplicação.

---

## Fluxos principais

### Upload de mídia
1. Cliente envia arquivo(s) via `POST /media/upload`.
2. Backend valida tipo e tamanho e faz stream para o B2.
3. Backend grava a mídia e um job de transcodificação no Postgres, na
   mesma transação.
4. Worker gera HLS (segmentos + `.m3u8`) com FFmpeg e sobe para o B2.

### Streaming
1. Cliente pede `GET /media/{id}/stream` com JWT.
2. Backend devolve a playlist `.m3u8`.
3. Player busca os segmentos, com suporte a HTTP Range Requests.

### Sala de cinema
1. Host cria a sala (`POST /rooms`), membros entram (`POST /rooms/{id}/join`).
2. Todos conectam em `WS /rooms/{id}/sync`.
3. Host emite play/pause/seek, servidor faz broadcast, membros ressincronizam
   quando a diferença passa do limite configurado.

---

## Variáveis de ambiente

| Variável | Uso |
|---|---|
| `HTTP_ADDR` | Endereço do servidor, padrão `:4222` |
| `DATABASE_URL` | Conexão do app com o Postgres |
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT` | Criação do container no `docker compose` |
