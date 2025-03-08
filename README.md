# Golang Websockets Microservices

This project is a **Golang-based WebSockets-powered microservices architecture** that enables real-time communication across different services. It includes several microservices for searching hashtags, places, users, and handling messaging, all orchestrated using **Docker Compose**.

## 📌 Project Structure

```
├── docker-compose.yml   # Docker Compose configuration
├── hashtags_search      # Hashtag search service
│   ├── controllers      # Request handlers for hashtags
│   ├── handlers         # WebSocket handlers
│   ├── models           # Data models
│   ├── modules          # Database and logging modules
│   ├── tests            # Unit tests
│   ├── build            # Docker build files
│   ├── go.mod, go.sum   # Go dependencies
│   └── search_hashtags.go # Service entry point
├── messenger_engine     # Real-time messaging service
│   ├── controllers      # Handles chat, messages, and broadcasts
│   ├── models           # Message and user models
│   ├── handlers         # WebSocket handlers
│   ├── modules          # Database and logger
│   ├── tests            # Unit tests
│   ├── build            # Docker build files
│   ├── go.mod, go.sum   # Go dependencies
│   └── messenger.go     # Service entry point
├── places_search        # Places search service
│   ├── controllers      # Request handlers for places
│   ├── handlers         # WebSocket handlers
│   ├── models           # Response models
│   ├── modules          # Database and logging modules
│   ├── tests            # Unit tests
│   ├── build            # Docker build files
│   ├── go.mod, go.sum   # Go dependencies
│   └── search_places.go # Service entry point
├── user_search          # User search service
│   ├── controllers      # Handles user search queries
│   ├── handlers         # WebSocket handlers
│   ├── models           # User models
│   ├── modules          # Database and logger
│   ├── tests            # Unit tests
│   ├── build            # Docker build files
│   ├── go.mod, go.sum   # Go dependencies
│   └── search.go        # Service entry point
```

## 🚀 Features

- **Real-time WebSockets Communication**: Enables instant messaging and live updates.
- **Microservices Architecture**: Each service operates independently for scalability.
- **Dockerized Setup**: Easy deployment with `docker-compose`.
- **PostgreSQL Database**: Centralized database for all services.
- **Environment Configuration**: Uses `.env` files to manage settings.
- **Efficient Resource Management**: Services have memory constraints for stability.

## 🛠️ Setup & Installation

### 1️⃣ Prerequisites
Ensure you have the following installed:
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### 2️⃣ Clone the Repository
```sh
 git clone https://github.com/your-repo/golang-websockets-microservices.git
 cd golang-websockets-microservices
```

### 3️⃣ Create an `.env` File
Each service loads environment variables from a `.env` file. Create one in the root directory with:
```sh
DATABASE_TYPE=postgres
DATABASE_USER=timofeyivankov
DATABASE_PASSWORD=Lovell32bd
DATABASE_NAME=my_database_new
SSL_MODE=False
DATABASE_HOST=localhost
```

### 4️⃣ Build and Run the Services
```sh
docker-compose up --build
```

### 5️⃣ Access the Services
- **Hashtag Search**: `http://localhost:8380`
- **Messenger Engine**: `http://localhost:8440`
- **Places Search**: `http://localhost:8285`
- **User Search**: `http://localhost:8280`

## 🧪 Running Tests
Each microservice includes test cases to validate functionality. To run tests inside a container, use:
```sh
docker exec -it <container_name> go test ./...
```
Or run locally with:
```sh
go test ./...
```

## 🔧 API Endpoints
Each microservice exposes different endpoints. Below are some common routes:

### Hashtag Search Service (`http://localhost:8380`)
- `GET /hashtags/search?query=<query>` - Search hashtags.
- `WS /hashtags/live` - WebSocket for real-time hashtag tracking.

### Messenger Engine (`http://localhost:8440`)
- `POST /messages/send` - Send a new message.
- `GET /messages/chat?chat_id=<id>` - Get chat messages.
- `WS /chat/connect` - WebSocket for live messaging.

### Places Search Service (`http://localhost:8285`)
- `GET /places/search?location=<lat,lon>` - Search for places.
- `WS /places/live` - WebSocket for real-time location updates.

### User Search Service (`http://localhost:8280`)
- `GET /users/search?name=<name>` - Search users.
- `WS /users/live` - WebSocket for real-time user updates.

## 📜 License
This project is licensed under the **MIT License**.

---
🔹 *Contributions are welcome! Feel free to submit a pull request.*
