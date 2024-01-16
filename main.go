package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
)

func main() {
    var err error
    // Define command line flags
    binaryNamePtr := flag.String("name", "", "name of the tool/keyword to run your tool in the cli (required)")
    mainFilePathPtr := flag.String("source", "", "optional: path to the main go source code file (using github.com/m-c-frank as remote)")
    targetDirPtr := flag.String("destination", "", "optional: target directory (relative to home, defaults to $HOME/bin)")

    // Parse command line flags
    flag.Parse()

    // Check if the required arguments are provided
    if *binaryNamePtr == "" {
        fmt.Println("Error: Missing required arguments.")
        flag.Usage()
        return
    }


    // collect the args
    binaryName := *binaryNamePtr
    mainFilePath := *mainFilePathPtr
    targetDir, err := getTargetPath(targetDirPtr)

    if !checkPath(binaryName) || *mainFilePathPtr == ""{
        cloneTool(binaryName)
    }

    // Use default target directory if not provided
    if targetDir == "" {
        targetDir, err = makeDirRelativeToHome("")
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
    }

    targetPath := filepath.Join(targetDir, binaryName)

    err = buildAndSetupGoBinary(binaryName, mainFilePath, targetDir, targetPath)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Build and PATH update completed successfully.")
    fmt.Println("Restart your Terminal or source ~/.bashrc")
    return
}

func getTargetPath(targetDirPtr *string) (string, error) {
    var err error
    if *targetDirPtr == "" {
        homeDir, err:= os.UserHomeDir()
        if err != nil {
            fmt.Println("Error: ", err)
            return "", err
        }
        targetDir := filepath.Join(homeDir, "bin")
        return targetDir, err
    }
    targetDir := *targetDirPtr
    return targetDir, err
}

func cloneTool(name string) error {
    cmd := exec.Command("git", "clone", "https://github.com/m-c-frank/" + name)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
            fmt.Println("Error: Unable to add file to Git staging.")
            return err
    }
    return nil
}

func checkPath(path string) bool {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    } else if err != nil {
        fmt.Println("Error checking path:", err)
        return false
    }
    return true
}


func buildAndSetupGoBinary(binaryName string, mainFilePath string, targetDir string, targetPath string) error {
    tempWorkDir, err := os.MkdirTemp("", "tempworkdir")

    binaryPath, err := buildBinary(mainFilePath, tempWorkDir, binaryName)
    if err != nil {
        return err
    }

    err = os.Rename(binaryPath, targetPath)
    if err != nil {
        return err
    }

    err = ensurePathInUserEnv(targetDir)
    if err != nil {
        return err
    }

    return err
}

// ensures the directory exists and returns the full path
func makeDirRelativeToHome(dirName string) (string, error) {
    usr, err := user.Current()
    if err != nil {
        fmt.Println("Error:", err)
        return "", err
    }

    // Ensure the ~/tools directory exists
    destinationPath := filepath.Join(usr.HomeDir, dirName)
    err = os.MkdirAll(destinationPath, os.ModePerm)
    if err != nil {
        fmt.Println("Error:", err)

        return "", err
    }

    return destinationPath, err
}

// buildAndAddToPath compiles the Go program and adds the binary to ~/tools, then adds ~/tools to the PATH.
func buildBinary(mainFilePath string, workDir string, binaryName string) (string, error) {
    outputBinary := filepath.Join(workDir, binaryName)

    // Compile the program
    cmd := exec.Command("go", "build", "-o", outputBinary, mainFilePath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        return "", err
    }

    return outputBinary, err
}

// make sure a given path is in user path or add it
func ensurePathInUserEnv(pathToAdd string) error {
    usr, err := user.Current()
    if err != nil {
        return err
    }

    pathIsInEnvironment, err := isPathInEnvironment(pathToAdd)
    if err != nil {
        return err
    }

    // Add pathToAdd to the PATH if not already present
    if !pathIsInEnvironment {
        profilePath := filepath.Join(usr.HomeDir, ".bashrc")
        file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            return err
        }
        defer file.Close()
        _, err = file.WriteString(fmt.Sprintf("\nexport PATH=$PATH:%s\n", pathToAdd))
        if err != nil {
            return err
        }
    }
    return err
}

// isPathInEnvironment checks if the given path is already in the PATH environment variable.
func isPathInEnvironment(path string) (bool, error) {
    value := os.Getenv("PATH")
    if value == "" {
        return false, errors.New("Error: failed getting $PATH")
    }

    paths := filepath.SplitList(value)
    for _, p := range paths {
        if p == path {
            return true, nil
        }
    }

    return false, nil
}

