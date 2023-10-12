package main

import (
	paymentpb "github.com/ReStorePUC/protobucket/payment"
	productpb "github.com/ReStorePUC/protobucket/product"
	pb "github.com/ReStorePUC/protobucket/user"
	"github.com/gin-contrib/cors"
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

	productConn, err := grpc.Dial("product:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer productConn.Close()
	prodC := productpb.NewProductClient(productConn)

	sRepo := repository.NewShop(db)
	sController := controller.NewShop(sRepo, c, prodC, pc)
	sHandler := handler.NewShop(sController)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowFiles:       true,
	}))

	router.POST("/private/request", sHandler.CreateRequest)
	router.PUT("/private/request/:id", sHandler.UpdateRequest)
	router.GET("/private/request/search/:storeID", sHandler.SearchRequest)
	router.GET("/private/request/profile/search/:profileID", sHandler.SearchProfileRequest)
	router.POST("/private/confirm-request/:paymentID", sHandler.ConfirmRequest)

	router.POST("/private/payment", sHandler.CreatePayment)
	router.PUT("/private/payment/:id", sHandler.UpdatePayment)
	router.GET("/private/payment/store/:storeID", sHandler.GetPayments)
	router.GET("/private/payment/search", sHandler.SearchPayments)

	router.Run(":8080")
}
