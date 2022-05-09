package image_previewer

import (
	"testing"
)

//картинка найдена в кэше;
//удаленный сервер не существует;
//удаленный сервер существует, но изображение не найдено (404 Not Found);
//удаленный сервер существует, но изображение не изображение, а скажем, exe-файл;
//удаленный сервер вернул ошибку;
//удаленный сервер вернул изображение;
//изображение меньше, чем нужный размер; и пр.

func TestOriginalImageInCache(t *testing.T) {
	//logger := log.With().Logger()
	//ctx := context.Background()
	//
	//svc := NewApp(
	//	logger,
	//	cache.NewCache(10),
	//	NewImageDownloader(logger),
	//	NewImageResizer(logger),
	//)
	//
	//tests := []struct {
	//	name    string
	//	ctx     context.Context
	//	width   int
	//	height  int
	//	url     string
	//	imgName string
	//}{
	//	{
	//		name:    "get_original_img_from_cache_case",
	//		ctx:     ctx,
	//		width:   256,
	//		height:  126,
	//		url:     "",
	//		imgName: "gopher_256x126_resized.jpg",
	//	},
	//}
	//
	//t.Run(tt.name, func(t *testing.T) {
	//	gotImg, err := id.Download(tt.ctx, ImageURL+tt.imgName)
	//	if err != nil {
	//		t.Errorf("Download() error = %v", err)
	//		return
	//	}
	//
	//	wantImg := loadImage(tt.imgName)
	//	if !reflect.DeepEqual(gotImg, wantImg) {
	//		t.Errorf("Download() gotImg = %v, want %v", gotImg, wantImg)
	//	}
	//})
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
