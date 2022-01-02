.PHONY: build clean

build:
	go build -ldflags "-X main.BuildDate=`date +%s` \
	  -X main.GitCommit=`git rev-parse --short HEAD` \
    -X main.Version=1.0.0 " \
		-o shell-compose

clean:
	rm -f shell-compose
