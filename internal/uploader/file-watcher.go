package uploader

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var DVRLastLine map[string]string

func fileChangeHandler(filePath string, streamMediaID string){
    time.Sleep(1 * time.Second)
    splitFileName := strings.Split(filePath, "\\")
    destName := streamMediaID + "/" +splitFileName[len(splitFileName) - 1]

    fileReader, err := os.Open(filePath)                        
    if err != nil {
        log.Printf("Failed to open %s for reading\n", filePath)
        return
    }

    defer fileReader.Close()

    if err = FileUploaderInstance.UploadFile(fileReader, destName); err != nil {
        log.Printf("Failed to upload %s: %v\n", filePath, err)
    } else {
        log.Printf("Uploaded %s as %s\n", filePath, destName)
    }

    if strings.HasSuffix(filePath, ".m3u8") {
        fileReader.Seek(0, io.SeekStart)
        WriteDVRPlaylist("./media/" + streamMediaID, splitFileName[len(splitFileName) - 1], fileReader)
    }
}

//TODO: should close the watcher when stream ends
func SetupFileWatcher(streamMediaID string) {
    mediaDir := "./media/" + streamMediaID
    DVRLastLine = make(map[string]string)

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
                    go fileChangeHandler(event.Name, streamMediaID)
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