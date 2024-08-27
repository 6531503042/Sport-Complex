package middlewarerepository

type (
	MiddlewareRepositoryService interface {

	}

	middlewarerepository struct {

	}
)

func NewMiddlewareRepository() MiddlewareRepositoryService {
	return &middlewarerepository{}
}