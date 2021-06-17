package serve

import (
	"betting/cmd/serve/bet"
	"betting/cmd/serve/table"
	"betting/storage/memory"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCmd associates the serve command with the instantiation of a Server.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "serve instantiates a server",
		Run:   StartServer,
	}
}

func StartServer(_ *cobra.Command, _ []string) {
	router := mux.NewRouter()

	tableStorage := memory.NewTableStorage()
	betStorage := memory.NewBetStorage()

	t := table.Load(router, tableStorage, betStorage)
	b := bet.Load(t, tableStorage, betStorage)

	log.Info("started server")

	if err := http.ListenAndServe(viper.GetString("port"), b); err != nil {
		log.Fatal(err)
	}
}
