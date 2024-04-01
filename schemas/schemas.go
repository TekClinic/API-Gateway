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
	ID   string `json:"id" binding:"required,min=1,max=100"`
	Type string `json:"type" binding:"required,min=1,max=100"`
}

// EmergencyContact implements EmergencyContact schema.
type EmergencyContact struct {
	Name      string `json:"name" binding:"required,min=1,max=100"`
	Closeness string `json:"closeness" binding:"required,min=1,max=100"`
	Phone     string `json:"phone" binding:"required,e164"`
}

// IDHolder implements IDHolder schema.
type IDHolder struct {
	ID int32 `json:"id"`
}

// PatientBase implements PatientBase schema.
type PatientBase struct {
	Name              string             `json:"name" binding:"required,min=1,max=100"`
	PersonalID        PersonalID         `json:"personal_id" binding:"required"`
	Gender            string             `json:"gender" binding:"oneof='unspecified' male female"`
	PhoneNumber       string             `json:"phone_number" binding:"omitempty,e164"`
	Languages         []string           `json:"languages" binding:"max=10"`
	BirthDate         string             `json:"birth_date" binding:"required,datetime=2006-01-02"`
	EmergencyContacts []EmergencyContact `json:"emergency_contacts" binding:"max=10,dive"`
	ReferredBy        string             `json:"referred_by" binding:"max=100"`
	SpecialNote       string             `json:"special_note" binding:"max=500"`
}

// Patient implements Patient schema.
type Patient struct {
	PatientBase
	ID     int32 `json:"id"`
	Active bool  `json:"active"`
	Age    int32 `json:"age"`
}

// AppointmentBase implements AppointmentBase schema.
type AppointmentBase struct {
	PatientID int32  `json:"patient_id"`
	DoctorID  int32  `json:"doctor_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// Appointment implements Appointment schema.
type Appointment struct {
	ID                int32  `json:"id"`
	PatientID         int32  `json:"patient_id"`
	DoctorID          int32  `json:"doctor_id"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	ApprovedByPatient bool   `json:"approved_by_patient"`
	Visited           bool   `json:"visited"`
}
