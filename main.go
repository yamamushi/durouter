package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/yamamushi/durouter/config"
	"github.com/yamamushi/durouter/config/db"
	"github.com/yamamushi/durouter/controllers"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Variables used for command line parameters
var (
	ConfPath string
)

func init() {
	// Read our command line options
	flag.StringVar(&ConfPath, "c", "durouter.conf", "Path to Router Config File - Default: durouter.conf")
	flag.Parse()

	_, err := os.Stat(ConfPath)
	if err != nil {
		log.Fatal("Config file is missing: ", ConfPath)

	}
}

func main() {

	fmt.Println("\n\n|| Starting durouter ||")
	log.SetOutput(ioutil.Discard)

	// Verify we can actually read our config file
	globalconfig, err := config.ReadConfig(ConfPath)
	if err != nil {
		fmt.Println("error reading config file at: ", ConfPath)
		return
	}

	dbmanager := db.NewDBManager(globalconfig)
	_, err = dbmanager.GetCollection(globalconfig.DBConfig.AccountColumn)
	if err != nil {
		fmt.Println("DB Connection error: " + err.Error())
		return
	}

	r := chi.NewRouter()
	r.Get("/", sayHello)

	verifyAuth := controllers.NewVerifyAuth(globalconfig)
	r.Get("/verifyauth/*", verifyAuth.Verify)
	serverstatus := controllers.NewServerStatus(globalconfig)
	r.Get("/serverstatus", serverstatus.GetStatus)
	r.Get("/serverschedule", serverstatus.GetSchedule)

	err = http.ListenAndServe(globalconfig.HostConfig.Host+":"+globalconfig.HostConfig.Port, r)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	fmt.Println("Completed")
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	_, _ = w.Write([]byte(message))
}
