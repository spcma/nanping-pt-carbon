set GOOS=linux
go build -o nanping-pt-carbon ./cmd/main.go
C:\Windows\System32\OpenSSH\scp.exe .\nanping-pt-carbon root@192.168.1.10:/home/yanping/app/nanping/yanping-carbon/tmp
