package handler

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"

	"maydere.com/opentel-labs/common-service/pb"
)

func (h *Handler) GetCountries(ctx context.Context, req *emptypb.Empty) (*pb.CountriesResponse, error) {
	fmt.Println("********************* GetCountries 1 ***********************************")
	ctx, span := otel.Tracer("common-service").Start(ctx, "GetCountries")
	defer span.End()

	fmt.Println("********************* GetCountries 2 ***********************************")

	dbResult, err := h.CountryStore.GetAllCountries(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fmt.Println("********************* GetCountries 3 ***********************************")

	var items []*pb.Country = make([]*pb.Country, len(dbResult))
	for i, dbItem := range dbResult {
		items[i] = &pb.Country{
			Id:           int32(dbItem.Id),
			Name:         dbItem.Name,
			EnglishName:  dbItem.EnglishName,
			IsoCode2:     dbItem.IsoCode2,
			IsoCode3:     dbItem.IsoCode3,
			CallingCode:  int32(dbItem.CallingCode),
			Published:    dbItem.Published,
			DisplayOrder: int32(dbItem.DisplayOrder),
		}
	}

	fmt.Println("********************* GetCountries 4 ***********************************")

	return &pb.CountriesResponse{Data: items}, nil
}

func (h *Handler) GetCountry(ctx context.Context, req *pb.Identifier) (*pb.CountryResponse, error) {
	ctx, span := otel.Tracer("common-service").Start(ctx, "GetCountry/{id}")
	defer span.End()

	dbItem, err := h.CountryStore.GetCountryById(ctx, int(req.Id))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &pb.Country{
		Id:           int32(dbItem.Id),
		Name:         dbItem.Name,
		EnglishName:  dbItem.EnglishName,
		IsoCode2:     dbItem.IsoCode2,
		IsoCode3:     dbItem.IsoCode3,
		CallingCode:  int32(dbItem.CallingCode),
		Published:    dbItem.Published,
		DisplayOrder: int32(dbItem.DisplayOrder),
	}

	return &pb.CountryResponse{Data: item}, nil
}
