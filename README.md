# Job Tracker API

API REST para gestionar postulaciones/“jobs” (crear, consultar, actualizar, eliminar y filtrar por estado).  
Construida en **Go** con **Gin**, persistencia en **PostgreSQL** (via GORM) y observabilidad con **OpenTelemetry**. El repositorio incluye un `compose.yml` con **Postgres (pgvector)** y **Grafana Alloy** como collector OTEL.

---

## Requisitos

- **Go** >= 1.25
- **Docker** y **Docker Compose** (recomendado para Postgres + Alloy)

---

## Configuración

1. Crea tu archivo `.env` a partir del ejemplo:
```bash 
cp .env.example .env
```

2. Ajusta variables (valores de ejemplo en `.env.example`):

- `PORT` (por defecto `8080`)
- Base de datos:
    - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- OpenTelemetry:
    - `OTEL_EXPORTER_OTLP_ENDPOINT` (por defecto `http://localhost:4318`)
    - `OTEL_EXPORTER_OTLP_PROTOCOL` (por defecto `http/protobuf`)
    - `OTEL_SERVICE_NAME`
    - `OTEL_RESOURCE_ATTRIBUTES`
- Grafana Alloy:
    - `GRAFANA_INSTANCE_ID` (placeholder)
- `GCLOUD_RW_API_KEY` (placeholder)

> Nota: usa placeholders para credenciales/keys y gestiona secretos con tu herramienta preferida.

---

## Levantar dependencias (Postgres + Alloy)

Desde la raíz del proyecto:
```bash 
docker compose up -d
``` 

Servicios expuestos por `compose.yml`:

- **Postgres**: `localhost:5432`
- **Alloy**:
  - OTEL gRPC: `localhost:4317`
  - OTEL HTTP: `localhost:4318`
  - UI/health: `http://localhost:12345`

---

## Ejecutar la API

Entrada principal:

- `cmd/api/main.go` (llama a `bootstrap.Start()`)

Ejecuta:
```bash 
go run ./cmd/api
```

La API debería escuchar en `http://localhost:${PORT}` (por defecto `http://localhost:8080`).

---

## Postman

Se incluye una colección lista para usar:

- `postman_collection.json`

Variables típicas:
- `BASE_URL` → `http://localhost:8080`

---

## Estructura del proyecto (alto nivel)

- `cmd/api` — main del servicio
- `internal/application` — casos de uso / servicios
- `internal/domain` — entidades y errores de dominio
- `internal/infrastructure` — handlers HTTP y repositorios
- `internal/bootstrap` — arranque (config, wiring, DB, OTEL)
- `docker-config` — configuración de Alloy
- `tests` — pruebas

---

## Troubleshooting

- **La API no conecta a Postgres**: revisa que `docker compose up -d` esté levantado y que `.env` tenga `DB_*` correctos.
- **OTEL no exporta**: confirma que Alloy esté healthy en `http://localhost:12345/-/healthy` y que `OTEL_EXPORTER_OTLP_ENDPOINT` apunte a `http://localhost:4318`.

---
