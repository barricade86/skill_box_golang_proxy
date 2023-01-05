package router

import (
	"github.com/go-chi/chi/v5"
	"webserver/internal/storage"
)

// RouterService
type RouterService struct {
	router      *chi.Mux
	dataStorage storage.Storage
}

// NewRouterService Creates and returns new instance of RouterService
func NewRouterService(dataStorage storage.Storage) *RouterService {
	return &RouterService{
		router:      chi.NewRouter(),
		dataStorage: dataStorage,
	}
}

// InitRoutes
func (rs *RouterService) Init() *chi.Mux {
	rs.router.Post("/user/create", rs.create)
	rs.router.Delete("/user/delete", rs.delete)
	rs.router.Post("/user/friends/add", rs.addFriendsForUser)
	rs.router.Get("/user/{userID}/friends/all", rs.getFriendsForUser)
	rs.router.Put("/user/{userID}/change", rs.updateAgeByUserId)

	return rs.router
}
