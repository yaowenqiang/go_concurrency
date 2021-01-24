package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
	"gopkg.in/gomail.v2"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

var stories []Story

var m *gomail.Message
var d *gomail.Dialer

func init() {
	hackerNewsClient = gophernews.NewClient()
	var err error

	redditSession, err = geddit.NewLoginSession("user", "pass", "gopher")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m = gomail.NewMessage()
	m.SetHeader("From", "user")
	m.SetHeader("To", "user")
	d = gomail.NewDialer("smtp.163.com", 465, "user", "123456")
}

func getHnStoryDetails(id int, c chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()
	story, err := hackerNewsClient.GetStory(id)
	if err != nil {
		return
	}
	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "hackerNews",
	}

	c <- newStory
}

//Story is the story struct of hacker news and reddit
type Story struct {
	title  string
	url    string
	author string
	source string
}

func newHnStories(c chan<- Story) {
	defer close(c)
	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, id := range changes.Items {
		wg.Add(1)
		go getHnStoryDetails(id, c, &wg)
	}
	wg.Wait()
}

func newRedditStories(c chan<- Story) {
	defer close(c)
	sort := geddit.PopularitySort(geddit.NewSubmissions)
	var listingOptions geddit.ListingOptions
	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range submissions {
		newStory := Story{
			title:  s.Title,
			url:    s.URL,
			author: s.Author,
			source: "Reddit /r/programming",
		}
		c <- newStory
	}
}

func outputToConsole(c <-chan Story) {
	for {
		s := <-c
		fmt.Printf("%s :%s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func outputToFile(c <-chan Story, file *os.File) {
	for {
		s := <-c
		fmt.Fprintf(file, "%s :%s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}
func main() {
	go func() {
		for {
			fmt.Println("Fetching new stories...")
			fromHn := make(chan Story, 8)
			fromReddit := make(chan Story, 8)
			toList := make(chan Story, 8)
			toFile := make(chan Story, 8)
			toConsole := make(chan Story, 8)
			toEmail := make(chan Story, 8)
			go newHnStories(fromHn)
			go newRedditStories(fromReddit)
			file, err := os.Create("stories.txt")

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			go outputToConsole(toConsole)
			go outputToFile(toFile, file)
			go sendMail(toEmail)

			hnOpen := true
			redditOpen := true

			for hnOpen || redditOpen {
				select {
				case story, open := <-fromHn:
					if open {
						toFile <- story
						toConsole <- story
						toEmail <- story
						toList <- story
					} else {
						hnOpen = false
					}
				case story, open := <-fromReddit:
					if open {
						toFile <- story
						toConsole <- story
						toList <- story
						toEmail <- story
					} else {
						redditOpen = false
					}
				}
			}
			fmt.Println("Done fetching new stories...")
			time.Sleep(30 * time.Second)
		}
	}()

	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	fmt.Println("start Listening and serving on 0.0.0.0:8881")
	if err := http.ListenAndServe(":8881", nil); err != nil {
		panic(err)
	}
}

func sendMail(c <-chan Story) {
	for {
		s := <-c
		// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", s.title)
		m.SetBody("text/html", "<a href='"+s.url+"'>"+s.title+"</a>")
		// m.Attach("/home/Alex/lolcat.jpg")
		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}
}

func outputStories(c <-chan Story) {
	for {
		s := <-c
		stories = append(stories, s)
	}
}

func topTen(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><bodh>"))

	form := "<form action='/search' method='get'>Search: <input type='text' name='q' > <input type='submit'></form>"
	w.Write([]byte(form))
	for i := len(stories) - 1; i >= 0 && len(stories)-i < 10; i-- {
		story := stories[i]
		w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
	}
	w.Write([]byte("</body></html>"))

}

func search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "Search parameter q is required to search.", http.StatusNotAcceptable)
		return
	}

	w.Write([]byte("<html><bodh>"))
	s := searchStories(query)
	if len(s) == 0 {
		w.Write([]byte(fmt.Sprintf("No results for query %s\n<br>", query)))
	} else {
		for _, story := range s {
			w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a><br>by %s on %s<br><br>", story.url, story.title, story.author, story.source)))
		}
	}
	w.Write([]byte("<a href='../'>Back</a>"))
	w.Write([]byte("</body></html>"))
}

func searchStories(query string) []Story {
	var foundstories []Story
	for _, story := range stories {
		if strings.Contains(strings.ToUpper(story.title), strings.ToUpper(query)) {
			foundstories = append(foundstories, story)
		}
	}
	return foundstories
}
