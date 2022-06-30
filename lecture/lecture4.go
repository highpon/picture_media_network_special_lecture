package lecture

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func Lecture4(inputDirPath, outputDirPath string) (retErr error) {
	if err := CheckExistDir(inputDirPath); err != nil {
		retErr = errors.Wrap(err, ErrDirNotFound)
		return
	}

	if err := CheckExistDir(outputDirPath); err != nil {
		retErr = errors.Wrap(err, ErrDirNotFound)
		return
	}

	inputFileLists, err := GetFileLists(inputDirPath)
	if err != nil {
		retErr = errors.Wrap(err, "failed GetFileLists")
		return err
	}

	fmt.Println(inputFileLists)

	for _, inputImagePath := range inputFileLists {
		inputImageFile, err := os.Open(inputImagePath)
		if err != nil {
			retErr = errors.Wrap(err, "failed to os.Open")
		}
		defer inputImageFile.Close()

		inputImg, _, err := image.Decode(inputImageFile)
		if err != nil {
			retErr = errors.Wrap(err, "failed image.Decode")
		}

		for _, quality := range []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100} {
			outputFileName := outputDirPath + strings.TrimSuffix(filepath.Base(inputImagePath), filepath.Ext(inputImagePath)) + "_q_" + strconv.Itoa(quality) + "." + "jpg"

			outputImg, err := os.Create(outputFileName)
			if err != nil {
				retErr = errors.Wrapf(err, "failed os.Create")
			}
			defer func() {
				if err := outputImg.Close(); err != nil {
					retErr = errors.Wrap(err, "failed outputFile.Close")
				}
			}()

			if err := jpeg.Encode(outputImg, inputImg, &jpeg.Options{Quality: quality}); err != nil {
				retErr = errors.Wrap(err, "failed jpeg.Encode")
				return
			}

			outputFile, err := os.Stat(outputFileName)
			if err != nil {
				retErr = errors.Wrap(err, "failed os.Stat")
				return
			}

			if _, err := outputImg.Seek(0, 0); err != nil {
				retErr = errors.Wrap(err, "failed outputImg.Seek(0,0)")
				return
			}

			imgSize, err := getImageSize(outputImg)
			if err != nil {
				retErr = errors.Wrap(err, "failed getImageSize")
				return
			}

			fmt.Printf("%s: BPP - %f\n", outputFileName, float64(8*outputFile.Size())/float64(imgSize))
		}
		fmt.Println("---------------------------------------------")
	}

	return nil
}

func getImageSize(imgReader io.Reader) (int, error) {
	img, _, err := image.Decode(imgReader)
	if err != nil {
		fmt.Println()
		return 0, errors.Wrap(err, "failed image.Decode")
	}

	return img.Bounds().Dx() * img.Bounds().Dy(), nil
}
