
# the -ldflags=-s remove the debug info making the executable much smaller
# -H=windowsgui no terminal console window
# Build for macos
go build -ldflags="-s -w" -o bin/ mcutie.go
# Build for windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o bin/mcutie.exe mcutie.go

#Compress the resulting files with UPX - not currently working with macos binaries under Bigsur
#upx mcutie.exe - windows 10 is false positive identifying this as a virus