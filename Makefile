# note: call scripts from /scripts
go-build-armv7l:
	GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-w -s -extldflags "-static"' -o exe -mod vendor -v ./cmd/v1
