package emqx

type Meta struct {
	Count   uint64 `json:"count"`
	Hasnext bool   `json:"hasnext"`
	Limit   uint16 `json:"limit"`
	Page    uint16 `json:"page"`
}
