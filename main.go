package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultArtifactsPath = "../../artifacts" //nolint:misspell
	defaultStacksPath    = "../../infrastructure/stacks"
)

var (
	artifactsPath string
	stacksPath    string
)

func main() {
	artifactsPathFlag := flag.String("artifacts_path", defaultArtifactsPath, //nolint:misspell
		"artifacts path is the equivalent of dists directory where each lambda name folder has a Go binary and a zip/rar/tar") //nolint:misspell
	stacksPathFlag := flag.String("stacks_path", defaultStacksPath,
		"stacks path is the XIATECH integrations folder where each stack consists of a lambda folder that has Go Lambda entrypoint")
	flag.Parse()

	artifactsPath, stacksPath = *artifactsPathFlag, *stacksPathFlag
	if strings.TrimSpace(artifactsPath) == "" || strings.TrimSpace(stacksPath) == "" {
		log.Println("[WARNING] no artifacts_path or no stacks_path is specified") //nolint:misspell
	}

	checkVersion()
	printASCII()
	makeArtifacts()
}

func checkVersion() {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalln("[ERROR]: ", err)
	}

	log.Println("[INFO]: ", string(out))
}

func printASCII() {
	log.Println(`
 _______  _______ __________________ _______ _________ _______  _______  _______ 
(  ___  )(  ____ )\__   __/\__   __/(  ____ \\__   __/(  ____ \(  ____ \(  ____ )
| (   ) || (    )|   ) (      ) (   | (    \/   ) (   | (    \/| (    \/| (    )|
| (___) || (____)|   | |      | |   | (__       | |   | |      | (__    | (____)|
|  ___  ||     __)   | |      | |   |  __)      | |   | |      |  __)   |     __)
| (   ) || (\ (      | |      | |   | (         | |   | |      | (      | (\ (   
| )   ( || ) \ \__   | |   ___) (___| )      ___) (___| (____/\| (____/\| ) \ \__
|/     \||/   \__/   )_(   \_______/|/       \_______/(_______/(_______/|/   \__/
	`)
}

func makeArtifacts() {
	absStacksPath, err := filepath.Abs(stacksPath)
	if err != nil {
		log.Fatalln("[ERROR] couldn't get absolute path of stacks")
	}
	absArtifactsPath, err := filepath.Abs(artifactsPath)
	if err != nil {
		log.Fatalln("[ERROR] couldn't get absolute path of artifacts") //nolint:misspell
	}
	err = filepath.WalkDir(absStacksPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("[WARNING] couldn't read file info for %v", d.Name())
			return err
		}

		if !d.IsDir() || d.Name() == "local" || d.Name() == "testdata" {
			return nil
		}

		pathValues := strings.Split(path, "/")
		for i := range pathValues {
			if pathValues[i] == "stacks" {
				if i+3 >= len(pathValues) {
					return nil
				}

				if pathValues[i+2] != "lambda" {
					return nil
				}
			}
		}
		lambdaName := pathValues[len(pathValues)-1]
		distPath := filepath.Join(absArtifactsPath, lambdaName)
		log.Println("[INFO] lambdaName:", lambdaName)
		log.Println("[INFO] path:", path)
		log.Println("[INFO] distPath:", distPath)

		err = buildBinary(path, distPath, lambdaName)
		if err != nil {
			log.Fatalln("[ERROR] ", err)
			return err
		}

		err = buildZip(distPath, lambdaName)
		if err != nil {
			log.Fatalln("[ERROR] ", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("[ERROR] failed to walk over directories under stacks path %v: %v \n", stacksPath, err)
	}
}

func buildBinary(currentPath, distPath, lambdaName string) error {
	err := os.MkdirAll(distPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to mkdir for path: %v: %v", distPath, err)
	}

	exe := filepath.Join(distPath, lambdaName)

	cmd := exec.Command("go", "build", "-o", exe, "main.go")
	cmd.Dir = currentPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		if len(stdout.String()) > 0 {
			log.Printf("[INFO] go build returns stdout %v\n", stdout.String())
		}
		if len(stderr.String()) > 0 {
			log.Printf("[WARNING] go build returns stderr %v\n", stderr.String())
		}
		return fmt.Errorf("failed to build for path %v: %v", distPath, err)
	}

	err = os.Chmod(distPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to chmod for path %v: %v", distPath, err)
	}
	return nil
}

func buildZip(distPath, lambdaName string) error {
	exe := filepath.Join(distPath, lambdaName)
	archiveFile, err := os.Create(fmt.Sprintf("%v.zip", exe))
	if err != nil {
		return fmt.Errorf("failed to create archive file for %v lambda: %v", lambdaName, err)
	}
	defer archiveFile.Close()
	w := zip.NewWriter(archiveFile)
	defer w.Close()

	execFile, err := os.Open(exe)
	if err != nil {
		return fmt.Errorf("failed to open exec file %v: %v", exe, err)
	}
	defer execFile.Close()

	writer, err := w.Create(lambdaName)
	if err != nil {
		return fmt.Errorf("failed to create space inside zip: %v", err)
	}
	if _, err = io.Copy(writer, execFile); err != nil {
		return fmt.Errorf("failed to put the executable file into the zip: %v", err)
	}
	return nil
}
