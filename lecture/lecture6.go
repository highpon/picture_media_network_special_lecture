package lecture

import (
	"bytes"
	"fmt"
	"image/color"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gonum.org/v1/plot/plotter"
)

type Lesson6FFMPEGParam struct {
	crf        int
	outputFile string
	GraphParam GraphParams
}

func Lecture6(inputDirPath, outputDirPath string) (retErr error) {
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

	PlotFontInit()
	p := PlotInit("RD曲線", "bit/sec(kb/s)", "PSNR")

	for _, inputFile := range inputFileLists {
		pts := make(plotter.XYs, 100)
		tmp := Lesson6FFMPEGParam{crf: 0, GraphParam: GraphParams{Name: "", LineColor: color.RGBA{R: 255, A: 255}, PointColor: color.RGBA{R: 255, A: 255}}}
		for i, crf := range []int{1, 10, 16, 23, 29, 35, 41, 47, 53, 59, 66, 71, 82, 100} {
			// for i, crf := range []int{71, 82, 100} {
			tmp.crf = crf
			pts[i].X, pts[i].Y, err = getPSMRFromFFMPEGWithYuv(tmp, inputFile, outputDirPath+"/out.mp4")
			fmt.Println(pts[i].X, pts[i].Y)
			if err != nil {
				retErr = errors.Wrap(err, "failed getPSMRFromFFMPEG")
				return err
			}
			fmt.Println("crf", crf)
		}
		CreatePlot(p, pts, tmp.GraphParam, outputDirPath+"/point.png", true)
	}

	return nil
}

func getPSMRFromFFMPEGWithYuv(input Lesson6FFMPEGParam, inputFile, outputFile string) (float64, float64, error) {
	fmt.Println("crf", input.crf)
	cmd := exec.Command("ffmpeg", "-y", "-r", strconv.Itoa(15), "-i", inputFile, "-vcodec", "libx264", "-preset", "fast", "-crf", strconv.Itoa(input.crf), "-tune", "psnr", "-psnr", outputFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return 0.0, 0.0, errors.Wrap(err, "failed to cmd.Run")
	}

	psnrRow := strings.Split(strings.Split(stderr.String(), "\n")[len(strings.Split(stderr.String(), "\n"))-2], " ")
	retval := make(plotter.XYs, 1)
	if s := strings.Split(psnrRow[9], ":"); s[0] != "Global" {
		return 0.0, 0.0, errors.Wrap(nil, "failed ffmpeg format")
	} else {
		retval[0].Y, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			return 0.0, 0.0, errors.Wrap(err, "failed strconv.ParseFloat")
		}
	}

	if s := strings.Split(psnrRow[10], ":"); s[0] != "kb/s" {
		return 0.0, 0.0, errors.Wrap(nil, "failed ffmpeg format")
	} else {
		retval[0].X, err = strconv.ParseFloat(s[1], 64)
		if err != nil {
			return 0.0, 0.0, errors.Wrap(err, "failed strconv.ParseFloat")
		}
	}
	return retval[0].X, retval[0].Y, nil
}
