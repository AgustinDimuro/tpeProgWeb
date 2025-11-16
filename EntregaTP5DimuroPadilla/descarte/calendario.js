document.addEventListener('DOMContentLoaded', async function() {
    const calendarEl = document.getElementById('calendario');

    if (!calendarEl) {
        console.error("No se encontró el elemento #calendario");
        return;
    }

    const calendar = new FullCalendar.Calendar(calendarEl, {
        initialView: 'dayGridMonth', // Vista inicial de mes
        locale: 'es', // Poner el calendario en español
        headerToolbar: {
            left: 'prev,next today',
            center: 'title',
            right: 'dayGridMonth,timeGridWeek,listWeek' // Vistas disponibles
        },
        buttonText: {
             today:    'Hoy',
             month:    'Mes',
             week:     'Semana',
             list:     'Lista'
        },

        // 'events' es la propiedad clave.
        // La usamos como una función para cargar nuestras reservas dinámicamente.
        events: async function(fetchInfo, successCallback, failureCallback) {
            try {
                // 1. Hacemos el fetch a nuestro endpoint existente
                const res = await fetch('/reservations'); // app.js usa ENDPOINT_LIST, pero aquí podemos usar la ruta directa
                if (!res.ok) {
                    throw new Error(`Error al obtener reservas: ${res.status}`);
                }
                const reservations = await res.json(); // Espera un array: [{id, cabin_id, fecha}, ...]

                // 2. Transformamos los datos
                const events = reservations.map(r => {
                    // Tu API devuelve la fecha completa (ej: "2025-10-28T00:00:00Z")
                    // FullCalendar la entiende, pero solo queremos la parte de la fecha.
                    const fecha = r.fecha.slice(0, 10); // Queda "2025-10-28"

                    return {
                        title: `Cabaña ${r.cabin_id ?? r.cabinID}`, // El texto que se muestra en el evento
                        start: fecha,      // El día del evento
                        allDay: true,      // Marca el día completo como ocupado
                        // (Opcional) Agregamos datos extra
                        extendedProps: {
                            cabinId: r.cabin_id ?? r.cabinID,
                            reservaId: r.id ?? r.ID
                        }
                    };
                });

                // 3. Entregamos los eventos transformados a FullCalendar
                successCallback(events);

            } catch (e) {
                console.error(e);
                failureCallback(e); // Informamos a FullCalendar que hubo un error
                alert("No se pudieron cargar las reservas en el calendario.");
            }
        },
        
        // (Opcional) Cambiar el color de los eventos
        eventColor: '#d64541', // Usa el color --red-fire de tu paleta
        
        // (Opcional) ¿Qué pasa al hacer clic en un evento?
        eventClick: function(info) {
            alert(`Reserva de Cabaña ${info.event.extendedProps.cabinId} el día ${info.event.start.toLocaleDateString()}`);
        }
    });

    // ¡Renderizar el calendario!
    calendar.render();
});