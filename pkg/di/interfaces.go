package di

import "go-api-posts/internal/user"

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserRepository interface {
	FindByEmail(email string) (*user.User, error)
}
