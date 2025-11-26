# ------------------------------------------------------
# 1.Login 
# ------------------------------------------------------

GET http://localhost:8080/login
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"

POST http://localhost:8080/login
[FormParams]
cabin_id: 2
password: admin
HTTP 303
[Captures]
auth_token: cookie "token"

# ------------------------------------------------------
# 2. Rutas de Usuario
# ------------------------------------------------------

GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Reservas" 

# ------------------------------------------------------
# 3. Gestión de Reservas
# ------------------------------------------------------

# CREAR RESERVA (POST)
POST http://localhost:8080/reservations
[Cookies]
token: {{auth_token}}
[FormParams]
cabin_id: 1
fecha: 2025-12-25
HTTP 200 
[Asserts]
body contains "2025-12-25"

# VERIFICAR CREACION
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body contains "2025-12-25"

# C. ACTUALIZAR RESERVA
POST http://localhost:8080/reservations/update
[Cookies]
token: {{auth_token}}
[FormParams]
cabin_id: 1
fecha: 2025-12-25
fecha_nueva: 2025-12-31
HTTP 200
[Asserts]
body not contains "2025-12-25"
body contains "2025-12-31"

# D. VERIFICAR ACTUALIZACIÓN
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body contains "2025-12-31"

# E. ELIMINAR RESERVA
POST http://localhost:8080/reservations/delete
[Cookies]
token: {{auth_token}}
[FormParams]
cabin_id: 1
fecha: 2025-12-31
HTTP 200

# F. VERIFICAR ELIMINACIÓN
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body not contains "2025-12-31"

# ------------------------------------------------------
# 4. Calendario
# ------------------------------------------------------

GET http://localhost:8080/calendario
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"

GET http://localhost:8080/calendario?month=1&year=2026
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body contains "Enero"
body contains "2026"

# ------------------------------------------------------
# 5. Logout
# ------------------------------------------------------

GET http://localhost:8080/logout
HTTP 303

# Verificar que ya no podemos entrar al home
GET http://localhost:8080/
HTTP 303
[Asserts]
header "Location" contains "/login"