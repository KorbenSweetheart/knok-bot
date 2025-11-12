package telegram

import (
	"log"
	"net/url"
	"read-reminder-bot/clients/telegram"
	"read-reminder-bot/lib/e"
	"read-reminder-bot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd" // rnd page
	HelpCmd  = "/help"
	StartCmd = "/start" // hi + help
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) { // add page: http://...
		// TODO: AddPage()
	}

	switch text {
	case RndCmd:
	case HelpCmd:
	case StartCmd:
	default:
	}

}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := NewMessageSender(chatID, p.tg) // with closure

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		// return p.tg.SendMessage(chatID, msgAlreadyExists)
		return sendMsg(msgAlreadyExists) // with closure
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	//if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
	if err := sendMsg(msgSaved); err != nil { // with closure
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random", err) }()

}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error { // Closure https://en.wikipedia.org/wiki/Closure_(computer_programming)
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != "" // won't work with links like yle.fi
}
