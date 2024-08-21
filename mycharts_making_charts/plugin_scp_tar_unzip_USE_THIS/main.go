package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

var username, key, port, chartPath, remotePath, host string
var helmBin string

func initialize() {
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.StringVar(&username, "u", "", "The remote server username")
	flag.StringVar(&key, "k", os.Getenv("HOME")+"/.ssh/id_rsa", "The SSH key")
	flag.StringVar(&port, "p", "22", "The remote server port")
	flag.StringVar(&remotePath, "r", "", "Path to the remote directory")
	flag.StringVar(&chartPath, "l", "", "Path to the chart directory")
	flag.StringVar(&host, "s", "", "The hostname or IP address")
	flag.Parse()
	if host == "" {
		flag.PrintDefaults()
		fmt.Println("Please supply the hostname or IP address to the remote host")
		os.Exit(2)
	}
	if username == "" {
		flag.PrintDefaults()
		fmt.Println("Please provide the username to connect to the remote host over SSH")
		os.Exit(2)
	}
	if remotePath == "" {
		flag.PrintDefaults()
		fmt.Println("Please provide the remote path to save the file")
		os.Exit(2)
	}
	if chartPath == "" {
		flag.PrintDefaults()
		fmt.Println("Please provide the path to the Helm chart")
		os.Exit(2)
	}
}
func main() {
	initialize()
	chartFile, err := Package(chartPath)
	if err != nil {
		log.Fatalf("Error while packaging the chart: %s", err)
		return
	}
	err = Upload(chartFile)
	if err != nil {
		log.Fatalf("Error while uploading the archive: %s", err)
		return
	}
	fmt.Printf("Success!\n")
}
func Package(chartPath string) (string, error) {
	if os.Getenv("HELM_BIN") != "" {
		helmBin = os.Getenv("HELM_BIN")
	} else {
		helmBin = "helm"
	}
	fmt.Printf("Packaging chart from %s\n", chartPath)
	cmd := exec.Command(helmBin, "package", chartPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	chartNameFullPath := strings.Split(out.String(), ":")[1]
	chartNameFullPath = strings.Trim(chartNameFullPath, "\n")
	chartNameFullPath = strings.Trim(chartNameFullPath, " ")
	return chartNameFullPath, nil
}
func Upload(filename string) error {
	if remotePath == "" {
		remotePath = fmt.Sprintf("/home/%s/", username)
	}
	if remotePath[len(remotePath)-1:] != "/" {
		remotePath = remotePath + "/"
	}
	clientConfig, _ := auth.PrivateKey(username, key, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(host+":"+port, &clientConfig)
	err := client.Connect()
	if err != nil {
		log.Fatal("Couldn't establish a connection to the remote server ", err)
		return err
	}
	// Open a file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open filename to upload: %s", err)
		return err
	}
	// Close client connection after the file has been copied
	defer client.Close()
	// Close the file after it has been copied
	defer f.Close()
	defer os.Remove(filename)
	baseFileName := filepath.Base(filename)
	fmt.Printf("Uploading %s to %s at %s@%s:%s\n", baseFileName, remotePath, username, host, port)
	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)
	err = client.CopyFile(f, remotePath+baseFileName, "0644")
	if err != nil {
		log.Fatalf("Could not upload the file to the remote server: %s", err)
		return err
	}
	fmt.Printf("Cleaning up\n")
	return nil
}
