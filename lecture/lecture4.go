package lecture

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gonum.org/v1/plot/plotter"
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

	PlotFontInit()
	p := PlotInit("RD曲線", "BPP(bit per pixel)", "PSNR")
	colors := []color.RGBA{color.RGBA{R: 255, A: 255}, color.RGBA{G: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{R: 128, B: 128, A: 255}, color.RGBA{G: 128, B: 128, A: 255}}

	for index, inputImagePath := range inputFileLists {
		inputImageFile, err := os.Open(inputImagePath)
		if err != nil {
			retErr = errors.Wrap(err, "failed to os.Open")
		}
		defer inputImageFile.Close()

		inputImg, _, err := image.Decode(inputImageFile)
		if err != nil {
			retErr = errors.Wrap(err, "failed image.Decode")
		}

		pts := make(plotter.XYs, 10)
		for i, quality := range []int{1, 10, 20, 30, 40, 50, 60, 70, 80, 90} {
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

			inputImageAbsPath, err := filepath.Abs(inputImagePath)
			if err != nil {
				retErr = errors.Wrap(err, "failed filapath.Abs")
				return
			}

			outputFileAbsPathName, err := filepath.Abs(outputFileName)
			if err != nil {
				retErr = errors.Wrap(err, "failed filapath.Abs")
				return
			}

			psnr, err := getPSNR(inputImageAbsPath, outputFileAbsPathName)
			if err != nil {
				retErr = errors.Wrap(err, "failed getPSNR")
			}

			pts[i].X, pts[i].Y = float64(8*outputFile.Size())/float64(imgSize), psnr
			fmt.Printf("%s: BPP - %f\n", outputFileName, float64(8*outputFile.Size())/float64(imgSize))
			fmt.Printf("PSNR: %f\n", psnr)
		}
		CreatePlot(p, pts, GraphParams{Name: strings.TrimSuffix(filepath.Base(inputImagePath), filepath.Ext(inputImagePath)), LineColor: colors[index], PointColor: colors[index]}, outputDirPath+strings.TrimSuffix(filepath.Base(inputImagePath), filepath.Ext(inputImagePath))+"_point.png", true)
		fmt.Println(strings.TrimSuffix(filepath.Base(inputImagePath), filepath.Ext(inputImagePath)))
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

func getPSNR(inputFilePath, outputFilePath string) (float64, error) {
	cmd := exec.Command("python3", "./lecture/psnr.py", "-inputPath", inputFilePath, "-outputPath", outputFilePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		// fmt.Printf("Stdout: %s\n", stdout.String())
		// fmt.Printf("Stderr: %s\n", stderr.String())
	}
	// fmt.Printf("Stdout: %s\n", stdout.String())
	psnr, err := strconv.ParseFloat(stdout.String(), 64)
	if err != nil {
		fmt.Println(err)
		return 0, errors.Wrap(err, "failed strconv.ParseFloat")
	}
	return psnr, err
}
