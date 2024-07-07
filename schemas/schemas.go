package schemas

// NamedAPIResourceList implements NamedAPIResourceList schema.
type NamedAPIResourceList struct {
	Count    int32              `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}

// NamedAPIResource implements NamedAPIResource schema.
type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PatientBase implements PatientBase schema.
type PatientBase struct {
	Name              string             `json:"name" binding:"required,min=1,max=100"`
	PersonalID        PersonalID         `json:"personal_id" binding:"required"`
	Gender            string             `json:"gender" binding:"oneof='unspecified' male female"`
	PhoneNumber       string             `json:"phone_number,omitempty" binding:"omitempty,e164"`
	Languages         []string           `json:"languages" binding:"max=10"`
	BirthDate         string             `json:"birth_date" binding:"required,datetime=2006-01-02"`
	EmergencyContacts []EmergencyContact `json:"emergency_contacts" binding:"max=10,dive"`
	ReferredBy        string             `json:"referred_by,omitempty" binding:"max=100"`
	SpecialNote       string             `json:"special_note,omitempty" binding:"max=500"`
}

// Patient implements Patient schema.
type Patient struct {
	PatientBase
	ID     int32 `json:"id"`
	Active bool  `json:"active"`
	Age    int32 `json:"age"`
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

type DoctorBase struct {
	Name         string   `json:"name" binding:"required,min=1,max=100"`
	Gender       string   `json:"gender" binding:"oneof='unspecified' male female"`
	PhoneNumber  string   `json:"phone_number" binding:"e164"`
	Specialities []string `json:"specialities" binding:"max=30"`
	SpecialNote  string   `json:"special_note,omitempty" binding:"max=500"`
}

// Doctor implements Doctor schema.
type Doctor struct {
	DoctorBase
	ID     int32 `json:"id"`
	Active bool  `json:"active"`
}

// AppointmentBase implements AppointmentBase schema.
type AppointmentBase struct {
	PatientID int32  `json:"patient_id,omitempty"`
	DoctorID  int32  `json:"doctor_id" binding:"required"`
	StartTime string `json:"start_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime   string `json:"end_time" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// Appointment implements Appointment schema.
type Appointment struct {
	AppointmentBase
	ID                int32 `json:"id"`
	ApprovedByPatient bool  `json:"approved_by_patient"`
	Visited           bool  `json:"visited"`
}

// IDHolder implements IDHolder schema.
type IDHolder struct {
	ID int32 `json:"id" binding:"required"`
}

// PatientIDHolder implements PatientIDHolder schema.
type PatientIDHolder struct {
	PatientID int32 `json:"patient_id" binding:"required"`
}

// ErrorResponse implements ErrorResponse schema.
type ErrorResponse struct {
	Message string `json:"message"`
}
