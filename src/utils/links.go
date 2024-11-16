package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// stackOverflowLinkQualifierRegex matches all anchor elements that meet the following conditions:
// * must be an anchor element
// * the anchor element must have a pathname beginning with /q or /questions
// * if there is a host, it must be stackoverflow.com or a subdomain
var stackOverflowLinkQualifierRegex = regexp.MustCompile(`<a\s[^>]*href="(?:https?://(?:www\.)?(?:\w+\.)*(?:stackoverflow|stackexchange|superuser|serverfault|askubuntu)\.com)?/(?:q|questions)/[^"]*"[^>]*>.*?</a>`)

func ReplaceStackOverflowLinks(html string) string {
	return stackOverflowLinkQualifierRegex.ReplaceAllStringFunc(html, func(match string) string {
		// Extract the href attribute value from the anchor tag
		hrefRegex := regexp.MustCompile(`href="([^"]*)"`)
		hrefMatch := hrefRegex.FindStringSubmatch(match)
		if len(hrefMatch) < 2 {
			return match
		}
		href := hrefMatch[1]

		// Parse the URL
		parsedUrl, err := url.Parse(href)
		if err != nil {
			return match
		}

		// Extract the host from the URL
		host := parsedUrl.Host
		parts := strings.Split(host, ".")

		// Initialize newPath with the original path
		newPath := parsedUrl.Path

		// Determine the new path based on the domain
		if strings.Contains(host, "askubuntu.com") {
			// If the host contains "askubuntu.com", use "/askubuntu/"
			if len(parts) > 2 && parts[0] != "askubuntu" {
				// Handle subdomains
				newPath = "/askubuntu/" + parts[0] + newPath
			} else {
				// No need to repeat "askubuntu" in the path
				newPath = "/askubuntu" + newPath
			}
		} else if strings.Contains(host, "serverfault.com") {
			// If the host contains "serverfault.com", use "/serverfault/"
			if len(parts) > 2 && parts[0] != "serverfault" {
				// Handle subdomains
				newPath = "/serverfault/" + parts[0] + newPath
			} else {
				// No need to repeat "serverfault" in the path
				newPath = "/serverfault" + newPath
			}
		} else if strings.Contains(host, "superuser.com") {
			// If the host contains "superuser.com", use "/superuser/"
			if len(parts) > 2 && parts[0] != "superuser" {
				// Handle subdomains
				newPath = "/superuser/" + parts[0] + newPath
			} else {
				// No need to repeat "superuser" in the path
				newPath = "/superuser" + newPath
			}
		} else if len(parts) > 2 {
			// For other subdomains, use "/exchange/"
			newPath = "/exchange/" + parts[0] + newPath
		}

		// Reconstruct the new URL
		newUrl := newPath + parsedUrl.RawQuery + parsedUrl.Fragment

		// Replace the href attribute value in the anchor tag
		return strings.Replace(match, hrefMatch[1], newUrl, 1)
	})
}
