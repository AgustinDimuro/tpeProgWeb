# Nombre del binario de la aplicación
APP_NAME=mi-app

# URL de la base de datos (coincide con las credenciales del docker run)
DB_URL="postgres://user:password@localhost:5432/TPEdb?sslmode=disable"

# Target por defecto: compilar
all: build

# Inicia el servidor con recarga automática usando Air
run:
	@air

# Genera código sqlc
generate:
	@echo "=> Generando código con sqlc..."
	@sqlc generate

# Construye el binario
build: generate
	@echo "=> Compilando aplicación..."
	@go build -o $(APP_NAME) .

# Limpia artefactos
clean:
	@rm -f $(APP_NAME)

# ---- DOCKER COMPOSE ----
docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

docker-logs:
	@docker compose logs -f db

# ---- FLUJO COMPLETO ----
start:
	@echo "🚀 Levantando Postgres con Docker..."
	@docker compose up -d
	@sleep 5 # esperar a que Postgres arranque
	@echo "🔄 Generando código con sqlc..."
	@sqlc generate
	@echo "▶️ Ejecutando la aplicación..."
	@go run main.go