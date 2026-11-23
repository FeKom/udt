# 🎬 Personal Media Server

Servidor de mídia pessoal e familiar com streaming de vídeo, upload de arquivos e sala de cinema sincronizada.

---

## Stack

| Camada | Tecnologia |
|---|---|
| Backend | Go (monolito modular) |
| Frontend | Next.js / TypeScript |
| Storage | Backblaze B2 |
| Banco de dados | PostgreSQL 18 |
| Tempo real | WebSockets |

---

## Arquitetura

Detalhes, regras de fronteira entre módulos e caminho de evolução em
[`ARCHITECTURE.md`](ARCHITECTURE.md). Plano de entregas em
[`ROADMAP.md`](ROADMAP.md).

- Usuários e famílias com histórico e salvos individuais
- Upload com stream do backend para o Backblaze B2
- Metadados salvos no PostgreSQL
- Streaming via HLS com suporte a múltiplos usuários simultâneos
- Sala de cinema sincronizada via WebSockets (watch party)
- Upload de um ou mais arquivos por vez

---

## Rodando localmente

Requer Go 1.26+ e Docker.

```bash
cp .env.example .env
make up      # sobe o Postgres na porta 5433 e aplica as migrations
make run     # API em http://localhost:4222
make test
```

Endpoints disponíveis:

| Rota | Descrição |
|---|---|
| `GET /healthz` | Liveness |
| `GET /readyz` | Readiness, checa o Postgres |
| `POST /users` | Cria usuário (`{"name","email"}`) |
| `GET /users?name=&email=&limit=` | Lista com filtros opcionais |
| `GET /users/{id}` | Busca por id |

Nova migration: `make migrate-new name=create_families`.

---

## Libs Go

### HTTP / Router
- **net/http** (stdlib, Go 1.22+) — https://pkg.go.dev/net/http

### Storage (Backblaze B2 — S3 compatible)
- **aws-sdk-go-v2** — https://github.com/aws/aws-sdk-go-v2
- **Backblaze S3-compatible API** — https://www.backblaze.com/docs/cloud-storage-s3-compatible-api

### Banco de dados
- **pgx** — https://github.com/jackc/pgx
- **squirrel** (query builder) — https://github.com/Masterminds/squirrel
- **golang-migrate** — https://github.com/golang-migrate/migrate

### Auth
- **golang-jwt** — https://github.com/golang-jwt/jwt

### WebSockets
- **gorilla/websocket** — https://github.com/gorilla/websocket
- **nhooyr.io/websocket** — https://pkg.go.dev/nhooyr.io/websocket *(alternativa mais moderna)*

### Variáveis de ambiente
- **godotenv** — https://github.com/joho/godotenv

### Transcoding / FFmpeg
- **ffmpeg-go** — https://github.com/u2takey/ffmpeg-go

---

## Leituras — Streaming de Vídeo

### HTTP Range Requests
Conceito central para seek de vídeo (avançar/voltar) sem baixar o arquivo inteiro.

- MDN — https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
- RFC 7233 (spec oficial) — https://www.rfc-editor.org/rfc/rfc7233

### HLS — HTTP Live Streaming
Protocolo usado pela maioria dos streamings. Suporta múltiplos usuários simultâneos e qualidade adaptativa.

- Apple HLS — https://developer.apple.com/documentation/http-live-streaming
- RFC 8216 (spec oficial) — https://www.rfc-editor.org/rfc/rfc8216

### FFmpeg
Necessário para transcodar vídeos para o formato HLS.

- Documentação oficial — https://ffmpeg.org/documentation.html

### Referências de empresas reais
- **AWS / Prime Video** — ABR Streaming (Adaptive Bitrate): https://aws.amazon.com/blogs/media/part-1-back-to-basics-with-adaptive-bitrate-video-streaming/
- **Twitch** — Transcoding em escala: https://blog.twitch.tv/en/2017/10/10/live-video-transmuxing-transcoding-f-fmpeg-vs-twitch-transcoder-part-i-489c1c125f28/

---

## Leituras — Sala de Cinema (Watch Party)

A sala sincroniza o estado de playback (play, pause, seek) entre todos os usuários via WebSocket. Um usuário é o host e os demais seguem seu estado.

- **Como o Disney+ construiu o GroupWatch** — https://medium.com/disney-streaming/how-we-built-the-disney-groupwatch-feature-6e9d2400a8fe
- **WebSockets — conceitos e arquitetura** — https://ably.com/topic/websockets

### Ponto de atenção
Latência diferente entre usuários pode causar dessincronização. Definir uma tolerância (ex: até X segundos de diferença antes de forçar resync) deve ser considerado na arquitetura.

---

## Deploy

- **Local** durante desenvolvimento
- **Oracle Cloud Free Tier** (VM ARM — 4 cores, 24GB RAM) para produção
- Acesso externo via **Tailscale** durante desenvolvimento local
