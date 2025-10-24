const API_BASE = "";
const ENDPOINT_LIST = `${API_BASE}/reservations`;

const lista = document.getElementById("lista-reservas");
const crearReservas = document.getElementById("form-crear-reserva");
const actualizarReservas = document.getElementById("form-actualizar-reserva");

async function getReservations() {
    try {
        const res = await fetch(`${API_BASE}/reservations`);
        if (!res.ok) throw new Error(`GET /reservations -> ${res.status}`);
            const data = await res.json(); // espera array
            lista.innerHTML = "";
            data.forEach(r => lista.appendChild(itemReserva(r)));
    } catch (e) {
        console.error(e);
    }
}

function itemReserva(r) {
    // Asumo que sqlc devuelve campos con nombres CabinID, Fecha, ID, etc.
    // Si Fecha es string ISO o fecha con zona: mostrala formateada.
    const li = document.createElement("li");
    li.style.margin = "6px 0";

    const fecha = typeof r.fecha === "string"
        ? r.fecha.slice(0,10)
        : new Date(r.fecha).toISOString().slice(0,10);

    li.textContent = `Cabin ${r.cabin_id ?? r.cabinID} — ${fecha}`;

    const btn = document.createElement("button");
    btn.textContent = "Eliminar";
    btn.style.marginLeft = "8px";
    btn.addEventListener("click", () => eliminarReserva(r.cabin_id ?? r.cabinID, fecha));

    li.appendChild(btn);
    return li;
}

crearReservas.addEventListener("submit", async (ev) => {
  ev.preventDefault();
  const fd = new FormData(crearReservas);
  const cabin_id = (fd.get("cabin_id") || "").trim();
  const fecha    = (fd.get("fecha") || "").trim();

  if (!cabin_id || !fecha) return alert("Completá Cabin ID y fecha (YYYY-MM-DD).");

  try {
    const url = `${API_BASE}/reservation?cabin_id=${encodeURIComponent(cabin_id)}&fecha=${encodeURIComponent(fecha)}`;
    const res = await fetch(url, { method: "POST" });
    if (!res.ok) 
        throw new Error(`POST /reservation -> ${res.status}`);
    await getReservations(); // refrescar lista completa
    crearReservas.reset();
  } catch (e) {
        console.error(e);
        alert("No se pudo crear la reserva.");
        crearReservas.reset();
  }
});

actualizarReservas.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    const fd = new FormData(actualizarReservas);
    const cabin_id   = (fd.get("cabin_id") || "").trim();
    const fecha      = (fd.get("fecha") || "").trim();
    const fechaNueva = (fd.get("fecha_nueva") || "").trim();

    if (!cabin_id || !fecha || !fechaNueva) {
        return alert("Completá Cabin ID, fecha actual y nueva fecha (YYYY-MM-DD).");
    }

    try {
        const url =
        `${API_BASE}/reservation?cabin_id=${encodeURIComponent(cabin_id)}&fecha=${encodeURIComponent(fecha)}&fecha_nueva=${encodeURIComponent(fechaNueva)}`;
        const res = await fetch(url, { method: "PUT" });
        if (!res.ok) 
            throw new Error(`PUT /reservation -> ${res.status}`);
        await getReservations();
        actualizarReservas.reset();
    } catch (e) {
        console.error(e);
        alert("No se pudo actualizar la reserva.");
        actualizarReservas.reset();
    }
});



async function eliminarReserva(cabin_id, fecha) {
  if (!confirm(`Eliminar reserva de cabin ${cabin_id} en ${fecha}?`)) return;
  try {
    const url = `${API_BASE}/reservation?cabin_id=${encodeURIComponent(cabin_id)}&fecha=${encodeURIComponent(fecha)}`;
    const res = await fetch(url, { method: "DELETE" });
    if (!res.ok) throw new Error(`DELETE /reservation -> ${res.status}`);
        await getReservations();
  } catch (e) {
        console.error(e);
        alert("No se pudo eliminar la reserva.");
  }
}

document.addEventListener("DOMContentLoaded", getReservations);