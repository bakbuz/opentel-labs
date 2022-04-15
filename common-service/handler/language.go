package handler

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"

	"maydere.com/opentel-labs/common-service/pb"
)

func (h *Handler) GetLanguages(ctx context.Context, req *emptypb.Empty) (*pb.LanguagesResponse, error) {
	ctx, span := otel.Tracer("common-service").Start(ctx, "GetLanguages")
	defer span.End()

	dbResult, err := h.LanguageStore.GetAllLanguages(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var items []*pb.Language = make([]*pb.Language, len(dbResult))
	for i, dbItem := range dbResult {
		items[i] = &pb.Language{
			Id:           int32(dbItem.Id),
			CultureCode:  dbItem.CultureCode,
			Name:         dbItem.Name,
			Rtl:          dbItem.Rtl,
			Published:    dbItem.Published,
			DisplayOrder: int32(dbItem.DisplayOrder),
		}
	}
	return &pb.LanguagesResponse{Data: items}, nil
}

func (h *Handler) GetLanguage(ctx context.Context, req *pb.Identifier) (*pb.LanguageResponse, error) {
	ctx, span := otel.Tracer("common-service").Start(ctx, "GetLanguage/{id}")
	defer span.End()

	dbItem, err := h.LanguageStore.GetLanguageById(ctx, int(req.Id))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	item := &pb.Language{
		Id:           int32(dbItem.Id),
		CultureCode:  dbItem.CultureCode,
		Name:         dbItem.Name,
		Rtl:          dbItem.Rtl,
		Published:    dbItem.Published,
		DisplayOrder: int32(dbItem.DisplayOrder),
	}

	return &pb.LanguageResponse{Data: item}, nil
}
