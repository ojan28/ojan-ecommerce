package controllers

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"ecommerce/database"
	"ecommerce/models"
	generate "ecommerce/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string{

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, givenPassword string)(bool,string){
	err :=	bcrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword))
	valid := true
	msg := ""

	if err!= nil {
        valid = false
        msg = "email or password is not valid"
    }
	return valid, msg
}

func SignUp() gin.HandlerFunc{

	return func(c *gin.Context){
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&user); err!=nil{
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

		validationErr := Validate.Struct(user)
		if validationErr!= nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
            return
        }

		UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
            c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
        }

		count, err = UserCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})	
		defer cancel()
		if err!= nil {
            log.Panic(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
            return
        }

		if count > 0 {
            c.JSON(http.StatusConflict, gin.H{"error": "phone number already used"})
            return
        }

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		User.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.RefreshToken = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address,0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr!= nil {
            log.Panic(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
            return
        }
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully signed in")
	}
}

func Login() gin.HandlerFunc{

	return func(c *gin.Context){
		context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var founduser models.User
		var user models.User
        if err := c.ShouldBindJSON(&user); err!= nil{
            c.JSON(http.StatusBadRequest, gin.H{"error": err})
            return
        }

		UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

		defer cancel()
		if!PasswordIsValid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			fmt.Println(msg)
            return
        }

		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusFound, founduser)
	}

}

func ProductViewerAdmin() gin.HandlerFunc{
	return func(c*gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InserOne(ctx, products
		if anyerr!= nil {
            log.Println(anyerr)
            c.JSON(http.StatusInternalServerError,  "something went wrong, please try after some time")
            return
        })
		defer cancel()
		c.JSON(http.StatusCreated, "product created successfully")
		}
	}

}

func SearchProduct() gin.HandlerFunc{
	return func(c *gin.Context) {

		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err: ProductCollection.Find(ctx, bson.D{{}})
		if err!= nil {
            log.Println(err)
            c.JSON(http.StatusInternalServerError,  "something went wrong, please try after some time")
            return
        }
		err = cursor.All(ctx, &productlist)
		if err!= nil {
            log.Println(err)
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil{
			log.Println(err)
            c.IntendedJSON(400, "invalid")
            return
		}
		defer cancel()
		c.IntendedJSON(200, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc{
	return func(c*gin.Context){
		var searchproducts []models.Product
		queryParam := c.Query("name")

	
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid search index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})
		if err!= nil {
			c.IntendedJSON(404,"something went wrong while fetching the data")
			return
		}

		searchquerydb.All(ctx, &searchproducts)
		if err!= nil {
			log.Println(err)
			c.IntendedJSON(400, "Invalid"
			return
		}
		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IntendedJSON(400, "Invalid Request")
			return
			}	

		defer cancel()
		c.IntendedJSON(200, searchproducts)	
		}
}