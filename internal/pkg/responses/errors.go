package responses

// Error is a JSON representation defined by https://datatracker.ietf.org/doc/html/rfc7807#section-3.1
type Error struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}
