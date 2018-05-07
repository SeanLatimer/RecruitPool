deps:
	go get -d

run:
	go run main.go

release:
	goreleaser

fixme:
	grep -rnwi --include *.go --include *.proto --exclude *.pb.go  --exclude .+_test.go "FIXME" .

todo:
	grep -rnwi --include *.go --include *.proto --exclude *.pb.go  --exclude .+_test.go "TODO" .

# Legacy code should be remove by the time of release
legacy:
	grep -rnwi --include *.go --include *.proto --exclude *.pb.go  --exclude .+_test.go "LEGACY" .
