package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/subash68/ate/ate_table_service/pkg/api/table"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"time"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type tableServiceServer struct {
	db *sql.DB
}

func NewTableServiceServer(db *sql.DB) table.TableServiceServer {
	return &tableServiceServer{db: db}
}

func (s *tableServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *tableServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error(), ctx, c)
	}

	return c, nil
}

func (s *tableServiceServer) Create(ctx context.Context, req *table.CreateRequest) (*table.CreateResponse, error) {

	if ctx == nil {
		return &table.CreateResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  true,
				Message: "user authentication failed: ",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &table.CreateResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
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

	// insert setting into database
	userID := ctx.Value("userId").(int64)

	l := uuid.New()
	qrCode := "table-" + strconv.FormatInt(userID, 10) + "-seats-" + strconv.FormatInt(req.Table.Seats, 10) + "-" + l.String()

	log.Print(req.Table)

	res, err := c.ExecContext(ctx, "INSERT INTO tables (`name`, `qrCode`, `seats`, `userId`, `isOpenToReservation`) VALUES (?, ?, ?, ?, ?)", req.Table.Name, qrCode, req.Table.Seats, userID, req.Table.IsOpenToReservation)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert setting information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created setting "+err.Error())
	}

	return &table.CreateResponse{
		Api: apiVersion,
		Id:  id,
		Error: &table.ResponseStatus{
			Status:  false,
			Message: "setting created successfully",
		},
	}, nil
}

func (s *tableServiceServer) ReadAll(ctx context.Context, req *table.ReadAllRequest) (*table.ReadAllResponse, error) {

	if ctx == nil {
		return &table.ReadAllResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  true,
				Message: "user authentication failed" + fmt.Sprintf("%v", ctx),
			},
		}, nil
	}

	log.Println("Step 1")

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &table.ReadAllResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}
	log.Println("Step 2")

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	// insert setting into database
	currentUser := ctx.Value("userId").(int64)
	log.Println("Step 3, ", currentUser)
	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `qrCode`, `seats`, `isOpenToReservation` FROM tables WHERE `userId` = ?;", currentUser)
	log.Println("Step 4")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from users -> "+err.Error())
	}
	log.Println("Step 5")
	data := []*table.Table{}

	log.Println("Step 7")
	for rows.Next() {
		td := new(table.Table)
		log.Println("Step 8")
		if err := rows.Scan(&td.Id, &td.Name, &td.QrCode, &td.Seats, &td.IsOpenToReservation); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from setting row-> "+err.Error())
		}
		log.Println(td)
		data = append(data, td)
	}
	log.Println("Step 9")

	_ = rows.Close()

	return &table.ReadAllResponse{
		Api:    apiVersion,
		Tables: data,
		Error: &table.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *tableServiceServer) Read(ctx context.Context, req *table.ReadRequest) (*table.ReadResponse, error) {

	if ctx == nil {
		return &table.ReadResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &table.ReadResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `qrCode`, `seats`, `isOpenToReservation` FROM tables WHERE `id` = ?;", req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from users -> "+err.Error())
	}
	log.Println("Step 1")
	if !rows.Next() {
		log.Println("Step 2")
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from categories -> "+err.Error())
		}
		log.Println("Step 3")
		return &table.ReadResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  true,
				Message: "setting not found",
			},
		}, nil
	}
	log.Println("Step 4")
	var td table.Table

	if err := rows.Scan(&td.Id, &td.Name, &td.QrCode, &td.Seats, &td.IsOpenToReservation); err != nil {
		log.Println("Step 4.1, ", err)
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from setting row-> "+err.Error())
	}
	log.Println("Step 5")

	_ = rows.Close()

	return &table.ReadResponse{
		Api:   apiVersion,
		Table: &td,
		Error: &table.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *tableServiceServer) Update(ctx context.Context, req *table.UpdateRequest) (*table.UpdateResponse, error) {

	if ctx == nil {
		return &table.UpdateResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &table.UpdateResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
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

	// insert setting into database

	modifiedAt := time.Now().UTC()

	res, err := c.ExecContext(ctx, "UPDATE tables SET `name` = ?, `seats` = ?, `isOpenToReservation` = ?, `modifiedAt` = ? WHERE `id` = ?",
		req.Table.Name, req.Table.Seats, req.Table.IsOpenToReservation, modifiedAt, req.Table.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to update setting information"+err.Error())
	}

	r, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for update rows count "+err.Error())
	}

	if r == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("setting with ID='%d' is not found",
			req.Table.Id))
	}

	return &table.UpdateResponse{
		Api:     apiVersion,
		Updated: r,
		Error: &table.ResponseStatus{
			Status:  false,
			Message: "setting updated successfully",
		},
	}, nil
}

func (s *tableServiceServer) Delete(ctx context.Context, req *table.DeleteRequest) (*table.DeleteResponse, error) {

	if ctx == nil {
		return &table.DeleteResponse{
			Api: apiVersion,
			Error: &table.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &table.DeleteResponse{
			Error: &table.ResponseStatus{
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
	res, err := c.ExecContext(ctx, "DELETE FROM tables WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete setting -> "+err.Error())
	}

	_, err = c.ExecContext(ctx, "DELETE FROM reservations WHERE `tableId`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete open time -> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("setting with ID='%d' is not found",
			req.Id))
	}

	return &table.DeleteResponse{
		Api:     apiVersion,
		Deleted: int32(rows),
		Error: &table.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
