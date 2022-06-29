package lecture

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	lz4 "github.com/pierrec/lz4/v4"
	"github.com/pkg/errors"
)

var (
	ErrDirNotFound = "directory not found"
)

func Lecture2(inputDirPath, outputDirPath string) error {
	if err := CheckExistDir(inputDirPath); err != nil {
		return errors.Wrap(err, ErrDirNotFound)
	}

	if err := CheckExistDir(outputDirPath); err != nil {
		return errors.Wrap(err, ErrDirNotFound)
	}

	inputFileLists, err := GetFileLists(inputDirPath)
	if err != nil {
		return err
	}

	for _, v := range inputFileLists {
		for _, compressAlgorithm := range []string{"zstd", "br", "lz4"} {
			fmt.Printf("圧縮アルゴリズム: %s\n", compressAlgorithm)
			baseOutputString := outputDirPath + filepath.Base(v)

			if err := createTar(v, baseOutputString); err != nil {
				return errors.Wrap(err, "failed createTar")
			}

			if err := compressFile(baseOutputString+".tar", compressAlgorithm); err != nil {
				return errors.Wrap(err, "failed compressFile")
			}

			baseFile, err := os.Stat(v)
			if err != nil {
				return errors.Wrap(err, "failed os.Stat")
			}

			outputFile, err := os.Stat(baseOutputString + ".tar." + compressAlgorithm)
			if err != nil {
				return errors.Wrap(err, "failed os.Stat")
			}

			fmt.Printf("ファイル名: %s\n", baseFile.Name())
			fmt.Printf("ファイルサイズ(byte): %d\n", baseFile.Size())
			fmt.Printf("圧縮後サイズ(byte): %d\n", outputFile.Size())
			fmt.Printf("割合: %f\n", float64(outputFile.Size())/float64(baseFile.Size()))
			fmt.Println("---------------------------------------------")
		}
	}

	return nil
}

func compressFile(inputFile, compressAlgorithm string) (retErr error) {
	dst, err := os.Create(inputFile + "." + compressAlgorithm)
	if err != nil {
		retErr = errors.Wrap(retErr, "failed os.Create")
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			retErr = errors.Wrap(retErr, "failed dst.Close")
		}
	}()

	f, err := os.Open(inputFile)
	if err != nil {
		retErr = errors.Wrap(retErr, "failed os.Open")
		return
	}
	defer f.Close()

	switch compressAlgorithm {
	case "zstd":
		if err := zstdCompress(f, dst); err != nil {
			retErr = errors.Wrap(retErr, "failed Compress")
			return
		}
	case "br":
		if err := brotliCompress(f, dst); err != nil {
			retErr = errors.Wrap(retErr, "failed brotliCompress")
			return
		}
	case "lz4":
		if err := lz4Compress(f, dst); err != nil {
			retErr = errors.Wrap(retErr, "failed lz4Compress")
			return
		}
	default:
		retErr = errors.Wrap(retErr, "failed undefined algorithm compressFile")
		return
	}

	return nil
}

func createTar(inputFile, outputDirPath string) (retErr error) {
	dst, err := os.Create(outputDirPath + ".tar")
	if err != nil {
		retErr = errors.Wrap(retErr, "failed os.Create")
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			retErr = errors.Wrap(retErr, "failed dst.Close")
		}
	}()

	tw := tar.NewWriter(dst)
	defer func() {
		if err := tw.Close(); err != nil {
			retErr = errors.Wrap(retErr, "failed tw.Close")
		}
	}()

	f, err := os.Open(inputFile)
	if err != nil {
		retErr = errors.Wrap(retErr, "failed os.Open")
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		retErr = errors.Wrap(retErr, "failed f.Stat")
		return
	}

	if err := tw.WriteHeader(&tar.Header{
		Name:    inputFile,
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
		Size:    stat.Size(),
	}); err != nil {
		retErr = errors.Wrap(retErr, "failed tw.WriteHeader")
		return
	}

	if _, err := io.Copy(tw, f); err != nil {
		retErr = errors.Wrap(retErr, "failed io.Copy")
		return
	}

	return nil
}

// Compress input to output.
func zstdCompress(in io.Reader, out io.Writer) error {
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return err
	}
	return enc.Close()
}

func brotliCompress(in io.Reader, out io.Writer) error {
	enc := brotli.NewWriter(out)
	_, err := io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return errors.Wrap(err, "failed brotliCompress")
	}
	return enc.Close()
}

func lz4Compress(in io.Reader, out io.Writer) error {
	enc := lz4.NewWriter(out)
	_, err := io.Copy(enc, in)
	if err != nil {
		enc.Close()
		return errors.Wrap(err, "failed lz4Compress")
	}
	return enc.Close()
}
