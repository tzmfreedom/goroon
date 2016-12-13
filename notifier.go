package main

import (
	"github.com/deckarep/gosx-notifier"
)

type Notifier struct {
}

func NewNotifier() *Notifier {
	n := &Notifier{}
	return n
}

func (n *Notifier) Notify(title string, subTitle string) error {
	//At a minimum specifiy a message to display to end-user.
	note := gosxnotifier.NewNotification(title)

	//Optionally, set a title
	note.Title = title

	//Optionally, set a subtitle
	note.Subtitle = subTitle

	//Optionally, set a sound from a predefined set.
	//note.Sound = gosxnotifier.Basso

	//Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
	//note.Group = "com.unique.yourapp.identifier"

	//Optionally, set a sender (Notification will now use the Safari icon)
	//note.Sender = "com.apple.Safari"

	//Optionally, specifiy a url or bundleid to open should the notification be
	//clicked.
	note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

	//Optionally, an app icon (10.9+ ONLY)
	//note.AppIcon = "gopher.png"

	//Optionally, a content image (10.9+ ONLY)
	//note.ContentImage = "gopher.png"

	//Then, push the notification
	err := note.Push()
	if err != nil {
		return err
	}
	return nil
}
