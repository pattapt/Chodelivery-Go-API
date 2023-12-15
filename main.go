package main

import (
	"fmt"
	"net/http"

	MerchantV1Main "gobasic/api/merchant"
	MerchantV1Auth "gobasic/api/merchant/auth"
	MerchantV1Store "gobasic/api/merchant/store"
	StoreV1Main "gobasic/api/store"
	StoreV1Address "gobasic/api/store/address"
	StoreV1Auth "gobasic/api/store/auth"
	StoreV1Cart "gobasic/api/store/cart"
	StoreV1Chat "gobasic/api/store/chat"
	StoreV1Merchant "gobasic/api/store/merchant"
	StoreV1Profile "gobasic/api/store/profile"
	StoreV1Transaction "gobasic/api/store/transaction"

	MerchantV1Chat "gobasic/api/merchant/chat"
	MerchantV1Product "gobasic/api/merchant/product"
	MerchantV1Transaction "gobasic/api/merchant/transaction"
	MerchantV1Util "gobasic/api/merchant/util"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Serve static files from the "cdn/" directory
	r.Static("/cdn", "./cdn")

	// Create routers for store and merchant paths
	r_store := r.Group("/api/store/v1")
	r_merc := r.Group("/api/merchant/v1")

	// Define routes for store and merchant
	// STORE SECTION
	r_store.POST("/auth/registerAccount", StoreV1Auth.RegisterAccount)     //FINISH
	r_store.POST("/auth/Login", StoreV1Auth.LoginAccount)                  //FINISH
	r_store.POST("/auth/InvokeAccessToken", StoreV1Auth.InvokeAccessToken) //FINISH

	r_store.GET("/profile/GetProfile", StoreV1Main.AccessTokenMiddleware, StoreV1Profile.GetProfile) //FINISH

	r_store.GET("/merchant/list", StoreV1Main.AccessTokenMiddleware, StoreV1Merchant.GetMerchantList)                                                                                              // FINISH                                                                                   //FINISH
	r_store.GET("/merchant/:merchantToken/info", StoreV1Main.AccessTokenMiddleware, StoreV1Main.GetMerchantId, StoreV1Merchant.GetMerchantInfo)                                                    //FINISH
	r_store.GET("/merchant/:merchantToken/category/all", StoreV1Main.AccessTokenMiddleware, StoreV1Main.GetMerchantId, StoreV1Merchant.GetAllCategories)                                           //FINISH
	r_store.GET("/merchant/:merchantToken/category/GetProductsFromCategory/:categoryToken", StoreV1Main.AccessTokenMiddleware, StoreV1Main.GetMerchantId, StoreV1Merchant.GetProductsByCategoryId) //FINISH
	r_store.GET("/merchant/:merchantToken/category/SearchProducts", StoreV1Main.AccessTokenMiddleware, StoreV1Main.GetMerchantId, StoreV1Merchant.SearchProducts)                                  //FINISH
	r_store.GET("/merchant/:merchantToken/product/:ProductToken/Detail", StoreV1Main.AccessTokenMiddleware, StoreV1Main.GetMerchantId, StoreV1Merchant.GetProductInfo)                             //FINISH

	r_store.POST("/cart/add/:ProductToken", StoreV1Main.AccessTokenMiddleware, StoreV1Cart.AddCart)      // FINISH
	r_store.POST("/cart/update/:CartToken", StoreV1Main.AccessTokenMiddleware, StoreV1Cart.UpdateCart)   // FINISH
	r_store.GET("/cart/list", StoreV1Main.AccessTokenMiddleware, StoreV1Cart.GetCart)                    // FINISH
	r_store.GET("/cart/:CartToken/Detail", StoreV1Main.AccessTokenMiddleware, StoreV1Cart.GetCartDetail) // FINISH
	r_store.POST("/cart/checkOut", StoreV1Main.AccessTokenMiddleware, StoreV1Cart.CheckOutCart)          // FINISH

	r_store.GET("/profile/address/list", StoreV1Main.AccessTokenMiddleware, StoreV1Address.GetAddress)                     // FINISH
	r_store.GET("/profile/address/:AddressToken/Detail", StoreV1Main.AccessTokenMiddleware, StoreV1Address.GetAddressInfo) // FINISH
	r_store.POST("/profile/address/Create", StoreV1Main.AccessTokenMiddleware, StoreV1Address.CreateAddress)               // FINISH
	r_store.POST("/profile/address/:AddressToken/Update", StoreV1Main.AccessTokenMiddleware, StoreV1Address.UpdateAddress) // FINISH

	r_store.GET("/transaction/list", StoreV1Main.AccessTokenMiddleware, StoreV1Transaction.GetTransaction)                         // FINISH
	r_store.GET("/transaction/:TransactionToken/Detail", StoreV1Main.AccessTokenMiddleware, StoreV1Transaction.GetTransactionInfo) // FINISH
	r_store.GET("/transaction/:TransactionToken/QR", StoreV1Main.AccessTokenMiddleware, StoreV1Transaction.GetQRCode)              // FINISH

	r_store.GET("/chat/:ChatToken/Message", StoreV1Main.AccessTokenMiddleware, StoreV1Chat.ValidateChat, StoreV1Chat.GetMessage)   // FINISH
	r_store.POST("/chat/:ChatToken/Send", StoreV1Main.AccessTokenMiddleware, StoreV1Chat.ValidateChat, StoreV1Chat.SendMessage)    // FINISH
	r_store.POST("/chat/:ChatToken/SendImage", StoreV1Main.AccessTokenMiddleware, StoreV1Chat.ValidateChat, StoreV1Chat.SendImage) // FINISH
	r_store.GET("/chat/:ChatToken/Image/:ImageToken", StoreV1Main.AccessTokenMiddleware, StoreV1Chat.ValidateChat, StoreV1Chat.ViewImage)
	r_store.GET("/Util/GetProvince", StoreV1Main.AccessTokenMiddleware, MerchantV1Util.GetDistrict)

	// MERCHANT SECTION
	r_merc.POST("/auth/registerAccount", MerchantV1Auth.RegisterAccount)     // FINISH
	r_merc.POST("/auth/Login", MerchantV1Auth.LoginAccount)                  // FINISH
	r_merc.POST("/auth/InvokeAccessToken", MerchantV1Auth.InvokeAccessToken) // FINISH

	r_merc.GET("/store/GetProfile", MerchantV1Main.AccessTokenMiddleware, MerchantV1Store.GetProfile)                  // FINISH
	r_merc.POST("/store/CreateMerchant", MerchantV1Main.AccessTokenMiddleware, MerchantV1Store.Createmerchant)         // FINISH
	r_merc.GET("/store/GetMerchant", MerchantV1Main.AccessTokenMiddleware, MerchantV1Store.GetAllMerchant)             // FINISH
	r_merc.GET("/store/:MerchantUUID/Detail", MerchantV1Main.AccessTokenMiddleware, MerchantV1Store.GetMerchantDetail) // FINISH
	r_merc.POST("/store/:MerchantUUID/Update", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Store.UpdateMerchant)
	r_merc.POST("/store/:MerchantUUID/UpdateImage", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Store.UpdateMerchantImage)
	r_merc.GET("/store/:MerchantUUID/Profile", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Store.GetMerchantDetail) // FINISH

	r_merc.GET("/store/:MerchantUUID/Summary", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Store.GetSummary)

	r_merc.GET("/store/:MerchantUUID/Orders", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Transaction.GetTransaction)                          // FINISH
	r_merc.GET("/store/:MerchantUUID/Orders/:OrderToken/Detail", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Transaction.GetTransactionDetail) // FINISH
	r_merc.POST("/store/:MerchantUUID/Orders/:OrderToken/Update", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Transaction.UpdateTransaction)   // FINISH

	r_merc.GET("/store/:MerchantUUID/Products", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Product.GetProduct)                                    // FINISH
	r_merc.POST("/store/:MerchantUUID/Products/Create", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Product.CreateProduct)                         // FINISH
	r_merc.POST("/store/:MerchantUUID/Products/:ProductToken/Update", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Product.UpdateProduct)           // FINISH
	r_merc.POST("/store/:MerchantUUID/Products/:ProductToken/UploadImage", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Product.UploadProductImage) // FINISH
	r_merc.GET("/store/:MerchantUUID/Products/:ProductToken/Detail", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Product.GetProductInfo)           // FINISH

	r_merc.GET("/store/:MerchantUUID/Category", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid)                        // SKIP
	r_merc.POST("/store/:MerchantUUID/Category/Create", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid)                // SKIP
	r_merc.POST("/store/:MerchantUUID/Category/:CategoryToken/Update", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid) // SKIP
	r_merc.GET("/store/:MerchantUUID/Category/:CategoryToken/Detail", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid)  // SKIP

	r_merc.GET("/store/:MerchantUUID/chat", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Chat.GetAllChat)                                                          // FINISH
	r_merc.GET("/store/:MerchantUUID/chat/:ChatToken/Message", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Chat.ValidateChat, MerchantV1Chat.GetMessage)          // FINISH
	r_merc.POST("/store/:MerchantUUID/chat/:ChatToken/Send", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Chat.ValidateChat, MerchantV1Chat.SendMessage)           // FINISH
	r_merc.POST("/store/:MerchantUUID/chat/:ChatToken/SendImage", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Chat.ValidateChat, MerchantV1Chat.SendImage)        // FINISH
	r_merc.GET("/store/:MerchantUUID/chat/:ChatToken/Image/:ImageToken", MerchantV1Main.AccessTokenMiddleware, MerchantV1Main.CheckMerchantValid, MerchantV1Chat.ValidateChat, MerchantV1Chat.ViewImage) // FINISH

	r_merc.GET("/Util/GetProvince", MerchantV1Main.AccessTokenMiddleware, MerchantV1Util.GetDistrict) // FINISH

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Wrong api path",
		})
	})

	fmt.Println("Server is running on :80")
	r.Run(":80")
}
