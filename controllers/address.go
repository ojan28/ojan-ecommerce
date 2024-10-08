package controllers

import(
 "context"
 "ecommerce/models"
 "github.com/gin-gonic/gin"
 "go.mongodb.org/mongo-driver/bson"
 "go.mongodb.org/mongo-driver/bson/primitive"
 "go.mongodb.org/mongo-driver/mongo"
)

func AddAddress()gin.HandlerFunc{
	return func (c *gin.Context){
		
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err!= nil {
            c.IntendedJSON(500, "Internal server error")
        }

		var addresses models.Address
		addresses.Address_id = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			C.IntendedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background, 100*time.Second)

		match_filter := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key:"_id", Value: address}}}}
		unwind := bson.D{{Key:"$unwind", Value:bson.D{primitive.E{Key:"path", Value:"$address"}}}}
		group := bson.D{{Key:"$group", Value: bson.D{primitive.E{Key:"_id", Value: "$address_id"}, {Key:"count", Value: bson.D{"$sum", Value: 1}}}}} 
		pointcursor , err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchfilter,unwind,group})
		if err!= nil {
            c.IntendedJSON(500, "Internal server error")
        }

		var addressinfo[]bson.M 
		pointcursor.All(ctx,&addressinfo); err != nil {
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
		count := address_no["count"]
		size = count.(int32)
		}
		if size < 2 {
          filter := bson.D{primitive.E{Key:"_id", Value: address}}
		  update := bson.D{{ Key:"$push", Value: bson.D{primitive.E{Key:"address", Value: addresses}}}}
		  _, err := UserCollection.UpdateOne(ctx, filter, update)
		  if err != nil {
			fmt.Println(err)
		  }
        }else{
			c.IntendedJSON (400, "Not Allowed")
		}

		defer cancel()
		ctx.Done()
	}

func EditHomeAddress()gin.HandlerFunc
{
	return func (c *gin.Context){
		user_id := c.Query("id")
        if user_id == "" {
            c.Header("Content-Type", "application/json")
            c.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
            c.Abort()
            return
        }
		user_id, err := primitive.ObjectIDFromHex(user_id) 
		if err != nil {
			c.IntendedJSON(500, "Internal server error"
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IntendedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel context = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{ Key:"$set", Value:bson.D{primitive.E{Key:"address.0.house_name", Value: editaddress.House}, {Key:"address.0.street_name", Value: editaddress.Street},{Key:"address.0.city", Value: editaddress.City},{Key:"address.0.pin_code", Value:editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err!= nil {
            c.JSON(500, "Something went wrong")
            return
        }
		defer cancel()
		ctx.Done()
		c.IntendedJSON(200, "Home Address Updated Successfully")
	}


}

func EditWorkAddress()gin.HandlerFunc{
	return func (c *gin.Context)
	{
		user_id := c.Query("id")
        if user_id == "" {
            c.Header("Content-Type", "application/json")
            c.JSON(http.StatusNotFound, gin.H{"error": "invalid"})
            c.Abort()
            return
        }
		user_id, err := primitive.ObjectIDFromHex(user_id) 
		if err != nil {
			c.IntendedJSON(500, "Internal server error"
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IntendedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel context = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{ Key:"$set", Value:bson.D{primitive.E{Key:"address.1.house_name", Value: editaddress.House}, {Key:"address.1.street_name", Value: editaddress.Street},{Key:"address.1.city", Value: editaddress.City},{Key:"address.1.pin_code", Value:editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err!= nil {
            c.JSON(500, "Something went wrong")
            return
        }
        defer cancel()
        ctx.Done()
        c.IntendedJSON(200, "Work Address Updated Successfully")
	}
}

func DeleteAddress()gin.HandlerFunc{

	return func (c *gin.Context){
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
            c.JSON(http.StatusNotFound, gin.H{"error": "invalid search index"})
			c.Abort()
            return
        }

		addresses := make([]models.Address, 0)
		user_id, err := primitive.ObjectIDFromHex(user_id) 
		if err != nil {
			c.IntendedJSON(500, "Internal server error"
		}
		var ctx, cancel context = context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
		update := bson.D{{ Key: "$set", Value: bson.D{primitive.E{Key:"address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err!= nil {
            c.JSON(404, "Wrong Command")
            return
        }
		defer cancel()
		ctx.Done()
		c.IntendedJSON(200, "Successfully Deleted"
	}

}
