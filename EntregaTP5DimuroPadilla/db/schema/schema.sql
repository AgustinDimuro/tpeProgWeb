-- db/schema.sql

CREATE TABLE IF NOT EXISTS cabins (
  id             SERIAL PRIMARY KEY,
  email_contact  TEXT    NOT NULL CHECK (position('@' in email_contact) > 1),
  phone_contact  TEXT    NOT NULL,
  password       TEXT    NOT NULL,
  role           TEXT    NOT NULL DEFAULT 'user', -- <--- ESTA ES LA LÍNEA QUE FALTA
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reservations (
  id         SERIAL PRIMARY KEY,
  cabin_id   INTEGER NOT NULL REFERENCES cabins(id) ON DELETE RESTRICT,
  fecha      DATE    NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_reservations_fecha UNIQUE (fecha)
);

-- Índices útiles
CREATE INDEX IF NOT EXISTS idx_reservations_cabin_id ON reservations(cabin_id);
CREATE INDEX IF NOT EXISTS idx_reservations_fecha ON reservations(fecha);