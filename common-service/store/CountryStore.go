package store

import (
	"context"
	"sync"

	"maydere.com/opentel-labs/common-service/pb"
)

type CountryStore struct {
	session string
}

var countryStoreInstance *CountryStore
var onceCountryStore sync.Once

func GetCountryStore(session string) *CountryStore {
	onceCountryStore.Do(func() {
		countryStoreInstance = &CountryStore{
			session: session,
		}
	})
	return countryStoreInstance
}

func (s *CountryStore) GetAllCountries(ctx context.Context) ([]*pb.Country, error) {
	data := make([]*pb.Country, 2)
	data[0] = &pb.Country{
		Id:           1,
		Name:         "Türkiye",
		EnglishName:  "Turkiye",
		IsoCode2:     "TR",
		IsoCode3:     "TUR",
		CallingCode:  90,
		Published:    true,
		DisplayOrder: 1,
	}
	data[1] = &pb.Country{
		Id:           2,
		Name:         "Amerika Birleşik Devletleri",
		EnglishName:  "United States",
		IsoCode2:     "US",
		IsoCode3:     "USA",
		CallingCode:  1,
		Published:    true,
		DisplayOrder: 2,
	}

	return data, nil
}

func (s *CountryStore) GetCountryById(ctx context.Context, id int) (*pb.Country, error) {
	data := &pb.Country{
		Id:           1,
		Name:         "Türkiye",
		EnglishName:  "Turkiye",
		IsoCode2:     "TR",
		IsoCode3:     "TUR",
		CallingCode:  90,
		Published:    true,
		DisplayOrder: 1,
	}

	return data, nil
}
