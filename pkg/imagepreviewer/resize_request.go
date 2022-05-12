package imagepreviewer

import "context"

type ResizeRequest struct {
	ctx     context.Context
	width   int
	height  int
	url     string
	headers map[string][]string
}

func NewResizeRequest(
	ctx context.Context,
	width int,
	height int,
	url string,
	headers map[string][]string,
) *ResizeRequest {
	return &ResizeRequest{ctx, width, height, url, headers}
}
