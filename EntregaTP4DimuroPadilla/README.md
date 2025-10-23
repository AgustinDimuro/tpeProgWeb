# TPEProgWebEntregaTP2

En este repositorio se encontrará la resolución de los incisos solicitados para la entrega referente al **Trabajo Práctico Especial de Programación Web** en el **Trabajo Práctico 2**.  

Los integrantes del grupo son:  
- Agustín Nicolás Dimuro  
- Tomás Agustín Padilla  

---

## Descripción de la base de datos

Actualmente la base de datos posee dos tablas:  

- Una tabla referente a las **cabañas** y la información asociada a ellas, como su identificador (número de cabaña), datos de contacto y credenciales de acceso a la plataforma.  
- La otra tabla almacena la información de las **reservas**, incluyendo el número de la cabaña que realiza la reserva y la fecha correspondiente.  

---

## Cómo inicializo la base de datos

1. Clonar el repositorio.  
2. Abrir una terminal y navegar hasta el directorio **`TPEProgWebEntregaTP2`** (podés usar `ls` para listar directorios y luego `cd` para entrar).  
3. Inicializar el docker, base de datos en PostgreSQL y generar el código necesario con sqlc mediante el comando:
   
   ```bash
   make start
   ```

   Luego de ejecutado este comando, se realizará la descarga de la imagen de PostgreSQL para el docker compose en caso de que no esté descargada en su dispositivo. Del mismo modo, se ejecutará el comando `sqlc generate` y se generará el código pertinente para las queries.  

   Por último, se ejecutará el `main.go` preparado para que pueda observar una prueba realizada sobre la base de datos en la cuál se creará una cabaña, se creará una reserva, se listarán tanto la reserva como la cabaña, y por último se preguntará si dada una fecha existe una reserva.  

---

## Requisitos previos

- [Go](https://go.dev/dl/) (versión 1.20 o superior recomendada)  
- Git (para clonar el repositorio)  
- SQLC (para poder generar el código)

Podés verificar si tenés Go instalado ejecutando en la terminal:

```bash
go version
```

---

## Consejo en caso de que el sqlc no esté instalado o genere error de path (127)

Los siguientes consejos serán útiles en caso de tener Ubuntu/Linux, LinuxMint o similares.

1. **Instalar sqlc**  
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. **Si al ejecutar se genera el error 127 o error de path, recomendamos:**  
   ```bash
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

3. **Verificar la instalación:**  
   ```bash
   sqlc version
   ```

4. **Confirmar dónde quedó instalado:**  
   ```bash
   go env GOPATH
   ```

   Esto debería mostrar `/home/tuUsuario/go`, lo que significa que `sqlc` está instalado en el path correcto.
