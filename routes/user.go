package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mend/db/tp"
	"mend/model/user"
)

// userSetup creates routes group for user handling
func userSetup(router *gin.Engine) {
	group := router.Group("/user")
	{
		group.GET("count", count())
		group.GET("", getAll())
		group.GET(":id", get())
		group.POST("", add())
		group.PUT(":id", update())
		group.DELETE(":id", remove())
	}
}

func count() gin.HandlerFunc {
	return func(c *gin.Context) {
		n := user.Count()
		c.JSON(http.StatusOK, gin.H{"count": n})
	}
}

func getAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := user.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, rows)
	}
}

// get returns data for selected user
func get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		row, err := user.GetWithID(tp.ID(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, row)
	}
}

// swagger:operation POST
// Adds new user to database
// ---
// Produces:
// - application/json
// Responses:
//  '200':
//      description: wszystko dobrze
func add() gin.HandlerFunc {
	return func(c *gin.Context) {
		var u user.User

		if err := c.BindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := u.AddToDatabase(); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"ID of new user": u.ID})
	}
}

// update updates data in db for selected user
func update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var u user.User
		if err = c.BindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		u.ID = tp.ID(id)
		if err = u.UpdateInDatabase(); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		c.JSON(http.StatusOK, "ok")
	}
}

// remove deletes data of selected user from db
func remove() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if err = user.DeleteFromDatabase(tp.ID(id)); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, "ok")
	}
}
