package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_cart_service/pkg/api/cart"
)

func (s *cartServiceServer) AddItem(ctx context.Context, req *cart.AddItemRequest) (*cart.AddItemResponse, error) {

	// Check if context is present
	log.Println("serivce: add item to cart")
	if ctx == nil {
		return &cart.AddItemResponse{
			Status: &cart.ResponseStatus{
				IsSuccess:    false,
				ErrorCode:    "E101",
				ErrorMessage: "User authentication failed",
			},
		}, nil
	}

	// currentUser := ctx.Value("userId").(int64)
	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	cartObj, _ := s.CheckIfCartExists(ctx)

	if cartObj.GetId() == 0 {
		cartObj, _ = s.CreateNewCart(ctx)
	}

	// add item to cart here
	result, _ := c.QueryContext(ctx,
		"call sp_ADD_NEW_ITEM(?, ?, ?)",
		cartObj.GetId(),
		req.Item.GetItemId(),
		req.Item.GetItemCount())

	defer result.Close()

	// update item with price from menu service

	// result should return item inserted
	Items := []*cart.CartItem{}
	for result.Next() {

		var id int32
		var itemId int32
		var itemCount int32

		scanerr := result.Scan(&id, &itemId, &itemCount)
		if scanerr != nil {
			log.Print(scanerr)
		}

		Items = append(Items, &cart.CartItem{
			Id: id, ItemId: itemId, ItemCount: itemCount, ItemPrice: s.GetPrice(ctx, itemId).DiscountedPrice * float32(itemCount),
		})
		log.Println(Items)
	}
	cartObj.Items = Items
	return &cart.AddItemResponse{
		Cart:   cartObj,
		Status: &cart.ResponseStatus{IsSuccess: true, SuccessMessage: "Item added successfully to cart."},
	}, nil
}

func (s *cartServiceServer) CheckIfCartExists(ctx context.Context) (*cart.Cart, error) {
	currentUser := ctx.Value("userId").(int64)
	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	row := c.QueryRowContext(ctx, "call sp_CHECK_CART_EXISTS(?)", currentUser)

	cart := &cart.Cart{
		Id:            0,
		Address:       "",
		Items:         []*cart.CartItem{},
		ShippingCost:  0,
		TotalItemCost: 0,
		Total:         0,
	}
	row.Scan(&cart)
	return cart, nil
}

func (s *cartServiceServer) CreateNewCart(ctx context.Context) (*cart.Cart, error) {
	currentUser := ctx.Value("userId").(int64)
	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	// Get current address here
	currentAddress := "	2217 Denver Avenue, Riverside, CA, California 92509"

	// order type will be default
	result := c.QueryRowContext(ctx, "call sp_CREATE_CART(?, ?)", currentUser, currentAddress)
	// This had to be fixed
	cart := &cart.Cart{
		Id:            0,
		Address:       currentAddress,
		ShippingCost:  0,
		TotalItemCost: 0,
		Total:         0,
	}
	result.Scan(&cart.Id)
	return cart, nil
}
