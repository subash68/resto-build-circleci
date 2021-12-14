package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/subash68/ate/ate_reservation_service/pkg/api/reservation"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type reservationServiceServer struct {
	db *sql.DB
}

func NewReservationServiceServer(db *sql.DB) reservation.ReservationServiceServer {
	return &reservationServiceServer{db: db}
}

type Reservation struct {
	Id          int64  `json:",omitempty" bson:"_id,omitempty"`
	Description string `json:"description" bson:"description"`
	Status      bool   `json:"status" bson:"status"`
	Order       int32  `json:"order" bson:"order"`
	User        int64  `json:"user,omitempty" bson:"user"`
}

func (s *reservationServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *reservationServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}

func (s *reservationServiceServer) Create(ctx context.Context, req *reservation.CreateRequest) (*reservation.CreateResponse, error) {

	if ctx == nil {
		return &reservation.CreateResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &reservation.CreateResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id` FROM reservations WHERE `tableId` = ? and (`fromReserved` = ? and `toReserved` = ?)", req.Reservation.TableId, req.Reservation.From, req.Reservation.To)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from reservation -> "+err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.AlreadyExists, "failed to retrieve data from reservation -> "+err.Error())
		}
		return nil, status.Error(codes.AlreadyExists, "failed to retrieve data from reservation -> "+err.Error())
	}

	res, err := c.ExecContext(ctx, "INSERT INTO reservations(`tableId`, `fromReserved`, `toReserved`, `reservedById`, `description`) VALUES (?, ?, ?, ?, ?)",
		req.Reservation.TableId, req.Reservation.From, req.Reservation.To, currentUser, req.Reservation.Description)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert reservation information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created reservation "+err.Error())
	}

	return &reservation.CreateResponse{
		Api: apiVersion,
		Id:  id,
		Error: &reservation.ResponseStatus{
			Status:  false,
			Message: "reservation created successfully",
		},
	}, nil
}

func (s *reservationServiceServer) ReadAll(ctx context.Context, req *reservation.ReadAllRequest) (*reservation.ReadAllResponse, error) {

	if ctx == nil {
		return &reservation.ReadAllResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &reservation.ReadAllResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id` FROM tables WHERE `userId` = ?;", currentUser)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from reservations -> "+err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from reservations -> "+err.Error())
	}
	// defer rows.Close()

	tablesIds := []int64{}
	data := []*reservation.Reservation{}
	for rows.Next() {
		id := int64(0)

		//TODO: few other items should also be added here
		if err := rows.Scan(&id); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
		}
		tablesIds = append(tablesIds, id)
	}

	_ = rows.Close()

	for _, v := range tablesIds {
		rows2, err := c.QueryContext(ctx, "SELECT `id`, `fromReserved`, `toReserved`, `reservedById`, `description` FROM reservations WHERE `tableId` = ?", v) // in (select `addonsId` from product_addons where `productId` = ?)", td.Id)

		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to select from reservations -> "+err.Error())
		}
		log.Println("Error Passed")
		reservations := []*reservation.Reservation{}
		for rows2.Next() {
			a := new(reservation.Reservation)
			if err := rows2.Scan(&a.Id, &a.From, &a.To, &a.ReservedById, &a.Description); err != nil {
				return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
			}
			log.Println("Passed", a.Id)
			a.TableId = v
			reservations = append(reservations, a)
		}
		data = append(data, reservations...)
	}

	return &reservation.ReadAllResponse{
		Api:          apiVersion,
		Reservations: data,
		Error: &reservation.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *reservationServiceServer) ReadAllTable(ctx context.Context, request *reservation.ReadAllTableRequest) (*reservation.ReadAllResponse, error) {
	if ctx == nil {
		return &reservation.ReadAllResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(request.Api); err != nil {
		return &reservation.ReadAllResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
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
	// currentUser := ctx.Value("userId").(int64)

	rows, err := c.QueryContext(ctx, "SELECT `id`, `fromReserved`, `toReserved`, `reservedById`, `description` FROM reservations WHERE `tableId` = ?", request.TableId)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from reservations -> "+err.Error())
	}
	defer rows.Close()

	data := []*reservation.Reservation{}
	for rows.Next() {
		td := new(reservation.Reservation)
		if err := rows.Scan(&td.Id, &td.From, &td.To, &td.ReservedById, &td.Description); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from reservation row-> "+err.Error())
		}
		td.TableId = request.TableId

		data = append(data, td)
	}

	return &reservation.ReadAllResponse{
		Api:          apiVersion,
		Reservations: data,
		Error: &reservation.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *reservationServiceServer) Read(ctx context.Context, req *reservation.ReadRequest) (*reservation.ReadResponse, error) {

	if ctx == nil {
		return &reservation.ReadResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &reservation.ReadResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id`, `tableId`, `fromReserved`, `toReserved`, `reservedById`, `description` FROM reservations WHERE `id` = ?", req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from reservations -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from reservations -> "+err.Error())
		}

		return &reservation.ReadResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "reservation not found",
			},
		}, nil
	}

	var data reservation.Reservation
	if err := rows.Scan(&data.Id, &data.TableId, &data.From, &data.To, &data.ReservedById, &data.Description); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from reservation row-> "+err.Error())
	}

	return &reservation.ReadResponse{
		Api:         apiVersion,
		Reservation: &data,
		Error: &reservation.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *reservationServiceServer) Update(ctx context.Context, req *reservation.UpdateRequest) (*reservation.UpdateResponse, error) {

	if ctx == nil {
		return &reservation.UpdateResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &reservation.UpdateResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id` FROM reservations WHERE `tableId` = ? and (`fromReserved` = ? and `toReserved` = ?)", req.Reservation.TableId, req.Reservation.From, req.Reservation.To)

	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "failed to select from reservations -> "+err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.AlreadyExists, "failed to retrieve data from reservations -> "+err.Error())
		}
		return nil, status.Error(codes.AlreadyExists, "failed to retrieve data from reservations -> "+err.Error())
	}

	res, err := c.ExecContext(ctx, "UPDATE reservations SET `tableId` = ?, `fromReserved` = ?, `toReserved` = ?, `reservedById` = ?, `description` = ?,`modifiedAt` = ? WHERE `id` = ?",
		req.Reservation.TableId, req.Reservation.From, req.Reservation.To, req.Reservation.ReservedById, req.Reservation.Description, modifiedAt, req.Reservation.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to update reservation information"+err.Error())
	}

	rows2, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for update rows count "+err.Error())
	}

	if rows2 == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("reservation with ID='%d' is not found",
			req.Reservation.Id))
	}

	return &reservation.UpdateResponse{
		Api:     apiVersion,
		Updated: rows2,
		Error: &reservation.ResponseStatus{
			Status:  false,
			Message: "reservation updated successfully",
		},
	}, nil
}

func (s *reservationServiceServer) Delete(ctx context.Context, req *reservation.DeleteRequest) (*reservation.DeleteResponse, error) {

	if ctx == nil {
		return &reservation.DeleteResponse{
			Api: apiVersion,
			Error: &reservation.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &reservation.DeleteResponse{
			Error: &reservation.ResponseStatus{
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
	res, err := c.ExecContext(ctx, "DELETE FROM reservations WHERE `id`=?", req.Id)
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

	return &reservation.DeleteResponse{
		Api:     apiVersion,
		Deleted: int32(rows),
		Error: &reservation.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
