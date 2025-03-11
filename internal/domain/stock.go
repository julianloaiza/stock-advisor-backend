package domain

// Stock representa la información y recomendaciones de un valor.
type Stock struct {
	Ticker     string
	TargetFrom string
	TargetTo   string
	Company    string
	Action     string
	Brokerage  string
	RatingFrom string
	RatingTo   string
	Time       string // Considera usar time.Time para manejar fechas
}
