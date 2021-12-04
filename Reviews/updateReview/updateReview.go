package main

import (
	"fmt"
	basic "go-backend/basic"
	"github.com/aws/aws-lambda-go/lambda"
)


func handler(request basic.Request)(basic.Response , error) {
	review_id := request.PathParameters["review_id"]
	err := basic.UpdateReview(request.Body , review_id) 
	if err != nil {
		return basic.Response{
			Body: err.Error(),
			StatusCode: 400,
		},nil
	}

	message := fmt.Sprintf("updated review_id : %s" , review_id)
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