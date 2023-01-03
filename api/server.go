package api

import (
	"fmt"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/token"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config       util.Config
	store        db.Store
	tokenCreator token.Maker
	router       *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenCreator, err := token.NewJWTCreator(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create token generator: %w", err)
	}

	server := &Server{
		config:       config,
		store:        store,
		tokenCreator: tokenCreator,
	}
	router := gin.Default()

	router.POST("/user/login", server.login)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/category", server.createCategory)
	router.GET("/category/:id", server.getCategory)
	router.GET("/category", server.listCategories)
	router.PATCH("/category/:id", server.updateCategory)
	router.DELETE("/category/:id", server.deleteCategory)

	router.POST("/muscle_group", server.createMuscleGroup)
	router.GET("/muscle_group/:name", server.getMuscleGroup)
	router.GET("muscle_group", server.listMuscleGroups)
	router.PATCH("muscle_group", server.updateMuscleGroup)
	router.DELETE("/muscle_group/:name", server.deleteMuscleGroup)

	router.POST("/exercise", server.createExercise)
	router.GET("/exercise/:name", server.getExercise)
	router.GET("/exercise", server.listExercises)
	router.GET("/exercise/group/:muscle_group", server.getMuscleGroupExercises)
	router.PATCH("/exercise/:name", server.updateExercise)
	router.DELETE("/exercise/:name", server.deleteExercise)

	router.POST("/workout/:user_id", server.createWorkout)
	router.GET("/workout/:workout_id", server.getWorkout)
	router.GET("/workout/history/:user_id", server.listWorkouts)
	router.PATCH("/workout/:workout_id", server.updateFinishTime)
	router.DELETE("/workout/:workout_id", server.deleteWorkout)

	router.POST("/lift", server.createLift)
	router.POST("/lift/:workout_id/:user_id", server.createLifts)
	router.GET("/lift/:id/:user_id", server.getLift)
	router.GET("/lift/history/:user_id", server.listLifts)
	router.GET("/lift/history/pr/:order_by/:user_id", server.listPRs)
	router.GET("/lift/pr/:exercise_name/:order_by/:user_id", server.listPRsByExercise)
	router.GET("/lift/pr/group/:muscle_group/:order_by/:user_id", server.listPRsByMuscleGroup)
	router.PATCH("/lift/:id", server.updateLift)
	router.DELETE("/lift/:id", server.deleteLift)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
