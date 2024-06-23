
//for windows
GOOS=windows go build -o client.exe -ldflags="-H windowsgui"
GOOS=windows go build -o server.exe -ldflags="-H windowsgui"