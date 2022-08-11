package main

import (
	"flag"
	"fmt"
	"log"
	"picture_media_network_special_lecture/lecture"
)

var lectureFlag int
var inputPath, outputPath string

func init() {
	flag.IntVar(&lectureFlag, "lecture", -1, "講義番号を入力")
	flag.StringVar(&inputPath, "inputPath", "", "Input file path")
	flag.StringVar(&outputPath, "outputPath", "", "Output file path")
}

func main() {
	// if err := lecture.Lecture2("./data_2/data_2/data/", "./data_2/data_2/output/"); err != nil {
	flag.Parse()
	switch lectureFlag {
	case 2:
		if err := lecture.Lecture2(inputPath, outputPath); err != nil {
			log.Println(err)
			return
		}
	case 3:
		if err := lecture.Lecture3(inputPath, outputPath); err != nil {
			log.Println(err)
			return
		}
	case 4:
		if err := lecture.Lecture4(inputPath, outputPath); err != nil {
			log.Println(err)
			return
		}
	case 5:
		if err := lecture.Lecture5(inputPath, outputPath); err != nil {
			log.Println(err)
			return
		}
	case 6:
		if err := lecture.Lecture6(inputPath, outputPath); err != nil {
			log.Println(err)
			return
		}
	default:
		fmt.Println("invalid lecture number")
	}
}
