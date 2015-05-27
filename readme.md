#goslacksender

There are so many go projects for Slack I couldn't choose between them, so I put together something very simple for my own purposes. This project simply allows you to post messages to a Slack channel asynchronously.

[![GoDoc](https://godoc.org/github.com/philpearl/goslacksender?status.svg)](https://godoc.org/github.com/philpearl/goslacksender)

[![Build Status](https://travis-ci.org/philpearl/goslacksender.svg)](https://travis-ci.org/philpearl/goslacksender)

## Usage

```go
import (
	"testing"

	"github.com/philpearl/goslacksender"
	"github.com/philpearl/ut"
)

const SLACK_URL = "https://hooks.slack.com/services/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"

var Sender = goslacksender.New(SLACK_URL)


func MyInterestingRoutine() {
	
	Sender.Text("Something interesting just happened!")
}
```