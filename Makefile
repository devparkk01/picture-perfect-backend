.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/addMovie addMovie/addMovie.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getMovie getMovie/getMovie.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/deleteMovie deleteMovie/deleteMovie.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/updateMovie updateMovie/updateMovie.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/addReview Reviews/addReview/addReview.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getReview Reviews/getReview/getReview.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/deleteReview Reviews/deleteReview/deleteReview.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/updateReview Reviews/updateReview/updateReview.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getReviews Reviews/getReviews/getReviews.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getMovies getMovies/getMovies.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
