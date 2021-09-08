package tests

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPDCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PDCLI Test Suite")
}

var executable string

var _ = BeforeSuite(func() {
	executableName := "pdcli"

	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}

	executable = filepath.Join("..", executableName)

	// Create the PDCLI binary in the project root directory
	cmd := exec.Command("go", "build", "-o", executable, "../cmd/pdcli")
	err := cmd.Run()

	Expect(err).ToNot(HaveOccurred(), "Error creating binary file for test suite", executable)
})

type TestCommand struct {
	args   []string
	config string
	env    map[string]string
	input  []byte
}

// TestResult contains the result of the CLI command executed.
type TestResult struct {
	configFile string
	configData []byte
	out        []byte
	err        []byte
	exitCode   int
}

// NewCommand creates a new CLI command runner.
func NewCommand() *TestCommand {
	return &TestCommand{
		env: map[string]string{},
	}
}

// GetStdIn gets the standard input for the CLI command.
func (tc *TestCommand) GetStdIn(value string) *TestCommand {
	tc.input = []byte(value)
	return tc
}

// Env sets an environment variable to the CLI command.
func (tc *TestCommand) Env(name, value string) *TestCommand {
	tc.env[name] = value
	return tc
}

// Args adds a set of arguments to the command.
func (tc *TestCommand) Args(values ...string) *TestCommand {
	tc.args = append(tc.args, values...)
	return tc
}

// ConfigFile returns the name of the temporary configuration file.
// It will return an empty string if the confif file doesn't exist.
func (tr *TestResult) ConfigFile() string {
	return tr.configFile
}

// ConfigString returns the content of the configuration file.
func (tr *TestResult) ConfigString() string {
	return string(tr.configData)
}

// Err returns the standard errour output of the test command.
func (tr *TestResult) ErrString() string {
	return string(tr.err)
}

// ExitCode returns the exit code of the test command.
func (tr *TestResult) ExitCode() int {
	return tr.exitCode
}

func (tc *TestCommand) Run() *TestResult {
	var err error

	// Create a temporary config directory
	tmpDir, err := ioutil.TempDir("", "pdcli-test-*.d")

	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	// Delete the temp dir after test suite is finished
	defer func() {
		err = os.RemoveAll(tmpDir)
		ExpectWithOffset(1, err).ToNot(HaveOccurred())
	}()

	// Create a temporary configuration file
	configFile := filepath.Join(tmpDir, "config.json")

	if tc.config != "" {
		err = ioutil.WriteFile(configFile, []byte(tc.config), 0600)
		ExpectWithOffset(1, err).ToNot(HaveOccurred())
	}

	// Parse the current environment into a map
	envMap := map[string]string{}
	for _, text := range os.Environ() {
		index := strings.Index(text, "=")
		var name string
		var value string
		if index > 0 {
			name = text[0:index]
			value = text[index+1:]
		} else {
			name = text
			value = ""
		}
		envMap[name] = value
	}

	// Add the environment variables
	for name, value := range tc.env {
		envMap[name] = value
	}

	// Add to the environment the variable that points to a configuration file
	envMap["PDCLI_CONFIG"] = configFile

	// Reconstruct the environment list
	envList := make([]string, 0, len(envMap))

	for name, value := range envMap {
		envList = append(envList, name+"="+value)
	}

	// Create the buffers
	inBuf := &bytes.Buffer{}

	// if standard input is provided
	if tc.input != nil {
		inBuf.Write(tc.input)
	}

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	// Create the pdcli command
	cmd := exec.Command(executable, tc.args...)
	cmd.Env = envList
	cmd.Stdin = inBuf
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	// Run the command
	err = cmd.Run()

	// Check if the configuration file exists
	_, err = os.Stat(configFile)

	if errors.Is(err, os.ErrNotExist) {
		configFile = ""
	} else if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}

	var configData []byte

	if configFile != "" {
		configData, err = ioutil.ReadFile(configFile)
		Expect(err).ToNot(HaveOccurred())
	}

	// The result rendered from the test command
	result := &TestResult{
		configFile: configFile,
		configData: configData,
		out:        outBuf.Bytes(),
		err:        errBuf.Bytes(),
		exitCode:   cmd.ProcessState.ExitCode(),
	}

	return result
}

var _ = AfterSuite(func() {
	err := os.Remove("../pdcli")
	Expect(err).ToNot(HaveOccurred())
})
