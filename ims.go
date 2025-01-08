package main

import (
	"fmt"

	imsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
)

func GetImageIdAndDisk(config *Config, client Client) (string, int32, error) {
	request := &imsMdl.ListImagesRequest{}
	request.Id = &config.ImageId
	response, err := client.ListImages(request)
	if err != nil {
		return "", 0, fmt.Errorf("error describing Images: %v", err)
	}
	if len(*response.Images) < 1 {
		return "", 0, fmt.Errorf("no image for given filters not found")
	}

	images := *response.Images
	if images[0].Status != imsMdl.GetImageInfoStatusEnum().ACTIVE {
		return "", 0, fmt.Errorf("%s invalid image status: %s", config.ImageId, images[0].Status)
	}
	return images[0].Id, images[0].MinDisk, nil
}
