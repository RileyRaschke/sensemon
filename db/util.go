package db

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const MAX_OBJECT_LENGTH int = 28

func passwordFromCommand(cmd string) (string, error) {
	res, err := exec.Command(cmd).Output()
	r := string(res)
	return strings.TrimSpace(r), err
}

func passwordFromShell() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	r := string(bytePassword)
	return strings.TrimSpace(r), err
}

func ValidName(t string) bool {
	tValid := regexp.MustCompile(`^[a-zA-Z0-9_]{1,30}$`)
	if tValid.MatchString(t) {
		return true
	}
	return false
}

func CleanName(name string) string {
	strVal := strings.Trim(name, " \n\t")
	if strVal == " " || strVal == "" {
		strVal = "S_"
	}
	if strVal[0] >= '0' && strVal[0] <= '9' {
		strVal = "N_" + strVal
	}
	if strings.ToUpper(strVal) == "DATE" {
		return "C_DATE"
	}
	strVal = strings.Replace(strVal, "'", "", -1)
	strVal = strings.Replace(strVal, "_#_", "num", -1)
	strVal = strings.Replace(strVal, "(", "", -1)
	strVal = strings.Replace(strVal, ")", "", -1)
	strVal = strings.Replace(strVal, "`", "_", -1)
	rmSpace := regexp.MustCompile(`((\s)|(-)|(–)|(_)|(\+)|(\*)|(\\)|(/)|(\|)|(#)|(%))+`)
	strVal = rmSpace.ReplaceAllString(
		strings.Replace(strVal, "'", "", -1), "_")

	if len(strVal) > MAX_OBJECT_LENGTH {
		strVal = string(strVal)[:MAX_OBJECT_LENGTH]
	}
	return strings.Trim(strings.ToUpper(strVal), "_")
}
