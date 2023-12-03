package commands

type TerraDaggerArgs []string
type DaggerEngineCMDs [][][]string

type TerraDaggerCMD struct {
	Binary                  string
	Command                 string
	Args                    TerraDaggerArgs
	OmitBinaryNameInCommand bool
}
