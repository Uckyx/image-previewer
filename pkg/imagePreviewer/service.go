package imagePreviewer

import "context"

func NewApp() {

}

type Service interface {
	GetResizedImage(ctx context.Context)
}
