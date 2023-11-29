package httpchi

import (
	"github.com/vlasashk/todo-manager/config"
	"github.com/vlasashk/todo-manager/internal/models/task"
	"log"
	"net/http"
)

func Run(db task.Repo, cfg config.AppCfg) {
	r := NewRouter(db)
	log.Fatalln(http.ListenAndServe(cfg.Host+":"+cfg.Port, r))
}
