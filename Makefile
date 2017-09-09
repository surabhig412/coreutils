
all: fmt cat tr test


cat: cat.go
	go build -o cat cat.go

tr: tr.go
	go build -o tr tr.go

test: cat_test tr_test

cat_test:
	go test -coverprofile cover.out cat.go cat_test.go

tr_test:
	go test -coverprofile cover.out tr.go tr_test.go

fmt:
	go fmt

clean:
	go clean
