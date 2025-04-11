package interfaces

//go:generate mockgen -source=services.go -destination=../../mocks/services/users_service.go -package=mockservices -exclude_interfaces=AuthService
type UsersService interface {
	UsersRepository
}

//go:generate mockgen -source=services.go -destination=../../mocks/services/auth_service.go -package=mockservices -exclude_interfaces=UsersService
type AuthService interface {
	AuthRepository
}
