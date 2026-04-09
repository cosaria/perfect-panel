package cmd_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestPhase7NoLegacyPackageImportsRemain(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位当前测试文件路径")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
	base := "github.com/perfect-panel/server/"
	legacyPrefixes := []string{
		base + "models",
		base + "modules",
		base + "adapter",
		base + "types",
	}

	offenders := make([]string, 0)
	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "bin" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fileNode, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		for _, imp := range fileNode.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")
			for _, legacy := range legacyPrefixes {
				if importPath == legacy || strings.HasPrefix(importPath, legacy+"/") {
					offenders = append(offenders, rel+": "+importPath)
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk server module: %v", err)
	}

	if len(offenders) > 0 {
		slices.Sort(offenders)
		t.Fatalf("检测到旧包导入残留（%d 项）:\n%s", len(offenders), strings.Join(offenders, "\n"))
	}
}
