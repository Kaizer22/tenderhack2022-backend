package entity

type ProfileType string

const (
	ProfileTypeAdmin    ProfileType = "admin"
	ProfileTypeProvider ProfileType = "provider"
	ProfileTypeCustomer ProfileType = "customer"
)

type Profile struct {
	OrganizationId   int64       `pg:"id,pk"`
	OrganizationName string      `pg:"org_name"`
	OrganizationType ProfileType `pg:"org_type"`
}

type ProfileData struct {
	OrganizationName string      `pg:"org_name"`
	OrganizationType ProfileType `pg:"org_type"`
}
