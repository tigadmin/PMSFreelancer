package controller

import (
	"errors"
	"net/http"
	"strconv"

	// "PMSFreelancer/config"
	// "PMSFreelancer/models"

	"github.com/gin-gonic/gin"
	"github.com/stevejo12/PMSFreelancer/config"
	"github.com/stevejo12/PMSFreelancer/models"
)

func userExperience(id string) ([]models.ExperienceReturnValue, error) {
	resp, err := config.DB.Query("SELECT * FROM experience WHERE user_id=?", id)

	if err != nil {
		return []models.ExperienceReturnValue{}, errors.New("Server unable to execute query to database")
	}

	allData := []models.ExperienceReturnValue{}

	for resp.Next() {
		var databaseData models.ExperienceTableResponse
		if err := resp.Scan(&databaseData.ID, &databaseData.Description, &databaseData.Place, &databaseData.Position, &databaseData.StartYear, &databaseData.EndYear, &databaseData.UserID); err != nil {
			return []models.ExperienceReturnValue{}, errors.New("Something is wrong with the database data")
		}

		var returnValue models.ExperienceReturnValue

		returnValue.ID = databaseData.ID
		returnValue.Position = databaseData.Position
		returnValue.Place = databaseData.Place
		returnValue.StartYear = databaseData.StartYear
		returnValue.EndYear = databaseData.EndYear
		returnValue.Description = databaseData.Description

		allData = append(allData, returnValue)
	}

	if resp.Err() != nil {
		return []models.ExperienceReturnValue{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

// GetOnlyUserExperience => Get Detail View for the User Experience
// GetOnlyUserExperience godoc
// @Summary User Experience
// @Produce json
// @Accept  json
// @Tags Experience
// @Param token header string true "Token Header"
// @Success 200 {object} models.ResponseOKGetUserExperience
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /userExperience [get]
func GetOnlyUserExperience(c *gin.Context) {
	id := idToken

	allUserExperience, err := userExperience(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    []models.ExperienceReturnValue{}})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "All User Experience data have been retrieved",
		"data":    allUserExperience})
}

// AddExperience => Add User Experience
// AddExperience godoc
// @Summary Adding User Experience
// @Produce json
// @Accept  json
// @Tags Experience
// @Param token header string true "Token Header"
// @Param Data body models.ExperienceParameters true "Data Format to add experience"
// @Success 200 {object} models.ResponseWithNoBody
// @Failure 500 {object} models.ResponseWithNoBody
// @Router /addExperience [post]
func AddExperience(c *gin.Context) {
	id := idToken

	var data models.ExperienceParameters

	err = c.BindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})
		return
	}

	if data.StartYear > data.EndYear {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Start year should be in the past compared to end year"})
		return
	}

	query := "INSERT INTO experience(place, position, starting_year, ending_year, user_id, description) VALUES"
	query = query + "(\"" + data.Place + "\", \"" + data.Position + "\"," + strconv.Itoa(data.StartYear) + ", " + strconv.Itoa(data.EndYear) + ", " + id + ", \"" + data.Description + "\")"

	_, err = config.DB.Exec(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server unable to execute query to database"})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Successfully Added Experience"})
}
