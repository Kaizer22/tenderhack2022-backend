package entity

import (
	"fmt"
)

type userRole string

const (
	AnonymRole   userRole = "anonym"
	CustomerRole userRole = "customer"
	AdminRole    userRole = "admin"
	ProviderRole userRole = "provider"
)

type Account struct {
	Id        int64    `pg:"id,pk" json:"id"`
	Login     string   `pg:"login,unique" json:"login"`
	Password  string   `pg:"password" json:"password"`
	Role      userRole `pg:"role" json:"role"`
	ProfileID int64    `pg:"profile_id" json:"profile_id"`
}

type AccountData struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Role     userRole `json:"role"`
}

func (a *Account) SetPassword(password string) error {

	if len(password) < 8 {
		return fmt.Errorf(
			"password for %s must be at least 8 characters",
			a.Login)
	}
	a.Password = password
	return nil
}

func (a *Account) InvalidPassword(password string) bool {
	if len(password) < 8 {
		return true
	}
	return a.Password != password
}

func (a *Account) InvalidRole(role string) bool {
	return role != string(AnonymRole) &&
		role != string(ProviderRole) &&
		role != string(AdminRole) &&
		role != string(CustomerRole)
}
