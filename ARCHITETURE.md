# Architecture

## Estilo Arquitetural

**Monolito Modular** com Bounded Contexts.

Cada módulo representa um contexto de domínio isolado com suas próprias responsabilidades. Módulos não acessam o código interno uns dos outros — comunicação acontece via interfaces explícitas. Essa estrutura permite extrair qualquer módulo como serviço independente no futuro sem refatoração.

---

## Estrutura de Pastas

```
/
├── cmd/
│   └── server/
│       └── main.go           ← entrypoint, inicializa dependências e rotas
├── internal/
│   ├── auth/
│   │   ├── handler.go        ← HTTP handlers
│   │   ├── service.go        ← regras de negócio
│   │   └── repository.go     ← acesso ao MongoDB
│   ├── family/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── media/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── transcoder.go     ← integração FFmpeg (HLS)
│   ├── stream/
│   │   ├── handler.go        ← serve segmentos HLS e Range Requests
│   │   └── service.go
│   ├── room/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── hub.go            ← gerencia conexões WebSocket da sala
│   └── middleware/
│       └── auth.go           ← valida JWT em rotas protegidas
├── pkg/
│   ├── storage/              ← cliente Backblaze B2 (S3-compatible)
│   └── database/             ← cliente MongoDB
├── migrations/               ← arquivos de migration (golang-migrate)
└── config/
    └── config.go             ← leitura de variáveis de ambiente
```

---

## Decisões Técnicas

### Linguagem
- **Go** — backend principal
- **Elixir** — candidato futuro para o módulo de sala de cinema (watch party), pelo modelo Actor via BEAM e Phoenix Channels

### Estilo de Código
- Estilo funcional — funções puras com dependências explícitas nos parâmetros
- Sem herança, sem OOP pesado
- Interfaces pequenas (1-2 métodos) seguindo o padrão idiomático Go
- Composição no lugar de herança

### HTTP
- **stdlib `net/http`** (Go 1.22+) — roteamento com parâmetros nativos via `{id}` e `r.PathValue()`
- Sem frameworks externos de roteamento
- Middleware como funções que envolvem `http.Handler`

```go
// Parâmetro de rota
mux.HandleFunc("GET /media/{id}", handler)
id := r.PathValue("id")

// Middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // valida JWT
        next.ServeHTTP(w, r)
    })
}
```

### Banco de Dados
- **MongoDB** — armazenamento de metadados (usuários, famílias, mídias, histórico, salvos, salas)
- **Driver oficial Go** (`go.mongodb.org/mongo-driver`) — sem ORM
- **golang-migrate** — gerenciamento de migrations

### Storage
- **Backblaze B2** via API S3-compatible
- **aws-sdk-go-v2** — SDK de acesso ao B2
- Upload recebe o arquivo no backend Go, faz stream pro B2
- Download de imagens via URL temporária do B2
- Vídeos servidos via HLS (segmentos no B2, playlist servida pelo backend)

### Auth
- **JWT** (`golang-jwt/jwt`) — autenticação stateless
- **bcrypt** (`golang.org/x/crypto`) — hash de senha
- Middleware valida token em todas as rotas protegidas

### Streaming de Vídeo
- Vídeos transcodados para **HLS** via FFmpeg após upload
- Segmentos `.ts` + playlist `.m3u8` armazenados no B2
- Backend serve a playlist e os segmentos com autenticação
- **HTTP Range Requests** para suporte a seek no player

### Sala de Cinema (Watch Party)
- **WebSockets** (`gorilla/websocket`) — sincronização em tempo real
- Host controla o estado (play, pause, seek)
- Servidor faz broadcast dos eventos para todos os membros da sala
- Tolerância de dessincronização configurável (threshold em segundos)

### Migrations
- **golang-migrate** — https://github.com/golang-migrate/migrate
- Arquivos em `/migrations/` versionados junto ao código

---

## Variáveis de Ambiente

```env
PORT=8080
MONGODB_URI=
MONGODB_DATABASE=

B2_ENDPOINT=
B2_BUCKET=
B2_KEY_ID=
B2_APPLICATION_KEY=

JWT_SECRET=
JWT_EXPIRATION_HOURS=

FFMPEG_PATH=
```

---

## Fluxos Principais

### Upload de Mídia
1. Cliente envia arquivo(s) via `POST /media/upload`
2. Backend valida tipo e tamanho
3. Backend faz stream do arquivo pro B2
4. FFmpeg transcoda o vídeo pra HLS (segmentos `.ts` + `.m3u8`)
5. Segmentos HLS sobem pro B2
6. Metadados salvos no MongoDB atomicamente

### Streaming de Vídeo
1. Cliente requisita `GET /media/:id/stream`
2. Backend valida JWT e retorna playlist `.m3u8`
3. Player requisita segmentos `.ts` individualmente
4. Backend valida JWT em cada segmento e serve com Range Request

### Sala de Cinema
1. Host cria sala via `POST /rooms`
2. Membros entram via `POST /rooms/:id/join`
3. Todos conectam via WebSocket em `WS /rooms/:id/sync`
4. Host emite eventos → servidor faz broadcast → membros atualizam player local
