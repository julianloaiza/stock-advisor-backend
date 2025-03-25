# Stock Advisor Backend

Stock Advisor Backend is a robust Go-based API for managing and querying stock market data, designed with clean architecture and hexagonal architecture principles.

![Swagger](capture.png)

## Features

- **RESTful API** for stock market data retrieval
- **Advanced Filtering**: Search and filter stocks by multiple criteria
- **Intelligent Recommendation Algorithm**: Score stocks based on target prices and ratings
- **Data Synchronization**: Efficient sync with external data sources
- **Database Agnostic**: Designed with GORM for flexible database support
- **Comprehensive Swagger Documentation**
- **Dependency Injection** using Uber FX
- **CORS Support**

## Technologies

- **Go 1.23+**
- **Echo Framework**
- **GORM**
- **PostgreSQL/CockroachDB**
- **Uber FX**
- **Swagger**
- **Testify**

## Requirements

- Go 1.23 or higher
- PostgreSQL or CockroachDB
- External Stock Data API (configured in `.env`)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/julianloaiza/stock-advisor-backend.git
cd stock-advisor-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Create and configure `.env` file:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Generate Swagger documentation:
```bash
swag init
```

## Running with Docker

You can run the application using Docker:

```bash
# Build the image
docker build -t stock-advisor-backend .

# Run the container
docker run -p 8080:8080 \
  -e DATABASE_URL=postgresql://user:password@host/database \
  -e STOCK_API_URL=https://api.example.com \
  stock-advisor-backend
  ...
```

### Full Deployment

For a complete application deployment, visit:
[julianloaiza/stock-advisor-deployment](https://github.com/julianloaiza/stock-advisor-deployment)

## Configuration

Configure the following in `.env`:
- `DATABASE_URL`: Database connection string
- `STOCK_API_URL`: External stock data API URL
- `STOCK_AUTH_TKN`: Authentication token for external API
- `SYNC_MAX_ITERATIONS`: Maximum sync iterations
- `SYNC_TIMEOUT`: Sync operation timeout
- `CORS_ALLOWED_ORIGINS`: Allowed CORS origins

You can also configure the recommendation algorithm using the `recommendation_factors.json` file.

## Running the Application

```bash
# Run the application
go run main.go
```

## Testing

```bash
# Run all tests
go test ./...

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API Documentation

Access Swagger documentation at:
`http://localhost:8080/swagger/index.html`

## Project Structure

```
└── 📁stock-advisor
    ├── 📁config               # Application configuration management
        └── config.go          # Loads and validates application configuration
    ├── 📁database             # Database connection setup
        └── database.go        # Establishes and manages database connection
    ├── 📁docs                 # Swagger documentation
    ├── 📁internal             # Core application logic
        ├── 📁domain           # Domain models and core entities
            └── stock.go       # Stock entity definition
        ├── 📁httpapi          # HTTP API layer
            ├── 📁handlers     # HTTP request handlers
                ├── handlers.go         # Base handler interface
                ├── 📁response          # API response utilities
                    └── response.go     # Standard API response structures
                └── 📁stocks            # Stock-specific handlers
                    ├── get.go          # GET stocks handler
                    ├── stocks.go       # Handler module configuration
                    └── sync.go         # Stock synchronization handler
            ├── httpapi.go             # HTTP API module configuration
            └── 📁middleware           # HTTP middleware
                └── cors.go            # CORS configuration
        ├── 📁repositories     # Data access layer
            ├── repositories.go        # Repository module configuration
            └── 📁stocks       # Stock-specific repositories
                ├── get.go             # Stock retrieval repository methods
                ├── stocks.go          # Repository module configuration
                └── sync.go            # Stock synchronization repository methods
        └── 📁services         # Business logic layer
            ├── 📁apiClient    # Client for external API communication
                ├── apiClient.go       # Client definitions and initialization
                └── get.go             # GET request implementation
            ├── services.go            # Services module configuration
            └── 📁stocks       # Stock-specific services
                ├── get.go             # Stock retrieval service logic
                ├── stocks.go          # Service module configuration
                ├── sync_parser.go     # Data transformation during synchronization
                ├── sync_recommendation.go # Recommendation scoring algorithm
                └── sync.go            # Stock synchronization service logic
    ├── recommendation_factors.json    # Recommendation algorithm configuration
    ├── .env                   # Environment configuration (local)
    ├── .env.example           # Example environment configuration
    ├── Dockerfile             # Docker container configuration
    ├── go.mod                 # Go module dependencies
    └── main.go                # Application entry point
```

## API Endpoints

- `GET /stocks`: Retrieve stocks with advanced filtering
- `POST /stocks/sync`: Synchronize stocks from external source
- `GET /swagger/*`: Swagger documentation

### GET /stocks Endpoint

#### Input Parameters (Query Params)
- `query` (optional): General search text
  - Searches in: ticker, company, brokerage, action, ratings
- `page` (optional): Page number 
  - Default value: 1
- `size` (optional): Number of records per page
  - Default value: 10
- `recommends` (optional): Order by recommendation score
  - Values: `true` or `false`
  - Default value: `false`
- `minTargetTo` (optional): Minimum target price
- `maxTargetTo` (optional): Maximum target price
- `currency` (optional): Price currency
  - Default value: "USD"

#### Example Request
```
GET /stocks?query=AAPL&page=1&size=10&recommends=true&minTargetTo=150&maxTargetTo=200&currency=USD
```

#### Successful Response (200 OK)
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
        "action": "upgraded by",
        "rating_from": "Hold",
        "rating_to": "Buy", 
        "target_from": 150,
        "target_to": 180,
        "currency": "USD",
        "recommend_score": 36.125
      }
    ],
    "total": 1000,
    "page": 1,
    "size": 10
  },
  "message": "Stock query successful"
}
```

### Recommendation Algorithm

The system calculates a `recommend_score` for each stock based on multiple factors:

1. **Percentage difference between target prices**: Higher increases receive higher scores
2. **Analyst ratings**: Upgrades to "Buy" and "Strong-Buy" are prioritized
3. **Action type**: Different scores are assigned to actions like "upgraded by", "target raised by", etc.
4. **Company and brokerage factors**: Configurable from `recommendation_factors.json`

This score allows sorting results when using the `recommends=true` parameter.

### POST /stocks/sync Endpoint

#### Input Parameters
```json
{
  "limit": 5  // Number of sync iterations
}
```

#### Constraints
- `limit` must be a positive integer
- Default value: 1
- Maximum configurable in server settings (default: 100)

#### Example Request
```json
{
  "limit": 5
}
```

#### Successful Response (200 OK)
```json
{
  "code": 200,
  "message": "Synchronization completed successfully"
}
```

#### Possible Errors
- 400 Bad Request: 
  - Invalid limit
  - Error reading request body
- 500 Internal Server Error: 
  - Error during synchronization with external API

#### Important Notes
- Each iteration updates approximately 10 stock records
- Synchronization COMPLETELY replaces existing data
- The operation cannot be undone once completed
- During synchronization, recommendation scores are calculated and stored in the database

## Data Flow

### Stock Query Flow
1. HTTP request arrives at the `GetStocks` handler
2. Handler validates and processes parameters
3. Stock service applies business logic
4. Repository performs database query
5. Results are transformed and returned to the client

### Synchronization Flow
1. HTTP request arrives at the `SyncStocks` handler
2. Stock service coordinates synchronization
3. API client fetches data from external source
4. Parser transforms data to internal format
5. Recommendation algorithm calculates scores
6. Repository replaces all data in the database