package api

import (
	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/muscle_group", server.createMuscleGroup)
	router.GET("/muscle_group/:group_name", server.getMuscleGroup)
	router.GET("muscle_group", server.listMuscleGroups)
	router.PATCH("muscle_group", server.updateMuscleGroup)
	router.DELETE("/muscle_group/:group_name", server.deleteMuscleGroup)

	router.POST("/exercise", server.createExercise)
	router.GET("/exercise/:exercise_name", server.getExercise)
	router.GET("/exercise", server.listExercises)
	router.GET("/exercise/group/:muscle_group", server.getMuscleGroupExercises)
	router.PATCH("/exercise/name", server.updateExerciseName)
	router.PATCH("/exercise/group", server.updateExerciseMuscleGroup)
	router.DELETE("/exercise/:exercise_name", server.deleteExercise)

	router.POST("/lift", server.createLift)
	router.GET("/lift/:id", server.getLift)
	router.GET("/lift", server.listLifts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
