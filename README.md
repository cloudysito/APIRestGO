# GameFaction REST API

A robust REST API built with **Go** and **MongoDB** for managing players and factions in a gaming system. Features token-based authentication, MongoDB integration, and a clean Repository pattern architecture.

## 🚀 Features

- ✅ **Player Management** - Full CRUD operations with ranking system
- ✅ **Faction System** - Create and manage player factions
- ✅ **Authentication** - Token-based middleware validation
- ✅ **Database** - MongoDB with aggregation pipelines
- ✅ **Robust Error Handling** - Validation at repository level
- ✅ **Partial Updates** - Update only specific fields without overwriting
- ✅ **Goroutines** - Asynchronous operations (MMR calculation)
- ✅ **Worker Pool** - Bounded concurrent MMR recalculation across all players
- ✅ **Graceful Shutdown** - Drains in-flight requests before exiting
- ✅ **Context Management** - Proper timeout handling

## 🛠 Tech Stack

- **Language:** Go 1.21+
- **Router:** chi (v5)
- **Database:** MongoDB
- **Libraries:**
  * `mongo-go-driver` - MongoDB client
  * `go-chi/chi` - HTTP router and middleware
  * `godotenv` - Environment configuration

## 📋 Prerequisites

- Go 1.21 or higher
- MongoDB 4.4+
- `.env` file with MongoDB URI and PORT

## ⚙️ Installation

1. **Clone the repository**

git clone https://github.com/cloudysito/APIRestGO
cd apirestgo


2. **Download dependencies**

go mod download


3. **Setup environment variables** Create a `.env` file in the root directory:

MONGO_URI=mongodb://localhost:27017
PORT=8080
SECRET_TOKEN=your-token-here


4. **Run the server**

go run main.go


Server will start at `http://localhost:8080`

## 📚 API Endpoints

### **Players**

| Method | Endpoint                    | Auth       | Description                                                |
| ------ | ---------------------------- | ---------- | ------------------------------------------------------------ |
| POST   | `/api/players`               | ✅ Required | Register a new player                                       |
| GET    | `/api/players`               | ❌ No       | Get all players                                              |
| GET    | `/api/players/{name}`        | ❌ No       | Get player by name                                           |
| GET    | `/api/stats`                 | ✅ Required | Get rank statistics                                          |
| PUT    | `/api/players/{name}/rank`   | ✅ Required | Update player rank                                           |
| DELETE | `/api/players/{name}`        | ✅ Required | Delete a player                                              |
| POST   | `/api/admin/recalculate-mmr` | ✅ Required | Recalculate MMR for all players (worker pool, concurrent)   |

### **Factions**

| Method | Endpoint                          | Auth       | Description              |
| ------ | ---------------------------------- | ---------- | -------------------------- |
| GET    | `/api/factions`                    | ❌ No       | Get all factions           |
| POST   | `/api/factions`                    | ✅ Required | Create new faction         |
| GET    | `/api/factions/{name}`             | ❌ No       | Get faction by name        |
| PUT    | `/api/factions/{name}`             | ✅ Required | Update faction (partial)   |
| DELETE | `/api/factions/{name}`             | ✅ Required | Delete faction             |

## 🏗 Project Structure

apirestgo/
├── main.go # Application entry point, routing and graceful shutdown
├── go.mod # Module definition
├── .env.example # Environment template
├── handlers/
│ ├── player_handler.go # Player HTTP handlers + MMR worker pool
│ └── faction_handler.go # Faction HTTP handlers
├── repository/
│ ├── mongo_repo.go # MongoDB player implementation
│ ├── faction_repo.go # MongoDB faction implementation
│ └── player_repo.go # Player repository interface
├── middleware/
│ └── auth.go # Token validation middleware
└── models/
├── player.go # Player data model
├── faction.go # Faction data model
└── stats.go # Statistics model


## 🎯 Architecture

This project uses the **Repository Pattern**:

HTTP Request
↓
Handler (Business Logic)
↓
Repository (Data Access)
↓
MongoDB (Persistence)


## ⚡ Concurrency

This project implements two real concurrency patterns, not just `go func()` calls:

- **Async MMR calculation on registration:** when a player registers, their initial MMR is calculated and persisted in a background goroutine using `context.Background()` (not the request context), so the write survives after the HTTP response is already sent.
- **Worker pool for bulk recalculation:** `POST /api/admin/recalculate-mmr` fans out work across a fixed pool of 5 workers using buffered channels and `sync.WaitGroup`, instead of spawning one goroutine per player. This bounds concurrent MongoDB writes regardless of how many players exist.
- **Graceful shutdown:** the server listens for `SIGINT`/`SIGTERM`, stops accepting new requests, and waits (up to a 10s timeout) for in-flight requests to finish before exiting — the same pattern used in production deployments (Docker, Kubernetes).

## 🔒 Error Handling

The API returns structured error responses:

```json
{
  "error": "Player not found."
}
```

All repository methods validate that operations actually succeeded:

- `UpdateRank()` checks `MatchedCount`
- `DeletePlayer()` checks `DeletedCount`
- `GetByName()` distinguishes "not found" from actual errors

## 🚀 Performance Features

- **Context Timeouts:** All DB operations timeout after 5 seconds
- **Goroutines:** Async MMR calculation on player registration
- **Worker Pool:** Bounded concurrency (5 workers) for bulk MMR recalculation
- **Graceful Shutdown:** Server drains in-flight requests before exiting on SIGINT/SIGTERM
- **Aggregation Pipelines:** Efficient rank statistics computation
- **Partial Updates:** MongoDB `$set` operator for atomic updates
