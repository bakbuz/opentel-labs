package handler

import (
	"maydere.com/opentel-labs/common-service/pb"
	"maydere.com/opentel-labs/common-service/store"
)

type Handler struct {
	pb.UnimplementedCommonServiceServer
	LanguageStore *store.LanguageStore
	CountryStore  *store.CountryStore
}

func NewHandler(dbSess string) *Handler {
	h := &Handler{
		LanguageStore: store.GetLanguageStore(dbSess),
		CountryStore:  store.GetCountryStore(dbSess),
	}

	return h
}
