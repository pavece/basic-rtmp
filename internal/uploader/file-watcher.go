package uploader

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var fileUploader *FileUploader = &FileUploader{} //Unique instance for every stream
type FileWatcher struct {
    watcher *fsnotify.Watcher
    dvrGenerator *DVRGenerator
}

func (w *FileWatcher) InitFileWatcher(streamMediaID string, videoHeight int) {
    mediaDir := os.Getenv("LOCAL_MEDIA_DIR") + "/" + streamMediaID

    watcher, err := fsnotify.NewWatcher()
    w.watcher = watcher
    if err != nil {
        log.Fatal(err)
    }

    fileUploader = fileUploader.NewUploader()
    w.dvrGenerator = NewDVRGenerator(videoHeight, streamMediaID, fileUploader)

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }

                if event.Op & fsnotify.Create == fsnotify.Create &&  (strings.HasSuffix(event.Name, ".ts") || strings.HasSuffix(event.Name, ".m3u8")){
                    go w.fileChangeHandler(event.Name, streamMediaID)
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

func (w *FileWatcher) fileChangeHandler(filePath string, streamMediaID string){
    time.Sleep(1 * time.Second)
    filename := filepath.Base(filePath)
    destName := streamMediaID + "/" + filename

    fileReader, err := os.Open(filePath)                        
    if err != nil {
        log.Printf("Failed to open %s for reading\n", filePath)
        return
    }

    defer fileReader.Close()

    if err = fileUploader.UploadFile(fileReader, destName); err != nil {
        log.Printf("Failed to upload %s: %v\n", filePath, err)
    } else {
        log.Printf("Uploaded %s as %s\n", filePath, destName)
    }

    if strings.HasSuffix(filePath, ".m3u8") {
        fileReader.Seek(0, io.SeekStart)
        w.dvrGenerator.WriteDVRPlaylist(os.Getenv("LOCAL_MEDIA_DIR") + "/" + streamMediaID, filename, fileReader)
    }
}