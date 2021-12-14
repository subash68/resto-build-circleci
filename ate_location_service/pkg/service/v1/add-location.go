package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_location_service/pkg/api/location"
)

// add new location -- this could be billing or permanent
func (s *locationServiceServer) AddLocation(ctx context.Context, req *location.AddLocationRequest) (*location.AddLocationResponse, error) {

	log.Println("Service: Location - add location")
	if ctx == nil {
		return &location.AddLocationResponse{
			Status: &location.ResponseStatus{
				IsSuccess:      false,
				SuccessMessage: "",
				ErrorCode:      "E101",
				ErrorMessage:   "User authentication failed",
				ErrorDetail:    "",
			},
		}, nil
	}

	c, err := s.connect(ctx)
	if err != nil {
		log.Printf("error connecting database : %v", err.Error())
		return nil, err
	}

	defer c.Close()

	//

	return &location.AddLocationResponse{
		Status: &location.ResponseStatus{
			IsSuccess:      false,
			SuccessMessage: "",
			ErrorCode:      "",
			ErrorMessage:   "",
			ErrorDetail:    "",
		},
	}, nil
}
