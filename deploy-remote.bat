set GOOS=linux
go build -o nanping-pt-carbon ./cmd/main.go
C:\Windows\System32\OpenSSH\scp.exe .\nanping-pt-carbon root@121.41.129.227:/home/app/nanping-carbon/tmp
