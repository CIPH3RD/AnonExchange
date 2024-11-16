package routes

import (
	"anonymousoverflow/config"
	"anonymousoverflow/src/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHome(c *gin.Context) {
	theme := utils.GetThemeFromEnv()
	c.HTML(200, "home.html", gin.H{
		"version": config.Version,
		"theme":   theme,
	})
}

type urlConversionRequest struct {
	URL string `form:"url" binding:"required"`
}

var coreRegex = regexp.MustCompile(`(?:https?://)?(?:www\.)?([^/]+)(/(?:questions|q|a)/.+)`)

// Will return `nil` if `rawUrl` is invalid.
func translateUrl(rawUrl string) string {
	coreMatches := coreRegex.FindStringSubmatch(rawUrl)
	if coreMatches == nil {
		return ""
	}

	// Extract domain and path from the URL
	domain := strings.TrimSpace(coreMatches[1])
	rest := coreMatches[2]

	// Clean up the domain by removing any leading or trailing whitespace
	domain = strings.ToLower(strings.Trim(domain, " "))

	// Define path prefix based on domain
	var pathPrefix string

	switch {
	case domain == "stackoverflow.com":
		pathPrefix = ""
	case domain == "askubuntu.com":
		pathPrefix = "/askubuntu"
	case domain == "serverfault.com":
		pathPrefix = "/serverfault"
	case domain == "superuser.com":
		pathPrefix = "/superuser"
	case strings.HasSuffix(domain, ".stackexchange.com"):
		subDomain := strings.TrimSuffix(domain, ".stackexchange.com")
		if subDomain == "" {
			return ""
		}
		// Full domain with dots should be used with "/exchange/" prefix
		pathPrefix = "/exchange/" + domain
		//	default:
		//		// Default behavior for other domains
		//		pathPrefix = "/exchange/" + domain
	}

	// Ensure proper formatting of the return string
	if pathPrefix == "" {
		return rest
	}
	if strings.HasPrefix(rest, "/") {
		return fmt.Sprintf("%s%s", pathPrefix, rest)
	}
	return fmt.Sprintf("%s/%s", pathPrefix, rest)
}

func PostHome(c *gin.Context) {
	body := urlConversionRequest{}

	if err := c.ShouldBind(&body); err != nil {
		c.HTML(400, "home.html", gin.H{
			"errorMessage": "Invalid request body",
		})
		return
	}

	translated := translateUrl(body.URL)

	if translated == "" {
		theme := utils.GetThemeFromEnv()
		c.HTML(400, "home.html", gin.H{
			"errorMessage": "Invalid stack overflow/exchange URL",
			"theme":        theme,
		})
		return
	}

	c.Redirect(302, translated)
}
