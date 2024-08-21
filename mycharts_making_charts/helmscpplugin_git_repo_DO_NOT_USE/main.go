package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

type Action int

const (
	Upload = iota
	Download
	Delete
	Init
)
const Protocol = "scp"

var key, chartPath string
var action Action
var helmBin string
var AllowedActions = []string{"init", "push", "delete"}

type URL struct {
	username string
	host     string
	port     string
	path     string
}
type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func detokenize(url string) (URL, error) {
	regex := `scp:\/\/(\w+)@(\d+\.\d+\.\d+\.\d+):?(\d+)?(.*)$`
	r := regexp.MustCompile(regex)
	if !r.MatchString(url) {
		return URL{}, errors.New("INVALID SCP URL")
	}
	m := r.FindAllStringSubmatch(url, -1)
	username := m[0][1]
	host := m[0][2]
	port := "22"
	if m[0][3] != "" {
		port = m[0][3]
	}
	remotePath := "/home/" + username + "/"
	if m[0][4] != "" {
		remotePath = m[0][4]
	}
	return URL{
		username: username,
		host:     host,
		port:     port,
		path:     remotePath,
	}, nil
}
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func getRepoURL(repo string) (string, error) {
	var repoURL string
	cmd := exec.Command(helmBin, "repo", "list", "-o", "json")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var repos []Repo
	err = json.Unmarshal(out.Bytes(), &repos)
	if err != nil {
		return "", fmt.Errorf("could not parse helm repo list output: %s", err)
	}
	for _, r := range repos {
		if r.Name == repo {
			repoURL = r.Url
			break
		}
	}
	return repoURL, nil
}
func initialize() (URL, error) {
	var url URL
	var err error
	if os.Getenv("SCP_KEY") != "" {
		key = os.Getenv("SCP_KEY")
	} else {
		key = fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	}
	if os.Getenv("HELM_BIN") != "" {
		helmBin = os.Getenv("HELM_BIN")
	} else {
		helmBin = "helm"
	}
	// if the plugin is used as a downloader plugin
	if len(os.Args) == 5 && !contains(AllowedActions, os.Args[1]) {
		url, err = detokenize(os.Args[4])
		if err != nil {
			return URL{}, errors.New("please make sure the URL is scp://username@host[:port]/path")
		}
		action = Download
	} else if os.Args[1] == "push" {
		url, err = detokenize(os.Args[3])
		if err != nil {
			return url, errors.New("please make sure the URL is scp://username@host[:port]/path")
		}
		chartPath = os.Args[2]
		action = Upload
	} else if os.Args[1] == "delete" {
		// We don't pass the URL when deleting charts
		action = Delete
	} else if os.Args[1] == "init" {
		url, err = detokenize(os.Args[2])
		if err != nil {
			return url, errors.New("please make sure the URL is scp://username@host[:port]/path")
		}
		action = Init
	} else {
		return URL{}, errors.New("incorrect arguments.\nUsage:\nhelmscp push /path/to/chart scp://username@hostname[:port]/path/to/remote\nOR\nhelmscp scp://username@hostname:port/path/to/chart")
	}
	return url, nil
}
func main() {
	url, err := initialize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if action == Upload {
		chartFile, err := Package(chartPath)
		if err != nil {
			log.Fatalf("Error while packaging the chart: %s", err)
			return
		}
		err = Scp(chartFile, url, Upload)
		if err != nil {
			log.Fatalf("Error while uploading the archive: %s", err)
			return
		}
		fmt.Printf("Success!\n")
	} else if action == Download {
		err = Scp("", url, Download)
		if err != nil {
			log.Fatalf("Error while downloading the asset: %s", err)
			return
		}
	} else if action == Delete {
		var version string
		mySet := flag.NewFlagSet("", flag.ExitOnError)
		mySet.StringVar(&version, "version", "", "Chart version")
		mySet.Parse(os.Args[3:])
		repoURL, err := getRepoURL(os.Args[len(os.Args)-1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		url, err = detokenize(repoURL)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		chartName := os.Args[2]
		err = delete(version, url, chartName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if action == Init {
		err = reindex(url)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Incorrect command usage")
		os.Exit(1)
	}
}
func Package(chartPath string) (string, error) {
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
func Scp(filename string, url URL, action Action) error {
	clientConfig, _ := auth.PrivateKey(url.username, key, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(url.host+":"+url.port, &clientConfig)
	err := client.Connect()
	if err != nil {
		log.Fatal("Couldn't establish a connection to the remote server ", err)
		return err
	}
	// Close client connection after the file has been copied
	defer client.Close()
	remoteFile := url.path
	baseFileName := filepath.Base(filename)
	if action == Upload {
		if !strings.HasSuffix(remoteFile, "/") {
			remoteFile = remoteFile + "/"
		}
		// Open a file
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Could not open %s: %s", filename, err)
			return err
		}
		defer f.Close()
		defer os.Remove(filename)
		fmt.Printf("Uploading %s to %s at %s@%s:%s\n", baseFileName, remoteFile, url.username, url.host, url.port)
		// Finaly, copy the file over
		// Usage: CopyFile(fileReader, remotePath, permission)
		err = client.CopyFile(f, remoteFile+baseFileName, "0644")
		if err != nil {
			return err
		}
		err = reindex(url)
		if err != nil {
			return err
		}
		fmt.Printf("Cleaning up\n")
		return nil
	} else if action == Download {
		// Must point to a file not a directory
		if strings.HasSuffix(remoteFile, "/") {
			return errors.New("remote path must be a file not a directory")
		}
		sshClient, err := ssh.Dial("tcp", url.host+":"+url.port, &clientConfig)
		if err != nil {
			return err
		}
		defer sshClient.Close()
		session, err := sshClient.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()
		if err := session.Run("stat " + remoteFile); err != nil {
			return fmt.Errorf("could not download %s", remoteFile)
		}
		err = client.CopyFromRemote(os.Stdout, remoteFile)
		if err != nil {
			return err
		}
		return nil
	} else if action == Delete {

	} else if action == Init {
		err = reindex(url)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
func delete(version string, url URL, chartName string) error {
	var err error
	clientConfig, _ := auth.PrivateKey(url.username, key, ssh.InsecureIgnoreHostKey())
	sshClient, err := ssh.Dial("tcp", url.host+":"+url.port, &clientConfig)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	remoteFile := fmt.Sprintf("%s/%s-%s.tgz", url.path, chartName, version)
	fmt.Printf("Deleting %s\n", remoteFile)
	if err := session.Run("rm -f " + remoteFile); err != nil {
		return fmt.Errorf("could not delete %s", remoteFile)
	}
	err = reindex(url)
	if err != nil {
		return err
	}
	return nil
}
func reindex(url URL) error {
	chartDir := url.path
	clientConfig, _ := auth.PrivateKey(url.username, key, ssh.InsecureIgnoreHostKey())
	sshClient, err := ssh.Dial("tcp", url.host+":"+url.port, &clientConfig)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("error while connecting to remote server: %s", err)
	}
	defer session.Close()
	fmt.Printf("Indexing %s\n", chartDir)
	if err := session.Run("helm repo index " + chartDir); err != nil {
		return fmt.Errorf("error while reindexing: %s", err)
	}
	return nil
}
