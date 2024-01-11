package terradagger

// type DataTransfer struct {
// 	SourcePath         string
// 	SourcePathAbs      string
// 	DestinationPath    string
// 	DestinationPathAbs string
// 	IsDir              bool
// }

type DataTransferExcluded struct {
	Files       []string
	Directories []string
}

type TransferToContainer struct {
	SourcePathInHostAbs        string
	DestinationPathInContainer string
}

type TransferToHost struct {
	SourcePathInContainer    string
	DestinationPathInHostAbs string
	BackupPathInHostAbs      string
}

type DataTransferToContainer struct {
	WorkDirPath string
	Files       []TransferToContainer
	Dirs        []TransferToContainer
}

type DataTransferToHost struct {
	WorkDirPath string
	Files       []TransferToHost
	Dirs        []TransferToHost
}

type DataBackupInHost struct {
	BackupPathInHostAbs string
	Files               []*DataBackupContent
	Dirs                []*DataBackupContent
	ID                  string
}

type DataBackupContent struct {
	SourcePathAbs      string
	DestinationPathAbs string
}

// type DataTransferToHost struct {
// 	WorkDirPath string
// 	Files       []string
// 	Dirs        []string
// }
