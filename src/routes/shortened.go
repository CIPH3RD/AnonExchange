package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func RedirectShortenedOverflowURL(c *gin.Context) {
	id := c.Param("id")
	answerId := c.Param("answerId")
	sub := c.Param("sub")

	// Define the allowed domains
	var stackDomains = map[string]bool{
		"stackoverflow.com": true,
		"askubuntu.com":     true,
		"superuser.com":     true,
		"serverfault.com":   true,
	}

	// Determine the domain
	var domain string
	if stackDomains[sub] {
		// Directly use the provided domain if it is allowed
		domain = sub
	} else if strings.Contains(sub, ".") {
		// Treat as a default Stack Exchange domain if not in the allowed list
		domain = "stackoverflow.com"
	} else if sub != "" {
		// For non-empty subdomains that are not in stackDomains, treat them as Stack Exchange
		domain = fmt.Sprintf("%s.stackexchange.com", sub)
	} else {
		// Default domain
		domain = "stackoverflow.com"
	}

	// Fetch the stack overflow URL
	client := resty.New()
	client.SetRedirectPolicy(
		resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}),
	)

	// Construct the URL to fetch
	urlToFetch := fmt.Sprintf("https://%s/a/%s/%s", domain, id, answerId)
	resp, err := client.R().Get(urlToFetch)
	if err != nil {
		c.HTML(400, "home.html", gin.H{
			"errorMessage": "Unable to fetch stack overflow URL",
		})
		return
	}

	if resp.StatusCode() != 302 {
		c.HTML(400, "home.html", gin.H{
			"errorMessage": fmt.Sprintf("Unexpected HTTP status from origin: %d", resp.StatusCode()),
		})
		return
	}

	// Get the redirect URL
	location := resp.Header().Get("Location")

	// Determine the correct path prefix based on the domain
	var pathPrefix string
	switch {
	case domain == "askubuntu.com":
		pathPrefix = "/askubuntu"
	case domain == "serverfault.com":
		pathPrefix = "/serverfault"
	case domain == "superuser.com":
		pathPrefix = "/superuser"
	case strings.HasSuffix(domain, ".stackexchange.com"):
		subDomain := strings.TrimSuffix(domain, ".stackexchange.com")
		pathPrefix = "/exchange/" + subDomain
	default:
		pathPrefix = "/exchange"
	}

	// Construct the redirect URL
	redirectPrefix := os.Getenv("APP_URL")
	if sub != "" {
		redirectPrefix += pathPrefix
	} else {
		redirectPrefix += "/exchange"
	}

	c.Redirect(302, fmt.Sprintf("%s%s", redirectPrefix, location))
}
