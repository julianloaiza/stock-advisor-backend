package response

// APIResponse define el formato unificado para las respuestas de la API.
// Incluye un código (por ejemplo, 200, 501, etc.), la data (que puede ser cualquier estructura),
// un mensaje descriptivo y, opcionalmente, detalles del error.
type APIResponse struct {
	Code    int         `json:"code"`              // Código HTTP o de aplicación, ej. 200, 501, etc.
	Data    interface{} `json:"data,omitempty"`    // La data, que puede ser una estructura simple o paginada.
	Message string      `json:"message,omitempty"` // Mensaje descriptivo.
	Error   string      `json:"error,omitempty"`   // Detalles del error, si lo hubiera.
}

// PaginatedData es una estructura para contener respuestas paginadas.
// Se coloca dentro de APIResponse.Data cuando la consulta requiere paginación.
type PaginatedData struct {
	Content interface{} `json:"content"` // El listado real de ítems.
	Total   int64       `json:"total"`   // Total de ítems disponibles en la consulta.
	Page    int         `json:"page"`    // Número de página actual.
	Size    int         `json:"size"`    // Tamaño de página (número de ítems por página).
}
