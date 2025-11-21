# TPEProgWebEntregaTP3y4

En este repositorio se encontrará la resolución de los incisos solicitados para la entrega referente al **Trabajo Práctico Especial de Programación Web** en el **Trabajo Práctico 5**.  

Los integrantes del grupo son:  
- Agustín Nicolás Dimuro  
- Tomás Agustín Padilla  

---

## Descripción de resoluciones implementadas para el Trabajo Práctico 5

### Cambios en la arquitectura
Como primer cambio estructural, decidimos adoptar la Onion Architecture con el objetivo de eliminar dependencias y reducir el fuerte acoplamiento de la logica de negocio con la tecnologias utilizadas, como la base de datos. Para lograr esto se establecieron las soguientes capas: Infraestructura, Aplicacion, Dominio y Entidades, asignando a cada una las responsabilidades necesarias para el correcto funcionamiento de la aplicacion.
Adicionalmente, modificamos todos los handlers (ubicados en la cada de Infraestructura), para que ya no manejen JSON y en cambio utilicen templ para implementar Server-Side Rendering. Para que esto sea posible se crearon archivos .templ en la seccion de views que determinan los componentes principales de la pagina web, los cuales son utilizados por los handlers previamente mencionados para poder generar dinamicamente el HTML completo que va a ser enviado al usuario.

---

### ¿Cómo pruebo la aplicación?

1. Clonar el repositorio.  
2. Abrir una terminal y navegar hasta el directorio **`TPEProgWebEntregaTP5`** (podés usar `ls` para listar directorios y luego `cd` para entrar).  
3. Inicializar el docker, base de datos en PostgreSQL y generar el código necesario con sqlc mediante el comando:
   
   ```bash
   make start
   ```

   Luego de ejecutado este comando, se realizará la descarga de la imagen de PostgreSQL para el docker compose en caso de que no esté descargada en su dispositivo. Del mismo modo, se ejecutará el comando `sqlc generate` y se generará el código pertinente para las queries.

   A su vez, luego de generado el sqlc se ejecutara el comando templ dentro del Makefile que tiene como trabajo principal ejecutar el comando templ generate para transformar los archivos .templ a .go y que puedan ser utilizados por la  aplicacion.  

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
## Requisitos previos

- [Go](https://go.dev/dl/) (versión 1.20 o superior recomendada)  
- Git (para clonar el repositorio)  
- SQLC (para poder generar el código)
- hurl
- templ