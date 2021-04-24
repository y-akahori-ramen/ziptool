package ziptool

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func archiveInternal(dst, src string) (bool, error) {
	fileCreated := false

	// 内部的には/のパスで扱う
	src = filepath.ToSlash(src)

	// 出力するファイルの保存先ディレクトリがなければ作成
	if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
		return fileCreated, err
	}

	dstZip, err := os.Create(dst)
	if err != nil {
		return fileCreated, err
	}
	defer dstZip.Close()

	fileCreated = true
	w := zip.NewWriter(dstZip)
	defer w.Close()

	// 末尾に/がついているディレクトリ指定は/を取り除く
	src = strings.TrimSuffix(src, "/")

	// アーカイブ対象が存在するか
	_, err = os.Stat(filepath.FromSlash(src))
	if os.IsNotExist(err) {
		return fileCreated, err
	}

	// srcの末尾には/がついていないため付与。皇族のTrimPrefixで/まで含めて取り除いたパスで取り除きたいため
	baseDir := path.Dir(src) + "/"

	err = filepath.Walk(filepath.FromSlash(src), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		pathName := strings.TrimPrefix(filepath.ToSlash(path), baseDir)
		f, err := w.Create(pathName)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		if err != nil {
			return err
		}

		return nil
	})

	return fileCreated, err
}

// Archive dstで指定したパスのファイル・ディレクトリをzipアーカイブしてsrcに出力する
func Archive(dst, src string) error {
	zipFileCreated, err := archiveInternal(dst, src)

	// アーカイブに失敗した場合に出力先のzipが作成されていればそのzipは不正データなので削除しておく
	if err != nil && zipFileCreated {
		os.Remove(dst)
	}

	return err
}

// Unarchive zipFileSrcで指定されたzipファイルをdstに解凍する
func Unarchive(dst, zipFileSrc string) error {
	r, err := zip.OpenReader(zipFileSrc)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Archiveで圧縮したファイルのパスは/で入っているため実行中プラットフォームの形式に変換
		dstFilePath := filepath.Join(dst, filepath.FromSlash(f.Name))

		// 出力するファイルの保存先ディレクトリがなければ作成
		if err := os.MkdirAll(filepath.Dir(dstFilePath), 0777); err != nil {
			return err
		}

		dstFile, err := os.Create(dstFilePath)
		if err != nil {
			return err
		}

		zipContentReader, err := f.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		io.Copy(dstFile, zipContentReader)

		zipContentReader.Close()
		dstFile.Close()
	}

	return nil
}
