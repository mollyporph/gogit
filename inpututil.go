package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", s)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
func getMfa() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("github 2step auth: ")
	mfa, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(mfa)
}
func getString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(input)
}
func getInt(prompt string) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	number, parseErr := strconv.Atoi(strings.TrimSpace(input))
	if parseErr != nil {
		log.Fatal(err)
	}
	return number
}
func getUsername() string {
	return getString("github username: ")
}
func getPassword() string {
	fd := os.Stdin.Fd()
	oldState, err := terminal.MakeRaw(int(fd))
	if err != nil {
		panic(err)
	}
	fmt.Printf("github password: ")
	password, err := terminal.ReadPassword(int(fd))
	terminal.Restore(int(fd), oldState)
	fmt.Println()
	return string(password)
}
