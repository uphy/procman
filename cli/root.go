package cli

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"github.com/uphy/procman/config"
	"github.com/uphy/procman/handlers"
	"github.com/uphy/procman/process"
)

func Execute() {
	command := cobra.Command{
		Use: "procman",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetUint32("port")
			processCommand := args
			env := config.NewEnv()
			fmt.Println(port, processCommand)

			authHandler := handlers.NewAuth(env.JWTSecret(), env.User(), env.Password())
			proc := process.New(processCommand)
			procHandler := handlers.NewProc(proc)

			e := echo.New()
			e.Static("/", "./static")
			e.POST("/api/auth/login", authHandler.Login)
			e.POST("/api/auth/info", authHandler.Info)
			e.POST("/api/process/start", procHandler.Start, authHandler.Authorized)
			e.POST("/api/process/stop", procHandler.Stop, authHandler.Authorized)
			e.POST("/api/process/restart", procHandler.Stop, authHandler.Authorized)
			e.POST("/api/process/status", procHandler.Status, authHandler.Authorized)
			e.Start(fmt.Sprintf(":%d", port))
		},
	}
	command.Flags().Uint32P("port", "p", 3000, "port number for web UI")

	if err := command.Execute(); err != nil {
		fmt.Println(err)
	}
}
