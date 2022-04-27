package imagePreviewer

import (
	"fmt"
	"image-previewer/pkg/cache"
	"net/http"
)

func NewApp(cache cache.Cache) Service {
	return &service{
		cache: cache,
	}
}

type Service interface {
	Fill(w http.ResponseWriter, r *http.Request) error
}

type service struct {
	cache cache.Cache
}

func (s *service) Fill(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Ура мы тут")
	//todo Проверить кеш
	//todo Скачать с урла
	//todo Обрезать картинку
	//todo Положить в кеш
	//todo Отдать клиенту

	return nil
}
