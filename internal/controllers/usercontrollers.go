package controllers

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"nagarjuna2323/books_api/internal/middlewares/authentication"
	hash "nagarjuna2323/books_api/internal/middlewares/hashpassword"
	L "nagarjuna2323/books_api/internal/middlewares/logger"
	mdl "nagarjuna2323/books_api/internal/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var usersCSV = "users.csv"

// SIGNUP CONTROLLER

func SignUpService(c *gin.Context) {
	var newUser mdl.User

	// Parse request JSON body into newUser struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Open or create the users CSV file
	file, err := os.OpenFile(usersCSV, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		L.BKSLog("E", "Error opening users file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	defer file.Close()

	// Check if the user already exists
	if userExists(newUser.Email, file) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Write headers if the file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		L.BKSLog("E", "Error getting file info", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	if fileInfo.Size() == 0 {
		writer := csv.NewWriter(file)
		headers := []string{"Email", "Password", "UserType"}
		if err := writer.Write(headers); err != nil {
			L.BKSLog("E", "Error writing headers", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		writer.Flush()
	}

	// Hash the password before storing it
	hashedPassword, err := hash.HashPassword(newUser.Password)
	if err != nil {
		L.BKSLog("E", "Error hashing password", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user entry in CSV format
	userRecord := []string{
		newUser.Email,
		hashedPassword,
		newUser.UserType,
	}

	// Write the user record to the CSV file
	writer := csv.NewWriter(file)
	// defer writer.Flush()

	if err := writer.Write(userRecord); err != nil {
		L.BKSLog("E", "Error writing user record to file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func userExists(email string, file *os.File) bool {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false // Assume user does not exist if unable to read CSV
	}
	for _, record := range records {
		if record[0] == email {
			return true // User already exists
		}
	}
	return false
}

// LOGIN CONTROLLER

func LogInService(c *gin.Context) {
	var signInReq mdl.SignInRequest

	// Parse request JSON body into signInReq struct
	if err := c.ShouldBindJSON(&signInReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Open the users CSV file
	file, err := os.Open(usersCSV)
	if err != nil {
		L.BKSLog("E", "Error opening users file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all user records from the CSV file
	userRecords, err := reader.ReadAll()
	if err != nil {
		L.BKSLog("E", "Error reading user records from file", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Find the user with the matching email
	var userRecord []string
	for _, record := range userRecords {
		if record[0] == signInReq.Email {
			userRecord = record
			break
		}
	}

	if len(userRecord) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the hashed password
	if err := hash.ComparePasswords(userRecord[1], signInReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := authentication.GenerateToken(mdl.User{
		Email:    userRecord[0],
		UserType: userRecord[2],
		// Add other user details as needed
	})
	if err != nil {
		L.BKSLog("E", "Error generating token", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// HOME CONTROLLER

func Home(c *gin.Context) {
	// Check user type (for demonstration purpose, assuming userType is extracted from JWT token)
	userType := c.GetString("userType") // Assuming userType is extracted from JWT token

	var (
		regularUserBooks []string
		adminUserBooks   []string
		err              error
	)

	// Read regular user books
	regularUserBooks, err = readBooksFromFile("regular_user.csv")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read regular user books"})
		return
	}

	// For admin users, read admin user books
	if userType == "admin" {
		adminUserBooks, err = readBooksFromFile("admin_user.csv")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read admin user books"})
			return
		}
	}

	// Combine regular user and admin user books
	allBooks := regularUserBooks
	if userType == "admin" {
		allBooks = append(allBooks, adminUserBooks...)
	}

	// Return the list of all books in the API response
	c.JSON(http.StatusOK, gin.H{"books": allBooks})
}

func readBooksFromFile(fileName string) ([]string, error) {
	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Extract book names from records
	var books []string
	for _, record := range records {
		books = append(books, record[0]) // Assuming the book name is in the first column of the CSV
	}

	return books, nil
}

// ADDBOOK CONTROLLER

func AddBook(c *gin.Context) {
	// Check user type (for demonstration purpose, assuming userType is extracted from JWT token)
	userType := c.GetString("userType") // Assuming userType is extracted from JWT token

	// Ensure only admin users can access this endpoint
	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can access this endpoint"})
		return
	}

	// Parse parameters
	bookName := c.PostForm("bookName")
	author := c.PostForm("author")
	publicationYearStr := c.PostForm("publicationYear")

	// Validate parameters
	if bookName == "" || author == "" || publicationYearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book Name, Author, and Publication Year are required"})
		return
	}

	// Check if publicationYearStr is a valid integer
	publicationYear, err := strconv.Atoi(publicationYearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Publication Year must be a number"})
		return
	}

	// Check if publication year is a valid year (not in future)
	currentYear := time.Now().Year()
	if publicationYear < 0 || publicationYear > currentYear {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Publication Year"})
		return
	}

	// Add the new book to regular_user.csv
	err = addBookToCSV("regular_user.csv", []string{bookName, author, publicationYearStr})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Book added successfully"})
}

func addBookToCSV(fileName string, book []string) error {
	// Open the CSV file with append mode
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if the file is empty (to determine if headers need to be written)
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// If the file is empty, write the headers
	if fileInfo.Size() == 0 {
		headers := []string{"Book Name", "Author", "Publication Year"}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}

	// Write the book record to the CSV file
	if err := writer.Write(book); err != nil {
		return err
	}

	// Flush and return any errors that occur during writing
	writer.Flush()
	return writer.Error()
}

// DELETE BOOK CONTROLLER

func DeleteBook(c *gin.Context) {
	// Check user type (for demonstration purpose, assuming userType is extracted from JWT token)
	userType := c.GetString("userType") // Assuming userType is extracted from JWT token

	// Ensure only admin users can access this endpoint
	if userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can access this endpoint"})
		return
	}

	// Parse parameter
	bookName := c.Query("bookName")
	if bookName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book Name parameter is required"})
		return
	}

	// Delete the book from regularUser.csv
	err := deleteBookFromCSV("regular_user.csv", bookName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func deleteBookFromCSV(fileName string, bookName string) error {
	// Open the CSV file
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Create a new CSV writer
	file, err = os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write records back to CSV except the one to be deleted
	var updatedRecords [][]string
	for _, record := range records {
		if strings.EqualFold(record[0], bookName) {
			continue // Skip the record to be deleted
		}
		updatedRecords = append(updatedRecords, record)
	}

	// Write updated records to CSV file
	if err := writer.WriteAll(updatedRecords); err != nil {
		return err
	}

	return nil
}
