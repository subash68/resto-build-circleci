package v1

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/subash68/ate/ate_addon_service/pkg/api/addon"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type addonServiceServer struct {
	db *sql.DB
}

func NewAddonServiceServer(db *sql.DB) addon.AddonServiceServer {
	return &addonServiceServer{db: db}
}

// type Addon struct {
// 	Id          int64 `json:",omitempty" bson:"_id,omitempty"`
// 	Name		string `json:"description" bson:"name"`
// 	Description string `json:"description" bson:"description"`
// 	Status      bool `json:"status" bson:"status"`
// 	Order       int32 `json:"order" bson:"order"`
// 	User        int64 `json:"user,omitempty" bson:"user"`
// }

func (s *addonServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *addonServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error(), ctx, c)
	}

	return c, nil
}

func (s *addonServiceServer) Create(ctx context.Context, req *addon.CreateRequest) (*addon.CreateResponse, error) {

	if ctx == nil {
		return &addon.CreateResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "user authentication failed: ",
			},
		}, nil
	}
	log.Println("Passed The Authentication: ", req.Api, req.Addon)
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &addon.CreateResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
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

	// insert addon into database
	currentUser := ctx.Value("userId").(int64)

	log.Print(req.Addon)

	res, err := c.ExecContext(ctx, "INSERT INTO addons (`name`, `price`, `userId`) VALUES (?, ?, ?)", req.Addon.Name, req.Addon.Price, currentUser)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert addon information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created addon "+err.Error())
	}

	return &addon.CreateResponse{
		Api: apiVersion,
		Id:  id,
		Error: &addon.ResponseStatus{
			Status:  false,
			Message: "addon created successfully",
		},
	}, nil
}

func (s *addonServiceServer) ReadAll(ctx context.Context, req *addon.ReadAllRequest) (*addon.ReadAllResponse, error) {

	if ctx == nil {
		return &addon.ReadAllResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "user authentication failed" + fmt.Sprintf("%v", ctx),
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &addon.ReadAllResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
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

	// insert addon into database
	currentUser := ctx.Value("userId").(int64)

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `price`, `userId` FROM addons WHERE `userId` = ?", currentUser)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from addons -> "+err.Error())
	}
	defer rows.Close()

	data := []*addon.Addon{}
	for rows.Next() {
		td := new(addon.Addon)

		if err := rows.Scan(&td.Id, &td.Name, &td.Price, &td.User); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from addon row-> "+err.Error())
		}

		data = append(data, td)
	}

	return &addon.ReadAllResponse{
		Api:    apiVersion,
		Addons: data,
		Error: &addon.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *addonServiceServer) Read(ctx context.Context, req *addon.ReadRequest) (*addon.ReadResponse, error) {

	if ctx == nil {
		return &addon.ReadResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &addon.ReadResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `price`, `userId` FROM addons WHERE `id` = ?", req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from addons -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from categories -> "+err.Error())
		}

		return &addon.ReadResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "addon not found",
			},
		}, nil
	}

	var data addon.Addon

	if err := rows.Scan(&data.Id, &data.Name, &data.Price, &data.User); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from addon row-> "+err.Error())
	}

	return &addon.ReadResponse{
		Api:   apiVersion,
		Addon: &data,
		Error: &addon.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *addonServiceServer) Update(ctx context.Context, req *addon.UpdateRequest) (*addon.UpdateResponse, error) {

	if ctx == nil {
		return &addon.UpdateResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &addon.UpdateResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
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

	// insert addon into database

	modifiedAt := time.Now().UTC()

	res, err := c.ExecContext(ctx, "UPDATE addons SET `name` = ?, `price` = ?, `modifiedAt` = ? WHERE `id` = ?", req.Addon.Name, req.Addon.Price, modifiedAt, req.Addon.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to update addon information"+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for update rows count "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("addon with ID='%d' is not found",
			req.Addon.Id))
	}

	return &addon.UpdateResponse{
		Api:     apiVersion,
		Updated: rows,
		Error: &addon.ResponseStatus{
			Status:  false,
			Message: "addon updated successfully",
		},
	}, nil
}

func (s *addonServiceServer) Delete(ctx context.Context, req *addon.DeleteRequest) (*addon.DeleteResponse, error) {

	if ctx == nil {
		return &addon.DeleteResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &addon.DeleteResponse{
			Error: &addon.ResponseStatus{
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
	res, err := c.ExecContext(ctx, "DELETE FROM addons WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete addon -> "+err.Error())
	}

	_, err = c.ExecContext(ctx, "DELETE FROM product_addons WHERE `addonsId`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete addon -> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("addon with ID='%d' is not found",
			req.Id))
	}

	return &addon.DeleteResponse{
		Api:     apiVersion,
		Deleted: int32(rows),
		Error: &addon.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *addonServiceServer) ReadByProduct(ctx context.Context, req *addon.ReadByProductRequest) (*addon.ReadByProductResponse, error) {

	if ctx == nil {
		return &addon.ReadByProductResponse{
			Api: apiVersion,
			Error: &addon.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &addon.ReadByProductResponse{
			Error: &addon.ResponseStatus{
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
	res, err := c.ExecContext(ctx, "DELETE FROM addons WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete addon -> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("addon with ID='%d' is not found",
			req.Id))
	}

	return &addon.ReadByProductResponse{
		Api:    apiVersion,
		Addons: nil, // addons should be sent from here
		Error: &addon.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
