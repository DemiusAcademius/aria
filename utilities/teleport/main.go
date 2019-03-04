package main

import (
	"os"
	"fmt"
	"log"
	"regexp"

	"github.com/fatih/color"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	executablePath := TeleportDir()
	homePath := UserHomeDir()
	fmt.Println("executable path: " + executablePath)
	fmt.Println("      home path: " + homePath)

	inventory := LoadInventory(executablePath + Separator + "inventory.yaml")

	var variables = map[string]Variable{
		"@HOME": Variable{Description: "User Home", Path: homePath},
		"@DESK": Variable{Description: "Desktop", Path: homePath + Separator + "Desktop"},
		"@DOCS": Variable{Description: "Documents", Path: homePath + Separator + "Documents"},
	}

	if len(os.Args) < 3 {
		usage(inventory, variables)
	}

	server := os.Args[1]
	source := os.Args[2]

	serverConfig, ok := inventory.Servers[server]
	if !ok {
		color.Red("not found server `%s` in inventory", server)
	}

	sshConfig := &ssh.ClientConfig{
		User: serverConfig.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(serverConfig.Password),
		},
	}

	Upload(&serverConfig, sshConfig, variables, inventory.Content, source, homePath, executablePath)

	color.Blue("OK")
}

// Connect create sftp session for server in serverConfig
func Connect(serverConfig *ServerConfig, sshConfig *ssh.ClientConfig) (connection *ssh.Client, sftpClient *sftp.Client) {
	url := serverConfig.URI + ":22"
	connection, err := ssh.Dial("tcp", url, sshConfig)
	if err != nil {
		log.Fatalf("Failed to dial to server `%s`: %s", url, err)
	}
	sftpClient, err = sftp.NewClient(connection)
	if err != nil {
		log.Fatalf("unable to start sftp subsystem: %V", url, err)
	}
	return connection, sftpClient
}

// Upload source file|dir from contentPath to server
func Upload(
	serverConfig *ServerConfig, sshConfig *ssh.ClientConfig,
	variables Dictionary,
	contentPath, source, homePath, executablePath string) {
	connection, sftpClient := Connect(serverConfig, sshConfig)
	defer connection.Close()
	defer sftpClient.Close()

	regex, _ := regexp.Compile(`@[a-zA-Z]+`)

	contentPath = expandTemplate(contentPath, regex, variables)

	if err := UploadExecutor(sftpClient, contentPath, source, "./"); err != nil {
		color.Red(fmt.Sprintf("%s", err))
	}
}
