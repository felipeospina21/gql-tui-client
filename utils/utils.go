package utils

import (
	"fmt"
	"log"
	"os"
)

func IsStringEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}

	return false
}

func LogToFile(s string) {
	f, err := os.OpenFile("debug.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", s)); err != nil {
		log.Println(err)
	}
}

func OverwriteEnvVars(s string) {
	f, err := os.Create(".env")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("\n %s", s))
}
