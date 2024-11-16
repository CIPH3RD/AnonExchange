package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"anonymousoverflow/src/types"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func RedirectShortenedOverflowURL(c *gin.Context) {
	id := c.Param("id")
	answerId := c.Param("answerId")
	sub := c.Param("sub")

	// fetch the stack overflow URL
	client := resty.New()
	client.SetRedirectPolicy(
		resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}),
	)

	// Initialize the domain
	domain := "www.stackoverflow.com"

	// Check if the sub parameter contains a dot (subdomain), indicating it's already a full domain
	if strings.Contains(sub, ".") {
		// If subdomain is provided (like "subdomain.stackexchange.com"), use it directly
		domain = sub
	} else if sub != "" {
		// If it's just a simple name (like "askubuntu"), look for it in the ExchangeDomains list
		for _, exchangeDomain := range types.ExchangeDomains {
			if strings.Contains(sub, exchangeDomain) {
				// If it's a valid exchange domain, format it as "example.com"
				domain = fmt.Sprintf("%s.%s.com", sub, exchangeDomain)
			}
		}
	}
	resp, err := client.R().Get(fmt.Sprintf("https://%s/a/%s/%s", domain, id, answerId))
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

	// get the redirect URL
	location := resp.Header().Get("Location")

	redirectPrefix := os.Getenv("APP_URL")
	if sub != "" {
		redirectPrefix += fmt.Sprintf("/%s", sub)
	}

	c.Redirect(302, fmt.Sprintf("%s%s", redirectPrefix, location))
}
