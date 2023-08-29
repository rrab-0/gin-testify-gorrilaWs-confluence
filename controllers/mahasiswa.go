package mahasiswa

import (
	"example/unit-test-hello-world/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Mahasiswa struct {
	NIM     string `json:"nim"`
	Nama    string `json:"nama"`
	Jurusan string `json:"jurusan"`
}

func Create(c *gin.Context) {
	var db = config.DB
	mahasiswa := new(Mahasiswa)

	if err := c.Bind(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// add a new mahasiswa
	_, err := db.Query(
		`
		INSERT INTO "mahasiswa"
		(nim, nama, jurusan, created_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`,
		mahasiswa.NIM, mahasiswa.Nama, mahasiswa.Jurusan,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa added successfully.",
	})
}

func Reads(c *gin.Context) {
	var db = config.DB

	query := `SELECT nim, nama, jurusan FROM mahasiswa`
	rows, err := db.Query(query)

	if err != nil {
		return
	}
	defer rows.Close()

	mahasiswas := []Mahasiswa{}

	for rows.Next() {
		var mahasiswa Mahasiswa
		if err := rows.Scan(&mahasiswa.NIM, &mahasiswa.Nama, &mahasiswa.Jurusan); err != nil {
			return
		}
		mahasiswas = append(mahasiswas, mahasiswa)
	}

	c.JSON(http.StatusOK, mahasiswas)
}

type MahasiswaId struct {
	ID int `form:"id"`
}

// grabs uri
func Read(c *gin.Context) {
	var db = config.DB
	id := new(MahasiswaId)

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	query := `SELECT nim, nama, jurusan FROM mahasiswa WHERE "id" = $1`
	row := db.QueryRow(query, id.ID)
	var (
		nim     string
		nama    string
		jurusan string
	)

	if err := row.Scan(&nim, &nama, &jurusan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	response := Mahasiswa{
		NIM:     nim,
		Nama:    nama,
		Jurusan: jurusan,
	}

	c.JSON(http.StatusOK, response)
}

// func Update(c *gin.Context) {
// 	var db = config.DB
// 	id := new(MahasiswaId)

// 	if err := c.ShouldBindUri(&id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 	}

// 	query := ``
// }

// grabs uri too
func Destroy(c *gin.Context) {
	var db = config.DB
	id := new(MahasiswaId)

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// delete a mahasiswa
	_, err := db.Query(
		`
		DELETE FROM mahasiswa WHERE id = $1`,
		id.ID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa deleted successfully.",
	})
}

// // grabs query string
// func Read(c *gin.Context) {
// 	var db = config.DB
// 	id := new(MahasiswaId)

// 	if err := c.ShouldBind(&id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 	}

// 	query := `SELECT nim, nama, jurusan FROM mahasiswa WHERE "id" = $1`
// 	row := db.QueryRow(query, id.ID)
// 	var (
// 		nim     string
// 		nama    string
// 		jurusan string
// 	)

// 	if err := row.Scan(&nim, &nama, &jurusan); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err,
// 		})
// 	}

// 	response := Mahasiswa{
// 		NIM:     nim,
// 		Nama:    nama,
// 		Jurusan: jurusan,
// 	}

// 	c.JSON(http.StatusOK, response)
// }
