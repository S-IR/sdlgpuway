package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)
func CopyDir(src, dst string) error {
    return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        rel, _ := filepath.Rel(src, path)
        target := filepath.Join(dst, rel)

        if d.IsDir() {
            return os.MkdirAll(target, 0o755)
        }

        in, err := os.Open(path)
        if err != nil { return err }
        defer in.Close()

        out, err := os.Create(target)
        if err != nil { return err }
        defer out.Close()

        _, err = io.Copy(out, in)
        return err
    })
}
