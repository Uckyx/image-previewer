package image_previewer

import "testing"

//картинка найдена в кэше;
//удаленный сервер не существует;
//удаленный сервер существует, но изображение не найдено (404 Not Found);
//удаленный сервер существует, но изображение не изображение, а скажем, exe-файл;
//удаленный сервер вернул ошибку;
//удаленный сервер вернул изображение;
//изображение меньше, чем нужный размер; и пр.

func TestImageInCache(t *testing.T) {

}

func TestNotFoundRemoteService(t *testing.T) {

}

func TestAllowedFormatImage(t *testing.T) {

}

func TestNotFoundImgInRemoteService(t *testing.T) {

}

func TestErrorInResponseRemoteService(t *testing.T) {

}

func TestRemoteServiceSuccessResponseImage(t *testing.T) {

}
