package lecture

import (
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"

	"github.com/pkg/errors"
)

func Lecture3(inputDirPath, outputDirPath string) (retErr error) {
	if err := CheckExistDir(inputDirPath); err != nil {
		retErr = errors.Wrap(err, ErrDirNotFound)
	}

	if err := CheckExistDir(outputDirPath); err != nil {
		retErr = errors.Wrap(err, ErrDirNotFound)
	}

	inputFileLists, err := GetFileLists(inputDirPath)
	if err != nil {
		retErr = errors.Wrap(err, "failed GetFileLists")
	}

	for _, inputImage := range inputFileLists {
		inputImageFile, err := os.Open(inputImage)
		if err != nil {
			retErr = errors.Wrap(err, "failed to os.Open")
		}
		defer inputImageFile.Close()

		inputImg, err := bmp.Decode(inputImageFile)
		fmt.Println("inputImageFile: ", inputImage)
		if err != nil {
			return errors.Wrap(err, "failed image.Decode")
		}

		baseFile, err := os.Stat(inputImage)
		if err != nil {
			return errors.Wrap(err, "failed os.Stat")
		}

		fmt.Printf("ファイル名: %s\n", baseFile.Name())
		fmt.Printf("ファイルサイズ(byte): %d\n", baseFile.Size())

		for _, encodeAlgorithm := range []string{"png", "gif", "jpeg"} {
			fmt.Printf("エンコードアルゴリズム: %s\n", encodeAlgorithm)
			baseOutputString := outputDirPath + strings.TrimSuffix(filepath.Base(inputImage), filepath.Ext(inputImage)) + "." + encodeAlgorithm

			outputImg, err := os.Create(baseOutputString)
			if err != nil {
				return errors.Wrapf(err, "failed os.Create")
			}
			defer func() {
				if err := outputImg.Close(); err != nil {
					retErr = errors.Wrap(err, "failed outputFile.Close")
				}
			}()

			switch encodeAlgorithm {
			case "jpeg", "jpg":
				jpeg.Encode(outputImg, inputImg, &jpeg.Options{})
			case "png":
				png.Encode(outputImg, inputImg)
			case "gif":
				gif.Encode(outputImg, inputImg, nil)
			default:
			}

			outputFile, err := os.Stat(baseOutputString)
			if err != nil {
				return errors.Wrap(err, "failed os.Stat")
			}
			fmt.Printf("圧縮後(%s)サイズ(byte): %d\n", encodeAlgorithm, outputFile.Size())
			fmt.Printf("割合: %f\n", float64(outputFile.Size())/float64(baseFile.Size()))
		}
		fmt.Println("---------------------------------------------")
	}
	return nil
}
