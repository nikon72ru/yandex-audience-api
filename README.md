# yandex-audience-api
#### Yandex Audience Api with Golang

##### official documentation: https://yandex.ru/dev/audience/doc/concept/about-docpage/

## Supported methods 
| Entity | Method | Support |
|--------|-------|-------|
| Pixels | [List of pixels](https://yandex.ru/dev/audience/doc/pixels/segments-docpage/) | :heavy_check_mark: | 
| Pixels | [Create pixel](https://yandex.ru/dev/audience/doc/pixels/create-docpage/) | :heavy_check_mark: | 
| Pixels | [Update pixel](https://yandex.ru/dev/audience/doc/pixels/edit-docpage/) | :heavy_check_mark: | 
| Pixels | [Remove pixel](https://yandex.ru/dev/audience/doc/pixels/delete-docpage/) | :heavy_check_mark: | 
| Pixels | [Restore pixel](https://yandex.ru/dev/audience/doc/pixels/undelete-docpage/) | :heavy_check_mark: | 
## Quickstart
``` golang
package main

import (
	"context"
	"fmt"
	"github.com/nikon72ru/yandex-audience-api/audience"
	"log"
)

func main() {
	//Creating audience client
	client, _ := audience.NewClient(context.Background())
	//Get all my segments
	allMySegments, _ := client.SegmentsList()
	_ = allMySegments

	//Create segment
	var segment = audience.CircleGeoSegment{
		BaseSegment: audience.BaseSegment{
			Name: "My new segment",
		},
		GeoSegmentType: "work",
		TimesQuantity:  20,
		PeriodLength:   30,
		Radius:         500,
		Points: []audience.Point{{
			Latitude:    65.534102,
			Longitude:   57.157753,
			Description: "random point",
		}}}
	if err := client.CreateCircleGeoSegment(&segment); err != nil {
		log.Fatal(err)
	}

	//Check ID of created segment
	fmt.Println(segment.ID)

	//Updating created segment
	segment.Name = "My updated segment"
	if err := client.UpdateSegment(segment.ID, &segment); err != nil {
		log.Fatal(err)
	}

	//Remove segment
	if err := client.RemoveSegment(segment.ID); err != nil {
		log.Fatal(err)
	}

}
```
-------------------------------------
## Token
### There are two ways to specify your token.
##### 1. As environment variable "YANDEX_AUDIENCE_TOKEN":
``` bash
export YANDEX_AUDIENCE_TOKEN=[YOUR TOKEN]
```
##### 2. In context:
``` golang
func main() {
	//Creating audience client
	client, _ := audience.NewClient(context.WithValue(context.Background(), "YANDEX_AUDIENCE_TOKEN", "[YOUR TOKEN]"))
	//Your another cool code
}
```
-------------------------------------
## Segment from file
### You can upload files with minimal buffering on server side
``` diff
- Notice! You need to save segment after uploading
```
``` golang
func main() {
	//Creating audience client
	client, _ := audience.NewClient(context.Background())

	var segment = audience.UploadingSegment{
		BaseSegment: audience.BaseSegment{
			Name: "segment from file",
		},
		Hashed:      false,
		ContentType: audience.Mac,
	}
	//upload file
	if err := client.CreateFileSegment(&segment, "./test-files/macs_for_uploads.csv"); err != nil {
		log.Fatal(err)
	}

	//Save uploaded segment
	if err := client.SaveUploadedSegment(&segment); err != nil {
		log.Fatal(err)
	}
}
```
CreateCSVSegment() and CreateReaderSegment() methods also supported
----------------------------------------
## Any questions?
Welcome to create issue!
