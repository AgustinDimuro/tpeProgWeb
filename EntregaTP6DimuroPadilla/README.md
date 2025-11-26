# TPEProgWebEntregaTP6

En este repositorio se encontrará la resolución de los incisos solicitados para la entrega referente al **Trabajo Práctico Especial de Programación Web** en el **Trabajo Práctico 6**.  

Los integrantes del grupo son:  
- Agustín Nicolás Dimuro  
- Tomás Agustín Padilla  

---

## Descripción de resoluciones implementadas para el Trabajo Práctico 6

### Cambios en la arquitectura
Como primer cambio estructural, decidimos adoptar la Onion Architecture con el objetivo de eliminar dependencias y reducir el fuerte acoplamiento de la logica de negocio con la tecnologias utilizadas, como la base de datos. Para lograr esto se establecieron las soguientes capas: Infraestructura, Aplicacion, Dominio y Entidades, asignando a cada una las responsabilidades necesarias para el correcto funcionamiento de la aplicacion.
Adicionalmente, modificamos todos los handlers (ubicados en la cada de Infraestructura), para que ya no manejen JSON y en cambio utilicen templ para implementar Server-Side Rendering. Para que esto sea posible se crearon archivos .templ en la seccion de views que determinan los componentes principales de la pagina web, los cuales son utilizados por los handlers previamente mencionados para poder generar dinamicamente el HTML completo que va a ser enviado al usuario.
Para cumplir con el objetivo principal de la última entrega, evolucionamos la aplicación utilizando HTMX para simplificar la interactividad del frontend, generando así una página dinámica que no requiere recargarse con cada actualización. Para lograrlo, implementamos la creación y actualización de reservas sin recargas mediante hx-post: al crear una reserva, el servidor devuelve únicamente el fragmento HTML con la lista actualizada en lugar de una redirección. Asimismo, aplicamos esta lógica a las eliminaciones; al borrar una reserva a través de hx-delete, la lista se actualiza eliminando el elemento correspondiente sin necesidad de refrescar la página completa.
---

### Aclaracion
Dejamos precargadas dos usuarios para asi poder probar la funcionalidad del sistema correctamente, estos dos usuarios son los siguientes:

Usuario
- ID = 1
- Password = 1234

Administrador:
- ID = 2
- Password = admin

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

   Por último, se ejecutará el `main.go` con datos precargados manualmente con el fin de facilitar la prueba del servicio. 
4. Para poder realizar las pruebas que preparamos para realizar un testeo  de la aplicación se puede ejecutar el siguiente comando en una consola distinta a la que el fue ejecutado el servidor:
   ```bash
   make hurltest
   ```
   Este comando realizará las pruebas que estan almacenadas en el archivo requests.hurl. Dentro de dicho archivo se encuentran las siguientes pruebas:
   - Carga de la página de Login y verificación de respuesta HTML.
   - Inicio de sesión exitoso y captura automática de la cookie de autenticación.
   - Acceso autorizado al Dashboard principal utilizando la sesión capturada.
   - Creación de una nueva reserva mediante envío de formulario.
   - Verificación de que la reserva creada aparece renderizada en el HTML.
   - Actualización de la fecha de la reserva existente.
   - Verificación de que la fecha vieja desaparece y la nueva se muestra correctamente.
   - Eliminación de la reserva del sistema.
   - Confirmación de que la reserva eliminada ya no se renderiza en el listado.
   - Carga correcta de la vista del calendario.
   - Navegación y filtrado del calendario por mes y año específicos.
   - Ejecución del logout e invalidación de la sesión.
   - Verificación de seguridad (bloqueo de acceso) al intentar entrar sin sesión.   

   Debería poder observar que al ejecutar este comando hurl le devuelve como resultado "Success". 

---
## Requisitos previos

- [Go](https://go.dev/dl/) (versión 1.20 o superior recomendada)  
- Git (para clonar el repositorio)  
- SQLC (para poder generar el código)
- hurl
- templ