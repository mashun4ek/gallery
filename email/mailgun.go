package email

import (
	"fmt"
	"net/url"

	"github.com/mailgun/mailgun-go"
)

var (
	welcomeText = `Hi there!

	Welcome to MyGravitation! 
	
	Best,
	Maria
	`

	welcomeHTML = `Hi there!
<br/>
Welcome to my awesome <a href="https://mygravitation.com>website!</a>.
<br/>
Stay tuned!
MyGravitation`
	welcomeSubject = "Welcome to MyGravitation.com"
	resetSubject   = "Instructions for resetting your password."
	resetBaseURL   = "https://mygravitation/reset"

	resetTextTmpl = `HELLO!

If you requested a password reset. please follow link below:

%s

If you are asked for a token. please use the following value:

%s

If you didn't request, ignore this email.

Best regards,
MyGravitation Support
`

	resetHTMLTmpl = `HELLO!<br/>
<br/>
If you requested a password reset. please follow link below:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token. please use the following value:<br/>
<br/>
%s
<br/>
If you didn't request, ignore this email.<br/>
<br/>
Best regards,<br/>
MyGravitation Support<br/>
`
)

func WithMailgun(domain, apiKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey)
		c.mg = mg
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(c *Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// from: "mashun4ek@gmail.com",
		from: "demo@sandbox0a473d6bb8c94265aa28ced67515768c.mailgun.org",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
