package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Job      string   `json:"job"`
	Email    string   `json:"email"`
	Country  string   `json:"country"`
	Company  string   `json:"company"`
	Browsers []string `json:"browsers"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintln(out, "found users:")
	var (
		count     int
		email     string
		isAndroid bool
		isMSIE    bool
	)

	scanner := bufio.NewScanner(file)
	browsers := make(map[string]bool, 114)
	user := User{}

	for scanner.Scan() {
		err = user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
			} else {
				continue
			}

			browsers[browser] = true
		}

		if isAndroid && isMSIE {
			email = strings.Replace(user.Email, "@", " [at] ", -1)
			fmt.Fprintf(out, "[%d] %s <%s>\n", count, user.Name, email)
		}

		count++
	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(browsers))
}
