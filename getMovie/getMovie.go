package main  


import (
	"encoding/json"
	basic "go-backend/basic"
	"github.com/aws/aws-lambda-go/lambda"
)



func handler(request basic.Request)(basic.Response , error) {
	movie_id := request.PathParameters["movie_id"]
	thisItem , err := basic.GetMovie(movie_id)

	if err != nil {
		return basic.Response{
			Body: "Movie not found",
			StatusCode : 400 , 
		}, nil 
	}

	marshalledItem , _ := json.Marshal(thisItem)

	return basic.Response{
		Body: string(marshalledItem), 
		StatusCode: 200,
		Headers : map[string]string{
			"Content-Type" : "application/json",
		},
	}, nil
	
}


func main() {
	lambda.Start(handler)
}
