package entity

type ProfileType string

const (
	ProfileTypeProvider ProfileType = "provider"
	ProfileTypeCustomer ProfileType = "customer"
)

// Profile example
type Profile struct {
	OrganizationId   int64       `pg:"id,pk"`
	OrganizationName string      `pg:"org_name"`
	OrganizationType ProfileType `pg:"org_type"`
}

// ProfileData example
type ProfileData struct {
	OrganizationName string      `pg:"org_name"`
	OrganizationType ProfileType `pg:"org_type"`
}
