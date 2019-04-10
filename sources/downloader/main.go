package main

import (
	"io"
	"path/filepath"
	"path"
	"os"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	executablePath := DownloaderDir()
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

	if err := Download(&serverConfig, sshConfig, source); err != nil {
		color.Red(fmt.Sprintf("%s", err))
		os.Exit(-1)
	}

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
		log.Fatalf("unable to start sftp subsystem: %v", err)
	}
	return connection, sftpClient
}

// Download source file from server to current dir
func Download(serverConfig *ServerConfig, sshConfig *ssh.ClientConfig,sourceFile string) error {
	connection, sftpClient := Connect(serverConfig, sshConfig)
	defer connection.Close()
	defer sftpClient.Close()

	srcF, err := sftpClient.OpenFile(sourceFile, os.O_RDONLY)
	if err != nil {
		return fmt.Errorf("can not open source file %s: %v", sourceFile, err)
	}

	destinationPath := path.Join(WorkingDir(), filepath.Base(sourceFile))
	dstF, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("can not create destination file %s: %v", destinationPath, err)
	}

	written, err := io.Copy(dstF, srcF)
	if err != nil {
		return fmt.Errorf("can not copy source file: %v", err)
	}

	c := color.New(color.FgMagenta)
	c.Print("  transmitted: ")
	color.Magenta("%v bytes", written)

	return nil
}
