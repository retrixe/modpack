package main

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzipFile(zipFile []byte, location string, exclude []string, include []string) error {
	r, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return err
	}
	includePresent := include != nil
	if includePresent && len(include) == 0 { // Don't bother extracting if no file is included.
		return nil
	}
	for _, f := range r.File {
		toContinue := includePresent       // If include is present, then continue by default.
		for _, excluded := range exclude { // If file is in exclude, then continue.
			if excluded == f.FileInfo().Name() {
				toContinue = true
			}
		}
		for _, included := range include { // If file is in include, then do not continue.
			if included == f.FileInfo().Name() {
				toContinue = false
			}
		}
		if toContinue {
			continue
		}
		fpath := filepath.Join(location, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(location)+string(os.PathSeparator)) {
			continue // "%s: illegal file path"
		}
		// Create folders.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		// Create parent folder of file if needed.
		err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		if err != nil {
			return err
		}
		// Open target file.
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		// Open file in zip.
		rc, err := f.Open()
		if err != nil {
			return err
		}
		// Copy file from zip to disk.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
		outFile.Close()
		rc.Close()
	}
	return nil
}
