NAME = lablog

all:
	make generate
	make format
	#make vet
	make test
	make build

generate:
	cd src/web/; go-bindata -pkg="web" html/

format:
	find . -name "*.go" -not -path './Godeps/*' -type f -exec gofmt -s=true -w=true {} \;
	find . -name "*.go" -not -path './Godeps/*' -type f -exec goimports -w=true {} \;

test:
	go test

build:
	go build -ldflags "-X main.buildTime `date +%s` -X main.buildVersion `git describe --always`" -o "$(NAME)"

clean:
	rm "$(NAME)"
	rm *.pprof
	rm *.pdf

install:
	cp "$(NAME)" /usr/local/bin

uninstall:
	rm "/usr/local/bin/$(NAME)"

callgraph:
	go tool pprof --pdf "$(NAME)" cpu.pprof > callgraph.pdf

memograph:
	go tool pprof --pdf "$(NAME)" mem.pprof > memograph.pdf

dependencies_save:
	godep save ./...

dependencies_restore:
	godep restore ./...

bench:
	mkdir -p benchmarks/`git describe --always`/
	go test -test.benchmem=true -test.bench . 2> /dev/null | tee benchmarks/`git describe --always`/`date +%s`

coverage:
	rm -f coverage.out
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o=/tmp/coverage.html

lint:
	golint ./...

vet:
	go vet
