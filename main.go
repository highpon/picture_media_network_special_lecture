package main

import (
	"log"
	"picture_media_network_special_lecture/lecture"
)

func main() {
	if err := lecture.Lecture2("./data_2/data_2/data/", "./data_2/data_2/output/"); err != nil {
		log.Println(err)
	}
}
