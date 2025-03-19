package stocks

import (
	"log"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// GetStocks obtiene los stocks filtrados, aplicando paginación en la base de datos.
func (r *repository) GetStocks(query string, page, size int, recommends bool, minTargetTo, maxTargetTo float64, currency string) ([]domain.Stock, int64, error) {
	var stocks []domain.Stock
	var total int64

	// Calculamos el offset para la paginación
	offset := (page - 1) * size

	// Construimos la consulta base
	dbQuery := r.buildBaseQuery(query, minTargetTo, maxTargetTo, currency)

	// Si se solicitan recomendaciones, ordenamos por el puntaje de recomendación en orden descendente
	if recommends {
		dbQuery = dbQuery.Order("recommend_score DESC")
	}

	// Contamos el total de registros sin paginar
	if err := dbQuery.Count(&total).Error; err != nil {
		log.Printf("Error contando registros: %v", err)
		return nil, 0, err
	}

	// Aplicamos paginación
	if err := dbQuery.
		Offset(offset).
		Limit(size).
		Find(&stocks).Error; err != nil {
		log.Printf("Error obteniendo stocks: %v", err)
		return nil, 0, err
	}

	return stocks, total, nil
}

// buildBaseQuery construye la consulta base con todos los filtros aplicados
func (r *repository) buildBaseQuery(query string, minTargetTo, maxTargetTo float64, currency string) *gorm.DB {
	// Preparar filtro de búsqueda
	likeQuery := "%" + query + "%"

	// Construir consulta base con filtro de texto
	dbQuery := r.db.Model(&domain.Stock{}).
		Where("ticker ILIKE ? OR company ILIKE ? OR brokerage ILIKE ? OR action ILIKE ? OR rating_from ILIKE ? OR rating_to ILIKE ?",
			likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery)

	// Aplicar filtro de currency
	if currency != "" {
		dbQuery = dbQuery.Where("currency = ?", currency)
	}

	// Aplicar filtros de target_to si están especificados
	if minTargetTo > 0 {
		dbQuery = dbQuery.Where("target_to >= ?", minTargetTo)
	}

	if maxTargetTo > 0 {
		dbQuery = dbQuery.Where("target_to <= ?", maxTargetTo)
	}

	return dbQuery
}
