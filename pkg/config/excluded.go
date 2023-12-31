package config

var excludedDirsDefault = []string{".git", ".terradagger",
	"dist/**", "node_modules/**", ".cache"} // TODO: Add conditionals to include some of these directories.

var excludedFilesDefault = []string{".gitignore"}
