.PHONY: build clean

build:
	go build -ldflags "-X main.BuildDate=`date +%s` \
	  -X main.GitCommit=`git rev-parse --short HEAD` \
    -X main.Version=`git describe --tags --abbrev=0`.`git rev-list \`git describe --tags --abbrev=0\`.. --count`" \
		-o shell-compose

clean:
	rm -f shell-compose
