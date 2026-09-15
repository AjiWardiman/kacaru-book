// Package router bertugas menghubungkan setiap URL endpoint ke handler yang bertanggung jawab menanganinya.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "tourism-booking/docs"
	"tourism-booking/internal/handler"
	"tourism-booking/internal/middleware"
	"tourism-booking/internal/model"
)

// Dependencies mengumpulkan semua handler yang dibutuhkan router.
type Dependencies struct {
	AuthHandler         *handler.AuthHandler
	DestinationHandler  *handler.DestinationHandler
	SlotHandler         *handler.SlotHandler
	BookingHandler      *handler.BookingHandler
	PaymentHandler      *handler.PaymentHandler
	OTPHandler          *handler.OTPHandler
	SettingHandler      *handler.SettingHandler
	UploadHandler       *handler.UploadHandler
	NotificationHandler *handler.NotificationHandler
	UserHandler         *handler.UserHandler
	FeatureHandler      *handler.FeatureHandler
	JWTSecret           string
}

// New menyusun semua route/endpoint aplikasi dan mengembalikan *gin.Engine.
func New(deps Dependencies) *gin.Engine {
	r := gin.Default()

	// Pasang middleware CORS
	r.Use(middleware.CORS())

	// ===== Endpoint Swagger UI =====
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Endpoint untuk melayani file statis yang diunggah
	r.Static("/uploads", "./uploads")

	// Endpoint untuk mengecek apakah server hidup
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "layanan berjalan normal"})
	})

	// ===== Endpoint AUTH (publik) =====
	auth := r.Group("/auth")
	auth.POST("/register", deps.AuthHandler.Register)
	auth.POST("/login", deps.AuthHandler.Login)
	auth.POST("/send-otp", deps.AuthHandler.SendOTP)

	// ===== Endpoint PUBLIK untuk melihat-lihat destinasi wisata =====
	r.GET("/destinations", deps.DestinationHandler.List)
	r.GET("/destinations/:id", deps.DestinationHandler.Get)
	r.GET("/destinations/:id/slots", deps.SlotHandler.ListByDestination)

	// ===== Endpoint FITUR AI & WEATHER =====
	r.GET("/weather", deps.FeatureHandler.GetWeather)
	r.POST("/recommendations", deps.FeatureHandler.GetAIRecommendations)

	// ===== Endpoint NOTIFIKASI & SSE STREAM (wajib login) =====
	notifications := r.Group("/notifications")
	notifications.Use(middleware.Auth(deps.JWTSecret))
	notifications.GET("", deps.NotificationHandler.List)
	notifications.PATCH("/:id/read", deps.NotificationHandler.MarkAsRead)
	notifications.PATCH("/read-all", deps.NotificationHandler.MarkAllAsRead)
	notifications.GET("/stream", deps.NotificationHandler.Stream)

	// ===== Endpoint BOOKING (wajib login) =====
	bookings := r.Group("/bookings")
	bookings.Use(middleware.Auth(deps.JWTSecret))
	bookings.POST("", deps.BookingHandler.Create)
	bookings.GET("/my-bookings", deps.BookingHandler.ListMyBookings)
	bookings.POST("/:id/pay", deps.PaymentHandler.CreatePayment)

	// ===== Endpoint USER PROFILE & SETTINGS (wajib login) =====
	usersGroup := r.Group("/users")
	usersGroup.Use(middleware.Auth(deps.JWTSecret))
	usersGroup.PUT("/profile", deps.UserHandler.UpdateSelfProfile)
	usersGroup.PUT("/password", deps.UserHandler.UpdateSelfPassword)

	// ===== Endpoint OTP (wajib login) =====
	otp := r.Group("/otp")
	otp.Use(middleware.Auth(deps.JWTSecret))
	otp.POST("/request", deps.OTPHandler.RequestOTP)
	otp.POST("/verify", deps.OTPHandler.VerifyOTP)

	// ===== Endpoint UPLOAD FILE UMUM (wajib login) =====
	r.POST("/upload", middleware.Auth(deps.JWTSecret), deps.UploadHandler.Upload)

	// ===== Endpoint WEBHOOK Xendit =====
	r.POST("/webhooks/xendit", deps.PaymentHandler.XenditWebhook)

	// ===== Endpoint ADMIN (wajib login DAN role harus admin) =====
	admin := r.Group("/admin")
	admin.Use(middleware.Auth(deps.JWTSecret), middleware.RequireRole(string(model.RoleAdmin)))
	admin.POST("/destinations", deps.DestinationHandler.Create)
	admin.POST("/destinations/upload", deps.UploadHandler.UploadDestination)
	admin.PUT("/destinations/:id", deps.DestinationHandler.Update)
	admin.DELETE("/destinations/:id", deps.DestinationHandler.Delete)
	admin.POST("/destinations/:id/slots", deps.SlotHandler.Create)
	admin.GET("/bookings", deps.BookingHandler.ListAllBookings)

	// ----- Sub-grup manajemen pengguna (users), khusus admin -----
	admin.GET("/users", deps.UserHandler.List)
	admin.GET("/users/:id", deps.UserHandler.Get)
	admin.POST("/users", deps.UserHandler.Create)
	admin.PUT("/users/:id", deps.UserHandler.Update)
	admin.DELETE("/users/:id", deps.UserHandler.Delete)

	// ----- Sub-grup pengaturan aplikasi (settings), khusus admin -----
	admin.GET("/settings", deps.SettingHandler.List)
	admin.GET("/settings/:key", deps.SettingHandler.Get)
	admin.PUT("/settings/:key", deps.SettingHandler.Update)

	// ----- Sub-grup laporan analytics & export (reports), khusus admin -----
	admin.GET("/reports/summary", deps.FeatureHandler.GetReportSummary)
	admin.GET("/reports/export/excel", deps.FeatureHandler.ExportReportExcel)
	admin.GET("/reports/export/pdf", deps.FeatureHandler.ExportReportPDF)

	return r
}
