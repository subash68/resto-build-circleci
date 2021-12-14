package v1

import (
	"context"

	"github.com/subash68/ate/ate_location_service/pkg/api/location"
)

func (s *locationServiceServer) GetCurrentLocation(ctx context.Context, req *location.VoidNoParams) (*location.CurrentLocationResponse, error) {

	return &location.CurrentLocationResponse{
		Status: &location.ResponseStatus{
			IsSuccess:      false,
			SuccessMessage: "",
			ErrorCode:      "",
			ErrorMessage:   "",
			ErrorDetail:    "",
		},
	}, nil
}
