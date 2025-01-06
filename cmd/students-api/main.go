package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/helloankitpandey/students-api/internal/config"
	"github.com/helloankitpandey/students-api/internal/http/handlers/student"
	"github.com/helloankitpandey/students-api/internal/storage/sqlite"
)

func main() {
	// fmt.Println("welcome to students-a pii")
	// 1.load config

	cfg := config.MustLoad()
	fmt.Println("cfg is loaded")
	// added some log to check




	// 2.database
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sql lite  is loaded")

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// for graphical user interface of database use =>> TablePlus





	// 3.setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	// getting students data by id
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	// getting all students data 
	router.HandleFunc("GET /api/students", student.GetList(storage))

	
	


	// 4.setup server

	server := http.Server{
		Addr:  cfg.Addr,
		Handler: router,
	}
	// fmt.Println("server started")

	// for checking address
	// fmt.Printf("server started %s", cfg.HTTPServer.Addr)
	
	// at place of printf we do slog =. structured log
	slog.Info("server started", slog.String("address", cfg.Addr))

	// if yha ongoing request chal rahi hai to jb server band krte hai to 
	// to uss program ko ussi moment pe stop kr dega
	// ongoing request fail ho jayegi
	// to yese nhi hona chahiye hame isse gracefully stud down krna hai
	
	// then what we can do => ongoing request ko hame complete krna hai uske bad hi hame server shut down hona chahiye hai
	// production ke andar gracefully shutdown jaruri hai


	// // isko ak alag go routine ke andar run krunga
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatal("failed to start server")
	// }

    // channels we use here
	done := make(chan os.Signal, 1)  // iske andar signal hongi 

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func()  {
		// isko ak alag go routine ke andar run krunga
	    err := server.ListenAndServe()
	    if err != nil {
		    log.Fatal("failed to start server")
	    }
	} ()

	<-done // yha done channel ko ham listen kar rhe hai
	// its means that => jb tk done channel ke ander ko signal jata nahi hai
	// tb tkk ham <-done yha pe block ho jayenge && iske aage hamra code badega nahi 
	// means -> main program continuosly run hote rhega && go rutine run hote rhega

	// after we write server stop logic
	slog.Info("shutting down the server")

	// we use timer for shuuting down the server using context
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	// called -> graceful shutdown
	slog.Info("server shutdown successfully")
}