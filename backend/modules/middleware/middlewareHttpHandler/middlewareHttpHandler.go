package middlewarehttphandler

import middlewareusecase "main/modules/middleware/middlewareUsecase"

type (
	middlewareHttpHandlerService interface {

	}

	middlewareHttpHandler struct {
		middlewareUsecase middlewareusecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareUsecase(middlewareUsecase middlewareusecase.MiddlewareUsecaseService) middlewareHttpHandlerService {
	return &middlewareHttpHandler{middlewareUsecase: }
}