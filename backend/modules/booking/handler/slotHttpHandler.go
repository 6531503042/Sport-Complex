package handler

import (
	"main/config"
	"main/modules/booking/usecase"
)

type(
	NewSlotHttpHandlerService interface {

	}

	slotsHttpHandler struct {
		cfg *config.Config
		slotsUsecase usecase.SlotUsecaseService
		
	}
)

func NewSlotHttpHandler(cfg *config.Config, slotUsecase usecase.SlotUsecaseService) NewSlotHttpHandlerService {
	return &slotsHttpHandler{cfg: cfg, slotsUsecase: slotUsecase}
}