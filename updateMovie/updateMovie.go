package main

import (
	"fmt"
	basic "go-backend/basic"
	"github.com/aws/aws-lambda-go/lambda"
)


func handler(request basic.Request)(basic.Response , error) {
	movie_id := request.PathParameters["movie_id"]
	err := basic.UpdateMovie(request.Body , movie_id) 
	if err != nil {
		return basic.Response{
			Body: err.Error(),
			StatusCode: 400,
		},nil
	}

	message := fmt.Sprintf("updated movie_id : %s" , movie_id)
	return basic.Response{
		Body: message,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type" : "application/json",
		},
	},nil


}


func main() {
	lambda.Start(handler)
}