package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=go"

type jobInfo struct {
	id       string
	title    string
	company  string
	location string
	summary  string
}

func (job jobInfo) getItems() []string {
	return []string{
		"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.company, job.location, job.summary,
	}
}

func main() {
	totalPages := getPages()
	var totalJobs []jobInfo
	c := make(chan []jobInfo)

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		totalJobs = append(totalJobs, <-c...)
	}

	writeToCSV(totalJobs)
}

func getPage(page int, c chan<- []jobInfo) {

	var jobs []jobInfo
	pageUrl := baseURL + "&start=" + strconv.Itoa(page*10)
	fmt.Println("Checking ", pageUrl)

	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, selection *goquery.Selection) {
		job := parsePage(selection)
		jobs = append(jobs, job)
	})

	c <- jobs
}

func getPages() int {

	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	fmt.Println(doc)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func parsePage(selection *goquery.Selection) jobInfo {
	id, _ := selection.Attr("data-jk")
	title := selection.Find(".title>a").Text()
	company := selection.Find(".company>a").Text()
	location := selection.Find(".location").Text()
	summary := selection.Find(".summary").Text()

	return jobInfo{
		id:       id,
		title:    title,
		company:  company,
		location: location,
		summary:  summary,
	}
}

func writeToCSV(jobs []jobInfo) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"Link", "Title", "Company", "Location", "Summary"})

	for _, job := range jobs {
		w.Write(job.getItems())
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request Failed in ", res.StatusCode)
	}
}
