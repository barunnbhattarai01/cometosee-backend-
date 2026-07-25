package common

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(to, subject, content string) error {
	from := "barunnbhattarai@gmail.com"
	password := os.Getenv("APP_PASSWORD")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>

<head>
<meta charset="UTF-8">
<title>%s</title>
</head>

<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,Helvetica,sans-serif;">

<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f6f9;padding:40px 0;">
<tr>
<td align="center">

<table width="650" cellpadding="0" cellspacing="0"
style="background:#ffffff;border-radius:14px;overflow:hidden;
box-shadow:0 8px 25px rgba(0,0,0,.08);">

<tr>
<td style="
background:linear-gradient(90deg,#ff6b00,#ff8f3d);
padding:28px;
text-align:center;
color:white;
font-size:30px;
font-weight:bold;">
 Cometosee
</td>
</tr>

<tr>
<td style="padding:40px;">

<h2 style="margin-top:0;color:#222;">
%s
</h2>

<div style="
font-size:16px;
line-height:30px;
color:#555;
">
%s
</div>

<div style="text-align:center;margin-top:40px;">

<a href="https://cometosee.vercel.app"
style="
background:#ff6b00;
padding:15px 35px;
color:white;
text-decoration:none;
border-radius:8px;
font-weight:bold;
font-size:16px;
display:inline-block;
">
Open Cometosee
</a>

</div>

<hr style="margin-top:45px;border:none;border-top:1px solid #eee;">

<p style="
font-size:13px;
color:#888;
text-align:center;
line-height:22px;
">

This is an automated email from
<strong>Cometosee</strong>.<br>

Please do not reply to this email.

</p>

</td>
</tr>

<tr>

<td style="
background:#fafafa;
padding:18px;
text-align:center;
font-size:12px;
color:#999;
">

© 2026 Cometosee. All Rights Reserved.

</td>

</tr>

</table>

</td>
</tr>
</table>

</body>
</html>
`, subject, subject, content)

	message := []byte(fmt.Sprintf(
		"From: Cometosee <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n%s",
		from,
		to,
		subject,
		html,
	))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)
}
