package uploader

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

//Watch for new .ts / .m3u8 in the media directories to opload them to the object store / CDN
func SetupFileWatcher(){

    _, err := os.Stat("./media")
    if os.IsNotExist(err) {
        os.Mkdir("./media", 0777)
    }


	watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, ".ts") {
                    log.Println("New segment created:", event.Name)
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("Watcher error:", err)
            }
        }
    }()

    err = watcher.Add("./media")
    if err != nil {
        log.Fatal("Failed to attach file watcher to media folder", err)
    }
}