# Roadmap

## Fase 0 — Fundação

- [x] Monolito modular em `cmd/api` + `internal/<modulo>`
- [x] PostgreSQL 18 via Docker Compose, migrations com golang-migrate
- [x] `pgx` + `squirrel` com `postgres.DBTX` (pool ou transação)
- [x] Middlewares de log e recover, shutdown gracioso
- [x] `GET /healthz` e `GET /readyz`
- [x] Módulo `identity` de referência: criar, buscar e listar usuários

---

## Fase 1 — Auth & Usuários

### 1.1 Usuário
- [ ] Tabela `users` (`id`, `name`, `email`, `created_at`) — feita na Fase 0; falta `password_hash`
- [ ] `POST /auth/register` — cadastro com senha hasheada (bcrypt)
- [ ] `POST /auth/login` — retorna JWT
- [ ] Middleware de autenticação (valida JWT em todas as rotas protegidas)

### 1.2 Família
- [ ] Tabelas `families` (`id`, `name`, `owner_id`, `invite_code`, `created_at`) e `family_members` (`family_id`, `user_id`, `role`)
- [ ] `POST /families` — cria família e gera `invite_code` único (owner vira membro automaticamente)
- [ ] `POST /families/join` — entra na família pelo `invite_code`
- [ ] `GET /families/me` — retorna família do usuário autenticado com lista de membros

---

## Fase 2 — Upload de Mídias

### 2.1 Modelo de mídia
- [ ] Modelo `Media` no Postgres (`id`, `title`, `type` (filme/série/anime/foto), `owner_id`, `family_id`, `b2_key`, `size`, `duration`, `thumbnail_key`, `created_at`)
- [ ] Modelo `Episode` para séries (`id`, `media_id`, `season`, `episode`, `title`, `b2_key`, `duration`)

### 2.2 Upload
- [ ] `POST /media/upload` — recebe arquivo(s), faz stream pro B2, salva metadados no Postgres
- [ ] Suporte a multipart upload (múltiplos arquivos simultâneos)
- [ ] Validação de tipo de arquivo (vídeo, imagem) e tamanho máximo
- [ ] Geração de thumbnail via FFmpeg após upload
- [ ] Transcodagem para HLS via FFmpeg (segmentos `.ts` + playlist `.m3u8`) e upload dos segmentos pro B2

### 2.3 Biblioteca
- [ ] `GET /media` — lista mídias da família com filtros (tipo, título)
- [ ] `GET /media/:id` — detalhes de uma mídia
- [ ] `DELETE /media/:id` — remove mídia do B2 e do Postgres (apenas owner)

---

## Fase 3 — Download & Streaming

### 3.1 Streaming HLS
- [ ] `GET /media/:id/stream` — retorna a playlist `.m3u8` da mídia
- [ ] `GET /media/:id/segments/:file` — serve os segmentos `.ts` com autenticação
- [ ] Suporte a HTTP Range Requests para seek no player

### 3.2 Fotos
- [ ] `GET /media/:id/download` — gera URL temporária do B2 para download direto de imagens

### 3.3 Histórico
- [ ] Modelo `History` (`user_id`, `media_id`, `progress_seconds`, `watched_at`)
- [ ] `POST /media/:id/progress` — salva progresso do usuário
- [ ] `GET /history` — retorna histórico do usuário autenticado

### 3.4 Salvos
- [ ] Modelo `Saved` (`user_id`, `media_id`, `saved_at`)
- [ ] `POST /media/:id/save` — salva mídia na lista do usuário
- [ ] `DELETE /media/:id/save` — remove dos salvos
- [ ] `GET /saved` — lista mídias salvas do usuário

---

## Fase 4 — Sala de Cinema (Watch Party)

### 4.1 Sala
- [ ] Modelo `Room` (`id`, `media_id`, `host_id`, `members[]`, `state`: play/pause/timestamp, `created_at`)
- [ ] `POST /rooms` — cria sala para uma mídia
- [ ] `POST /rooms/:id/join` — entra na sala
- [ ] `DELETE /rooms/:id` — fecha a sala (apenas host)

### 4.2 Sincronização em tempo real
- [ ] `WS /rooms/:id/sync` — conexão WebSocket por membro
- [ ] Eventos: `play`, `pause`, `seek`, `member_joined`, `member_left`
- [ ] Host emite eventos → servidor faz broadcast pra todos os membros
- [ ] Lógica de tolerância de dessincronização (threshold de X segundos antes de forçar resync)
