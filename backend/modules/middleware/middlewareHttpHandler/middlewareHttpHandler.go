package middlewarehttphandler

import middlewareusecase "main/modules/middleware/middlewareUsecase"

type (
	MiddlewareHttpHandlerService interface {

	}

	MiddlewareHttpHandler struct {
		middlewareUsecase middlewareusecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareUsecase(middlewareUsecase middlewareusecase.MiddlewareUsecaseService) MiddlewareHttpHandlerService {
	return &MiddlewareHttpHandler{middlewareUsecase}
}