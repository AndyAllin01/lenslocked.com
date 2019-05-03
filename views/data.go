package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	AlertMsgGeneric = "Something went wrong. Please try again and contact us if the problem persists"
)

//Alert used to render bootstrap alert messages
type Alert struct {
	Level   string
	Message string
}

//Data is the top level structure that views expect
//to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}
