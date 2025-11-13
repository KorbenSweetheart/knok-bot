package telegram

const msgHelp = `I can save articles for later. And offer them to read on request.

To save the page, just send me a link to it.
e.g. https://go.dev/

You can request an article by sending this command:

/rnd - to get a random page from your list.
⚠️ After that, this page will be removed from your list!`

const msgHello = "Hi there! 👋\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤯"
	msgNoSavedPages   = "You don't have saved pages 😩"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You already have this page in your list 🤓"
)
