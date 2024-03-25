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

// PersonalID implements PersonalID schema.
type PersonalID struct {
	ID   string `json:"id" validate:"required,min=1,max=100"`
	Type string `json:"type" validate:"required,min=1,max=100"`
}

// EmergencyContact implements EmergencyContact schema.
type EmergencyContact struct {
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Closeness string `json:"closeness" validate:"required,min=1,max=100"`
	Phone     string `json:"phone" validate:"e164"`
}

// PatientBase implements PatientBase schema.
type PatientBase struct {
	Name              string             `json:"name" validate:"required,min=1,max=100"`
	PersonalID        PersonalID         `json:"personal_id" validate:"required"`
	Gender            string             `json:"gender" validate:"oneof='unspecified' male female"`
	PhoneNumber       string             `json:"phone_number" validate:"e164"`
	Languages         []string           `json:"languages" validate:"max=10"`
	BirthDate         string             `json:"birth_date" validate:"required,datetime=2000-12-25"`
	EmergencyContacts []EmergencyContact `json:"emergency_contacts" validate:"max=10"`
	ReferredBy        string             `json:"referred_by" validate:"max=100"`
	SpecialNote       string             `json:"special_note" validate:"max=500"`
}

// Patient implements Patient schema.
type Patient struct {
	PatientBase
	ID     int32 `json:"id"`
	Active bool  `json:"active"`
	Age    int32 `json:"age"`
}
