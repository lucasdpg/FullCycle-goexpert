package usecase

import (
	"github.com/lucasdpg/FullCycle-goexpert/Clean-Architecture/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

//func (c *CreateOrderUseCase) ListExecute() (OrderOutputDTO, error) {

//dto := OrderOutputDTO{
//	ID:         order.ID,
//	Price:      order.Price,
//	Tax:        order.Tax,
//	FinalPrice: order.Price + order.Tax,
//}
//
//return dto, nil
//}
