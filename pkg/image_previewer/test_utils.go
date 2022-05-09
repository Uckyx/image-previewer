package image_previewer

import (
	"bufio"
	"fmt"
	"os"
)

const ImageURL = "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/"
const OriginalImgName = "_gopher_original_1024x504.jpg"
const ResizedImgName = "gopher_256x126_resized.jpg"

func loadImage(imgName string) []byte {
	fileToBeUploaded := "./image_test/" + imgName
	file, err := os.Open(fileToBeUploaded)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return bytes
}
