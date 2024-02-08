package terraform

var (
	tfDefaultCacheDir            = ".terraform"
	tfDefaultStateFileName       = "terraform.tfstate"
	tfDefaultStateBackupFileName = "terraform.tfstate.backup"
	tfDefaultLockFileName        = ".terraform.lock.hcl"
)

type TfConfig interface {
	GetCacheDir() string
	GetStateFileName() string
	GetStateBackupFileName() string
	GetLockFileName() string
}

type tfCfg struct{}

func (t *tfCfg) GetCacheDir() string {
	return tfDefaultCacheDir
}

func (t *tfCfg) GetStateFileName() string {
	return tfDefaultStateFileName
}

func (t *tfCfg) GetStateBackupFileName() string {
	return tfDefaultStateBackupFileName
}

func (t *tfCfg) GetLockFileName() string {
	return tfDefaultLockFileName
}
