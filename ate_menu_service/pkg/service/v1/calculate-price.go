package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
)

type ProductInfo struct {
	name            string
	price           float32
	discount        float32
	discountType    int32
	discountedPrice float32
}

func (s *menuServiceServer) CalculatePrice(ctx context.Context, req *menu.CalculatePriceRequest) (*menu.CalculatePriceResponse, error) {

	if ctx == nil {
		return &menu.CalculatePriceResponse{
			Status: &menu.ResponseStatus{
				Status:  false,
				Message: "",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &menu.CalculatePriceResponse{
			Status: &menu.ResponseStatus{
				Status:  false,
				Message: "Unsupported api version",
			},
		}, nil
	}

	productPrice := s.CalculateProductPrice(ctx, int32(req.GetProduct()))

	return &menu.CalculatePriceResponse{
		Name:            productPrice.name,
		Price:           productPrice.price,
		Discount:        productPrice.discount,
		DiscountType:    productPrice.discountType,
		DiscountedPrice: productPrice.discountedPrice,
		Status:          &menu.ResponseStatus{Status: true, Message: ""},
	}, nil
}

func (s *menuServiceServer) CalculateProductPrice(ctx context.Context, product int32) *ProductInfo {
	//
	c, err := s.connect(ctx)

	if err != nil {
		return &ProductInfo{} // return nothing here
	}

	defer c.Close()

	// get product price and discount price
	result := c.QueryRowContext(ctx, "call sp_GET_PRODUCT_PRICE(?)", product)
	// TODO: Fix error here
	productInfo := &ProductInfo{}

	err = result.Scan(&productInfo.name, &productInfo.price, &productInfo.discount, &productInfo.discountType)
	if err != nil {
		log.Print("Error while reading product information")
		log.Print(err.Error())
		return &ProductInfo{}
	}

	productInfo.discountedPrice = ApplyPriceDiscount(productInfo.price, productInfo.discount, productInfo.discountType)
	return productInfo

}

func ApplyPriceDiscount(price float32, discount float32, discountType int32) float32 {

	log.Print("printing discount price")
	log.Print(price)
	log.Print(discount)
	log.Print(discountType)

	discountedPrice := float32(0.0)
	if discountType == 1 {
		discountedPrice = float32(price - price*(discount/100))
	} else if discountType == 2 {
		discountedPrice = price - discount
	}

	return discountedPrice
}
