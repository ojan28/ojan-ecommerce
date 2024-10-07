package controllers 

import (

	"time"
	"context"
	"errors"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application{
	return &Application{
        prodCollection: prodCollection,
        userCollection: userCollection,
    }
}

func (app *Application)AddToCart() gin.HandlerFunc{
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
            log.Println("Product id is empty"
           
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User id is empty"
            _ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
            return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil {
            log.Println(err)
            _ = c.AbortWithError(http.StatusInternalServerError)
            return
        }

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, product, app.userCollection, productID, userQueryID)
		if err!= nil {
			c.IntendedJSON(http.StatusInternalServerError, err)
		}
		c.IntendedJSON(http.StatusOK, gin.H{"message": "Product succesfully added to cart"})
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc{

	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
            log.Println("Product id is empty"
           
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User id is empty"
            _ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
            return
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil {
            log.Println(err)
            _ = c.AbortWithError(http.StatusInternalServerError)
            return
        }

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil {
			c.IntendedJSON(http.StatusInternalServer, err)
			return
		}
		c.IntendedJSON(http.StatusOK, gin.H{"message": "Product succesfully removed from cart"})
	}
}

func GetItemFromCart() gin.HandlerFunc
{
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
            c.Abort()
            return
        }

		user_id, _ := primitve.ObjectIDFromHex(user_id) 
		
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: user_id }}).Decode(&filledcart)
		
		if err != nil {
			log.Println(err)
			c.IntendedJSON(500, "not found")
			return
		}

		filter_match := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key:"_id", Value:user_id}}}}
		unwind := bson.D{{Key:"$unwind", Value:bson.D{primitive.E{Key:"path", Value:"$usercart"}}}}
		grouping := bson.D{{Key:"$group", Value:bson.D{primitive.E{Key:"_id",Value:"$_id"}, {Key:"total", Value:bson.D{primitive.E{Key:"$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err := pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing{
			c.IntendedJSON(200, json["total"])
			c.IntendedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc{
		return func(c *gin.Context) {
			userQueryID := c.Query("id")
			if userQueryID == "" {
				log.Panicln("user Id is empty")
				_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			}
			context.WithTimeout(context.Background(), 100*time.Second)
			
			defer cancel()
			database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
			if err := nil {
				c.IntendedJSON(http.StatusInternalServer, err)
                return
			}
			c.IntendedJSON(http.StatusOK, gin.H{"message": "Successfully bought the items from cart"})
		}
}

func (app *Application)InstantBuy() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
            log.Println("Product id is empty"
           
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("User id is empty"
            _ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
            return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil {
            log.Println(err)
            _ = c.AbortWithError(http.StatusInternalServerError)
            return
        }

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
        if err!= nil {
            c.IntendedJSON(http.StatusInternalServer, err)
            return
        }
        c.IntendedJSON(http.StatusOK, gin.H{"message": "Successfully placed the order "})
	}
}