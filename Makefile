TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=local
NAMESPACE=kffl
NAME=totp
BINARY=terraform-provider-${NAME}
VERSION=0.1.2
OS != go env GOHOSTOS
ARCH != go env GOHOSTARCH
OS_ARCH=${OS}_${ARCH}

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   