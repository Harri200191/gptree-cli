package internal

import (
    "io/fs"
    "path/filepath"
    "strings"
)

func WalkDir(root string, ignore []string) ([]string, error) {
    var files []string
    ignoreSet := make(map[string]struct{})

    for _, dir := range ignore {
        ignoreSet[filepath.Join(root, dir)] = struct{}{}
    }

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        for ignored := range ignoreSet {
            if strings.HasPrefix(path, ignored) {
                return nil
            }
        }

        if !d.IsDir() && filepath.Ext(path) != "" {
            files = append(files, path)
        }
        return nil
    })

    return files, err
}
