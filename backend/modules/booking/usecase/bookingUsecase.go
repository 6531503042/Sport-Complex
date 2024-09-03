package usecase

import "main/modules/booking/repository"

type (
	BookingUsecaseService interface {

	}

	bookingUsecase struct {
		bookingRepository repository.BookingRepositoryService
	}
)

func NewBookingUsecase(bookingRepository repository.BookingRepositoryService) BookingUsecaseService {
	return &bookingUsecase{
		bookingRepository: bookingRepository,
	}
}



