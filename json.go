package main

import "net/http"

func respondWithError(w http.ResponseWrite,code int, msg string){
	if code > 499{
		log.println("Responding with 5XX error:",msg)
	}
	type errResponse struct{
		Error string `json:"error"`
	}
	respondWithJSON(w,code,err)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	data,err := json.Marshal(payload)
	if err != nil{
		log.printf("Failed to marshal JSON response: %v",payload)
		w.writeHeader(500)
		return
	}
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(code )
	w.Write(data)
}