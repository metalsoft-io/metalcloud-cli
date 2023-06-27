package configuration

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
)

var (
	Version string
	Date    string
	Commit  string
	BuiltBy string
)

func IsAdmin() bool {
	return os.Getenv("METALCLOUD_ADMIN") == "true"
}

func ReadInputFromPipe() ([]byte, error) {

	reader := bufio.NewReader(GetStdin())
	var content []byte

	for {
		input, err := reader.ReadByte()
		if err != nil && err == io.EOF {
			break
		}
		content = append(content, input)
	}

	return content, nil
}

func ReadInputFromFile(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ConsoleIOChannel represents an IO channel, typically stdin and stdout but could be anything
type ConsoleIOChannel struct {
	Stdin  io.Reader
	Stdout io.Writer
}

var consoleIOChannelInstance ConsoleIOChannel

var once sync.Once

// GetConsoleIOChannel returns the console channel singleton
func GetConsoleIOChannel() *ConsoleIOChannel {
	once.Do(func() {

		consoleIOChannelInstance = ConsoleIOChannel{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
		}
	})

	return &consoleIOChannelInstance
}

// GetStdout returns the configured output channel
func GetStdout() io.Writer {
	return GetConsoleIOChannel().Stdout
}

// GetStdin returns the configured input channel
func GetStdin() io.Reader {
	return GetConsoleIOChannel().Stdin
}

// SetConsoleIOChannel configures the stdin and stdout to be used by all io with
func SetConsoleIOChannel(in io.Reader, out io.Writer) {
	channel := GetConsoleIOChannel()
	channel.Stdin = in
	channel.Stdout = out
}

// GetUserEmail returns the API key's owner
func GetUserEmail() string {
	return os.Getenv("METALCLOUD_USER_EMAIL")
}
