package main

import (
	paymentpb "github.com/ReStorePUC/protobucket/payment"
	pb "github.com/ReStorePUC/protobucket/user"
	"github.com/gin-gonic/gin"
	"github.com/restore/shop/config"
	"github.com/restore/shop/controller"
	"github.com/restore/shop/handler"
	"github.com/restore/shop/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	config.Init()
	dbCfg := config.NewDBConfig()

	db, err := repository.Init(dbCfg)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("user:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserClient(conn)

	paymentConn, err := grpc.Dial("payment:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer paymentConn.Close()
	pc := paymentpb.NewPaymentClient(paymentConn)

	sRepo := repository.NewShop(db)
	sController := controller.NewShop(sRepo, c, pc)
	sHandler := handler.NewShop(sController)

	router := gin.Default()

	router.POST("/private/request", sHandler.CreateRequest)
	router.PUT("/private/request/:id", sHandler.UpdateRequest)
	router.GET("/private/request/search/:storeID", sHandler.SearchRequest)

	router.POST("/private/payment", sHandler.CreatePayment)
	router.PUT("/private/payment/:id", sHandler.UpdatePayment)
	router.GET("/private/payment/store/:storeID", sHandler.GetPayments)
	router.GET("/private/payment/search", sHandler.SearchPayments)

	router.Run(":8080")
}
