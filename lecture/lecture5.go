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

type Lesson5FFMPEGParam struct {
	Bf         int
	Gop        int
	GraphParam GraphParams
}

func Lecture5(inputDirPath, outputDirPath string) (retErr error) {
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

	tmp := []Lesson5FFMPEGParam{
		{Bf: 0, Gop: 0, GraphParam: GraphParams{Name: "Iフレームのみ", LineColor: color.RGBA{R: 255, A: 255}, PointColor: color.RGBA{R: 255, A: 255}}},
		{Bf: 0, Gop: 5, GraphParam: GraphParams{Name: "I,Pフレーム (Gop: 5)", LineColor: color.RGBA{B: 255, A: 255}, PointColor: color.RGBA{B: 255, A: 255}}},
		{Bf: 0, Gop: 50, GraphParam: GraphParams{Name: "I,Pフレーム(Gop: 50)", LineColor: color.RGBA{G: 255, A: 255}, PointColor: color.RGBA{G: 255, A: 255}}},
		{Bf: 8, Gop: 5, GraphParam: GraphParams{Name: "I,P,Bフレーム(Gop: 5, bf: 8)", LineColor: color.RGBA{R: 128, G: 128, A: 255}, PointColor: color.RGBA{R: 128, G: 128, A: 255}}},
		{Bf: 8, Gop: 50, GraphParam: GraphParams{Name: "I,P,Bフレーム(Gop: 50, bf: 8)", LineColor: color.RGBA{B: 128, G: 128, A: 255}, PointColor: color.RGBA{B: 128, G: 128, A: 255}}},
	}

	for _, inputFile := range inputFileLists {
		for _, param := range tmp {
			pts := make(plotter.XYs, 5)
			for i, crf := range []int{1, 10, 23, 35, 51} {
				fmt.Println("-------------------------------------------------")
				if param.Bf == 0 && param.Gop == 0 {
					fmt.Println("only I frame", crf)
				} else if param.Bf == 0 && (param.Gop == 5 || param.Gop == 50) {
					fmt.Println("I P frame", crf)
				} else if param.Bf == 8 && (param.Gop == 5 || param.Gop == 50) {
					fmt.Println("I P B frame", crf)
				}

				pts[i].X, pts[i].Y, err = getPSMRFromFFMPEG(inputFile, outputDirPath+"/out.avi", param.Gop, param.Bf, crf)
				if err != nil {
					retErr = errors.Wrap(err, "failed getPSMRFromFFMPEG")
					return err
				}
			}
			CreatePlot(p, pts, param.GraphParam, outputDirPath+"/point.png", true)
		}
	}

	return nil
}

func getPSMRFromFFMPEG(inputFilePath, outputFilePath string, gop, bf, crf int) (float64, float64, error) {
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFilePath, "-g", strconv.Itoa(gop), "-bf", strconv.Itoa(bf), "-crf", strconv.Itoa(crf), "-vcodec", "libx264", outputFilePath, "-loglevel", "info", "-psnr")
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
