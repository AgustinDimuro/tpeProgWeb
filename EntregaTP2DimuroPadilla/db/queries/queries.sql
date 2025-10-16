--                              CRUD de cabins

-- name: CreateCabin :one
INSERT INTO cabins (email_contact, phone_contact, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCabin :one
SELECT * FROM cabins WHERE id = $1;

-- name: ListCabins :many
SELECT * FROM cabins ORDER BY id;

-- name: UpdateCabin :one
UPDATE cabins
SET email_contact = $2,
    phone_contact = $3,
    password = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCabin :exec
DELETE FROM cabins WHERE id = $1;

--                          CRUD de reservations

-- name: CreateReservation :one
INSERT INTO reservations (cabin_id, fecha)
VALUES ($1, $2)
RETURNING *;

-- name: GetReservation :one
SELECT * FROM reservations WHERE id = $1;

-- name: GetReservationByFecha :one
SELECT * FROM reservations WHERE fecha = $1;

-- name: ListReservations :many
SELECT r.id, r.cabin_id, r.fecha, r.created_at,
       c.email_contact, c.phone_contact
FROM reservations r
JOIN cabins c ON c.id = r.cabin_id
ORDER BY r.fecha DESC;

-- name: ListReservationsByCabin :many
SELECT * FROM reservations
WHERE cabin_id = $1
ORDER BY fecha DESC;

-- name: UpdateReservationCabin :one
UPDATE reservations
SET cabin_id = $2
WHERE id = $1
RETURNING *;

-- name: UpdateReservationFecha :one
UPDATE reservations
SET fecha = $2
WHERE id = $1
RETURNING *;

-- name: DeleteReservation :exec
DELETE FROM reservations WHERE id = $1;

--                              Consultas
-- name: IsFechaDisponible :one
SELECT NOT EXISTS (
    SELECT 1 FROM reservations WHERE fecha = $1
    ) AS disponible;