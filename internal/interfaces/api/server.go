package api

type Server interface {
	RegisterServerRoutes()
	StartServer() error
}
