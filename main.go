package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/88250/gulu"
	"github.com/parnurzeal/gorequest"
)

var logger = gulu.Log.NewLogger(os.Stdout)


const (
	githubUserName = "Achuan-2"
	hacpaiUserName = "Achuan-2"
)

func main() {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://ld246.com/api/v2/user/"+hacpaiUserName+"/events?size=8").Timeout(7*time.Second).
		Set("User-Agent", "Profile Bot; +https://github.com/"+githubUserName+"/"+githubUserName).EndStruct(&result)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Fatalf("fetch events failed: %+v, %s", errors, data)
	}
	if 0 != result["code"].(float64) {
		logger.Fatalf("fetch events failed: %s", data)
	}

	buf := &bytes.Buffer{}
	buf.WriteString("\n\n")
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	updated := time.Now().In(cstSh).Format("2006-01-02 15:04:05")
	buf.WriteString("### Recent updates in Liandi \n\n Last Update TimeοΌ`" + updated + "`\n\nπ Posts &nbsp; π¬ Comments &nbsp; π£ Replies &nbsp; π Gossip &nbsp; β­οΈ Follow &nbsp; π Like &nbsp; π Thank &nbsp; π° Reward &nbsp; π Collection\n\n")
	for _, event := range result["data"].([]interface{}) {
		evt := event.(map[string]interface{})
		operation := evt["operation"].(string)
		title := evt["title"].(string)
		typ := evt["type"].(string)
		var emoji string
		switch typ {
		case "article":
			emoji = "π"
		case "comment":
			emoji = "π¬"
		case "comment2":
			emoji = "π£"
		case "breezemoon":
			emoji = "π"
			title = operation
		case "vote-article":
			emoji = "ππ"
		case "vote-comment":
			emoji = "ππ¬"
		case "vote-comment2":
			emoji = "ππ£"
		case "vote-breezemoon":
			emoji = "ππ"
			title = operation
		case "reward-article":
			emoji = "π°π"
		case "thank-article":
			emoji = "ππ"
		case "thank-comment":
			emoji = "ππ¬"
		case "accept-comment":
			emoji = "βπ¬"
		case "thank-comment2":
			emoji = "ππ£"
		case "thank-breezemoon":
			emoji = "ππ"
			title = operation
		case "follow-user":
			emoji = "β­οΈπ¨βπ»"
		case "follow-tag":
			emoji = "β­οΈπ·οΈ"
		case "collect-article":
			emoji = "ππ"
		}

		url := evt["url"].(string)
		content := evt["content"].(string)
		buf.WriteString("* " + emoji + " [" + title + "](" + url + ")\n\n" + "  > " + content + "\n")
	}
	buf.WriteString("\n\n")

	fmt.Println(buf.String())

	readme, err := ioutil.ReadFile("README.md")
	if nil != err {
		logger.Fatalf("read README.md failed: %s", data)
	}

	startFlag := []byte("<!--events start -->")
	beforeStart := readme[:bytes.Index(readme, startFlag)+len(startFlag)]
	newBeforeStart := make([]byte, len(beforeStart))
	copy(newBeforeStart, beforeStart)
	endFlag := []byte("<!--events end -->")
	afterEnd := readme[bytes.Index(readme, endFlag):]
	newAfterEnd := make([]byte, len(afterEnd))
	copy(newAfterEnd, afterEnd)
	newReadme := append(newBeforeStart, buf.Bytes()...)
	newReadme = append(newReadme, newAfterEnd...)
	if err := ioutil.WriteFile("README.md", newReadme, 0644); nil != err {
		logger.Fatalf("write README.md failed: %s", data)
	}
}
