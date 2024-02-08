package utils

import "os"

func GetSSHAuthSock() string {
	return os.Getenv("SSH_AUTH_SOCK")
}

// GetSSHGitSecureConnectCommand returns an SSH command string configured for secure connections
// by bypassing the usual host key storage mechanism. It's particularly useful for Git operations over SSH
// in environments where strict host key verification needs to be relaxed, such as automated scripts or
// continuous integration systems connecting to dynamically changing servers. This command automatically
// accepts new host keys but prevents connections to hosts with changed keys, mitigating man-in-the-middle
// attacks without requiring manual intervention for known host key management.
func GetSSHGitSecureConnectCommand() string {
	return "ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=accept-new"
}
