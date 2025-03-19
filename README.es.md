# Stock Advisor Backend

Stock Advisor Backend es una API robusta desarrollada en Go para gestionar y consultar datos del mercado de valores, diseñada con principios de arquitectura limpia y arquitectura hexagonal.

![Swagger](capture.png)

## Características

- **API RESTful** para recuperación de datos de acciones
- **Filtrado Avanzado**: Búsqueda y filtrado de acciones por múltiples criterios
- **Algoritmo de Recomendación Inteligente**: Puntuación de acciones basada en precios objetivo y calificaciones
- **Sincronización de Datos**: Sincronización de acciones desde fuentes de datos externas
- **Base de Datos Agnóstica**: Diseñada con GORM para soporte de bases de datos flexible
- **Documentación Swagger Completa**
- **Inyección de Dependencias** usando Uber FX
- **Soporte CORS**

## Tecnologías

- **Go 1.23+**
- **Framework Echo**
- **GORM**
- **PostgreSQL**
- **Uber FX**
- **Swagger**
- **Testify**

## Requisitos

- Go 1.23 o superior
- PostgreSQL
- API externa de datos de acciones (configurada en `.env`)

## Instalación

1. Clonar el repositorio:
```bash
git clone https://github.com/julianloaiza/stock-advisor-backend.git
cd stock-advisor-backend
```

2. Instalar dependencias:
```bash
go mod download
```

3. Crear y configurar archivo `.env`:
```bash
cp .env.example .env
# Editar .env con tu configuración
```

4. Generar documentación Swagger:
```bash
swag init
```

## Configuración

Configurar lo siguiente en `.env`:
- `DATABASE_URL`: Cadena de conexión a PostgreSQL
- `STOCK_API_URL`: URL de la API externa de datos de acciones
- `STOCK_API_KEY`: Clave de autenticación de la API
- `SYNC_MAX_ITERATIONS`: Máximo de iteraciones de sincronización
- `SYNC_TIMEOUT`: Tiempo de espera de la operación de sincronización
- `CORS_ALLOWED_ORIGINS`: Orígenes permitidos para CORS

## Ejecutando la Aplicación

```bash
# Ejecutar la aplicación
go run main.go
```

## Pruebas

```bash
# Ejecutar todas las pruebas
go test ./...

# Generar reporte de cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Documentación de la API

Acceder a la documentación Swagger en:
`http://localhost:8080/swagger/index.html`

## Estructura del Proyecto

```
└── 📁stock-advisor
    ├── 📁config               # Gestión de configuración de la aplicación
        └── config.go          # Carga y valida la configuración de la aplicación
    ├── 📁database             # Configuración de conexión a base de datos
        └── database.go        # Establece y gestiona la conexión a base de datos
    ├── 📁docs                 # Documentación Swagger
        ├── docs.go            # Documentación Swagger generada
        ├── swagger.json       # Especificación Swagger en JSON
        └── swagger.yaml       # Especificación Swagger en YAML
    ├── 📁internal             # Lógica central de la aplicación
        ├── 📁domain           # Modelos de dominio y entidades centrales
            └── stock.go       # Definición de entidad Stock
        ├── 📁httpapi          # Capa de API HTTP
            ├── 📁handlers     # Manejadores de solicitudes HTTP
                ├── handlers.go         # Interfaz base de manejadores
                ├── 📁response          # Utilidades de respuesta API
                    └── response.go     # Estructuras de respuesta API estándar
                └── 📁stocks            # Manejadores específicos de stocks
                    ├── get.go          # Manejador GET de stocks
                    ├── get_test.go     # Pruebas para manejador GET
                    ├── stocks.go       # Configuración y construcción del módulo de manejadores de stocks
                    ├── sync.go         # Manejador de sincronización de stocks
                    └── sync_test.go    # Pruebas para manejador de sincronización
            ├── httpapi.go             # Configuración del módulo de API HTTP
            └── 📁middleware           # Middleware HTTP
                └── cors.go            # Configuración de CORS
        ├── 📁repositories     # Capa de acceso a datos
            ├── repositories.go        # Configuración del módulo de repositorios
            └── 📁stocks       # Repositorios específicos de stocks
                ├── get.go             # Métodos de recuperación de stocks
                ├── get_test.go        # Pruebas de recuperación de stocks
                ├── stocks.go          # Configuración y construcción del módulo de repositorios de stocks
                ├── sync.go            # Métodos de sincronización de stocks
                └── sync_test.go       # Pruebas de métodos de sincronización
        └── 📁services         # Capa de lógica de negocio
            ├── services.go            # Configuración del módulo de servicios
            └── 📁stocks       # Servicios específicos de stocks
                ├── get.go             # Lógica de recuperación de stocks
                ├── get_test.go        # Pruebas de servicio de recuperación
                ├── recommendation.go  # Algoritmo de recomendación de stocks
                ├── recommendation_test.go # Pruebas del algoritmo de recomendación
                ├── stocks.go          # Configuración y construcción del módulo de servicios de stocks
                ├── sync.go            # Lógica de sincronización de stocks
                └── sync_test.go       # Pruebas de servicio de sincronización
    ├── .env                   # Configuración de entorno (local)
    ├── .env.example           # Ejemplo de configuración de entorno
    ├── .gitignore             # Archivo de ignorados de Git
    ├── Dockerfile             # Configuración de contenedor Docker
    ├── go.mod                 # Dependencias del módulo Go
    ├── go.sum                 # Versiones exactas de dependencias
    └── main.go                # Punto de entrada de la aplicación
```

## Endpoints de la API

- `GET /stocks`: Recuperar stocks con filtrado avanzado
- `POST /stocks/sync`: Sincronizar stocks desde fuente externa
- `GET /swagger/*`: Documentación Swagger

### Endpoint GET /stocks

#### Parámetros de Entrada (Parámetros de Consulta)
- `query` (opcional): Texto de búsqueda general
  - Busca en: ticker, compañía, casa de bolsa, acción, calificaciones
- `page` (opcional): Número de página 
  - Valor por defecto: 1
- `size` (opcional): Número de registros por página
  - Valor por defecto: 10
- `recommends` (opcional): Ordenar por puntuación de recomendación
  - Valores: `true` o `false`
  - Valor por defecto: `false`
- `minTargetTo` (opcional): Precio objetivo mínimo
- `maxTargetTo` (opcional): Precio objetivo máximo
- `currency` (opcional): Moneda de los precios
  - Valor por defecto: "USD"

#### Ejemplo de Solicitud
```
GET /stocks?query=AAPL&page=1&size=10&recommends=true&minTargetTo=150&maxTargetTo=200&currency=USD
```

#### Respuesta Exitosa (200 OK)
```json
{
  "code": 200,
  "data": {
    "content": [
      {
        "id": 1054506709730689025,
        "ticker": "AAPL",
        "company": "Apple Inc.",
        "brokerage": "Goldman Sachs",
        "action": "actualizado por",
        "rating_from": "Mantener",
        "rating_to": "Comprar", 
        "target_from": 150,
        "target_to": 180,
        "currency": "USD",
        "time": "2025-02-26T19:30:06.366255-05:00"
      }
    ],
    "total": 1000,
    "page": 1,
    "size": 10
  },
  "message": "Consulta de acciones exitosa"
}
```

### Endpoint POST /stocks/sync

#### Parámetros de Entrada
```json
{
  "limit": 5  // Número de iteraciones de sincronización
}
```

#### Restricciones
- `limit` debe ser un número entero positivo
- Valor por defecto: 1
- Máximo configurable en la configuración del servidor (por defecto: 100)

#### Ejemplo de Solicitud
```json
{
  "limit": 5
}
```

#### Respuesta Exitosa (200 OK)
```json
{
  "code": 200,
  "message": "Sincronización completada exitosamente"
}
```

#### Posibles Errores
- 400 Bad Request: 
  - Límite inválido
  - Error al leer el cuerpo de la solicitud
- 500 Internal Server Error: 
  - Error durante la sincronización con la API externa

#### Notas Importantes
- Cada iteración actualiza aproximadamente 10 registros de acciones
- La sincronización REEMPLAZA COMPLETAMENTE los datos existentes
- La operación no se puede deshacer una vez completada