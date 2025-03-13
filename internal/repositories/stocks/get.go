package stocks

import (
	"github.com/julianloaiza/stock-advisor/internal/domain"
)

// GetStocks obtiene los stocks filtrados, aplicando paginación en la base de datos.
func (r *repository) GetStocks(query string, minTargetTo, maxTargetTo float64, page, size int) ([]domain.Stock, int64, error) {
	var stocks []domain.Stock
	var total int64
	likeQuery := "%" + query + "%"
	offset := (page - 1) * size

	// Construcción de la consulta con filtros
	dbQuery := r.db.Model(&domain.Stock{}).
		Where("ticker ILIKE ? OR company ILIKE ? OR brokerage ILIKE ? OR action ILIKE ? OR rating_from ILIKE ? OR rating_to ILIKE ?",
			likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery)

	// Aplicamos filtro por rango de Target To si se especifica
	if minTargetTo > 0 {
		dbQuery = dbQuery.Where("target_to >= ?", minTargetTo)
	}
	if maxTargetTo > 0 {
		dbQuery = dbQuery.Where("target_to <= ?", maxTargetTo)
	}

	// Contamos el total de registros sin paginar
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicamos paginación directamente en la base de datos
	if err := dbQuery.Offset(offset).Limit(size).Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}
