package http

type CreateLinkRequest struct {
	LinkValue string `json:"linkValue"`
	TTL       int64  `json:"ttl"`
}
