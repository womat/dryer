set GOARCH=arm
set GOOS=linux
go build -o ..\bin\dryer ..\cmd\main.go ..\cmd\config.go ..\cmd\loadsave.go ..\cmd\server.go ..\cmd\webservice.go

set GOARCH=386
set GOOS=windows
go build -o ..\bin\dryer.exe ..\cmd\main.go ..\cmd\config.go ..\cmd\loadsave.go ..\cmd\server.go ..\cmd\webservice.go