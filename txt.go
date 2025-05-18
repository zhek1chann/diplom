package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type FileEntry struct {
	Path     string
	Contents string
}

func main() {
	root := "."
	outputFile := "project_summary.txt"
	var entries []FileEntry
	var structure strings.Builder

	// Сканируем .go файлы
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "vendor") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			relPath, _ := filepath.Rel(root, path)
			dir := filepath.Dir(relPath)
			entries = append(entries, FileEntry{
				Path:     relPath,
				Contents: readFile(path),
			})
			addToStructure(&structure, dir, filepath.Base(path))
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	// Пишем основной файл
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("📂 PROJECT STRUCTURE (.go only):\n\n")
	writer.WriteString(structure.String())
	writer.WriteString("\n\n📄 GO FILES CONTENT:\n\n")

	for _, entry := range entries {
		writer.WriteString(fmt.Sprintf("%s\n", entry.Path))
		writer.WriteString(strings.Repeat("-", 80) + "\n")
		writer.WriteString(entry.Contents + "\n\n")
	}
	writer.Flush()

	fmt.Println("✅ Project summary written to", outputFile)

	// Разделение на две части
	splitFileIntoTwo(outputFile, "project_summary_part1.txt", "project_summary_part2.txt")
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("// ERROR reading %s: %v", path, err)
	}
	return string(content)
}

func addToStructure(sb *strings.Builder, dir, file string) {
	levels := strings.Split(dir, string(os.PathSeparator))
	indent := ""
	for i, level := range levels {
		if level == "." {
			continue
		}
		indent = strings.Repeat("│   ", i)
		sb.WriteString(fmt.Sprintf("%s├── /%s\n", indent, level))
	}
	indent = strings.Repeat("│   ", len(levels))
	sb.WriteString(fmt.Sprintf("%s└── %s\n", indent, file))
}

func splitFileIntoTwo(input, output1, output2 string) {
	// Читаем все строки
	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	mid := len(lines) / 2

	// Пишем первую половину
	err = os.WriteFile(output1, []byte(strings.Join(lines[:mid], "\n")), 0644)
	if err != nil {
		panic(err)
	}

	// Пишем вторую половину
	err = os.WriteFile(output2, []byte(strings.Join(lines[mid:], "\n")), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("📝 File split into:", output1, "and", output2)
}
