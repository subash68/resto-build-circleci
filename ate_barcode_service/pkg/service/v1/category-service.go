package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/subash68/ate/ate_category_service/pkg/api/category"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type categoryServiceServer struct {
	db *sql.DB
}

func NewCategoryServiceServer(db *sql.DB) category.CategoryServiceServer {
	return &categoryServiceServer{db: db}
}

type Category struct {
	Id          int64  `json:",omitempty" bson:"_id,omitempty"`
	Description string `json:"description" bson:"description"`
	Status      bool   `json:"status" bson:"status"`
	Order       int32  `json:"order" bson:"order"`
	User        int64  `json:"user,omitempty" bson:"user"`
}

func (s *categoryServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *categoryServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}

func (s *categoryServiceServer) Create(ctx context.Context, req *category.CreateRequest) (*category.CreateResponse, error) {

	if ctx == nil {
		return &category.CreateResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &category.CreateResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	// insert reservation into database
	currentUser := ctx.Value("userId").(int64)

	res, err := c.ExecContext(ctx, "INSERT INTO categories(`name`, `description`, `status`, `categoryOrder`, `userId`) VALUES (?, ?, ?, ?, ?)", req.Category.Name, req.Category.Description, 1, req.Category.Order, currentUser)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert reservation information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created reservation "+err.Error())
	}

	return &category.CreateResponse{
		Api: apiVersion,
		Id:  id,
		Error: &category.ResponseStatus{
			Status:  false,
			Message: "reservation created successfully",
		},
	}, nil
}

func (s *categoryServiceServer) ReadAll(ctx context.Context, req *category.ReadAllRequest) (*category.ReadAllResponse, error) {

	/*if ctx == nil {
		return &reservation.ReadAllResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}*/

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &category.ReadAllResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	// insert reservation into database
	currentUser := ctx.Value("userId").(int64)

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `description`, `status`, `categoryOrder`, `userId` FROM categories WHERE `userId` = ?", currentUser)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
	}
	defer rows.Close()

	data := []*category.Category{}
	for rows.Next() {
		td := new(category.Category)

		if err := rows.Scan(&td.Id, &td.Name, &td.Description, &td.Status, &td.Order, &td.User); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from reservation row-> "+err.Error())
		}

		data = append(data, td)
	}

	return &category.ReadAllResponse{
		Api:        apiVersion,
		Categories: data,
		Error: &category.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *categoryServiceServer) Read(ctx context.Context, req *category.ReadRequest) (*category.ReadResponse, error) {

	if ctx == nil {
		return &category.ReadResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &category.ReadResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `description`, `status`, `categoryOrder`, `userId` FROM categories WHERE `id` = ?", req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from categories -> "+err.Error())
		}

		return &category.ReadResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  true,
				Message: "reservation not found",
			},
		}, nil
	}

	var data category.Category

	if err := rows.Scan(&data.Id, &data.Name, &data.Description, &data.Status, &data.Order, &data.User); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from reservation row-> "+err.Error())
	}

	return &category.ReadResponse{
		Api:      apiVersion,
		Category: &data,
		Error: &category.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *categoryServiceServer) Update(ctx context.Context, req *category.UpdateRequest) (*category.UpdateResponse, error) {

	if ctx == nil {
		return &category.UpdateResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &category.UpdateResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	// insert reservation into database

	modifiedAt := time.Now().UTC()

	res, err := c.ExecContext(ctx, "UPDATE categories SET `name` = ?, `description` = ?, `status` = ?, `categoryOrder` = ?, `modifiedAt` = ? WHERE `id` = ?", req.Category.Name, req.Category.Description, req.Category.Status, req.Category.Order, modifiedAt, req.Category.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to update reservation information"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for update rows count "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("reservation with ID='%d' is not found",
			req.Category.Id))
	}

	return &category.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
		Error: &category.ResponseStatus{
			Status:  false,
			Message: "reservation updated successfully",
		},
	}, nil
}

func (s *categoryServiceServer) Delete(ctx context.Context, req *category.DeleteRequest) (*category.DeleteResponse, error) {

	if ctx == nil {
		return &category.DeleteResponse{
			Api: apiVersion,
			Error: &category.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &category.DeleteResponse{
			Error: &category.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// delete ToDo
	res, err := c.ExecContext(ctx, "DELETE FROM categories WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete reservation -> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("reservation with ID='%d' is not found",
			req.Id))
	}

	return &category.DeleteResponse{
		Api:     apiVersion,
		Deleted: int32(rows),
		Error: &category.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
