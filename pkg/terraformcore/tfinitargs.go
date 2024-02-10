package terraformcore

type InitArgsOptions struct {
	// NoColor is a flag to disable colors in terraform output
	NoColor bool
	// BackendConfigFile is the path to the backend config file
	BackendConfigFile string
	// Upgrade is a flag to upgrade the modules and plugins
	Upgrade bool
}

type InitArgs interface {
	GetArgNoColor() []string
	GetArgBackendConfigFile() []string
	GetArgUpgrade() []string
}

func (ti *InitArgsOptions) GetArgNoColor() []string {
	arg := []string{"-no-color"}
	if ti.NoColor {
		return arg
	}

	return []string{}
}

func (ti *InitArgsOptions) GetArgBackendConfigFile() []string {
	arg := []string{"-backend-config", ti.BackendConfigFile}
	if ti.BackendConfigFile != "" {
		return arg
	}

	return []string{}
}

func (ti *InitArgsOptions) GetArgUpgrade() []string {
	arg := []string{"-upgrade"}
	if ti.Upgrade {
		return arg
	}

	return []string{}
}
