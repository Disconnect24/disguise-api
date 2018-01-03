package disguise

// The name will be the mlid.
type Accounts struct {
	Passwd  string
	Mlchkid string
}

// Stored mail. We don't need a mail ID as Cloud Datastore does that for us.
// (Incomplete key assigns an automatic name)
type Mail struct {
	SenderID    string
	Body        string
	RecipientID string
	Sent        bool
}
