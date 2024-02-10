package terraformcore

var (
	tfDefaultCacheDir            = ".terraform"
	tfDefaultStateFileName       = "terraform.tfstate"
	tfDefaultStateBackupFileName = "terraform.tfstate.backup"
	tfDefaultLockFileName        = ".terraform.lock.hcl"
)

type TerraformConfig interface {
	GetCacheDir() string
	GetStateFileName() string
	GetStateBackupFileName() string
	GetLockFileName() string
}

type tfConfig struct{}

func (t *tfConfig) GetCacheDir() string {
	return tfDefaultCacheDir
}

func (t *tfConfig) GetStateFileName() string {
	return tfDefaultStateFileName
}

func (t *tfConfig) GetStateBackupFileName() string {
	return tfDefaultStateBackupFileName
}

func (t *tfConfig) GetLockFileName() string {
	return tfDefaultLockFileName
}
