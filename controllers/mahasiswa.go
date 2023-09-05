package controllers

import (
	"example/unit-test-hello-world/config"
	"net/http"

	"example/unit-test-hello-world/localWebsocket"

	"github.com/gin-gonic/gin"
)

type MahasiswaController struct{}

type Mahasiswa struct {
	NIM     string `json:"nim"`
	Nama    string `json:"nama"`
	Jurusan string `json:"jurusan"`
}

type SuccessMessage struct {
	Message string
}

// Post a new Mahasiswa godoc
// @Summary      Create a new mahasiswa
// @Description  Create a new mahasiswa with their NIM, Nama, and Jurusan.
// @Tags         mahasiswa
// @Accept       json
// @Produce      json
// @Param        Mahasiswa body Mahasiswa true "Mahasiswa need to have NIM, Nama, and Jurusan"
// @Success      200  {object} SuccessMessage
// @Failure      400  {object} SuccessMessage
// @Failure      500  {object} SuccessMessage
// @Router       /api/v1/mahasiswa/ [post]
func (m *MahasiswaController) Create(c *gin.Context) {
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

	localWebsocket.Writer("Mahasiswa added successfully.")
	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa added successfully.",
	})
}

type AllMahasiswaResponse struct {
	Mahasiswa []Mahasiswa `json:"mahasiswa"`
}

// Get ALl Mahasiswa godoc
// @Summary      Returns all mahasiswa
// @Description  Returns all mahasiswa
// @Tags         mahasiswa
// @Accept       json
// @Produce      json
// @Success      200  {object} AllMahasiswaResponse
// @Failure      400  {object} SuccessMessage
// @Failure      500  {object} SuccessMessage
// @Router       /api/v1/mahasiswa/ [get]
func (m *MahasiswaController) Reads(c *gin.Context) {
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
		if err := rows.Scan(
			&mahasiswa.NIM, &mahasiswa.Nama, &mahasiswa.Jurusan); err != nil {
			return
		}
		mahasiswas = append(mahasiswas, mahasiswa)
	}

	localWebsocket.Writer("Get all mahasiswa successfully.")
	c.JSON(http.StatusOK, mahasiswas)
}

type MahasiswaId struct {
	ID int `uri:"id"`
}

// Get Mahasiswa by ID godoc
// @Summary      Returns one mahasiswa
// @Description  Returns one mahasiswa
// @Tags         mahasiswa
// @Accept       json
// @Produce      json
// @Param		 id path int true "Mahasiswa ID"
// @Success      200  {object} Mahasiswa
// @Failure      400  {object} SuccessMessage
// @Failure      500  {object} SuccessMessage
// @Router       /api/v1/mahasiswa/{id} [get]
func (m *MahasiswaController) Read(c *gin.Context) {
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
			"error": err.Error(),
		})
	}

	response := Mahasiswa{
		NIM:     nim,
		Nama:    nama,
		Jurusan: jurusan,
	}

	localWebsocket.Writer("Get mahasiswa by ID successfully.")
	c.JSON(http.StatusOK, response)
}

// Update Mahasiswa by ID godoc
// @Summary      Updates one mahasiswa
// @Description  Updates one mahasiswa
// @Tags         mahasiswa
// @Accept       json
// @Produce      json
// @Param		 id path int true "Mahasiswa ID"
// @Success      200  {object} SuccessMessage
// @Failure      400  {object} SuccessMessage
// @Failure      500  {object} SuccessMessage
// @Router       /api/v1/mahasiswa/{id} [patch]
func (m *MahasiswaController) Update(c *gin.Context) {
	var db = config.DB
	id := new(MahasiswaId)

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// get
	getQuery := `
	SELECT nim, nama, jurusan 
	FROM mahasiswa 
	WHERE "id" = $1`
	row := db.QueryRow(getQuery, id.ID)
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

	mahasiswa := new(Mahasiswa)
	if err := c.Bind(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if mahasiswa.NIM == "" {
		mahasiswa.NIM = nim
	}

	if mahasiswa.Nama == "" {
		mahasiswa.Nama = nama
	}

	if mahasiswa.Jurusan == "" {
		mahasiswa.Jurusan = jurusan
	}

	// update
	updateQuery := `
	UPDATE mahasiswa
	SET nim = $1, nama = $2, jurusan = $3
	WHERE id = $4`

	_, err := db.Query(updateQuery,
		mahasiswa.NIM,
		mahasiswa.Nama,
		mahasiswa.Jurusan,
		id.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	localWebsocket.Writer("Mahasiswa updated successfully.")
	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa updated successfully.",
	})
}

// Delete Mahasiswa by ID godoc
// @Summary      Deletes one mahasiswa
// @Description  Deletes one mahasiswa
// @Tags         mahasiswa
// @Accept       json
// @Produce      json
// @Param		 id path int true "Mahasiswa ID"
// @Success      200  {object} SuccessMessage
// @Failure      400  {object} SuccessMessage
// @Failure      500  {object} SuccessMessage
// @Router       /api/v1/mahasiswa/{id} [delete]
func (m *MahasiswaController) Destroy(c *gin.Context) {
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

	localWebsocket.Writer("Mahasiswa deleted successfully.")
	c.JSON(http.StatusOK, gin.H{
		"message": "Mahasiswa deleted successfully.",
	})
}
