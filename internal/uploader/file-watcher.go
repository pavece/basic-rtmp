package uploader

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

//TODO: add segments to the DVR playlist
//Watch for new .ts / .m3u8 in the media directories to opload them to the object store / CDN
func SetupFileWatcher() {
    mediaDir := "./media"

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

                if event.Op & fsnotify.Create == fsnotify.Create &&  (strings.HasSuffix(event.Name, ".ts") || strings.HasSuffix(event.Name, ".m3u8")){
                    log.Printf("File detected: %s\n", event.Name)

                    go func(filePath string) {
                        time.Sleep(1 * time.Second)
                        destName := strings.Split(filePath, "\\")[1]

                        if err := FileUploaderInstance.UploadFile(filePath, destName); err != nil {
                            log.Printf("Failed to upload %s: %v\n", filePath, err)
                        } else {
                            log.Printf("Uploaded %s as %s\n", filePath, destName)
                        }
                    }(event.Name)
                }

            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("Watcher error:", err)
            }
        }
    }()

    if err := watcher.Add(mediaDir); err != nil {
        log.Fatal("failed to attach watcher to media folder:", err)
    }

    log.Println("Watching", mediaDir, "for new .ts / .m3u8 files...")
}