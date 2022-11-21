package gotranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/language"
)

type GoogleTranslate struct {
	Host   string
	Client *http.Client
	Target language.Tag
	Source language.Tag
}

type Sentences struct {
	Trans string `json:"trans"`
	Orig  string `json:"orig"`
}

type Response struct {
	Sentences []Sentences `json:"sentences"`
}

type GoogleTranslateOpt func(c *GoogleTranslate)

func WithTarget(tg language.Tag) GoogleTranslateOpt {
	return func(g *GoogleTranslate) {
		g.Target = tg
	}
}

func NewGoogleTranslate(opts ...GoogleTranslateOpt) *GoogleTranslate {
	client := &GoogleTranslate{
		Host:   "https://translate.googleapis.com",
		Client: http.DefaultClient,
		Source: language.English,
		Target: language.Vietnamese,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *GoogleTranslate) Translate(source string, opts ...GoogleTranslateOpt) (string, error) {
	if len(strings.TrimSpace(source)) == 0 {
		return source, nil
	}
	for _, opt := range opts {
		opt(c)
	}
	resp, err := c.Client.Get(fmt.Sprintf("%s/translate_a/single?client=gtx&dj=1&dt=t&sl=%s&tl=%s&q=%s", c.Host, c.Source, c.Target, url.QueryEscape(source)))
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}
	res := Response{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if len(res.Sentences) == 0 {
		return source, nil
	}
	return res.Sentences[0].Trans, nil
}
