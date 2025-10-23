# ---------------------------------
# Pruebas del Endpoint /cabin (SINGULAR)
# ---------------------------------

# 1. Obtener la cabaña con ID 1
GET http://localhost:8080/cabin?id=1
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.email_contact" == "contacto@ejemplo.com"


# 2. Actualizar la cabaña con ID 1
PUT http://localhost:8080/cabin?id=1&emailContact=nuevo@mail.com&phoneContact=111222&password=newpass
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.email_contact" == "nuevo@mail.com"
jsonpath "$.phone_contact" == "111222"


# 3. Verificar que la cabaña se actualizó correctamente
GET http://localhost:8080/cabin?id=1
HTTP 200
[Asserts]
jsonpath "$.email_contact" == "nuevo@mail.com"


# 4. Intentar obtener una cabaña que no existe
GET http://localhost:8080/cabin?id=999
HTTP 500


# ---------------------------------
# Pruebas del Endpoint /reservation (SINGULAR)
# ---------------------------------

# 5. Crear una nueva reserva para la cabaña 1
POST http://localhost:8080/reservation?cabin_id=1&fecha=2025-10-30
HTTP 200
[Asserts]
jsonpath "$.cabin_id" == 1
jsonpath "$.fecha" contains "2025-10-30T00:00:00Z"


# 6. Obtener la reserva recién creada
GET http://localhost:8080/reservation?fecha=2025-10-30
HTTP 200
[Asserts]
jsonpath "$.cabin_id" == 1


# 7. Actualizar la fecha de la reserva
PUT http://localhost:8080/reservation?cabin_id=1&fecha=2025-10-30&fecha_nueva=2025-11-01
HTTP 200
[Asserts]
jsonpath "$.cabin_id" == 1
jsonpath "$.fecha" contains "2025-11-01T00:00:00Z"


# 8. Eliminar la reserva (con su nueva fecha)
DELETE http://localhost:8080/reservation?cabin_id=1&fecha=2025-11-01
HTTP 200

# 9. Verificar que la reserva fue eliminada
GET http://localhost:8080/reservation?fecha=2025-11-01
HTTP 404

# ---------------------------------
# Pruebas de Errores de Input
# ---------------------------------

# 10. Probar un formato de fecha inválido
GET http://localhost:8080/reservation?fecha=30-10-2025
HTTP 404


# 11. Probar un ID de cabaña no numérico
PUT http://localhost:8080/cabin?id=abc&emailContact=test@mail.com&phoneContact=123&password=pass
HTTP 400


# 12. Probar método no permitido en /cabin
POST http://localhost:8080/cabin
HTTP 405

# ---------------------------------
# Prueba para /reservations (PLURAL)
# ---------------------------------

# 13. Listar TODAS las reservas (esta sí es plural)
GET http://localhost:8080/reservations
HTTP 200
[Asserts]
jsonpath "$.[0].cabin_id" > 0 # Comprueba que al menos hay una reserva en la lista