package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/namnv2496/go-coffee-shop-demo/internal/cache"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderId int32, customerId int32) (domain.OrderDtoRes, error)
	CreateOrder(ctx context.Context, orders []domain.OrderItem, customerId int32) error
	SubmitOrder(ctx context.Context, customerId int32) (int32, error)
	UpdateStatusOrder(ctx context.Context, orderId int32, status int32) error
	ClearAllOrderEOD(ctx context.Context) error
}

type orderService struct {
	Cache        cache.Client
	OrderRepo    repo.OrderRepo
	CustomerRepo repo.CustomerRepo
}

func NewOrderService(
	orderRepo repo.OrderRepo,
	customerRepo repo.CustomerRepo,
	cache cache.Client,
) OrderService {
	return &orderService{
		OrderRepo:    orderRepo,
		CustomerRepo: customerRepo,
		Cache:        cache,
	}
}

func (s orderService) GetOrder(
	ctx context.Context,
	orderId int32,
	customerId int32,
) (domain.OrderDtoRes, error) {

	var res domain.OrderDtoRes
	var orderElement []domain.OrderDto
	var element domain.OrderDto
	var customer domain.Customer
	var orderItems []domain.OrderItem
	var err error
	var orders []domain.Order

	if orderId != 0 {
		// get by orderId
		order, err := s.OrderRepo.GetOrderById(ctx, orderId)
		if err != nil {
			fmt.Println("Error when get order")
		}

		orderItems, err = s.OrderRepo.GetOrderItem(ctx, []int32{orderId})
		if err != nil {
			fmt.Println("Error when get order")
		}

		customer, err = s.CustomerRepo.GetCustomer(ctx, order.Customer_id)
		if err != nil {
			fmt.Println("Cannot get customer information")
		}

		element.Customer = customer
		element.Order = order
		element.OrderItems = orderItems
		orderElement = append(orderElement, element)
	} else if customerId != 0 {
		// get by customerId
		customer, err = s.CustomerRepo.GetCustomer(ctx, customerId)
		if err != nil {
			fmt.Println("Cannot get customer information")
		}

		orders, err = s.OrderRepo.GetOrderByCustomerId(ctx, customerId)
		if err != nil {
			fmt.Println("Error when get order")
		}

		tmp, _ := s.getOrderByIds(ctx, orders, customer)
		orderElement = append(orderElement, tmp...)
	} else {
		// get all orders
		orders, err = s.OrderRepo.GetOrders(ctx)
		if err != nil {
			fmt.Println("Error when get order")
		}
		customers, err := s.CustomerRepo.GetCustomers(ctx)
		if err != nil {
			fmt.Println("Error when get customer")
		}
		customerMap := map[int32]domain.Customer{}
		for _, customer := range customers {
			customerMap[customer.Id] = customer
		}
		for _, order := range orders {
			tmp, _ := s.getOrderByIds(ctx, []domain.Order{order}, customerMap[order.Id])
			orderElement = append(orderElement, tmp...)
		}
	}
	res.Orders = orderElement
	return res, nil
}

func (s orderService) getOrderByIds(
	ctx context.Context,
	orders []domain.Order,
	customer domain.Customer,
) ([]domain.OrderDto, error) {
	var element domain.OrderDto
	var orderElement []domain.OrderDto

	orderIds := make([]int32, 0)
	for _, order := range orders {
		orderIds = append(orderIds, order.Id)
	}
	orderItems, err := s.OrderRepo.GetOrderItem(ctx, orderIds)
	if err != nil {
		fmt.Println("Error when get order")
	}

	orderMap := map[int32][]domain.OrderItem{}
	for _, orderItem := range orderItems {
		if value, exists := orderMap[orderItem.OrderId]; exists {
			value = append(value, orderItem)
			orderMap[orderItem.OrderId] = value
		} else {
			orderMap[orderItem.OrderId] = []domain.OrderItem{orderItem}
		}
	}
	for _, order := range orders {
		element.Customer = customer
		element.Order = order
		element.OrderItems = orderMap[order.Id]
		orderElement = append(orderElement, element)
	}
	return orderElement, nil
}

func (s orderService) CreateOrder(
	ctx context.Context,
	orders []domain.OrderItem,
	customerId int32,
) error {

	data, exist := s.Cache.Get(ctx, mq.REDIS_KEY_ORDER)
	newData, ok := s.convertToDTO(orders, customerId)
	if ok != nil {
		return status.Error(codes.Internal, "failed to convert orders to newData template")
	}
	if exist != nil {
		fmt.Println("[TEST] add new order to redis")
		// check valid data
		// TO DO
		// add new to redis
		json, ok := json.Marshal(newData)
		if ok != nil {
			return status.Error(codes.Internal, "failed to marshall data into cache")
		}
		s.Cache.Set(ctx, mq.REDIS_KEY_ORDER, json)
		return nil
	}
	fmt.Println("[TEST] update old order in redis")
	var orderItemMap = make(map[int32](map[int32]mq.RedisOrder))

	var jsonData []mq.RedisOrderDTO
	publicKey, _ := data.(string)
	if err := json.Unmarshal([]byte(publicKey), &jsonData); err != nil {
		return status.Error(codes.Internal, "failed to unmarshal data")
	}

	for _, customer := range jsonData {
		var value = make(map[int32]mq.RedisOrder)
		for _, orderItem := range customer.RedisOrders {
			value[orderItem.ItemId] = orderItem
		}
		orderItemMap[customer.CustomerId] = value
	}

	for _, customer := range newData {
		v, ok := orderItemMap[customer.CustomerId]
		if ok {
			// exist on map
			for _, orderItem := range customer.RedisOrders {
				item, ok := v[orderItem.ItemId]
				if ok {
					// update quantity and price
					item.Quantity = item.Quantity + orderItem.Quantity
					item.Price = item.Price + orderItem.Quantity*orderItem.Price
					if item.Status == mq.Done {
						item.Status = mq.Processing
					} else if item.Status == mq.Canceled {
						item.Status = mq.Processing
						item.Quantity = orderItem.Quantity
						item.Price = orderItem.Quantity * orderItem.Price
					}
					v[orderItem.ItemId] = item
				} else {
					v[orderItem.ItemId] = orderItem
				}
			}
		} else {
			// new customer
			var newCustomer = make(map[int32]mq.RedisOrder)
			for _, orderItem := range customer.RedisOrders {
				newCustomer[orderItem.ItemId] = orderItem
			}
			orderItemMap[customer.CustomerId] = newCustomer
		}
	}

	var redisData []mq.RedisOrderDTO
	for key, value := range orderItemMap {
		var newRedisOrders []mq.RedisOrder
		for itemId, orderItem := range value {
			newRedisOrders = append(newRedisOrders, mq.RedisOrder{
				ItemId:   itemId,
				Quantity: orderItem.Quantity,
				Price:    orderItem.Price,
				Status:   orderItem.Status,
				Finished: orderItem.Finished,
			})
		}
		redisData = append(redisData, mq.RedisOrderDTO{
			CustomerId:  key,
			RedisOrders: newRedisOrders,
		})
	}

	json, ok := json.Marshal(redisData)
	if ok != nil {
		return status.Error(codes.Internal, "failed to marshall data into cache")
	}
	s.Cache.Set(ctx, mq.REDIS_KEY_ORDER, json)
	return nil
}

func (s orderService) convertToDTO(
	orders []domain.OrderItem,
	customerId int32,
) ([]mq.RedisOrderDTO, error) {

	var res []mq.RedisOrderDTO
	var redisOrder []mq.RedisOrder
	for _, order := range orders {
		var element mq.RedisOrder
		element.ItemId = order.ItemId
		element.Quantity = order.Quantity
		element.Price = order.Price
		redisOrder = append(redisOrder, element)
	}
	res = append(res, mq.RedisOrderDTO{
		CustomerId:  customerId,
		RedisOrders: redisOrder,
	})
	return res, nil
}

func (s orderService) SubmitOrder(
	ctx context.Context,
	customerId int32,
) (int32, error) {
	data, exist := s.Cache.Get(ctx, mq.REDIS_KEY_ORDER)
	if exist != nil {
		return 0, status.Error(codes.Internal, "failed to update status of not exist customerId")
	}
	var jsonData []mq.RedisOrderDTO
	jsonString, _ := data.(string)
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		return 0, status.Error(codes.Internal, "failed to unmarshal data")
	}

	for i, redisOrderDTO := range jsonData {
		if redisOrderDTO.CustomerId == customerId {
			var orders []domain.OrderItem
			for _, redisOrder := range redisOrderDTO.RedisOrders {
				if redisOrder.Status == mq.Done {
					orders = append(orders, domain.OrderItem{
						ItemId: redisOrder.ItemId,
						// only save finished item volume
						Quantity: redisOrder.Finished,
						Price:    redisOrder.Price,
					})
				}
			}
			// remove order in redis
			jsonData = append(jsonData[:i], jsonData[i+1:]...)
			json, ok := json.Marshal(jsonData)
			if ok != nil {
				return 0, status.Error(codes.Internal, "failed to marshall data into cache")
			}
			s.Cache.Set(ctx, mq.REDIS_KEY_ORDER, json)
			return s.OrderRepo.CreateOrder(ctx, orders, customerId)
		}
	}
	return 0, status.Error(codes.Internal, "order of customerId is not exist")
}

func (s orderService) UpdateStatusOrder(
	ctx context.Context,
	orderId int32,
	status int32,
) error {

	return s.OrderRepo.UpdateStatusOrder(ctx, orderId, status)
}

func (s orderService) ClearAllOrderEOD(
	ctx context.Context,
) error {

	var jsonData = make([]mq.RedisOrderDTO, 0)
	json, ok := json.Marshal(jsonData)
	if ok != nil {
		return status.Error(codes.Internal, "failed to marshall data into cache")
	}
	s.Cache.Set(ctx, mq.REDIS_KEY_ORDER, json)
	return nil
}
