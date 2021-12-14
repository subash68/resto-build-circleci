package v1

import (
	"context"

	"github.com/subash68/ate/ate_location_service/pkg/api/location"
)

// add new location -- this could be billing or permanent
func (s *locationServiceServer) UpdateCurrentLocation(ctx context.Context, req *location.CurrentLocationRequest) (*location.CurrentLocationResponse, error) {

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
