# TPEProgWebEntregaTP1

En este repositorio se encontrará la resolución de los incisos solicitados para la entrega referente al **Trabajo Práctico Especial de Programación Web** en el **Trabajo Práctico 1**.  

Los integrantes del grupo son:  
- Agustín Nicolás Dimuro  
- Tomás Agustín Padilla  

---

## Descripción de la aplicación

Esta aplicación web tiene como objetivo organizar, en una página en línea, las reservas de un quincho ubicado dentro de un camping de cabañas.  
Las reservas se realizan por día completo y están asociadas al número de cabaña que efectúa la solicitud.  

Cada cabaña cuenta con:  
- Información de identificación (número de cabaña).  
- Datos de contacto (correo electrónico, número de teléfono).  
- Credenciales de acceso a la plataforma (ID y contraseña).  

Al iniciar sesión, cada cabaña debe ingresar su **ID** y **contraseña** para acceder al menú de reservas.  

---

## Cómo iniciar el servidor

1. Clonar el repositorio.  
2. Abrir una terminal y navegar hasta el directorio **`servidorTPE`** (podés usar `ls` para listar directorios y luego `cd` para entrar).  
3. Ejecutar el servidor con:  

   ```bash
   go run main.go

4. Acceder a la aplicación desde el navegador en:
http://localhost:8080/

## Credenciales de prueba

Número de cabaña: 1  
Contraseña: 1234  

## Requisitos previos

Antes de iniciar el servidor, asegurate de contar con lo siguiente instalado en tu sistema:

- [Go](https://go.dev/dl/) (versión 1.20 o superior recomendada)  
- Navegador web actualizado (Chrome, Firefox, Edge, etc.)  
- Git (para clonar el repositorio)  

Podés verificar si tenés Go instalado ejecutando en la terminal:

```bash
go version


