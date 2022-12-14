package test

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
	"github.com/vincentvanpro/customer-web-app/internals/controller"
	"github.com/vincentvanpro/customer-web-app/internals/database"
	"github.com/vincentvanpro/customer-web-app/internals/models"
	"gorm.io/gorm"
)

type TestSuiteEnv struct {
	suite.Suite
	db *gorm.DB
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {
	database.Connect()
	suite.db = database.GetDB()
}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	database.ClearTable()
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(TestSuiteEnv))
}

func (suite *TestSuiteEnv)Test_GetCustomers_EmptyResult() {
	req, w := setGetCustomersRouter(suite.db)
	a := suite.Assert()
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
	body, err := ioutil.ReadAll(w.Body)
		if err != nil {
		a.Error(err)
	}

	actual := models.Customer{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	expected := models.Customer{}
	a.Equal(expected, actual)
}

func setGetCustomersRouter(db *gorm.DB) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	r.GET("/users", controller.GetCustomers)
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func (suite *TestSuiteEnv)Test_GetCustomerByID_OK() {
	a := suite.Assert()
	
	newCustomer, err := insertTestCustomer(suite.db)
	if err != nil {
		a.Error(err)
	}

	req, w := setGetCustomerRouter(suite.db, "/users/1")

	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
	body, err := ioutil.ReadAll(w.Body)
		if err != nil {
		a.Error(err)
	}

	actual := models.Customer{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	actual.Model = gorm.Model{}

	expected := newCustomer
	expected.Model = gorm.Model{}
	a.Equal(expected.ID, actual.ID)
	a.Equal(expected.FirstName, actual.FirstName)
	a.Equal(expected.LastName, actual.LastName)
	a.Equal(expected.Email, actual.Email)
	a.Equal(expected.Gender, actual.Gender)
	a.Equal(expected.Address, actual.Address)
	a.WithinDuration(expected.DateOfBirth, actual.DateOfBirth, 10*time.Second)
}

func setGetCustomerRouter(db *gorm.DB, url string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	r.GET("/users/:id", controller.GetCustomerById)
	req, err := http.NewRequest(http.MethodGet, "/users/1", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func insertTestCustomer(db *gorm.DB) (models.Customer, error){
	b := models.Customer{FirstName: "Test", LastName: "Test", Email: "a@email.com", DateOfBirth: time.Now(), Gender: "Male", }

	if err := db.Create(&b).Error; err != nil {
		return b, err
	}

	return b, nil
}

func (suite *TestSuiteEnv)Test_CreateCustomer_OK() {
	a := suite.Assert()
	customers := []models.Customer{}
	suite.db.Find(&customers)
	a.Equal(len(customers), 0)

	var jsonStr = []byte(`{
		"first_name": "Test",
		"last_name": "Test",
		"email": "a@email.com",
		"date_of_birth": "2002-12-12T17:04:05+02:00",
		"gender": "Male",
		"address": ""
	}`)
	req, w := setCreateCustomerRouter(suite.db, jsonStr)

	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusCreated, w.Code, "HTTP request status code error")
	suite.db.Find(&customers)
	a.Equal(len(customers), 1)
}

func (suite *TestSuiteEnv)Test_CreateCustomer_AgeLess18() {
	a := suite.Assert()

	var jsonStr = []byte(`{
		"first_name": "Test",
		"last_name": "Test",
		"email": "a@email.com",
		"date_of_birth": "2021-12-12T17:04:05+02:00",
		"gender": "Male",
		"address": ""
	}`)
	req, w := setCreateCustomerRouter(suite.db, jsonStr)

	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
}

func (suite *TestSuiteEnv)Test_CreateCustomer_AgeGreater60() {
	a := suite.Assert()

	var jsonStr = []byte(`{
		"first_name": "Test",
		"last_name": "Test",
		"email": "a@email.com",
		"date_of_birth": "1960-12-12T17:04:05+02:00",
		"gender": "Male",
		"address": ""
	}`)
	req, w := setCreateCustomerRouter(suite.db, jsonStr)

	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
}

//Custom validation (allowed from 18 till 60 years)
var ageRestrictions validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		var age = today.Year() - date.Year();
		if age <= 0 || age < 18 || age > 60 {
			return false
		}
	}
	return true
}

func setCreateCustomerRouter(db *gorm.DB, jsonData []byte) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	//Custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("ageRestrictions", ageRestrictions)
	}
	r.POST("/users", controller.CreateCustomerTest)
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func (suite *TestSuiteEnv)Test_ModifyExistingCustomer_OK() {
	a := suite.Assert()
	
	newCustomer, err := insertTestCustomer(suite.db)
	if err != nil {
		a.Error(err)
	}

	var jsonStr = []byte(`{
		"ID": 1,
		"first_name": "Test",
		"last_name": "Modified",
		"email": "b@email.com",
		"date_of_birth": "2002-12-12T17:04:05+02:00",
		"gender": "Male",
		"address": ""
	}`)

	req, w := setModifyCustomerRouter(suite.db, "/users/1", jsonStr)

	a.Equal(http.MethodPatch, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
	body, err := ioutil.ReadAll(w.Body)
		if err != nil {
		a.Error(err)
	}

	actual := models.Customer{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	actual.Model = gorm.Model{}

	expected := newCustomer
	expected.Model = gorm.Model{}
	a.Equal(expected.ID, actual.ID)
	a.Equal(expected.FirstName, actual.FirstName)
	a.NotEqual(expected.LastName, actual.LastName)
	a.NotEqual(expected.Email, actual.Email)
	a.Equal(expected.Gender, actual.Gender)
	a.Equal(expected.Address, actual.Address)
}

func (suite *TestSuiteEnv)Test_ModifyExistingCustomer_AgeLess18() {
	a := suite.Assert()
	
	_, err := insertTestCustomer(suite.db)
	if err != nil {
		a.Error(err)
	}

	
	var jsonStr = []byte(`{
		"ID": 1,
		"date_of_birth": "2021-12-12T17:04:05+02:00",
	}`)

	req, w := setModifyCustomerRouter(suite.db, "/users/1", jsonStr)

	a.Equal(http.MethodPatch, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
}

func (suite *TestSuiteEnv)Test_ModifyExistingCustomer_AgeGreater60() {
	a := suite.Assert()
	
	_, err := insertTestCustomer(suite.db)
	if err != nil {
		a.Error(err)
	}

	
	var jsonStr = []byte(`{
		"ID": 1,
		"date_of_birth": "1960-12-12T17:04:05+02:00",
	}`)

	req, w := setModifyCustomerRouter(suite.db, "/users/1", jsonStr)

	a.Equal(http.MethodPatch, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
}

func (suite *TestSuiteEnv)Test_ModifyNonExistingCustomer_Fail() {
	a := suite.Assert()
	
	_, err := insertTestCustomer(suite.db)
	if err != nil {
		a.Error(err)
	}

	
	var jsonStr = []byte(`{
		"ID": 10,
		"first_name": "Test",
	}`)

	req, w := setModifyCustomerRouter(suite.db, "/users/10", jsonStr)

	a.Equal(http.MethodPatch, req.Method, "HTTP request method error")
	a.Equal(http.StatusBadRequest, w.Code, "HTTP request status code error")
}

func setModifyCustomerRouter(db *gorm.DB, url string, jsonData []byte) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	//Custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("ageRestrictions", ageRestrictions)
	}
	r.PATCH("/users/:id", controller.ModifyCustomerTest)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}
