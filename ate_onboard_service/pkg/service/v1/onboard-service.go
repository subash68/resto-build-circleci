package v1

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/subash68/ate/ate_onboard_service/pkg/api/onboard"
	"github.com/subash68/ate/ate_onboard_service/pkg/helper"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type onboardServiceServer struct {
	db *sql.DB
}

func NewOnboardServiceServer(db *sql.DB) onboard.OnboardServiceServer {
	return &onboardServiceServer{db: db}
}

type User struct {
	Id       int64   `json:",omitempty" bson:"_id,omitempty"`
	Fullname string  `json:"fullname" bson:"fullname"`
	Email    string  `json:"email" bson:"email"`
	Phone    *string `json:"phone" bson:"phone"`
	Password string  `json:"password" bson:"password"`
	UserType uint8   `json:"usertype" bson:"usertype"`
}

func (s *onboardServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Error(codes.Unimplemented, "unsupported API version.")
		}
	}

	return nil
}

func (s *onboardServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}

func (s *onboardServiceServer) Login(ctx context.Context, req *onboard.LoginUserRequest) (*onboard.LoginUserResponse, error) {

	if err := s.checkAPI(req.Api); err != nil {
		return &onboard.LoginUserResponse{
			Error: &onboard.ResponseStatus{
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

	fmt.Println(req.Email)
	rows, err := c.QueryContext(ctx, "SELECT `id`, `fullname`, `email`, `password`, `phone`, `type` FROM users WHERE `email` = ?", req.Email)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from users -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from users -> "+err.Error())
		}
		return &onboard.LoginUserResponse{
			Error: &onboard.ResponseStatus{
				Status:  false,
				Message: "user not registered",
			},
		}, nil
	}

	var user User
	if err := rows.Scan(&user.Id, &user.Fullname, &user.Email, &user.Password, &user.Phone, &user.UserType); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from user row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple users with email : %s", req.Email))
	}

	log.Print(user)

	//compare password here

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &onboard.LoginUserResponse{
			Api: apiVersion,
			// User: &user,
			Error: &onboard.ResponseStatus{
				Status:  false,
				Message: "authentication failed",
			},
		}, nil
	}

	generatedToken := helper.GenerateToken(user.Id, user.Fullname, user.Email, user.UserType)

	return &onboard.LoginUserResponse{
		Api:   apiVersion,
		Token: generatedToken,
		User: &onboard.UserResponse{
			Id:       user.Id,
			Fullname: user.Fullname,
			Email:    user.Email,
			Type:     int32(user.UserType),
		},
		Error: &onboard.ResponseStatus{
			Status:  false,
			Message: "user not registered",
		},
	}, nil
}

func (s *onboardServiceServer) Register(ctx context.Context, req *onboard.RegisterUserRequest) (*onboard.RegisterUserResponse, error) {

	if err := s.checkAPI(req.Api); err != nil {
		return &onboard.RegisterUserResponse{
			Api: req.Api,
			Error: &onboard.ResponseStatus{
				Status:  true,
				Message: "Api version unsupported",
			},
		}, nil
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "updated at field has invalid format "+err.Error())
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.User.Password), 8)

	//TODO: Execute procedure here
	res, err := c.ExecContext(ctx, "INSERT INTO users(`fullname`, `email`, `phone`, `password`, `type`) VALUES (?, ?, ?, ?, ?)",
		req.User.Fullname, req.User.Email, req.User.Phone, string(hashedPassword), req.User.Type)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert user information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created user "+err.Error())
	}

	return &onboard.RegisterUserResponse{
		Api: req.Api,
		Error: &onboard.ResponseStatus{
			Status:  false,
			Message: "User registered successfully",
		},
		UserId: id,
	}, nil
}

func (s *onboardServiceServer) NotificationRegister(ctx context.Context, req *onboard.NotificationUserRequest) (*onboard.NotificationUserResponse, error) {
	if ctx == nil {
		return &onboard.NotificationUserResponse{
			Api: apiVersion,
			Error: &onboard.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	if err := s.checkAPI(req.Api); err != nil {
		return &onboard.NotificationUserResponse{
			Error: &onboard.ResponseStatus{
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

	currentUser := ctx.Value("userId").(int64)

	_, err = c.ExecContext(ctx, "UPDATE users SET notificationToken = ? where id = ?", req.NotificationToken, currentUser)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to Update Cart information"+err.Error())
	}

	return &onboard.NotificationUserResponse{
		Api: req.Api,
		Error: &onboard.ResponseStatus{
			Status:  false,
			Message: "User registered successfully",
		},
		UserId: currentUser,
	}, nil

}
