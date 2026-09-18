# SmartHouse — Cómo correr el proyecto
 
## Requisitos
- Tener Docker y Docker Compose instalados
## 1. Crear el archivo `.env` en la raíz del proyecto
 
```env
POSTGRES_HOST=db
POSTGRES_USER=user
POSTGRES_PASSWORD=mysecretpassword
POSTGRES_DB=gorm
POSTGRES_PORT=5432
```
 
## 2. Levantar todo
 
```bash
docker compose up --build
```
 
Esto levanta el Front, el Middleware, el Load Balancer, las instancias de cada microservicio (back, servicios, inmuebles), las cuales son 2 instancias por microservicio, y la base de datos.
 
## 3. Abrir la app
 
```
http://localhost:3000
```
 
## Cómo apagar servicios
 
**Apagar todo:**
```bash
docker compose down
```
 
**Apagar y borrar los datos de la base de datos:**
```bash
docker compose down -v
```

