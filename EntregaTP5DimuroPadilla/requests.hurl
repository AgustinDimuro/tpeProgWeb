# ------------------------------------------------------
# 1. FLUJO DE AUTENTICACIÓN (Login)
# ------------------------------------------------------

# Accedemos al formulario de login (GET) para verificar que carga
GET http://localhost:8080/login
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"

# Enviamos las credenciales (POST) - Simulando el submit del form
# Usamos cabin_id=1 y password="secreta" definidos en tu main.go
POST http://localhost:8080/login
[FormParams]
cabin_id: 1
password: secreta
HTTP 303
[Captures]
# CAPTURAMOS LA COOKIE: Es vital para las siguientes peticiones protegidas
auth_token: cookie "token"

# ------------------------------------------------------
# 2. RUTAS PROTEGIDAS (Usuario)
# ------------------------------------------------------

# Acceder al Dashboard (Home) con la cookie capturada
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
# Verificamos que cargó algo típico del dashboard, ej: "Reservas"
body contains "Reservas" 

# ------------------------------------------------------
# 3. GESTIÓN DE RESERVAS (Crear, Actualizar, Borrar)
# ------------------------------------------------------

# A. CREAR RESERVA (POST)
# Nota: Usamos una fecha futura lejana para evitar conflicto con validaciones
POST http://localhost:8080/reservations
[Cookies]
token: {{auth_token}}
[FormParams]
cabin_id: 1
fecha: 2025-12-25
HTTP 303 
# Esperamos 303 See Other porque tu handler hace un http.Redirect al final

# B. VERIFICAR CREACIÓN (GET)
# Vamos al home y buscamos si la fecha aparece en el HTML renderizado
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body contains "2025-12-25"

# C. ACTUALIZAR RESERVA (POST a /reservations/update)
POST http://localhost:8080/reservations/update
[Cookies]
token: {{auth_token}}
[FormParams]
cabin_id: 1
fecha: 2025-12-25
fecha_nueva: 2025-12-31
HTTP 303

# D. VERIFICAR ACTUALIZACIÓN (GET)
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
# La fecha vieja ya no debería estar
body not contains "2025-12-25"
# La fecha nueva sí debería estar
body contains "2025-12-31"

# E. ELIMINAR RESERVA (POST a /reservations/delete)
POST http://localhost:8080/reservations/delete
[Cookies]
token: {{auth_token}}
[FormParams]
cabin_id: 1
fecha: 2025-12-31
HTTP 303

# F. VERIFICAR ELIMINACIÓN (GET)
GET http://localhost:8080/
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body not contains "2025-12-31"

# ------------------------------------------------------
# 4. CALENDARIO
# ------------------------------------------------------

GET http://localhost:8080/calendario
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"

# Probar navegación del calendario (query params)
GET http://localhost:8080/calendario?month=1&year=2026
[Cookies]
token: {{auth_token}}
HTTP 200
[Asserts]
body contains "Enero"
body contains "2026"

# ------------------------------------------------------
# 5. ADMINISTRACIÓN (Admin)
# ------------------------------------------------------

# Listar reservas (Admin)
GET http://localhost:8080/admin/reservations
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"

# Actualizar Cabaña (PUT)
# Nota: Aunque HTML puro no soporta PUT en forms sin JS, 
# tu main.go escucha explícitamente PUT en "/admin/cabins".
# Hurl sí puede enviar PUT.
PUT http://localhost:8080/admin/cabins
[FormParams]
id: 1
email_contact: admin@test.com
phone_contact: 9999999
password: nueva_password
HTTP 303

# ------------------------------------------------------
# 6. LOGOUT
# ------------------------------------------------------

GET http://localhost:8080/logout
HTTP 303
# Debería redirigir al login

# Verificar que ya no podemos entrar al home
GET http://localhost:8080/
# Importante: Aquí no enviamos la cookie porque el navegador la habría borrado
# o enviamos una caducada. Hurl no guarda estado entre requests salvo que usemos variables.
# Sin la cookie, el middleware debería redirigirnos al login (303)
HTTP 303
[Asserts]
header "Location" contains "/login"