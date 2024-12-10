package rss

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	File            string
	ExpectedTitle   string
	ExpectedContent string
}

func TestParseUTF16LE(t *testing.T) {
	testCases := []testCase{
		testCase{
			File:            "testdata/utf16-le",
			ExpectedTitle:   "Vålerenga slo ut rival i NM",
			ExpectedContent: "Vålerenga gikk seirende ut av naboduellen mot Lyn med 5-2, men det så lenge ut til å bli en frustrerende aften for de blåkledde.",
		},
		testCase{
			File:            "testdata/utf8",
			ExpectedTitle:   "Launching Feeder Dashboard",
			ExpectedContent: `<figure data-orig-width="1296" data-orig-height="900" class="tmblr-full"><img src="https://78.media.tumblr.com/ccf077ef7a27cf6799bf1bb393f9749d/tumblr_inline_pa9o4eALT71qleij7_540.gif" alt="image" data-orig-width="1296" data-orig-height="900"/></figure><p>After half a year in development, and 4 months in Beta. We are proud to release Feeder Dashboard live. Together with a release offer of 50% off.</p><p><b>Create your own information terminal</b></p><p>Enable up to 10 columns, or decks, for specific feeds or folders, and watch as new posts come in real-time. Customize the layout in the way that suits your needs to create a unique dashboard.</p><figure data-orig-width="1818" data-orig-height="1458" class="tmblr-full"><img src="https://78.media.tumblr.com/76d04c7d74406f29119cf2c0658230b6/tumblr_inline_pa9o8s8AVY1qleij7_540.png" alt="image" data-orig-width="1818" data-orig-height="1458"/></figure><p><b>Made for professionals </b></p><p>Feeder Dashboard is targeted to professionals who need to get a good overview of large quantities of information. Are you a job seeker that need an overview of different job searches on site? Or an analyst who uses Feeder for keeping track of Google Alerts, stocks, news? Then Feeder Dashboard is the perfect addition to your daily workflow.</p><p><b>Launch Discount - 50% off!</b></p><p>For the release of Feeder Dashboard, we&rsquo;re offering a 50% discount on the first 3 months. If you choose to pay for a year upfront, you get a 50% discount for the full year.</p><figure data-orig-width="2216" data-orig-height="1742" class="tmblr-full"><img src="https://78.media.tumblr.com/927747fcf1cfb0f194ab3d44bb2f3a58/tumblr_inline_pa9o9b2Toc1qleij7_540.jpg" alt="image" data-orig-width="2216" data-orig-height="1742"/></figure><p><b>Developed with the Norwegian News Agency NTB</b> </p><p>Working together we developed a solution to consolidate their information sources into a dashboard interface. Helping them increase speed and agility in news uptake. Editors each have an instance of Feeder set up with feeds from their respective branch. When new items come in, through sharing options they distribute those items to reporters to use as sources.</p><p><b>Can Feeder help you?</b></p><p>We love developing custom solutions for business and industries. Contact us to discuss how Feeder can improve your company&rsquo;s workflow.</p><p><b>Thank you</b></p><p>We are so thankful for all the beta testers who’s feedback made Feeder Dashboard even better. Our goal is to make the best experience there can be.</p><p>You can read more about Feeder Dashboard on our website: <a href="https://feeder.co/dashboard">feeder.co/dashboard</a> or our launch blog post: <a href="https://feeder.co/blog">feeder.co/blog</a><br/></p>`,
		},
	}

	for _, testCase := range testCases {
		input, headers := parseTestCase(testCase.File)
		feed, err := Parse(input, ParseOptions{
			ResponseHeaders: headers,
		})
		if err != nil {
			t.Error("Should not error on parsing", err)
			continue
		}

		if feed.Items[0].Title.String() != testCase.ExpectedTitle {
			t.Error("Incorrect title")
		}
		if feed.Items[0].Summary != testCase.ExpectedContent {
			t.Errorf("Incorrect content")
		}
	}
}

func parseTestCase(fileName string) ([]byte, http.Header) {
	contents, _ := os.ReadFile(fileName + ".response")
	headers, _ := os.ReadFile(fileName + ".headers")

	responseHeaders := http.Header{}

	headerLines := strings.Split(string(headers), "\r\n")
	for _, headerLine := range headerLines {
		values := strings.Split(headerLine, ": ")
		if len(values) > 1 {
			responseHeaders.Add(values[0], values[1])
		}
	}

	return []byte(contents), responseHeaders
}

// Use this when developing more tests
func downloadURLAsTestCase(url string, fileName string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln("Could not download", url, "due to", err)
	}

	headerFile, _ := os.Create("testdata/" + fileName + ".headers")
	defer headerFile.Close()

	responseFile, _ := os.Create("testdata/" + fileName + ".response")
	defer responseFile.Close()

	headerWriter := bufio.NewWriter(headerFile)
	responseWriter := bufio.NewWriter(responseFile)

	err = response.Header.Write(headerWriter)
	if err != nil {
		log.Fatalln("Could not write headers", err)
	}

	headerWriter.Flush()

	bodyBytes, _ := io.ReadAll(response.Body)
	responseWriter.Write(bodyBytes)

	responseWriter.Flush()

	fmt.Println("Done!")
}
