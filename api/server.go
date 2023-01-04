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

	server.buildRoutes()
	return server, nil
}

func (server *Server) buildRoutes() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.POST("/user/login", server.login)

	authRouter := router.Group("/").Use(authenticationMiddleware(server.tokenCreator))

	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccounts)
	authRouter.DELETE("/accounts/:id", server.deleteAccount)

	authRouter.POST("/category", server.createCategory)
	authRouter.GET("/category/:id", server.getCategory)
	authRouter.GET("/category", server.listCategories)
	authRouter.PATCH("/category/:id", server.updateCategory)
	authRouter.DELETE("/category/:id", server.deleteCategory)

	authRouter.POST("/muscle_group", server.createMuscleGroup)
	authRouter.GET("/muscle_group/:name", server.getMuscleGroup)
	authRouter.GET("muscle_group", server.listMuscleGroups)
	authRouter.PATCH("muscle_group", server.updateMuscleGroup)
	authRouter.DELETE("/muscle_group/:name", server.deleteMuscleGroup)

	authRouter.POST("/exercise", server.createExercise)
	authRouter.GET("/exercise/:name", server.getExercise)
	authRouter.GET("/exercise", server.listExercises)
	authRouter.GET("/exercise/group/:muscle_group", server.getMuscleGroupExercises)
	authRouter.PATCH("/exercise/:name", server.updateExercise)
	authRouter.DELETE("/exercise/:name", server.deleteExercise)

	authRouter.POST("/workout/:user_id", server.createWorkout)
	authRouter.GET("/workout/:workout_id", server.getWorkout)
	authRouter.GET("/workout/history/:user_id", server.listWorkouts)
	authRouter.PATCH("/workout/:workout_id", server.updateFinishTime)
	authRouter.DELETE("/workout/:workout_id", server.deleteWorkout)

	authRouter.POST("/lift", server.createLift)
	authRouter.POST("/lift/:workout_id/:user_id", server.createLifts)
	authRouter.GET("/lift/:id/:user_id", server.getLift)
	authRouter.GET("/lift/history/:user_id", server.listLifts)
	authRouter.GET("/lift/history/pr/:order_by/:user_id", server.listPRs)
	authRouter.GET("/lift/pr/:exercise_name/:order_by/:user_id", server.listPRsByExercise)
	authRouter.GET("/lift/pr/group/:muscle_group/:order_by/:user_id", server.listPRsByMuscleGroup)
	authRouter.PATCH("/lift/:id", server.updateLift)
	authRouter.DELETE("/lift/:id", server.deleteLift)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
