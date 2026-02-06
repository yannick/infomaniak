package api

// Response wraps every Infomaniak API response.
type Response[T any] struct {
	Result string     `json:"result"`
	Data   T          `json:"data,omitempty"`
	Error  *ErrorBody `json:"error,omitempty"`
}

// ErrorBody contains error details from the API.
type ErrorBody struct {
	Code        string        `json:"code"`
	Description string        `json:"description"`
	Errors      []ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail provides granular validation error context.
type ErrorDetail struct {
	Code        string            `json:"code"`
	Description string            `json:"description"`
	Context     map[string]string `json:"context,omitempty"`
}

// Domain represents a domain returned by the Infomaniak API.
type Domain struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	TLD       string         `json:"tld"`
	IsPremium bool           `json:"is_premium"`
	CreatedAt int64          `json:"created_at"`
	ExpiresAt int64          `json:"expires_at"`
	Options   DomainOptions  `json:"options"`
	Contacts  DomainContacts `json:"contacts"`
}

// DomainOptions holds optional flags for a domain.
type DomainOptions struct {
	DNSAnycast      bool `json:"dns_anycast"`
	RenewalWarranty bool `json:"renewal_warranty"`
	DomainPrivacy   bool `json:"domain_privacy"`
	DNSSEC          bool `json:"dnssec"`
}

// DomainContacts contains all contact types for a domain.
type DomainContacts struct {
	Owner   *Contact `json:"owner,omitempty"`
	Admin   *Contact `json:"admin,omitempty"`
	Tech    *Contact `json:"tech,omitempty"`
	Billing *Contact `json:"billing,omitempty"`
}

// Contact represents a domain contact record.
type Contact struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Phone       string `json:"phone"`
	Fax         string `json:"fax"`
	Email       string `json:"email"`
	IsValidated bool   `json:"is_validated"`
	ValidatedAt *int64 `json:"validated_at"`
	CreatedAt   int64  `json:"created_at"`
}

// UpdateNameserversInput is the request body for updating nameservers.
type UpdateNameserversInput struct {
	Nameservers          []string `json:"nameservers"`
	VerifyNSAvailability bool     `json:"verify_ns_availability"`
}
