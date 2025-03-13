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

// recommendationScore calcula una puntuación para un stock que sirve para determinar su recomendación.
// La lógica es sumar la diferencia entre target_to y target_from y aplicar un bonus basado en el rating.
func recommendationScore(s domain.Stock) float64 {
	diff := s.TargetTo - s.TargetFrom
	var bonus float64

	// Asignamos bonus base según rating_to (convertido a minúsculas)
	switch strings.ToLower(s.RatingTo) {
	case "buy":
		bonus = 10
	case "strong-buy":
		bonus = 15
	case "outperform":
		bonus = 12
	case "market perform":
		bonus = 5
	case "neutral", "hold":
		bonus = 5
	case "sector perform", "sector weight":
		bonus = 3
	case "overweight":
		bonus = 8
	case "equal weight":
		bonus = 2
	case "underweight":
		bonus = -5
	case "sell":
		bonus = -10
	default:
		bonus = 0
	}

	// Bonus adicional para cambios positivos: si rating_from es negativo ("sell" o "underweight")
	// y rating_to es muy positivo ("buy", "strong-buy" o "outperform"), añadimos +15.
	ratingFrom := strings.ToLower(s.RatingFrom)
	ratingTo := strings.ToLower(s.RatingTo)
	if (ratingFrom == "sell" || ratingFrom == "underweight") &&
		(ratingTo == "buy" || ratingTo == "strong-buy" || ratingTo == "outperform") {
		bonus += 15
	}

	return diff + bonus
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
