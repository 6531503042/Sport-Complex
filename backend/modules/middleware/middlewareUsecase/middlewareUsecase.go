package middlewareusecase

import (
	middlewareRepository "main/modules/middleware/middlewareRepository"
)

type (
	MiddlewareUsecaseService interface {
	}

	middlewareUsecase struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoryService
	}
	
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewareRepository}
}
