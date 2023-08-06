package controller

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "errorr occured while listing order Item"})

		}
		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrderItems)

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderId = c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError.gin.H{"error": "error occured while listing item by order"})
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (orderItems []primitive.M, err error) {

}

func GetOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background, 100*time.Second)
		defer cancel()

		orderItemId = c.Param("order_item_id")
		var orderItem = models.OrderItem
		err := Collection.FindOne(ctx, bson.M{"order_item_id", orderItemId}).decode(orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errpor": "error while fetching item"})
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background, 100*time.Second)

		var orderItemPack OrderItemPack
		var orderItem models.OrderItem
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		order.Order_Date = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemToBeInserted := []interface{}{}
		order.Table_id = OrderItemPack.Table_id
		order_id = OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.orderItems {
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItemPack)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItemPack.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItemPack.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItemPack.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemToBeInserted = append(orderItemToBeInserted, orderItem)
		}

		insertedOrderItem, err := orderItemCollection.InsertMany(ctx, orderItemToBeInserted)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, insertedOrderItem)

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem
		var orderItemId = c.Param("order_item_id")

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		// if orderItem.Order_id != nil {
		// 	err:= orderCollection.FindOne(ctx, bson.M{"order_id": orderItem.Order_id}).Decode(&order)
		// 	defer cancel()
		// 	if err!= nil{
		// 		msg := fmt.Sprintf("message:Menu was not found")
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
		// 		return
		// 	}
		// 	updateObj = append(updateObj, bson.E{"table_id", order.Table_id})
		// }

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", orderItem.Unit_price})
		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})
		
		upsert := true

		filter := bson.M{"order_item_id":orderItemId}
		opt := options.UpdateOptions{
			Upsert:&upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "OrderItem update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
