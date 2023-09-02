package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"example/unit-test-hello-world/config"
	mahasiswa "example/unit-test-hello-world/controllers"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func SetUpDB() (*sql.DB, error) {
	var (
		DB_NAME     = os.Getenv("DB_NAME")
		DB_USER     = os.Getenv("DB_USER")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_HOST     = os.Getenv("DB_HOST")
		DB_PORT     = os.Getenv("DB_PORT")
	)

	connStr := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	if err = db.Ping(); err != nil {
		return db, err
	}

	return db, err
}

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("ERROR: Could not load .env")
		os.Exit(1)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestHomepage(t *testing.T) {
	assert := assert.New(t)

	app := SetUpRouter()
	app.GET("/", Homepage)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	expected := `{"message":"Service is up and running."}`
	actual, _ := io.ReadAll(w.Body)

	assert.Equal(expected, string(actual))
	assert.Equal(http.StatusOK, w.Code)
}

type TestCase struct {
	path               string
	expectedStatus     int
	expectedUpdatedNIM string
}

var dummyMahasiswa = mahasiswa.Mahasiswa{
	NIM:     "321",
	Nama:    "dummy",
	Jurusan: "dummy",
}

func seedMahasiswa(db *sql.DB) error {
	_, err := db.Query(
		`
		INSERT INTO "mahasiswa"
		(nim, nama, jurusan, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`,
		dummyMahasiswa.NIM, dummyMahasiswa.Nama, dummyMahasiswa.Jurusan,
	)

	if err != nil {
		return err
	}

	return err
}

func cleanSeedMahasiswa(db *sql.DB) error {
	_, err := db.Query(
		`
		DELETE FROM mahasiswa 
		WHERE nim = $1 
		AND nama = $2 
		AND jurusan = $3`,
		dummyMahasiswa.NIM, dummyMahasiswa.Nama, dummyMahasiswa.Jurusan,
	)

	if err != nil {
		return err
	}

	return err
}

// steps:
// - create new dummy
// - read dummy,
// - if 200, success
func TestReads(t *testing.T) {
	assert := assert.New(t)
	testCase := &TestCase{
		path:           "/api/v1/mahasiswa/",
		expectedStatus: http.StatusOK,
	}

	db, err := SetUpDB()
	config.DB = db
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	err = seedMahasiswa(db)
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	req, _ := http.NewRequest(http.MethodGet, testCase.path, nil)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Request.URL.Path = testCase.path

	mahasiswa.Reads(c)
	assert.Equal(testCase.expectedStatus, rec.Code)

	t.Cleanup(func() {
		err := cleanSeedMahasiswa(db)
		if err != nil {
			t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
		}
	})
}

// steps:
// - create a dummy mahasiswa
// - read dummy
// - read the code
// - if 200, success
// - delete dummy
func TestRead(t *testing.T) {
	assert := assert.New(t)
	testCase := &TestCase{
		path:           "/api/v1/mahasiswa/:id",
		expectedStatus: http.StatusOK,
	}

	db, err := SetUpDB()
	config.DB = db
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	err = seedMahasiswa(db)
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	req, _ := http.NewRequest(http.MethodGet, testCase.path, nil)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.AddParam("id", "1")
	c.Request = req
	c.Request.URL.Path = testCase.path

	mahasiswa.Read(c)
	assert.Equal(testCase.expectedStatus, rec.Code)

	t.Cleanup(func() {
		err := cleanSeedMahasiswa(db)
		if err != nil {
			t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
		}
	})
}

// steps:
// - create a dummy mahasiswa
// - read the code
// - if 200, success
// - delete dummy
func TestCreate(t *testing.T) {
	assert := assert.New(t)
	testCase := &TestCase{
		path:           "/api/v1/mahasiswa/",
		expectedStatus: http.StatusOK,
	}

	db, err := SetUpDB()
	config.DB = db
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	newMahasiswa := mahasiswa.Mahasiswa{
		NIM:     dummyMahasiswa.NIM,
		Nama:    dummyMahasiswa.Nama,
		Jurusan: dummyMahasiswa.Jurusan,
	}

	reqBody, _ := json.Marshal(newMahasiswa)
	req, _ := http.NewRequest(
		http.MethodPost,
		testCase.path,
		bytes.NewBuffer(reqBody),
	)
	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Request.URL.Path = testCase.path

	mahasiswa.Create(c)
	assert.Equal(testCase.expectedStatus, rec.Code)

	t.Cleanup(func() {
		err := cleanSeedMahasiswa(db)
		if err != nil {
			t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
		}
	})
}

func seedMahasiswaWithReturn(db *sql.DB) *sql.Row {
	row := db.QueryRow(
		`
		INSERT INTO "mahasiswa"
		(nim, nama, jurusan, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP) 
		RETURNING id`,
		dummyMahasiswa.NIM, dummyMahasiswa.Nama, dummyMahasiswa.Jurusan,
	)

	return row
}

// steps:
// - create dummy mahasiswa
// - delete dummy
// - if 200, success
func TestDelete(t *testing.T) {
	assert := assert.New(t)
	testCase := &TestCase{
		path:           "/api/v1/mahasiswa/:id",
		expectedStatus: http.StatusOK,
	}

	db, err := SetUpDB()
	config.DB = db
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	row := seedMahasiswaWithReturn(db)
	var id int
	if err := row.Scan(&id); err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	req, _ := http.NewRequest(http.MethodDelete, testCase.path, nil)
	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.AddParam("id", strconv.Itoa(id))
	c.Request = req
	c.Request.URL.Path = testCase.path

	mahasiswa.Destroy(c)
	assert.Equal(testCase.expectedStatus, rec.Code)
}

type UpdateMahasiswaTest struct {
	NIM string
}

func checkUpdatedMahasiswaNIM(db *sql.DB, id int) *sql.Row {
	row := db.QueryRow(
		`
		SELECT nim FROM "mahasiswa"
		WHERE id = $1`,
		id,
	)

	return row
}

func cleanSeedUpdatedMahasiswa(db *sql.DB, id int) error {
	_, err := db.Query(
		`
		DELETE FROM mahasiswa
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	return err
}

// steps:
//   - create dummy mahasiswa
//   - update it
//   - check response code and updated value
//   - if res code and update value is equal
//     to expected updated value, success
//   - delete dummy
func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	testCase := &TestCase{
		path:               "/api/v1/mahasiswa/:id",
		expectedStatus:     http.StatusOK,
		expectedUpdatedNIM: "123",
	}

	db, err := SetUpDB()
	config.DB = db
	if err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	row := seedMahasiswaWithReturn(db)
	var id int
	if err := row.Scan(&id); err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	newlyUpdatedMahasiswa := UpdateMahasiswaTest{
		NIM: testCase.expectedUpdatedNIM,
	}
	reqBody, _ := json.Marshal(newlyUpdatedMahasiswa)
	req, _ := http.NewRequest(http.MethodPatch, testCase.path, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.AddParam("id", strconv.Itoa(id))
	c.Request = req
	c.Request.URL.Path = testCase.path

	mahasiswa.Update(c)

	updatedRow := checkUpdatedMahasiswaNIM(db, id)
	var NIM string
	if err := updatedRow.Scan(&NIM); err != nil {
		t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
	}

	assert.Equal(testCase.expectedStatus, rec.Code)
	assert.Equal(testCase.expectedUpdatedNIM, NIM)

	t.Cleanup(func() {
		err := cleanSeedUpdatedMahasiswa(db, id)
		if err != nil {
			t.Errorf("ERROR: Expected error to be %v, got %v", nil, err)
		}
	})
}
