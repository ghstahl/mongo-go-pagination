package main

import (
	"context"
	"fmt"
	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Product struct
type Product struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Quantity float64            `json:"qty" bson:"qty"`
	Price    float64            `json:"price" bson:"price"`
}

func main() {
	// Establishing mongo db connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// Example for Normal Find query
	filter := bson.M{}
	var limit int64 = 10
	var page int64 = 1
	collection := client.Database("myaggregate").Collection("stocks")
	projection := bson.D{
		{"name", 1},
		{"qty", 1},
	}
	// Querying paginated data
	// If you want to do some complex sort like sort by score(weight) for full text search fields you can do it easily
	// sortValue := bson.M{
	//		"$meta" : "textScore",
	//	}
	// paginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("score", sortValue)...
	var products []Product
	paginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Sort("qty", -1).Select(projection).Filter(filter).Decode(&products).Find()
	if err != nil {
		panic(err)
	}

	// paginated data is in paginatedData.Data
	// pagination info can be accessed in  paginatedData.Pagination
	//// if you want to marshal data to your defined struct
	// print ProductList
	fmt.Printf("Norm Find Data: %+v\n", products)

	// print pagination data
	fmt.Printf("Normal find pagination info: %+v\n", paginatedData.Pagination)

	//Example for Aggregation

	//match query
	match := bson.M{"$match": bson.M{"qty": bson.M{"$gt": 10}}}

	//group query
	projectQuery := bson.M{"$project": bson.M{"_id": 1, "qty": 1}}

	// you can easily chain function and pass multiple query like here we are passing match
	// query and projection query as params in Aggregate function you cannot use filter with Aggregate
	// because you can pass filters directly through Aggregate param
	aggPaginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Aggregate(match, projectQuery)
	if err != nil {
		panic(err)
	}

	var aggProductList []Product
	for _, raw := range aggPaginatedData.Data {
		var product *Product
		if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
			aggProductList = append(aggProductList, *product)
		}

	}

	// print ProductList
	fmt.Printf("Aggregate Product List: %+v\n", aggProductList)

	// print pagination data
	fmt.Printf("Aggregate Pagination Data: %+v\n", aggPaginatedData.Pagination)
}
