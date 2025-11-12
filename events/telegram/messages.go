package telegram

const msgHelp = `I can save and keep your pages. Also, I can offer them to read.

To save the page, just send me a link to it.

To get a random page from your list, send me the command /rnd.
Caution: After that, this page will be removed from your list!`

const msgHello = "Hi there! 👋\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤯"
	msgNoSavedPages   = "You don't have saved pages 😩"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You already have this page in your list 🤓"
)
