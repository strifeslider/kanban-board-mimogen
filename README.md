# KanBan Board Service - MimoGen

Board management microservice for KanBan SaaS platform.

## Features

- Workspace CRUD
- Board management
- Column CRUD with reordering
- WebSocket for real-time updates
- Membership management

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/workspaces | Create workspace |
| GET | /api/v1/workspaces | List workspaces |
| GET | /api/v1/workspaces/:id | Get workspace |
| POST | /api/v1/boards | Create board |
| GET | /api/v1/boards/:id | Get board |
| POST | /api/v1/columns | Create column |
| PUT | /api/v1/columns/reorder | Reorder columns |
| WS | /ws?token=jwt | WebSocket |

## Tech Stack

- Go 1.22
- Chi Router
- PostgreSQL
- gorilla/websocket

## Quick Start

```bash
export DB_HOST=localhost
export DB_NAME=kanban_board
export JWT_SECRET=your-secret

go run ./cmd/main.go
```

---

<details>
<summary><strong>IMPORTANT NOTICE</strong></summary>

<br>

**This repository was entirely generated using [MiMoCode](https://github.com/xiaomi/mimocode) - an AI-powered coding assistant by Xiaomi.**

All code, tests, documentation, and infrastructure configuration in this repository were created through AI-assisted development with MiMoCode.

---

**Этот репозиторий был полностью сгенерирован с помощью [MiMoCode](https://github.com/xiaomi/mimocode) - AI-ассистента для программирования от Xiaomi.**

Весь код, тесты, документация и инфраструктурная конфигурация в этом репозитории были созданы с помощью AI-ассистированной разработки MiMoCode.

</details>
