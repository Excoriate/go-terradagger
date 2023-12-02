package config

const MountPathPrefixInDagger = "/mnt"

var ExcludedDirsDefault = []string{".git/", ".gitignore", ".terra-dagger/",
	"dist/**", "node_modules/**"} // TODO: Add conditionals to include some of these directories.

var ExcludedDirsTerraform = []string{".terraform/**", ".terraform-lock.hcl", ".terraform-cache/**"}
