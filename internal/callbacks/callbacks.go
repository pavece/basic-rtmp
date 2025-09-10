package callbacks

import (
	"fmt"

	"github.com/google/uuid"
)

/*
	This package contains all the callbacks you can modify.
*/

//Use this to mark streams as completed or run custom logic related to stream ends
func OnStramEnd(streamId int, mediaId string){
	fmt.Printf("Stream with id %d and media id %s has ended \n", streamId, mediaId)
}



//Create your own validation logic
func ValidateStreamKey(streamKey string) error{
	return nil
}

/* 
	ID returned here will be used to upload media to object store
	Create your own logic for generating them (this function is called after stream key validation)

	Recommended: get user from database via streamKey, add a new video record and return the record's id here
*/
func GenerateMediaId(streamKey string) (string, error){
	return uuid.NewString(), nil
}