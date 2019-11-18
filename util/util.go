package util

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func CapitalizeFirstAll(str string) string {
	split := strings.Split(str, " ")

	for i, s := range split {
		if len(s) > 0 {
			split[i] = strings.Join([]string{strings.ToUpper(string(s[0])), s[1:]}, "")
		}
	}

	return strings.Join(split, " ")
}

func GetUserHomeDirectory() (string, error) {
	var homeDir string

	user, err := user.Current()
	if err != nil {
		return homeDir, err
	}

	homeDir = user.HomeDir

	return homeDir, nil
}

func JoinPaths(paths ...string) string {
	sep := string(os.PathSeparator)

	return filepath.Clean(strings.Join(paths, sep))
}

func PadDigit(digit int, length int) string {
	var builder strings.Builder
	digitStr := strconv.Itoa(digit)

	for i := 0; i < length-len(digitStr); i++ {
		builder.WriteString("0")
	}

	builder.WriteString(fmt.Sprintf("%d", digit))

	return builder.String()
}

func EscapeSpaces(str string) string {
	return strings.Join(strings.Split(str, " "), "\\ ")
}

func InsertTemplateData(templateString string, data interface{}) (string, error) {
	var builder strings.Builder

	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"padDigit": PadDigit,
		"sep": func() string {
			return string(os.PathSeparator)
		},
		"home": GetUserHomeDirectory,
		"homePath": func(path string) (string, error) {
			homeDir, err := GetUserHomeDirectory()
			if err != nil {
				return "", err
			}

			return JoinPaths(homeDir, path), nil
		},
		"escapeSpaces": EscapeSpaces,
	}).Parse(templateString)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&builder, data)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func GetYesOrNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [Y/N]: ", prompt)
	text, _ := reader.ReadString('\n')

	switch strings.TrimSpace(text) {
	case "Y", "y", "1":
		return true
	default:
		return false
	}
}

func GetOption(prompt string, options []string) int {
	fmt.Println(prompt)

	for i, option := range options {
		fmt.Printf("%d. %s", i+1, option)
	}

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	choice, _ := strconv.Atoi(input)

	return choice - 1
}

func MoveFile(src string, dest string, prompt bool) error {
	var err error
	move := true

	if prompt {
		move = GetYesOrNo(fmt.Sprintf("Move file\n'%s'\nto\n'%s'?\n", src, dest))
	}

	if move {
		err = os.MkdirAll(filepath.Dir(dest), 0644)
		if err == nil {
			err = os.Rename(src, dest)
		}
	}

	return err
}

func CombineStringArrays(arrs ...[]string) []string {
	ret := make([]string, 0)

	for _, arr := range arrs {
		ret = append(ret, arr...)
	}

	return ret
}
