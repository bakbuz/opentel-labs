package store

import (
	"context"
	"sync"

	"maydere.com/opentel-labs/common-service/pb"
)

type LanguageStore struct {
	session string
}

var languageStoreInstance *LanguageStore
var onceLanguageStore sync.Once

func GetLanguageStore(session string) *LanguageStore {
	onceLanguageStore.Do(func() {
		languageStoreInstance = &LanguageStore{
			session: session,
		}
	})
	return languageStoreInstance
}

func (s *LanguageStore) GetAllLanguages(ctx context.Context) ([]*pb.Language, error) {
	data := make([]*pb.Language, 2)
	data[0] = &pb.Language{
		Id:           1,
		CultureCode:  "tr",
		Name:         "Türkçe",
		Rtl:          false,
		Published:    true,
		DisplayOrder: 1,
	}
	data[1] = &pb.Language{
		Id:           2,
		CultureCode:  "en",
		Name:         "English",
		Rtl:          false,
		Published:    true,
		DisplayOrder: 2,
	}
	return data, nil
}

func (s *LanguageStore) GetLanguageById(ctx context.Context, id int) (*pb.Language, error) {
	data := &pb.Language{
		Id:           1,
		CultureCode:  "tr",
		Name:         "Türkçe",
		Rtl:          false,
		Published:    true,
		DisplayOrder: 1,
	}

	return data, nil
}
