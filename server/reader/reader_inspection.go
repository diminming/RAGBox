package reader

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	docx "github.com/fumiama/go-docx"
)

func NewInspctRptReader(filePath string) (*InspectionReader, error) {
	resolvedPath, err := resolveReadablePath(filePath)
	if err != nil {
		slog.Error("文件不存在或无法访问", "error", err, "filePath", filePath)
		return nil, err
	}

	return &InspectionReader{filePath: resolvedPath}, nil
}

func resolveReadablePath(filePath string) (string, error) {
	candidates := buildPathCandidates(filePath)

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("path not found on %s, tried: %s", runtime.GOOS, strings.Join(candidates, " | "))
}

func buildPathCandidates(filePath string) []string {
	seen := map[string]struct{}{}
	add := func(dst *[]string, p string) {
		p = strings.TrimSpace(p)
		if p == "" {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		seen[p] = struct{}{}
		*dst = append(*dst, p)
	}

	candidates := make([]string, 0, 4)
	add(&candidates, filePath)
	add(&candidates, strings.ReplaceAll(filePath, "\\", "/"))

	if runtime.GOOS != "windows" {
		if mntPath, rootPath, ok := convertWindowsPathToUnix(filePath); ok {
			add(&candidates, mntPath)
			add(&candidates, rootPath)
		}
	}

	return candidates
}

func convertWindowsPathToUnix(path string) (string, string, bool) {
	if len(path) < 3 || path[1] != ':' {
		return "", "", false
	}
	if path[2] != '\\' && path[2] != '/' {
		return "", "", false
	}

	drive := path[0]
	if drive >= 'A' && drive <= 'Z' {
		drive = drive - 'A' + 'a'
	}
	if drive < 'a' || drive > 'z' {
		return "", "", false
	}

	rest := strings.TrimLeft(path[2:], "\\/")
	rest = strings.ReplaceAll(rest, "\\", "/")

	// Try common mount styles in Linux/WSL environments.
	return fmt.Sprintf("/mnt/%c/%s", drive, rest), fmt.Sprintf("/%c/%s", drive, rest), true
}

type InspectionReader struct {
	filePath string
	FileName string
}

func detectHeadingLevel(para *docx.Paragraph) (int, bool) {
	if para == nil || para.Properties == nil || para.Properties.Style == nil {
		return 0, false
	}

	// Prefer style-based detection because it is explicit and stable in most docx files.
	style := strings.ToLower(strings.TrimSpace(para.Properties.Style.Val))
	style = strings.ReplaceAll(style, " ", "")

	switch style {
	case "heading1", "标题1", "标题一":
		return 1, true
	case "heading2", "标题2", "标题二":
		return 2, true
	case "heading3", "标题3", "标题三":
		return 3, true
	}

	return 0, false
}

func renderTable(table *docx.Table) string {
	if table == nil {
		return ""
	}

	rows := make([]string, 0, len(table.TableRows))
	for _, row := range table.TableRows {
		if row == nil {
			continue
		}

		cells := make([]string, 0, len(row.TableCells))
		for _, cell := range row.TableCells {
			cells = append(cells, renderTableCell(cell))
		}
		rows = append(rows, strings.Join(cells, "\t"))
	}

	return strings.Join(rows, "\n")
}

func renderTableCell(cell *docx.WTableCell) string {
	if cell == nil {
		return ""
	}

	parts := make([]string, 0, len(cell.Paragraphs)+len(cell.Tables))
	for _, para := range cell.Paragraphs {
		if para == nil {
			continue
		}

		text := strings.TrimSpace(para.String())
		if text != "" {
			parts = append(parts, text)
		}
	}
	for _, nestedTable := range cell.Tables {
		text := strings.TrimSpace(renderTable(nestedTable))
		if text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, " ")
}

func (r *InspectionReader) Read() ([]string, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		slog.Error("打开文件失败", "error", err, "filePath", r.filePath)
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		slog.Error("读取文件信息失败", "error", err, "filePath", r.filePath)
		return nil, err
	}
	r.FileName = info.Name()

	doc, err := docx.Parse(file, info.Size())
	if err != nil {
		slog.Error("解析文件失败", "error", err, "filePath", r.filePath)
		return nil, err
	}

	headingStack := make([]string, 10) // 用于存储当前的标题层级，最多支持10级标题
	infoStack := []string{}
	contentStack := []string{}
	depth := -1
	curr := -1
	// 遍历文档主体中的段落，识别标题并打印文本。
	for _, item := range doc.Document.Body.Items {

		switch node := item.(type) {
		case *docx.Paragraph:
			if level, ok := detectHeadingLevel(node); ok {
				// slog.Debug("检测到标题", "level", level, "text", node.String())
				curr = level - 1
				if curr > depth {
					depth = curr
				}
				headingStack[curr] = node.String()
			} else {
				if depth > -1 {
					content := fmt.Sprintf("[%s] [%s] %s", strings.Join(infoStack, " | "), strings.Join(headingStack[:curr+1], " > "), node.String())
					contentStack = append(contentStack, content)
					// slog.Debug("段落内容", "content", content)
				} else {
					infoStack = append(infoStack, node.String())
				}
			}
		case *docx.Table:
			content := renderTable(node)
			if depth > -1 {
				content := fmt.Sprintf("[%s] [%s] %s", strings.Join(infoStack, " | "), strings.Join(headingStack[:curr+1], " > "), content)
				// slog.Debug("段落内容", "content", content)
				contentStack = append(contentStack, content)
			} else {
				infoStack = append(infoStack, content)
			}
		}

	}
	pos := 0
	previous := ""
	section := strings.Builder{}
	sections := []string{}

	for idx, content := range contentStack {
		segments := strings.Split(content, " > ")
		segment := segments[len(segments)-1]
		segments = strings.Split(segment, "] ")
		segment = segments[0]

		if idx > 0 && segment != previous {

			previous = segment
			for pos < idx {
				section.WriteString(contentStack[pos])
				section.WriteString("\n")
				pos++
			}
			slog.Debug("获取到Block", "section", section.String())
			sections = append(sections, section.String())

			section = strings.Builder{}
			section.WriteString(content)
			section.WriteString("\n")

			pos = idx + 1
		} else {
			if previous == "" {
				contentStack[idx] = content
			} else {
				contentStack[idx] = segments[1]
			}

		}

	}

	slog.Debug("获取到Block", "section", section.String())
	sections = append(sections, section.String())
	return sections, nil

}
