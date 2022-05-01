package imagePreviewer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"

	"image-previewer/pkg/cache"
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
	vars := mux.Vars(r)
	rawUrl := vars["imageUrl"]
	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		return errors.New("не удалось сконвертировать строку с параметром height в число")
	}
	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		return errors.New("не удалось сконвертировать строку с параметром width в число")
	}

	//todo Проверить кеш На наличие оригинальной картинки

	//Вынести в сервис imageDownloader
	if len(rawUrl) == 0 {
		return errors.New("передан пустой урл, для скачивания картинки")
	}

	parsedUrl, err := url.Parse(rawUrl)
	if parsedUrl.Scheme != "" {
		return errors.New("ссылка не должна содержать схему http или https")
	}

	imageUrl := fmt.Sprintf("https://%s", rawUrl)
	fmt.Println(imageUrl)

	filePath, err := download(imageUrl)
	if err != nil {
		fmt.Println(err)

		return err
	}
	//---------------------------------

	//todo проверить кеш на наличие обрезанной картинки

	//Вынести в image resizer
	src, err := imaging.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	currentTimeStamp, err := fmt.Println(time.Now().Format(time.RFC850))
	if err != nil {
		return errors.New("не удалось получить текущее время для названия файла")
	}

	resizedImageName := fmt.Sprintf("image_%d_resized.jpg", currentTimeStamp)
	img := imaging.Resize(src, width, height, imaging.Lanczos)
	err = imaging.Save(img, resizedImageName)
	if err != nil {
		return errors.New("не удалось сохранить обрезанную картинку")
	}

	//-----------------------------------

	//todo Положить в кеш
	//todo Отдать клиенту

	return nil
}

func download(downloadUrl string) (filePath string, err error) {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	currentTimeStamp, err := fmt.Println(time.Now().Format(time.RFC850))
	if err != nil {
		return "", errors.New("не удалось получить текущее время для названия файла")
	}

	filePath = fmt.Sprintf("image_%d.jpg", currentTimeStamp)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return filePath, err
}
