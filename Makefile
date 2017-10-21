
all: fmt cat tr wc head test

head: head.go
	go build -o head head.go

cat: cat.go
	go build -o cat cat.go

tr: tr.go
	go build -o tr tr.go

wc: wc.go
	go build -o wc wc.go

test: cat_test tr_test

cat_test:
	go test -cover cat.go cat_test.go

tr_test:
	go test -cover tr.go tr_test.go

fmt:
	go fmt

clean:
	go clean
