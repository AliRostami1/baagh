package core

var events = newEventRegistry()

func Subscribe(fns ...EventHandler) {
	events.AddEventListener(fns...)
}
