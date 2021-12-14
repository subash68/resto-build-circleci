package v1

import (
	"context"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/subash68/ate/ate_setting_service/pkg/api/setting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type settingServiceServer struct {
	db *sql.DB
}

func NewSettingServiceServer(db *sql.DB) setting.SettingServiceServer {
	return &settingServiceServer{db: db}
}

// type Setting struct {
// 	Id          int64 `json:",omitempty" bson:"_id,omitempty"`
// 	Name		string `json:"description" bson:"name"`
// 	Description string `json:"description" bson:"description"`
// 	Status      bool `json:"status" bson:"status"`
// 	Order       int32 `json:"order" bson:"order"`
// 	User        int64 `json:"user,omitempty" bson:"user"`
// }

func (s *settingServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *settingServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error(), ctx, c)
	}

	return c, nil
}

func (s *settingServiceServer) Create(ctx context.Context, req *setting.CreateRequest) (*setting.CreateResponse, error) {

	if ctx == nil {
		return &setting.CreateResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
				Status:  true,
				Message: "user authentication failed: ",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &setting.CreateResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
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
	_ = ctx.Value("userId").(int64)

	log.Print(req.Setting)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Setting.Password), 8)

	res, err := c.ExecContext(ctx, "INSERT INTO users (`email`, `password`) VALUES (?, ?)", req.Setting.Email, string(hashedPassword))
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert setting information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created setting "+err.Error())
	}

	return &setting.CreateResponse{
		Api: apiVersion,
		Id:  id,
		Error: &setting.ResponseStatus{
			Status:  false,
			Message: "setting created successfully",
		},
	}, nil
}

func (s *settingServiceServer) ReadAll(ctx context.Context, req *setting.ReadAllRequest) (*setting.ReadAllResponse, error) {

	if ctx == nil {
		return &setting.ReadAllResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
				Status:  true,
				Message: "user authentication failed" + fmt.Sprintf("%v", ctx),
			},
		}, nil
	}

	log.Println("Step 1")

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &setting.ReadAllResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
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
	rows, err := c.QueryContext(ctx, "SELECT `id`, `fullname`, `email`, `phone`, `type`, `cuisine`, `status`, `everyday`, `profileImageId`, `shopLogoId`, `shopBannerId`, `isVeg`, `mealService`, `partyCatering`, `deliveryTakeAway`, `delivery`, `freeDelivery`, `offerType`, `offer`, `offerAmount`, `maxDeliveryTime`, `description`, `location`, `locLongitude`, `locLatitude` FROM users WHERE `id` = ?", currentUser)
	log.Println("Step 4")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from users -> "+err.Error())
	}
	log.Println("Step 5")
	data := []*setting.Setting{}

	log.Println("Step 7")
	for rows.Next() {
		td := new(setting.Setting)
		log.Println("Step 8")
		if err := rows.Scan(&td.Id, &td.Fullname, &td.Email, &td.Phone, &td.Type, &td.Cuisine, &td.Status, &td.Everyday, &td.ProfileImageUrl, &td.ShopLogoUrl, &td.ShopBannerUrl, &td.IsVeg, &td.MealService, &td.PartyCatering, &td.DeliveryTakeAway, &td.Delivery, &td.FreeDelivery, &td.OfferType, &td.Offer, &td.OfferAmount, &td.MaxDeliveryTime, &td.Description, &td.Location, &td.LocLongitude, &td.LocLatitude); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from setting row-> "+err.Error())
		}
		log.Println(td)
		data = append(data, td)
	}
	log.Println("Step 9")

	_ = rows.Close()

	for _, v := range data {
		log.Println(c)
		rows2, err := c.QueryContext(ctx, "SELECT `dayName`, `fromOpen`, `toOpen` FROM open_time WHERE `userId`= ?", v.Id)

		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
		}
		openings := []*setting.OpenTime{}
		for rows2.Next() {
			a := new(setting.OpenTime)
			if err := rows2.Scan(&a.DayName, &a.OpenFrom, &a.OpenTo); err != nil {
				return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
			}
			openings = append(openings, a)
		}
		v.Opening = openings
		_ = rows2.Close()
	}

	return &setting.ReadAllResponse{
		Api:      apiVersion,
		Settings: data,
		Error: &setting.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *settingServiceServer) Read(ctx context.Context, req *setting.ReadRequest) (*setting.ReadResponse, error) {

	if ctx == nil {
		return &setting.ReadResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &setting.ReadResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
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

	rows, err := c.QueryContext(ctx, "SELECT `id`, `fullname`, `email`, `phone`, `type`, `cuisine`, `status`, `everyday`, `profileImageId`, `shopLogoId`, `shopBannerId`, `isVeg`, `mealService`, `partyCatering`, `deliveryTakeAway`, `delivery`, `freeDelivery`, `offerType`, `offer`, `offerAmount`, `maxDeliveryTime`, `description`, `location`, `locLongitude`, `locLatitude` FROM users WHERE `id` = ?", req.Id)

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
		return &setting.ReadResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
				Status:  true,
				Message: "setting not found",
			},
		}, nil
	}
	log.Println("Step 4")
	var td setting.Setting

	if err := rows.Scan(&td.Id, &td.Fullname, &td.Email, &td.Phone, &td.Type, &td.Cuisine, &td.Status, &td.Everyday, &td.ProfileImageUrl, &td.ShopLogoUrl, &td.ShopBannerUrl, &td.IsVeg, &td.MealService, &td.PartyCatering, &td.DeliveryTakeAway, &td.Delivery, &td.FreeDelivery, &td.OfferType, &td.Offer, &td.OfferAmount, &td.MaxDeliveryTime, &td.Description, &td.Location, &td.LocLongitude, &td.LocLatitude); err != nil {
		log.Println("Step 4.1, ", err)
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from setting row-> "+err.Error())
	}
	log.Println("Step 5")

	_ = rows.Close()
	log.Println(c)
	rows2, err := c.QueryContext(ctx, "SELECT `dayName`, `fromOpen`, `toOpen` FROM open_time WHERE `userId`= ?", td.Id)
	log.Println("Step 6")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
	}
	openings := []*setting.OpenTime{}
	for rows2.Next() {
		a := new(setting.OpenTime)
		if err := rows2.Scan(&a.DayName, &a.OpenFrom, &a.OpenTo); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
		}
		openings = append(openings, a)
	}
	td.Opening = openings
	_ = rows2.Close()

	return &setting.ReadResponse{
		Api:     apiVersion,
		Setting: &td,
		Error: &setting.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}

func (s *settingServiceServer) Update(ctx context.Context, req *setting.UpdateRequest) (*setting.UpdateResponse, error) {

	if ctx == nil {
		return &setting.UpdateResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &setting.UpdateResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
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

	res, err := c.ExecContext(ctx, "UPDATE users SET `fullname` = ?, `email` = ?, `password` = ?, `phone` = ?, `type` = ?, `cuisine` = ?, `status` = ?, `everyday` = ?, `profileImageId` = ?, `shopLogoId` = ?, `shopBannerId` = ?, `isVeg` = ?, `mealService` = ?, `partyCatering` = ?, `deliveryTakeAway` = ?, `delivery` = ?, `freeDelivery` = ?, `offerType` = ?, `offer` = ?, `offerAmount` = ?, `maxDeliveryTime` = ?, `description` = ?, `location` = ?, `locLongitude` = ?, `locLatitude` = ?, `modifiedAt` = ? WHERE `id` = ?",
		req.Setting.Fullname, req.Setting.Email, req.Setting.Password, req.Setting.Phone, req.Setting.Type, req.Setting.Cuisine, req.Setting.Status, req.Setting.Everyday, req.Setting.ProfileImageUrl, req.Setting.ShopLogoUrl, req.Setting.ShopBannerUrl, req.Setting.IsVeg, req.Setting.MealService, req.Setting.PartyCatering, req.Setting.DeliveryTakeAway, req.Setting.Delivery, req.Setting.FreeDelivery, req.Setting.OfferType, req.Setting.Offer, req.Setting.OfferAmount, req.Setting.MaxDeliveryTime, req.Setting.Description, req.Setting.Location, req.Setting.LocLongitude, req.Setting.LocLatitude, modifiedAt, req.Setting.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to update setting information"+err.Error())
	}

	r, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for update rows count "+err.Error())
	}

	_, err = c.ExecContext(ctx, "DELETE FROM open_time WHERE `userId`=?", req.Setting.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete open time -> "+err.Error())
	}

	for i := 0; i < len(req.Setting.Opening); i++ {
		addonsResponse, err := c.ExecContext(ctx, "INSERT INTO open_time (`dayName`, `fromOpen`, `toOpen`) VALUES (?, ?, ?)", req.Setting.Opening[i].DayName, req.Setting.Opening[i].OpenFrom, req.Setting.Opening[i].OpenTo)

		if err != nil {
			return nil, status.Errorf(codes.Unknown, "failed to insert menu information"+err.Error())
		}

		addonsId, err := addonsResponse.LastInsertId()
		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve id for created menu "+err.Error())
		}

		log.Print(addonsId)
	}

	if r == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("setting with ID='%d' is not found",
			req.Setting.Id))
	}

	return &setting.UpdateResponse{
		Api:     apiVersion,
		Updated: r,
		Error: &setting.ResponseStatus{
			Status:  false,
			Message: "setting updated successfully",
		},
	}, nil
}

func (s *settingServiceServer) Delete(ctx context.Context, req *setting.DeleteRequest) (*setting.DeleteResponse, error) {

	if ctx == nil {
		return &setting.DeleteResponse{
			Api: apiVersion,
			Error: &setting.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &setting.DeleteResponse{
			Error: &setting.ResponseStatus{
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
	res, err := c.ExecContext(ctx, "DELETE FROM users WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete setting -> "+err.Error())
	}

	_, err = c.ExecContext(ctx, "DELETE FROM open_time WHERE `userId`=?", req.Id)
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

	return &setting.DeleteResponse{
		Api:     apiVersion,
		Deleted: int32(rows),
		Error: &setting.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
