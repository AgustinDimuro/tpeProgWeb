# TPEProgWebEntregaTP3y4

En este repositorio se encontrará la resolución de los incisos solicitados para la entrega referente al **Trabajo Práctico Especial de Programación Web** en el **Trabajo Práctico 3 y en el Trabajo Práctico 4**.  

Los integrantes del grupo son:  
- Agustín Nicolás Dimuro  
- Tomás Agustín Padilla  

---

## Descripción de resoluciones implementadas para el Trabajo Práctico 3

### Conexión a la base de datos
El primer paso que realizamos fue conectar nuestra aplicación con la base de datos creada dirante la resolución del Trabajo Práctico 2.
Para ello agregamos funciones encargadas de realizar las operaciones de creación, modificación, borrado y lectura de datos tanto para la tabla de usuarios como para la tabla de reservas. Estas funciones se comunican con la base de datos a través de las funciones creadas por SQLC. Para que el usuario de nuestra aplicación pueda realizar las dichas acciones sobre la base de datos, creamos handlers que se encargan de procesar las solictudes o los datos que llegan por HTTP mediante los metodos GET, PUT y DELETE. 

---

### ¿Cómo pruebo la aplicación?

1. Clonar el repositorio.  
2. Abrir una terminal y navegar hasta el directorio **`TPEProgWebEntregaTP4`** (podés usar `ls` para listar directorios y luego `cd` para entrar).  
3. Inicializar el docker, base de datos en PostgreSQL y generar el código necesario con sqlc mediante el comando:
   
   ```bash
   make start
   ```

   Luego de ejecutado este comando, se realizará la descarga de la imagen de PostgreSQL para el docker compose en caso de que no esté descargada en su dispositivo. Del mismo modo, se ejecutará el comando `sqlc generate` y se generará el código pertinente para las queries.  

   Por último, se ejecutará el `main.go` preparado para que pueda observar una prueba realizada sobre la base de datos en la cuál se creará una cabaña, se creará una reserva, se listarán tanto la reserva como la cabaña, y por último se preguntará si dada una fecha existe una reserva.  
4. Para poder realizar las pruebas que preparamos para realizar un testeo  de la aplicación se puede ejecutar el siguiente comando en una consola distinta a la que el fue ejecutado el servidor:
   ```bash
   make hurl
   ```
   Este comando realizará las pruebas que estan almacenadas en el archivo requests.hurl. Dentro de dicho archivo se encuentran las siguientes pruebas:
      - Obtención de una cabaña que ya está cargada en el sistema y chequeo de que sea la cabaña pedida.
      - Actualización de datos de la cabaña.
      - Verificación que las modificaciones previas fueron realizadas de forma exitosa.
      - Intento de obtener cabaña inexistente.
      - Creación de nueva reserva.
      - Obtención de la reserva recién creada.
      - Actualización sobre el atributo fecha de la reserva creada recientemente.
      - Eliminación de la reserva.
      - Verificación que la reserva se eliminó correctamente.
      - Verificación de que la aplicación detecte un incorrecto formato en la fecha de una reserva.
      - Verificación de que la aplicación detecte un incorrecto formato en el ID de una cabaña.
      - Verificación de que la aplicación detecte un metodo incorrecto al intentar conectarse a un endpoint.
      - Lista todas las reservas a menos que no haya ninguna.

   Si desea ver las mismas pruebas pero en el modo de testeo del comando hurl, puede ejecutar:
   
   ```bash
   make hurltest
   ```
   Debería poder observar que al ejecutar este comando hurl le devuelve como resultado "Success". 
---

## Descripción de resoluciones implementadas para el Trabajo Práctico 4

### Estructura HTML
Al acceder a la página web se podran observar tres secciones. En la primer sección se puede encontrar tanto el título de la página, como un link de redireccionamiento hacia el calendario donde se podran ver las reservas. Para el caso de la segunda sección, es la encargada de implementar la creación de reservas. Por último, la tercer sección se pueden realizar modificaciones sobre las fechas de reservas existentes. A su vez, se pueden observar en forma de lista las reservas ya creadas, las cuales se actualizará dinámicamente para el caso de que se agrege o modifique una reserva.

### Comunicación entre la API y JavaScript

Para la comunicación entre el frontend (HTML) y la API REST, se utiliza el archivo `app.js`. Este script se encarga de manejar toda la interactividad de la página y las solicitudes de datos.

* **Obtención de datos (Read):** Al cargar la página, se ejecuta la función `getReservations`. Esta función realiza una petición `GET` al endpoint `/reservations` para obtener el listado completo de reservas y las muestra dinámicamente en la lista.
* **Creación de reservas (Create):** El formulario "form-crear-reserva" es manejado por un *event listener*. Al enviarlo, se capturan los datos (`cabin_id` y `fecha`) y se realiza una petición `POST` al endpoint `/reservation` para crear la nueva reserva.
* **Actualización de reservas (Update):** De forma similar, el formulario "form-actualizar-reserva" envía una petición `PUT` al endpoint `/reservation` para modificar la fecha de una reserva existente.
* **Eliminación de reservas (Delete):** Cada reserva en la lista tiene un botón "Eliminar" propio. Al hacer clic, la función `eliminarReserva` ejecuta una petición `DELETE` al endpoint `/reservation` para borrarla.

Todas las operaciones se realizan de forma asíncrona usando `async/await` con la API `fetch`. Después de cada operación exitosa (crear, actualizar o eliminar), se vuelve a llamar a `getReservations()` para refrescar la lista de reservas en el HTML, asegurando que el usuario siempre vea los datos actualizados.

### ¿Cómo pruebo la aplicación?
Si ya fue iniciado el servidor no es necesario realizar nada adicional. En caso contrario, se pueden seguir los pasos mencionados en la sección de mismo nombre dentro de las Descripciones de resoluciones implemmentadas para el Trabajo Práctico 3. Tener en cuenta que no es necesario ejecutar el "make hurl".
Una vez inicializado el servidor, acceda detro de su navegador al siguiente link "http://localhost:8080/". Allí verá la aplicación descripta anteriormente. Otro punto a tener en cuenta es que las cabañas que ya estan cargadas son la cabaña de ID 1 e ID 2, si desea probar con otra cabaña deberá crearla previamente.

---
## Requisitos previos

- [Go](https://go.dev/dl/) (versión 1.20 o superior recomendada)  
- Git (para clonar el repositorio)  
- SQLC (para poder generar el código)
- hurl

---


