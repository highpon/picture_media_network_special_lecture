package lecture

import (
	"fmt"
	"image/jpeg"
	"image/png"
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

	for _, inputImage := range inputFileLists {
		inputImageFile, err := os.Open(inputImage)
		if err != nil {
			retErr = errors.Wrap(err, "failed to os.Open")
		}
		defer inputImageFile.Close()

		inputImg, err := png.Decode(inputImageFile)
		if err != nil {
			retErr = errors.Wrap(err, "failed image.Decode")
		}

		for _, quality := range []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100} {
			baseOutputString := outputDirPath + strings.TrimSuffix(filepath.Base(inputImage), filepath.Ext(inputImage)) + "_q_" + strconv.Itoa(quality) + "." + "jpg"

			outputImg, err := os.Create(baseOutputString)
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
		}

	}

	return nil
}
