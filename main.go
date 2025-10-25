package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	inputDir      = flag.String("input", ".", "Directory path to process")
	dryRun        = flag.Bool("dry-run", false, "Preview operations without executing")
	season        = flag.Int("season", 0, "Season number (required, no default)")       // Must be specified
	showName      = flag.String("show", "", "Anime series name (required, no default)") // Must be specified
	keepOtherTags = flag.Bool("keep-other-tags", true, "Preserve existing tags")
)

type NameInfo struct {
	RawName   string
	Ext       string
	Episode   int // -1 indicates no episode number found
	OtherTags []string
}

var (
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

	episodeRegex = regexp.MustCompile(`- (\d{1,2}) |\[(\d{1,2})\]|E(\d{1,2})||EP(\d{1,2})|第(\d{1,2})集`)
)

func main() {

	flag.Parse()

	missing := []string{}
	if *season == 0 {
		missing = append(missing, "-season")
	}
	if *showName == "" {
		missing = append(missing, "-show")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "error: Missing required parameters: %s\n", strings.Join(missing, ", "))
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	names, err := getDirFileNames(*inputDir)
	if err != nil {
		fmt.Printf("error: Failed to read directory contents: %v", err)
		os.Exit(1)
	}

	if len(names) == 0 {
		fmt.Println("no matching files found")
		return
	}

	for _, name := range names {
		oldPath := filepath.Join(*inputDir, name)
		info := getNameInfo(name)

		if info.Episode == -1 {
			fmt.Printf("warning: cannot extract episode number, skipping: %s\n", name)
			continue
		}

		newName := buildScrapedName(*showName, *season, info)
		if newName == "" {
			fmt.Printf("warning: failed to generate new filename, skipping: %s\n", name)
			continue
		}

		newPath := filepath.Join(*inputDir, newName)

		if oldPath == newPath {
			fmt.Printf("success: filename already correct, no change needed: %s\n", name)
			continue
		}

		if *dryRun {
			fmt.Printf("info: dry-run: move %s -> %s\n", name, newName)
		} else {
			if err := safeMove(oldPath, newPath); err != nil {
				fmt.Printf("error: failed to move file %s -> %s: %v\n", name, newName, err)
			} else {
				fmt.Printf("success: moved %s -> %s\n", name, newName)
			}
		}
	}
}

func getDirFileNames(dirPath string) ([]string, error) {
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

		if isHiddenFile(name) {
			continue
		}

		if !isAllowExt(filepath.Ext(name)) {
			continue
		}

		names = append(names, name)
	}

	return names, nil
}

func isAllowExt(ext string) bool {
	return AllowExts[ext]
}

func getNameInfo(rawName string) *NameInfo {

	info := &NameInfo{
		RawName:   rawName,
		Ext:       filepath.Ext(rawName),
		Episode:   extractEpisode(rawName),
		OtherTags: extractOtherTags(rawName),
	}

	return info
}

func extractOtherTags(rawName string) []string {

	tags := make([]string, 0)

	re := regexp.MustCompile(`\[.*?\]`)
	matches := re.FindAllString(rawName, -1)

	tags = append(tags, matches...)

	return tags
}

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

func buildScrapedName(showName string, season int, info *NameInfo) string {

	if info.Episode == -1 {
		return ""
	}

	otherTagsStr := ""
	if *keepOtherTags && len(info.OtherTags) > 0 {
		otherTagsStr = " " + strings.Join(info.OtherTags, " ")
	}

	seasonStr := padNumber(fmt.Sprintf("%d", season), 2)
	episodeStr := padNumber(fmt.Sprintf("%d", info.Episode), 2)

	return fmt.Sprintf("%s S%sE%s%s%s", showName, seasonStr, episodeStr, otherTagsStr, info.Ext)
}

func padNumber(num string, width int) string {
	if len(num) >= width {
		return num
	}
	return strings.Repeat("0", width-len(num)) + num
}

func safeMove(src, dst string) error {
	if fileExists(dst) {
		return fmt.Errorf("error: destination file already exists")
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isHiddenFile(path string) bool {
	filename := filepath.Base(path)
	return strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "~")
}
