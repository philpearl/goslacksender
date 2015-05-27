// goslacksender is a very simple API to send messages to a Slack channel.
// Messages are queued to a background go routine so the call does not block.
package goslacksender

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	// The size of the queue to the background goroutine
	channel_size int = 100
)

// Message encodes a slack message. Use it with Sender.Queue(). If you just want to
// send a text message use Sender.Text()
type Message struct {
	Text      string `json:"text"`
	Username  string `json:"username,omitempty"`
	IconUrl   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Channel   string `json:"channel,omitempty"`
}

// Sender represents an instance of a slack sender. Create one via New()
type Sender interface {
	// Queue() queues a message to be send to Slack in the background
	Queue(msg Message)
	// Text() sends a simple text message to the default Slack channel
	Text(text string)
	// Close() closes the sender and waits for any queued messages to be flushed out
	Close()
}

type senderImpl struct {
	url     string       // The url to post events too, including project details
	channel chan Message // For queuing events to the background
	done    chan bool    // For clean exiting
}

/*
New() creates a new Sender.

This creates a background goroutine to aggregate and send your events.

Get the url to pass in by going to
https://<teamname>.slack.com/services/new/incoming-webhook and following the
instructions.
*/
func New(url string) Sender {
	sender := &senderImpl{
		url:     url,
		channel: make(chan Message, channel_size),
		done:    make(chan bool),
	}
	go sender.run()
	return sender
}

/*
Queue events to be sent to Slack
*/
func (sender *senderImpl) Queue(msg Message) {
	sender.channel <- msg
}

// Text() sends a text message to Slack
func (sender *senderImpl) Text(text string) {
	sender.Queue(Message{Text: text})
}

/*
Close the sender and wait for queued events to be sent
*/
func (sender *senderImpl) Close() {
	// Closing the channel signals the background thread to exit
	close(sender.channel)
	// Wait for the background thread to signal it has flushed all events and exited
	<-sender.done
}

func (sender *senderImpl) send(msg Message) {
	// Convert data to JSON
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(msg)

	rsp, err := http.Post(sender.url, "application/json", &buf)
	if err != nil {
		log.Printf("goslacksender: Failed to post msg to Slack.  %v\n", err)
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		log.Printf("goslacksender: Failure return for slack msg post.  %d, %s\n", rsp.StatusCode, rsp.Status)
	}
}

func (sender *senderImpl) run() {
	for msg := range sender.channel {
		sender.send(msg)
	}

	// Indicate that this thread is over
	sender.done <- true
	log.Printf("goslacksender: Slack sender exited\n")
}
