Go-Mail
=======

Just a simple little app for sending mail.<br />

Build:
```bash
go build mail.go
```

Usage:
```bash
echo "This is a test" | ./mail -from email@example.org -to email2@example.org -subject "This is an Example" -server smtp.google.com:465 -user username -pass secret
```
