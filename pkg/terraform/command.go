package terraform

import "github.com/Excoriate/go-terradagger/pkg/terradagger"

type TFCommands interface {
	Init(args TFCommandArgs) error
	//Validate(args TFCommandArgs) error
	//Plan(args TFCommandArgs) error
	//Apply(args TFCommandArgs) error
	//Destroy(args TFCommandArgs) error
}

type TFCommandArgs interface {
	//GetArgs() []string
}

type TFCommand struct {
	tfCfg     TfConfig
	tfOptions TfGlobalOptions
	td        *terradagger.TD
}
