package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	secretKey = []byte("secret_key") // Kunci rahasia untuk JWT
)

// User merupakan struktur untuk tabel User
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Age       uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Photo merupakan struktur untuk tabel Photo
type Photo struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Caption   string
	PhotoURL  string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Comment merupakan struktur untuk tabel Comment
type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	PhotoID   uint      `gorm:"not null"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// SocialMedia merupakan struktur untuk tabel SocialMedia
type SocialMedia struct {
	ID             uint      `gorm:"primaryKey"`
	Name           string    `gorm:"not null"`
	SocialMediaURL string    `gorm:"not null"`
	UserID         uint      `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

// CustomClaims merupakan struktur yang akan disimpan di JWT
type CustomClaims struct {
	UserID uint
	jwt.StandardClaims
}

func main() {
	// Membuat koneksi database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrasi tabel
	db.AutoMigrate(&User{}, &Photo{}, &Comment{}, &SocialMedia{})

	// Membuat router Gin
	router := gin.Default()

	// Endpoint untuk registrasi pengguna
	router.POST("/register", register)

	// Endpoint untuk login
	router.POST("/login", login)

	// Middleware untuk memeriksa token
	router.Use(authMiddleware())

	// Endpoint untuk mengakses data tabel Photo
	router.GET("/photos", getPhotos)

	// Endpoint untuk mengakses data tabel Comment
	router.GET("/comments", getComments)

	// Endpoint untuk mengakses data tabel SocialMedia
	router.GET("/socialmedia", getSocialMedia)

	// Endpoint untuk menambah komentar
	router.POST("/comments", addComment)

	// Endpoint untuk menambah data SocialMedia
	router.POST("/socialmedia", addSocialMedia)

	// Jalankan server
	router.Run(":8080")
}

// Middleware untuk memeriksa token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Validasi token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set userID dalam context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

// Handler untuk registrasi pengguna
func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validasi umur
	if user.Age <= 8 {
		c.JSON(400, gin.H{"error": "Age must be greater than 8"})
		return
	}

	// Enkripsi password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)

	// Simpan pengguna ke database
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}

// Handler untuk login
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Cari pengguna berdasarkan email
	db := c.MustGet("db").(*gorm.DB)
	var existingUser User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Periksa password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserID: existingUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token berlaku selama 1 hari
		},
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": tokenString})
}

// Handler untuk mendapatkan semua foto
func getPhotos(c *gin.Context) {
	var photos []Photo
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&photos)
	c.JSON(200, photos)
}

// Handler untuk mendapatkan semua komentar
func getComments(c *gin.Context) {
	var comments []Comment
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&comments)
	c.JSON(200, comments)
}

// Handler untuk mendapatkan semua media sosial
func getSocialMedia(c *gin.Context) {
	var socialMedia []SocialMedia
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&socialMedia)
	c.JSON(200, socialMedia)
}

// Handler untuk menambah komentar
func addComment(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	comment.UserID = userID.(uint)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(201, gin.H{"message": "Comment added successfully"})
}

// Handler untuk menambah media sosial
func addSocialMedia(c *gin.Context) {
	var socialMedia SocialMedia
	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	socialMedia.UserID = userID.(uint)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&socialMedia).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create social media"})
		return
	}

	c.JSON(201, gin.H{"message": "Social media added successfully"})
}
