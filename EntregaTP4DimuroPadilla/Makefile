# Nombre del binario de la aplicación
APP_NAME=mi-app

# URL de la base de datos
DB_URL="postgres://user:password@localhost:5432/TPEdb?sslmode=disable"

# Target por defecto
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
	@echo "=> Verificando estado del contenedor de la base..."
	@if [ "$$(docker ps -q -f name=mydb)" ]; then \
		echo "✅ Contenedor 'mydb' ya está corriendo."; \
	elif [ "$$(docker ps -aq -f name=mydb)" ]; then \
		echo "⚙️ Contenedor 'mydb' existe pero está detenido. Iniciando..."; \
		docker start mydb >/dev/null; \
	else \
		echo "🚀 No existe contenedor, iniciando con docker compose..."; \
		docker compose up -d; \
	fi


docker-down:
	@docker compose down

docker-logs:
	@docker compose logs -f db

wait-for-db: docker-up
	@echo "=> Esperando que la base de datos esté lista..."
	@until docker exec mydb pg_isready -U user -d mydb; do \
		sleep 1; \
	done
	@echo "✅ Base de datos lista!"

reset:
	docker stop mydb
	docker rm mydb

hurl:
	hurl --test requests.hurl

# ---- FLUJO COMPLETO ----
start: build wait-for-db 
	@echo "=> Iniciando aplicación..."
	@./$(APP_NAME)
