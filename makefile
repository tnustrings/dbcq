all: bin/dbcq

bin/dbcq: dbcq.ct cmd/dbcq/main.ct
	ct dbcq.ct
	ct cmd/dbcq/main.ct
	#mv main.go cmd/dbcq
	go build -o bin/dbcq cmd/dbcq/main.go
