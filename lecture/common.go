package lecture

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
)

func CheckExistDir(path string) error {
	if d, err := os.Stat(path); os.IsNotExist(err) || d.IsDir() {
		return err
	}
	return nil
}

func GetFileLists(inputPath string) ([]string, error) {
	findList := []string{}
	depth := strings.Count(inputPath, string(os.PathSeparator))

	err := filepath.WalkDir(inputPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}
		if info.IsDir() {
			return nil
		}
		if depth < strings.Count(path, string(os.PathSeparator)) {
			return nil
		}

		findList = append(findList, path)
		return nil
	})
	return findList, err
}

func PlotInit(title, xLabel, yLabel string) *plot.Plot {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	p.Add(plotter.NewGrid())

	return p
}

func PlotFontInit() {
	// download font from debian
	const url = "http://http.debian.net/debian/pool/main/f/fonts-ipafont/fonts-ipafont_00303.orig.tar.gz"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("could not download IPA font file: %+v", err)
	}
	defer resp.Body.Close()

	ttf, err := untargz("IPAfont00303/ipam.ttf", resp.Body)
	if err != nil {
		log.Fatalf("could not untar archive: %+v", err)
	}

	fontTTF, err := opentype.Parse(ttf)
	if err != nil {
		log.Fatal(err)
	}
	mincho := font.Font{Typeface: "Mincho"}
	font.DefaultCache.Add([]font.Face{
		{
			Font: mincho,
			Face: fontTTF,
		},
	})
	if !font.DefaultCache.Has(mincho) {
		log.Fatalf("no font %q!", mincho.Typeface)
	}
	plot.DefaultFont = mincho
	plotter.DefaultFont = mincho
}

func untargz(name string, r io.Reader) ([]byte, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not create gzip reader: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("could not find %q in tar archive", name)
			}
			return nil, fmt.Errorf("could not extract header from tar archive: %v", err)
		}

		if hdr == nil || hdr.Name != name {
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, tr)
		if err != nil {
			return nil, fmt.Errorf("could not extract %q file from tar archive: %v", name, err)
		}
		return buf.Bytes(), nil
	}
}
