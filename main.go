package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dongwlin/rr/internal/renamer"

	"github.com/dongwlin/rr/internal/style"
)

var (
	inputDir      = flag.String("input", ".", "Directory path to process")
	dryRun        = flag.Bool("dry-run", false, "Preview operations without executing")
	season        = flag.Int("season", 0, "Season number (required, no default)")       // Must be specified
	showName      = flag.String("show", "", "Anime series name (required, no default)") // Must be specified
	keepOtherTags = flag.Bool("keep-other-tags", true, "Preserve existing tags")
	noColor       = flag.Bool("no-color", false, "Disable colored output")
)

func main() {

	flag.Parse()

	// honor -no-color flag by disabling styled output
	if *noColor {
		style.DisableColor()
	}

	missing := []string{}
	if *season == 0 {
		missing = append(missing, "-season")
	}
	if *showName == "" {
		missing = append(missing, "-show")
	}

	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "%s: missing required parameters: %s\n",
			style.RenderError("error"),
			strings.Join(missing, ", "),
		)
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	names, err := renamer.GetDirFileNames(*inputDir)
	if err != nil {
		fmt.Printf("%s: failed to read directory contents: %v",
			style.Error.Render("error"),
			err,
		)
		os.Exit(1)
	}

	if len(names) == 0 {
		fmt.Printf("%s: no matching files found\n",
			style.RenderWarning("warning"),
		)
		return
	}

	for _, name := range names {
		oldPath := filepath.Join(*inputDir, name)
		info := renamer.GetNameInfo(name)

		if info.Episode == -1 {
			fmt.Printf("%s: cannot extract episode number, skipping: %s\n",
				style.RenderWarning("warning"), name,
			)
			continue
		}

		newName := renamer.BuildScrapedName(*showName, *season, info, *keepOtherTags)
		if newName == "" {
			fmt.Printf("%s: failed to generate new filename, skipping: %s\n",
				style.RenderWarning("warning"),
				name,
			)
			continue
		}

		newPath := filepath.Join(*inputDir, newName)

		if oldPath == newPath {
			fmt.Printf("%s: filename already correct, no change needed: %s\n",
				style.RenderSuccess("success"),
				name,
			)
			continue
		}

		if *dryRun {
			fmt.Printf("%s: dry-run\n",
				style.RenderInfo("info"),
			)
			fmt.Printf("  from: %s\n", name)
			fmt.Printf("  to  : %s\n", newName)
			fmt.Println()
		} else {
			if err := renamer.SafeMove(oldPath, newPath); err != nil {
				fmt.Printf("%s: failed to move\n",
					style.RenderError("error"),
				)
				fmt.Printf("  from: %s\n", name)
				fmt.Printf("  to  : %s\n", newName)
				fmt.Printf("  cause: %v\n", err)
				fmt.Println()
			} else {
				fmt.Printf("%s: moved\n",
					style.RenderSuccess("success"),
				)
				fmt.Printf("  from: %s\n", name)
				fmt.Printf("  to  : %s\n", newName)
				fmt.Println()
			}
		}
	}
}
