package telegram

import (
	"context"
	"errors"
	"io/fs"
	"knok-bot/clients/telegram"
	"knok-bot/lib/e"
	"knok-bot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd" // rnd page
	HelpCmd  = "/help"
	StartCmd = "/start" // hi + help
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username) // to collect logs and what cmds people use.

	// Creating router
	if isAddCmd(text) { // add page: http://...
		return p.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := NewMessageSender(ctx, chatID, p.tg) // with closure just to try out closure

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		// return p.tg.SendMessage(chatID, msgAlreadyExists)
		return sendMsg(msgAlreadyExists) // with closure
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	//if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
	if err := sendMsg(msgSaved); err != nil { // with closure
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	// TODO: add check username folder, if don't exists, sendmsg that no links saved. // WARNING: maybe we already did it.
	if errors.Is(err, fs.ErrNotExist) {
		return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(ctx, chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) sendHelp(ctx context.Context, chatId int) error {
	return p.tg.SendMessage(ctx, chatId, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatId int) error {
	return p.tg.SendMessage(ctx, chatId, msgHello)
}

func NewMessageSender(ctx context.Context, chatID int, tg *telegram.Client) func(string) error { // Closure https://en.wikipedia.org/wiki/Closure_(computer_programming)
	return func(msg string) error {
		return tg.SendMessage(ctx, chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != "" // won't work with links like yle.fi
}
