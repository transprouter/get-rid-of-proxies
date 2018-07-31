all: build test specs

install:
	go get -t -v ./...
	pip install --quiet --upgrade ansible docker dnspython
	bundle install --gemfile specs/Gemfile

build:
	CGO_ENABLED=0 go build

test:
	go test -v -coverprofile=coverage.txt -covermode=atomic \
		./proxy

acceptance-infra:
	ansible-playbook specs/acceptance-infra.yml

specs: acceptance-infra
	cd specs/ && cucumber

clean:
	rm -f coverage.txt
	rm -f transprouter
	rm -f *.retry specs/*.retry
