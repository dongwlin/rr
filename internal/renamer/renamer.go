package renamer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// NameInfo holds extracted metadata from a filename.
type NameInfo struct {
	RawName   string
	Ext       string
	Episode   int // -1 indicates no episode number found
	OtherTags []string
}

var (
	// AllowExts defines the file extensions we process.
	AllowExts = map[string]bool{
		// video
		".mp4": true,
		".mkv": true,
		".avi": true,
		// subtitle
		".srt": true,
		".ass": true,
		".ssa": true,
		".idx": true,
	}

	tagRegex = regexp.MustCompile(`\[.*?\]`)

	episodeRegex = regexp.MustCompile(`` +
		`- (\d{1,2}) ` +
		`|\[(\d{1,2})\]` +
		`|E(\d{1,2})` +
		`|EP(\d{1,2})` +
		`|第(\d{1,2})集`,
	)
)

// GetDirFileNames returns non-hidden files from dirPath that match allowed extensions.
func GetDirFileNames(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if IsHiddenFile(name) {
			continue
		}

		if !IsAllowExt(filepath.Ext(name)) {
			continue
		}

		names = append(names, name)
	}

	return names, nil
}

// IsAllowExt reports whether the extension is allowed.
func IsAllowExt(ext string) bool {
	return AllowExts[ext]
}

// GetNameInfo extracts structured NameInfo from a raw filename.
func GetNameInfo(rawName string) *NameInfo {
	info := &NameInfo{
		RawName:   rawName,
		Ext:       filepath.Ext(rawName),
		Episode:   extractEpisode(rawName),
		OtherTags: extractOtherTags(rawName),
	}

	return info
}

// extractOtherTags finds and returns all bracketed tags in the provided
// filename. Tags are substrings that match the package-level `tagRegex`, for
// example "[Baha]" or "[1080p]". The returned slice preserves the order in
// which tags appear in the input. If no tags are found an empty slice is
// returned (not nil). Note: tags are not de-duplicated or normalized.
func extractOtherTags(rawName string) []string {
	tags := make([]string, 0)
	matches := tagRegex.FindAllString(rawName, -1)
	tags = append(tags, matches...)
	return tags
}

// extractEpisode attempts to locate an episode number in the provided
// filename using the package-level `episodeRegex`. The regex contains several
// alternative capture groups to support formats like:
//   - "- 04 "  (dash-space-number-space)
//   - "[02]"
//   - "E07" / "EP12"
//   - "第03集"
//
// The function examines the capture groups in order and returns the first
// non-empty numeric capture converted to an int. If no match is found, or if
// the captured text cannot be parsed as an integer, the function returns -1.
func extractEpisode(fileName string) int {
	match := episodeRegex.FindStringSubmatch(fileName)
	if len(match) == 0 {
		return -1
	}

	var episodeStr string
	for i := 1; i < len(match); i++ {
		if match[i] != "" {
			episodeStr = match[i]
			break
		}
	}

	episode, err := strconv.Atoi(episodeStr)
	if err != nil {
		return -1
	}

	return episode
}

// BuildScrapedName constructs the target filename. keepOtherTags controls whether
// existing bracketed tags are preserved in the result.
func BuildScrapedName(showName string, season int, info *NameInfo, keepOtherTags bool) string {
	if info.Episode == -1 {
		return ""
	}

	otherTagsStr := ""
	if keepOtherTags && len(info.OtherTags) > 0 {
		otherTagsStr = " " + strings.Join(info.OtherTags, " ")
	}

	seasonStr := padNumber(fmt.Sprintf("%d", season), 2)
	episodeStr := padNumber(fmt.Sprintf("%d", info.Episode), 2)

	return fmt.Sprintf("%s S%sE%s%s%s", showName, seasonStr, episodeStr, otherTagsStr, info.Ext)
}

// PadNumber returns num left-padded with zeros to the given width.
func padNumber(num string, width int) string {
	if len(num) >= width {
		return num
	}
	return strings.Repeat("0", width-len(num)) + num
}

// SafeMove moves src to dst. It tries os.Rename first; if that fails it copies
// the content and removes the source. It returns an error if dst already exists.
func SafeMove(src, dst string) error {
	if FileExists(dst) {
		return fmt.Errorf("destination file already exists")
	}

	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if info, err := os.Stat(src); err == nil {
		os.Chmod(dst, info.Mode())
	}

	return os.Remove(src)
}

// FileExists reports whether the path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsHiddenFile reports whether the filename is hidden (starts with '.' or '~').
func IsHiddenFile(path string) bool {
	filename := filepath.Base(path)
	return strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "~")
}
