package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("Request is Failed")

func main() {

	results := map[string]string{}

	c := make(chan error)

	urls := []string{
		"https://www.ruliweb.com/",
		"https://www.naver.com/",
		"https://www.amazon.co.jp/",
		"https://nomadcoders.co/go-for-beginners/lectures/1524",
		"https://www.hahwul.com/2019/11/18/how-to-fix-xcrun-error-after-macos-update/",
		"https://watcha.com/",
	}

	for _, url := range urls {
		go hitUrl(url, c)
	}

	for i := 0; i < len(urls); i++ {

		if <-c != nil {
			results[urls[i]] = "Failed"
		} else {
			results[urls[i]] = "OK"
		}

		fmt.Println(urls[i], results[urls[i]])
	}
}

func hitUrl(url string, c chan error) {

	fmt.Println("Check URL : ", url)

	req, err := http.Get(url)

	if err != nil || req.StatusCode >= 400 {
		c <- errRequestFailed
	}

	c <- nil
}
