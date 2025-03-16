package response

// APIResponse define el formato unificado para las respuestas de la API.
type APIResponse struct {
	Code    int         `json:"code"`              // Código HTTP
	Data    interface{} `json:"data,omitempty"`    // Datos de respuesta
	Message string      `json:"message,omitempty"` // Mensaje descriptivo
	Error   string      `json:"error,omitempty"`   // Detalles del error
}

// PaginatedData estructura para respuestas paginadas.
type PaginatedData struct {
	Content interface{} `json:"content"` // Lista de ítems
	Total   int64       `json:"total"`   // Total de ítems disponibles
	Page    int         `json:"page"`    // Número de página actual
	Size    int         `json:"size"`    // Ítems por página
}

// NewSuccess crea una respuesta exitosa.
func NewSuccess(code int, data interface{}, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

// NewError crea una respuesta de error.
func NewError(code int, message, errorDetail string) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Error:   errorDetail,
	}
}

// NewPaginated crea una respuesta paginada.
func NewPaginated(content interface{}, total int64, page, size int) PaginatedData {
	return PaginatedData{
		Content: content,
		Total:   total,
		Page:    page,
		Size:    size,
	}
}
