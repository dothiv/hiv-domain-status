package hivdomainstatus

import (
	_ "crypto/sha512"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var CLICKCOUNTER_SCRIPT = "//dothiv-registry.appspot.com/static/clickcounter.min.js"

type DomainCheckResult struct {
	Domain         string
	DnsOk          bool
	Addresses      []string
	URL            *url.URL
	bodyFile       string
	body           []byte
	StatusCode     int
	ScriptPresent  bool
	SaveBody       bool
	IframePresent  bool
	IframeTarget   string
	IframeTargetOk bool
	Valid          bool
}

func NewDomainCheckResult(domain string) (checkResult *DomainCheckResult) {
	checkResult = new(DomainCheckResult)
	checkResult.Domain = domain
	checkResult.URL, _ = url.Parse("http://www." + checkResult.Domain + "/")
	return
}

func (checkResult *DomainCheckResult) IsHivDomain() bool {
	return strings.ToLower(checkResult.Domain[len(checkResult.Domain)-3:]) == "hiv"
}

func (checkResult *DomainCheckResult) IsWWW() bool {
	var wwwMatch = regexp.MustCompile(`^https?://www\.+`)
	return wwwMatch.MatchString(checkResult.URL.String())
}

func (checkResult *DomainCheckResult) Check() (err error) {
	checkResult.Valid = true
	if checkResult.IsHivDomain() {
		err = checkResult.dnsCheck()
		if err != nil {
			checkResult.Valid = false
			return
		}
	}
	err = checkResult.fetch()
	if err != nil {
		if checkResult.IsWWW() {
			// Fetch page without www (if not present)
			checkResult.URL, _ = url.Parse("http://" + checkResult.Domain + "/")
			err = checkResult.Check()
			return
		}
		checkResult.Valid = false
		return
	}
	if !checkResult.IsHivDomain() {
		return
	}
	err = checkResult.checkClickCounter()
	if err != nil {
		checkResult.Valid = false
		return
	}
	err = checkResult.checkIframe()
	if err != nil {
		checkResult.Valid = false
		return
	}
	if len(checkResult.IframeTarget) > 0 {
		redirectUrl, redirectUrlErr := url.Parse(checkResult.IframeTarget)
		if redirectUrlErr != nil {
			checkResult.Valid = false
			err = redirectUrlErr
			return
		}
		if len(redirectUrl.Scheme) == 0 {
			redirectUrl.Scheme = checkResult.URL.Scheme
			checkResult.IframeTarget = redirectUrl.String()
		}
		redirectChecker := NewDomainCheckResult(redirectUrl.Host)
		redirectChecker.URL = redirectUrl
		redirectChecker.SaveBody = false
		redirectCheckErr := redirectChecker.Check()
		if redirectCheckErr != nil {
			checkResult.IframeTargetOk = false
			checkResult.Valid = false
			err = redirectCheckErr
			return
		} else {
			checkResult.IframeTargetOk = true
		}
	}
	return
}

var lookupHost = func(domain string) (addresses []string, err error) {
	addresses, err = net.LookupHost(domain)
	if err != nil {
		return
	}
	sort.Strings(addresses)
	return
}

// checks the DNS
func (checkResult *DomainCheckResult) dnsCheck() (err error) {
	checkResult.Addresses, err = lookupHost(checkResult.Domain)
	if err != nil {
		return
	}
	checkResult.DnsOk = true
	return
}

func TimeoutDialer(timeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, timeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(timeout))
		return conn, nil
	}
}

var timeout = time.Duration(5 * time.Second)

// fetches an URL and saves it as a temp file
// then opens it
func (checkResult *DomainCheckResult) fetch() (err error) {
	log.Printf("[%s] Fetching %s\n", checkResult.Domain, checkResult.URL)
	var response *http.Response
	url := checkResult.URL.String()
	transport := http.Transport{
		Dial: TimeoutDialer(time.Duration(5 * time.Second)),
	}
	client := http.Client{
		Transport: &transport,
	}
	response, err = client.Get(url)
	if err != nil {
		return
	}
	var newUrl = response.Request.URL
	if newUrl.String() != checkResult.URL.String() {
		log.Printf("[%s] Redirect to: %s\n", checkResult.Domain, newUrl)
		checkResult.URL = newUrl
	}
	checkResult.body, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return
	}

	if checkResult.SaveBody {
		tmpFile, tmpFileErr := ioutil.TempFile(os.TempDir(), checkResult.Domain+"-check")
		if tmpFileErr == nil {
			defer tmpFile.Close()
			tmpFile.Write(checkResult.body)
			checkResult.bodyFile = tmpFile.Name()
			log.Printf("[%s] Saved body to %s\n", checkResult.Domain, checkResult.bodyFile)
		} else {
			log.Printf("ERROR: Failed to save body to temp file.")
		}
	}

	checkResult.StatusCode = response.StatusCode
	log.Printf("[%s] Status %d\n", checkResult.Domain, checkResult.StatusCode)

	if checkResult.StatusCode != http.StatusOK {
		err = fmt.Errorf("Failed to load '%s': %s", checkResult.URL, checkResult.body)
		return
	}
	return
}

var scriptTagMatch = regexp.MustCompile(`<script[^>]+>`)
var iframeTagMatch = regexp.MustCompile(`<iframe[^>]+>`)
var srcAttributeMatch = regexp.MustCompile(`(src="([^"]+)"|src='([^']+)'|src=([^ ]+) )`)
var idAttributeMatch = regexp.MustCompile(`(id="clickcounter-target-iframe"|src='clickcounter-target-iframe'|src=clickcounter-target-iframe\W)`)

// Checks if the click-counter code snipped is installed
func (checkResult *DomainCheckResult) checkClickCounter() (err error) {
	allScripts := scriptTagMatch.FindAllSubmatch(checkResult.body, -1)
	for _, scriptTag := range allScripts {
		srcAttribute := srcAttributeMatch.FindSubmatch(scriptTag[0])
		if srcAttribute != nil {
			if string(srcAttribute[2]) == CLICKCOUNTER_SCRIPT || string(srcAttribute[3]) == CLICKCOUNTER_SCRIPT || string(srcAttribute[4]) == CLICKCOUNTER_SCRIPT {
				checkResult.ScriptPresent = true
			}
		}
	}
	if checkResult.ScriptPresent {
		log.Printf("[%s] click-counter script installed\n", checkResult.Domain)
	} else {
		err = fmt.Errorf("click-counter script not installed")
		return
	}
	return
}

// Checks if a click-counter iframe is used and the redirect works
func (checkResult *DomainCheckResult) checkIframe() (err error) {
	allIframes := iframeTagMatch.FindAllSubmatch(checkResult.body, -1)
	for _, iframeTag := range allIframes {
		idAttribute := idAttributeMatch.FindSubmatch(iframeTag[0])
		if idAttribute != nil {
			srcAttribute := srcAttributeMatch.FindSubmatch(iframeTag[0])
			checkResult.IframePresent = true
			keys := []int{2, 3, 4}
			for i := range keys {
				if len(srcAttribute[keys[i]]) > 0 {
					checkResult.IframeTarget = string(srcAttribute[keys[i]])
				}
			}
		}
	}
	if checkResult.IframePresent {
		log.Printf("[%s] iframe present\n", checkResult.Domain)
		if len(checkResult.IframeTarget) > 0 {
			log.Printf("[%s] iframe src: %s\n", checkResult.Domain, checkResult.IframeTarget)
		} else {
			err = fmt.Errorf("iframe has no src")
			return
		}
	}
	return
}

func CheckDomain(config *Config, domain string) (checkResult *DomainCheckResult, err error) {
	checkResult = NewDomainCheckResult(domain)
	err = checkResult.Check()
	if !checkResult.Valid {
		log.Printf("[%s] PROBLEM: %s\n", checkResult.Domain, err.Error())
	} else {
		log.Printf("[%s] A-OK\n", checkResult.Domain)
	}
	return
}
