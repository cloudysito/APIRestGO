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
- ✅ **Context Management** - Proper timeout handling

## 🛠 Tech Stack

- **Language:** Go 1.21+
- **Database:** MongoDB
- **Libraries:** 
  - `mongo-go-driver` - MongoDB client
  - `godotenv` - Environment configuration
  - `net/http` - Standard HTTP package

## 📋 Prerequisites

- Go 1.21 or higher
- MongoDB 4.4+
- `.env` file with MongoDB URI and PORT

## ⚙️ Installation

1. **Clone the repository**
```bash
git clone https://github.com/cloudysito/APIRestGO
cd apirestgo
```

2. **Download dependencies**
```bash
go mod download
```

3. **Setup environment variables**
Create a `.env` file in the root directory:
```env
MONGO_URI=mongodb://localhost:27017
PORT=8080
```

4. **Run the server**
```bash
go run main.go
```

Server will start at `http://localhost:8080`

## 📚 API Endpoints

### **Players**

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/register` | ✅ Required | Register a new player |
| GET | `/api/players` | ❌ No | Get all players |
| GET | `/api/player/{name}` | ❌ No | Get player by name |
| GET | `/api/stats` | ✅ Required | Get rank statistics |
| PUT | `/api/player/rank/{name}` | ✅ Required | Update player rank |
| DELETE | `/api/player/delete/{name}` | ✅ Required | Delete a player |

### **Factions**

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/factions` | ❌ No | Get all factions |
| POST | `/api/factions` | ✅ Required | Create new faction |
| GET | `/api/faction?name={name}` | ❌ No | Get faction by name |
| PUT | `/api/faction/update?name={name}` | ✅ Required | Update faction (partial) |
| DELETE | `/api/faction/delete?name={name}` | ✅ Required | Delete faction |


## 🏗 Project Structure

```
apirestgo/
├── main.go                 # Application entry point
├── go.mod                  # Module definition
├── .env.example            # Environment template
├── handlers/
│   ├── player_handler.go   # Player HTTP handlers
│   └── faction_handler.go  # Faction HTTP handlers
├── repository/
│   ├── mongo_repo.go       # MongoDB player implementation
│   ├── faction_repo.go     # MongoDB faction implementation
│   └── player_repo.go      # Player repository interface
├── middleware/
│   └── auth.go             # Token validation middleware
└── models/
    ├── player.go           # Player data model
    ├── faction.go          # Faction data model
    └── stats.go            # Statistics model
```

## 🎯 Architecture

This project uses the **Repository Pattern**:

```
HTTP Request
    ↓
Handler (Business Logic)
    ↓
Repository (Data Access)
    ↓
MongoDB (Persistence)
```

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
- **Aggregation Pipelines:** Efficient rank statistics computation
- **Partial Updates:** MongoDB `$set` operator for atomic updates

## 👨‍💻 Author

Built as a portfolio project demonstrating Go REST API best practices.
