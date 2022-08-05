package object

type (
	// struct of request url query
	Query struct {
		// is only media
		OnlyMedia string
		// max id
		MaxID string
		// since id
		SinceID string
		// limit
		Limit string
	}
)
