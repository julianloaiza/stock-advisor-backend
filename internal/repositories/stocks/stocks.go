package stocks

import (
	"log"
	"sort"
	"strings"

	"github.com/julianloaiza/stock-advisor/internal/domain"
	"gorm.io/gorm"
)

// Repository define las operaciones disponibles para manejar stocks.
type Repository interface {
	ReplaceAllStocks(stocks []domain.Stock) error
	// GetStocks busca stocks cuyos atributos contengan el string query (case-insensitive)
	// en los campos ticker, company, brokerage, action, rating_from y rating_to.
	// Aplica paginación a la consulta, devolviendo además el total de registros encontrados.
	// Si recommends es true, se reordena la data según un algoritmo de recomendación.
	GetStocks(query string, page, size int, recommends bool) ([]domain.Stock, int64, error)
}

type repository struct {
	db *gorm.DB
}

// New crea una nueva instancia del repositorio de stocks.
func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ReplaceAllStocks reemplaza la data existente en la tabla Stock por la nueva data,
// utilizando una transacción para asegurar atomicidad.
func (r *repository) ReplaceAllStocks(stocks []domain.Stock) error {
	log.Println("Reemplazando todos los stocks existentes")
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&domain.Stock{}).Error; err != nil {
			return err
		}
		return tx.Create(&stocks).Error
	})
}

// GetStocks busca stocks que coincidan con la query en los campos string indicados.
// Si recommends es false, se aplica paginación en la consulta SQL; si es true, se
// obtiene la data completa, se ordena en memoria usando nuestro algoritmo de recomendación
// y luego se aplica la paginación.
func (r *repository) GetStocks(query string, page, size int, recommends bool) ([]domain.Stock, int64, error) {
	var stocks []domain.Stock
	var total int64
	offset := (page - 1) * size
	likeQuery := "%" + query + "%"

	dbQuery := r.db.Model(&domain.Stock{}).
		Where("ticker ILIKE ? OR company ILIKE ? OR brokerage ILIKE ? OR action ILIKE ? OR rating_from ILIKE ? OR rating_to ILIKE ?",
			likeQuery, likeQuery, likeQuery, likeQuery, likeQuery, likeQuery)

	if !recommends {
		// Realiza el conteo y la consulta con paginación en SQL.
		if err := dbQuery.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		if err := dbQuery.Offset(offset).Limit(size).Find(&stocks).Error; err != nil {
			return nil, 0, err
		}
		return stocks, total, nil
	}

	// Si se solicita recomendación, obtenemos todos los registros que coinciden.
	var allMatching []domain.Stock
	if err := dbQuery.Find(&allMatching).Error; err != nil {
		return nil, 0, err
	}
	total = int64(len(allMatching))

	// Ordenamos la data usando nuestro algoritmo de recomendación.
	sort.Slice(allMatching, func(i, j int) bool {
		return recommendationScore(allMatching[i]) > recommendationScore(allMatching[j])
	})

	// Aplicamos paginación manualmente.
	if offset > len(allMatching) {
		return []domain.Stock{}, total, nil
	}
	end := offset + size
	if end > len(allMatching) {
		end = len(allMatching)
	}
	stocks = allMatching[offset:end]
	return stocks, total, nil
}

// recommendationScore calcula un puntaje de recomendación para una acción
// basada en varios factores, incluyendo la diferencia porcentual entre
// el precio objetivo inicial y final, la diferencia absoluta,
// el rating final, la transición de rating y la acción realizada.
// Devuelve un puntaje flotante que representa la recomendación.
func recommendationScore(s domain.Stock) float64 {
	// Diferencia porcentual
	var percentDiff float64
	if s.TargetFrom != 0 {
		percentDiff = ((s.TargetTo - s.TargetFrom) / s.TargetFrom) * 100
	} else {
		percentDiff = s.TargetTo
	}

	// Diferencia absoluta
	absoluteDiff := s.TargetTo - s.TargetFrom

	// Bonus por magnitud absoluta (USD)
	var absoluteBonus float64
	switch {
	case absoluteDiff >= 100:
		absoluteBonus = 10
	case absoluteDiff >= 50:
		absoluteBonus = 7
	case absoluteDiff >= 20:
		absoluteBonus = 5
	case absoluteDiff >= 10:
		absoluteBonus = 3
	case absoluteDiff >= 5:
		absoluteBonus = 2
	default:
		absoluteBonus = 0
	}

	// Bonus según el rating_to
	var ratingBonus float64
	switch strings.ToLower(s.RatingTo) {
	case "buy":
		ratingBonus = 10
	case "strong-buy":
		ratingBonus = 15
	case "outperform":
		ratingBonus = 12
	case "market perform", "neutral", "hold", "in-line":
		ratingBonus = 5
	case "sector perform", "sector weight":
		ratingBonus = 3
	case "overweight":
		ratingBonus = 8
	case "equal weight":
		ratingBonus = 2
	case "underweight":
		ratingBonus = -5
	case "sell":
		ratingBonus = -10
	case "underperform", "weak-sell":
		ratingBonus = -8
	default:
		ratingBonus = 0
	}

	// Bonus adicional por transición positiva o neutral a positiva
	ratingFrom := strings.ToLower(s.RatingFrom)
	ratingTo := strings.ToLower(s.RatingTo)
	var transitionBonus float64
	switch {
	case (ratingFrom == "sell" || ratingFrom == "underweight" || ratingFrom == "weak-sell") &&
		(ratingTo == "buy" || ratingTo == "strong-buy" || ratingTo == "outperform"):
		transitionBonus = 15
	case (ratingFrom == "neutral" || ratingFrom == "hold" || ratingFrom == "in-line") &&
		(ratingTo == "buy" || ratingTo == "strong-buy" || ratingTo == "outperform"):
		transitionBonus = 5
	case (ratingFrom == "sell" || ratingFrom == "underweight" || ratingFrom == "weak-sell") &&
		(ratingTo == "market perform" || ratingTo == "neutral" || ratingTo == "hold" || ratingTo == "in-line"):
		transitionBonus = 3
	case (ratingFrom == "neutral" || ratingFrom == "hold" || ratingFrom == "in-line") &&
		(ratingTo == "market perform" || ratingTo == "neutral" || ratingTo == "hold" || ratingTo == "in-line"):
		transitionBonus = 1
	default:
		transitionBonus = 0
	}

	// Bonus según la acción realizada
	var actionBonus float64
	switch strings.ToLower(s.Action) {
	case "target raised by":
		actionBonus = 5
	case "upgraded by":
		actionBonus = 7
	case "initiated by":
		actionBonus = 3
	case "reiterated by":
		actionBonus = 2
	case "target lowered by":
		actionBonus = -5
	case "downgraded by":
		actionBonus = -7
	default:
		actionBonus = 0
	}

	// Puntaje final equilibrado
	score := percentDiff + absoluteBonus + ratingBonus + transitionBonus + actionBonus

	return score
}

// Nota: El algoritmo de recomendación actual calcula en tiempo real un "score" para cada acción,
// combinando la diferencia entre target_to y target_from con un bonus basado en el cambio de rating.
// Esta estrategia permite ordenar los stocks de mejor a peor recomendados según la información de la base de datos.
//
// Sin embargo, para mejorar aún más la calidad de las recomendaciones, se podrían implementar algunas mejoras:
//
// 1. Análisis por bróker:
//    - Se podría analizar el historial de recomendaciones de cada bróker (brokerage) para determinar cuáles
//      tienen un desempeño superior (por ejemplo, aquellas firmas que consistentemente actualizan a "Buy" o
//      suben significativamente el precio objetivo).
//    - De esta forma, se podría asignar un bonus adicional o mayor peso a los stocks recomendados por brókeres
//      con un historial comprobado de aciertos.
//
// 2. Lista de empresas reconocidas:
//    - Mantener una lista predefinida de empresas reconocidas (por ejemplo, de gran capitalización o con buen
//      historial) y, al comparar, asignar un bonus extra a los stocks de dichas empresas.
//    - Esto puede ayudar a filtrar o priorizar acciones que, por su reputación, tengan un mayor potencial de inversión.
//
// 3. Pre-cálculo del score:
//    - Otra opción sería calcular el score de recomendación en el momento en que se puebla el dataset (o en
//      procesos batch periódicos) y almacenarlo en la base de datos.
//    - Así, las consultas futuras podrían simplemente ordenar por el score pre-calculado, lo que mejoraría la
//      eficiencia de la búsqueda y la recomendación.
//
// Dado que el reto pedía específicamente desarrollar un algoritmo en tiempo real para la recomendación,
// se implementó la solución actual. Estas mejoras son sugerencias para futuras optimizaciones y enriquecimientos
// de la lógica de recomendación.
