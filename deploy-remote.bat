set GOOS=linux
go build -o nanping-pt-carbon ./cmd/main.go
C:\Windows\System32\OpenSSH\scp.exe .\nanping-pt-carbon root@115.236.162.150:/home/yanping/app/nanping/yanping-carbon/tmp
