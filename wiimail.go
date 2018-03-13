package disguise

import (
	//"encoding/base64"
	"fmt"
	"strings"
)

const CRLF = "\r\n"

func FormulateMail(from string, to string, body string, potentialImage []byte) (string, error) {
	boundary := GenerateBoundary()

	// Set up headers and set up first boundary with body.
	// The body could be empty: that's fine, it'll have no value
	// (compared to nil) and the Wii will ignore that section.
	mailContent := fmt.Sprint("From: ", from, CRLF,
		"Subject: ", "PC Wii Mail", CRLF,
		"To: ", to, CRLF,
		"MIME-Version: 1.0", CRLF,
		"Content-Type: MULTIPART/mixed; BOUNDARY=", `"`, boundary, `"`, CRLF,
		CRLF,
		"--", boundary, CRLF,
		"Content-Type: TEXT/plain; CHARSET=utf-8",  CRLF,
		"Content-Description: wiimail", CRLF,
		CRLF,
		body,
		strings.Repeat(CRLF, 3),
		"--", boundary,
	)

	// If there's an attachment, we need to factor that in.
	// Otherwise we're done.
	if potentialImage == nil {
		return fmt.Sprint(mailContent, "--"), nil
	}

	// We're assuming once we're called that this is a jpeg.
	// Go ahead and convert its binary to base64.

	return "", nil
	//
	//	ext := mime.TypeByExtension(filepath.Ext(filename))
	//	if ext == "" {
	//		ext = "text/plain"
	//	}
	//
	//	h := textproto.MIMEHeader{}
	//	h.Add("Content-Type", ext)
	//	h.Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	//	h.Add("Content-Transfer-Encoding", "base64")
	//	newpart, err := multiw.CreatePart(h)
	//	if err != nil {
	//		return err
	//	}
	//	buf := bytes.NewBuffer([]byte{})
	//	bcdr := NewBase64Email(buf, base64.StdEncoding)
	//	if _, err = io.Copy(bcdr, file); err != nil {
	//		return err
	//	}
	//	if err = bcdr.Close(); err != nil {
	//		return err
	//	}
	//	if _, err = io.Copy(newpart, buf); err != nil {
	//		return err
	//	}
	//}
	//return multiw.Close()
}