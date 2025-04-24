package service

import (
	"fmt"
	"go-project/logger"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func FetchContent(url string, _type string) (string, error) {
	maxSectionLen := 1200

	var doc *goquery.Document
	var err error

	if _type == "http" {
		// 使用传统的HTTP请求获取内容
		var resp *http.Response
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		// 解析文档
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", err
		}
	} else {
		// 使用无头浏览器获取内容
		doc, err = initializePageWhitRod(url)
		if err != nil {
			return "", err
		}
	}

	var sections []string
	var sectionContent string

	doc.Find("h1, h2, h3, h4, h5, h6, p, strong, span").Each(func(i int, s *goquery.Selection) {
		// // 检查元素的style属性
		if style, exists := s.Attr("style"); exists && strings.Contains(style, "display:none") {
			// 如果style包含display:none，则跳过该元素
			return
		}

		// 去除前后空格
		text := strings.TrimSpace(s.Text())
		if text != "" {
			if !strings.Contains(sectionContent, text) {
				sectionContent += text + "\n"
			}
		}

		s.Find("img").Each(func(j int, img *goquery.Selection) {
			imgSrc, exists := img.Attr("src")
			if exists {
				sectionContent += fmt.Sprintf("<img src=\"%s\">\n", imgSrc)
			}
		})

		// 分段阈值
		if len(sectionContent) >= maxSectionLen {
			sections = append(sections, sectionContent)
			sectionContent = ""
		}
	})

	// 如果最后的sectionContent不为空，添加到sections
	if sectionContent != "" {
		sections = append(sections, strings.TrimRight(sectionContent, "\n"))
	}

	// return sections, nil
	// Convert sections to UTF-8 encoding to handle Chinese characters properly
	result := strings.Join(sections, "\n")

	if !utf8.ValidString(result) {
		// Try to decode as GBK if not valid UTF-8
		gbkBytes := []byte(result)
		utf8Bytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(gbkBytes)
		if err == nil {
			result = string(utf8Bytes)
		}
	}

	return result, nil
}

func initializePageWhitRod(url string) (*goquery.Document, error) {
	logger.Logger.Debug("rod初始化")
	l := launcher.MustNewManaged("ws://localhost:7317").Headless(true)
	browser := rod.New().Client(l.MustClient()).MustConnect()
	logger.Logger.Debug("rod初始化成功")

	defer browser.MustClose()

	page := browser.MustPage(url)

	page.MustWaitLoad()

	// 2秒超时，等待dom稳定性变化阈值95%
	err := page.WaitDOMStable(2*time.Second, 0.95)

	if err != nil {
		fmt.Println("Error waiting for DOMContentLoaded event")
	}

	content, err := page.HTML()
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func FetchLinksWithHttp(urlStr string) ([]string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	urlStr = strings.TrimRight(urlStr, "/")

	links := make(map[string]struct{})
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "link") {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link := strings.TrimSpace(a.Val)
					if !regexp.MustCompile(`\.(css|js|ico)$|(?i)javascript`).MatchString(link) {
						if link == "" || strings.HasPrefix(link, "//") || (regexp.MustCompile(`//+`).MatchString(link) && !strings.HasPrefix(link, urlStr)) {
							return
						}
						if regexp.MustCompile(`^[:a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(link) {
							return
						}
						if strings.HasSuffix(link, "#") || regexp.MustCompile(`#([a-zA-Z0-9]+)$`).MatchString(link) {
							return
						}
						if strings.Contains(link, "/img/") {
							return
						}

						baseURL, err := url.Parse(urlStr)
						if err != nil {
							return
						}

						baseDomain := baseURL.Scheme + "://" + baseURL.Host

						if !strings.HasPrefix(link, "http") {
							relativeURL, _ := url.Parse(link)

							absoluteURL := baseURL.ResolveReference(relativeURL)
							link = absoluteURL.String()
						}

						if strings.HasPrefix(link, baseDomain) {
							links[link] = struct{}{}
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	if len(links) == 0 {
		links[urlStr] = struct{}{}
	}

	sortedLinks := make([]string, 0, len(links))
	for link := range links {
		sortedLinks = append(sortedLinks, link)
	}

	sort.Strings(sortedLinks)

	return sortedLinks, nil
}

func FetchLinksWithRod(urlStr string) ([]string, error) {
	doc, err := initializePageWhitRod(urlStr)
	if err != nil {
		return nil, err
	}

	urlStr = strings.TrimRight(urlStr, "/")

	links := make(map[string]struct{})

	doc.Find("a, link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			href = strings.TrimSpace(href)
			if !regexp.MustCompile(`\.(css|js|ico)$|(?i)javascript`).MatchString(href) {
				if href == "" || strings.HasPrefix(href, "//") || (regexp.MustCompile(`//+`).MatchString(href) && !strings.HasPrefix(href, urlStr)) {
					return
				}
				if regexp.MustCompile(`^[:a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(href) {
					return
				}
				if strings.HasSuffix(href, "#") || regexp.MustCompile(`#([a-zA-Z0-9]+)$`).MatchString(href) {
					return
				}
				if strings.Contains(href, "/img/") {
					return
				}

				baseURL, err := url.Parse(urlStr)
				if err != nil {
					return
				}

				baseDomain := baseURL.Scheme + "://" + baseURL.Host

				if !strings.HasPrefix(href, "http") {
					relativeURL, _ := url.Parse(href)

					absoluteURL := baseURL.ResolveReference(relativeURL)
					href = absoluteURL.String()
				}

				if strings.HasPrefix(href, baseDomain) {
					links[href] = struct{}{}
				}
			}
		}
	})

	if len(links) == 0 {
		links[urlStr] = struct{}{}
	}

	sortedLinks := make([]string, 0, len(links))
	for link := range links {
		sortedLinks = append(sortedLinks, link)
	}

	sort.Strings(sortedLinks)

	return sortedLinks, nil
}
