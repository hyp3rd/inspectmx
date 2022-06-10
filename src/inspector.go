package inspectmx

type Inspector struct {
	Email    string  `json:"email" validate:"required,email,min=6,max=64" san:"min=6,max=64,trim,lower"`
	Provider *string `json:"provider,omitempty"`
	Valid    bool    `json:"valid,omitempty"`
	Status   string  `json:"status,omitempty"`
	// IPAddress string `json:"ip_address"`
}
