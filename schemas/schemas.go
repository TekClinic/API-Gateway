package schemas

// NamedAPIResource implements NamedAPIResource schema.
type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NamedAPIResourceList implements NamedAPIResourceList schema.
type NamedAPIResourceList struct {
	Count    int32              `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

// ErrorResponse implements ErrorResponse schema.
type ErrorResponse struct {
	Message string `json:"message"`
}
