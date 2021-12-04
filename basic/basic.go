package basic

import (
	"encoding/json"
	"fmt"
	"strconv"

	// "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

const MOVIE_TABLE = "movie-table"
const REVIEW_TABLE = "review-table"
type Response events.APIGatewayProxyResponse 
type Request events.APIGatewayProxyRequest

type Movie_item struct{
	Movie_id string `json:"movie_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Genre []string `json:"genre"`
	StarCast []string `json:"starCast"`
	Director []string `json:"director"`
	Poster string `json:"poster"`
}

type Review_item struct {
	Review_id string `json:"review_id"`
	Movie_id string `json:"movie_id"`
	Username string `json:"username"`
	Headline string `json:"headline"`
	Rating int `json:"rating"`
	Review string `json:"review"`
}


func createSession() (*dynamodb.DynamoDB) {
	mySession := session.Must(session.NewSession())
	return dynamodb.New(mySession) 
}


func AddMovie(body string )(Movie_item , error) {
	svc := createSession()

	// unMarshall the request body
	var thisItem Movie_item
	json.Unmarshal([]byte(body), &thisItem)

	thisItem.Movie_id = uuid.NewString()


	// Marshall the Item into a Map DynamoDB can deal with
	attributeValues , err := dynamodbattribute.MarshalMap(thisItem)
	if err != nil {
		// fmt.Println("Got error marshalling map:")
		// fmt.Println(err.Error())
		return thisItem, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(MOVIE_TABLE),
	}
	_, err = svc.PutItem(input)
	return thisItem, err

}


func GetMovie(movie_id string)(Movie_item , error) {
	svc := createSession()
	
	output , err := svc.GetItem(&dynamodb.GetItemInput{
		TableName : aws.String(MOVIE_TABLE), 
		Key: map[string]*dynamodb.AttributeValue{
			"movie_id": {
				S: aws.String(movie_id),
			},
		},
	})
	
	if err != nil {
		return Movie_item{} , nil 
	}
	
	var thisItem Movie_item 
	
	dynamodbattribute.UnmarshalMap(output.Item , &thisItem)
	
	return thisItem , nil 
}

func DeleteMovie(movie_id string)(error){
	svc := createSession()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"movie_id" :{
				S: aws.String(movie_id),
			},
		},
		TableName : aws.String(MOVIE_TABLE), 
	}

	_ , err := svc.DeleteItem(input)
	if err != nil {
		return err 
	}
	return nil 
}

func UpdateMovie(body string, movie_id string)(error){
	svc := createSession()

	var thisItem Movie_item
	json.Unmarshal([]byte(body) , &thisItem) 

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(MOVIE_TABLE),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(thisItem.Title),
			},
			":d": {
				S: aws.String(thisItem.Description),
			},
			":p" : {
				S : aws.String(thisItem.Poster),
			},

		},

		Key: map[string]*dynamodb.AttributeValue{
			"movie_id": {
				S: aws.String(movie_id),
			},
		},
		UpdateExpression: aws.String("set title = :t , description = :d , poster = :p"),
		ReturnValues: aws.String("UPDATED_NEW"),

	}
	_ , err := svc.UpdateItem(input)

	if err != nil {
		return err
	}
	return nil 
}

// Reviews related functions

func AddReview(body string)(Review_item , error) {
	svc := createSession()
	var thisItem Review_item 
	json.Unmarshal([]byte(body) , &thisItem) 

	thisItem.Review_id = uuid.NewString() 
	
	attributeValues , err := dynamodbattribute.MarshalMap(thisItem)
	if err!= nil {
		return thisItem , err
	}

	input := &dynamodb.PutItemInput{
		Item : attributeValues , 
		TableName : aws.String(REVIEW_TABLE),
	}

	_ , err = svc.PutItem(input)
	return thisItem, err

}

func GetReview(review_id string)(Review_item , error) {

	svc := createSession()
	output , err := svc.GetItem(&dynamodb.GetItemInput{
		TableName : aws.String(REVIEW_TABLE), 
		Key: map[string]*dynamodb.AttributeValue{
			"review_id": {
				S: aws.String(review_id),
			},
		},
	})
	
	if err != nil {
		return Review_item{} , nil 
	}
	
	var thisItem Review_item 
	
	dynamodbattribute.UnmarshalMap(output.Item , &thisItem)
	return thisItem , nil 

}

func DeleteReview(review_id string)(error){
	svc := createSession()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"review_id" :{
				S: aws.String(review_id),
			},
		},
		TableName : aws.String(REVIEW_TABLE), 
	}

	_ , err := svc.DeleteItem(input)
	if err != nil {
		return err 
	}
	return nil 

}

func UpdateReview(body string, review_id string)( error) {
	svc := createSession()

	var thisItem Review_item
	json.Unmarshal([]byte(body) , &thisItem) 

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(REVIEW_TABLE),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":h": {
				S: aws.String(thisItem.Headline),
			},
			":r": {
				N : aws.String(strconv.FormatInt(int64(thisItem.Rating) , 10)),
			},
			":re" : {
				S : aws.String(thisItem.Review),
			},

		},

		Key: map[string]*dynamodb.AttributeValue{
			"review_id": {
				S: aws.String(review_id),
			},
		},
		UpdateExpression: aws.String("set headline = :h , rating = :r , review = :re"),
		ReturnValues: aws.String("UPDATED_NEW"),

	}
	_ , err := svc.UpdateItem(input)

	if err != nil {
		return err 
	}
	return nil 


}


func GetReviews(movie_id string)([]Review_item , error){
	svc := createSession()
	fetchedItems := []Review_item{}

	//  create the filter expression 
	filterExp := expression.Name("movie_id").Equal(expression.Value(movie_id))
	// create the projection, i.e, the attributes that we want 
	projection := expression.NamesList((expression.Name("review_id")), expression.Name("movie_id"), 
		expression.Name("username") ,expression.Name("headline"), expression.Name("rating"),
		expression.Name("review"),
	) 

	expr , err := expression.NewBuilder().WithFilter(filterExp).WithProjection(projection).Build()

	if err != nil{
		fmt.Println("Go error building expression")
		return fetchedItems , err
	}

	//  build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:        expr.Filter(),
		ProjectionExpression:    expr.Projection(),
		TableName:   aws.String(REVIEW_TABLE),
	}

	result , err := svc.Scan(params)
	if err != nil{
		return fetchedItems , err 
	}

	for _ , item := range result.Items{
		newItem := Review_item{}

		_ = dynamodbattribute.UnmarshalMap(item , &newItem)
		fetchedItems = append(fetchedItems , newItem)

	}

	return  fetchedItems , nil 


}

func GetMovies()([]Movie_item , error){
	svc := createSession()
	fetchedItems := []Movie_item{}

	//  create the filter expression 
	// filterExp := expression.Name("movie_id").Equal(expression.Value(movie_id))
	// create the projection, i.e, the attributes that we want 
	// projection := expression.NamesList((expression.Name("review_id")), expression.Name("movie_id"), 
	// 	expression.Name("username") ,expression.Name("headline"), expression.Name("rating"),
	// 	expression.Name("review"),
	// ) 

	// expr , err := expression.NewBuilder().Build()

	// if err != nil{
	// 	fmt.Println("Go error building expression")
	// 	return fetchedItems , err
	// }

	//  build the query input parameters
	params := &dynamodb.ScanInput{
		// ExpressionAttributeNames: expr.Names(),
		// ExpressionAttributeValues: expr.Values(),
		// FilterExpression:        expr.Filter(),
		// ProjectionExpression:    expr.Projection(),
		TableName:   aws.String(MOVIE_TABLE),
	}

	result , err := svc.Scan(params)
	if err != nil{
		return fetchedItems , err 
	}

	for _ , item := range result.Items{
		newItem := Movie_item{}

		_ = dynamodbattribute.UnmarshalMap(item , &newItem)
		fetchedItems = append(fetchedItems , newItem)

	}

	return  fetchedItems , nil 	


}