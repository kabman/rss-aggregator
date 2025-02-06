package main

import "net/http"

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Name string `json:"name"`

	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error passing json: %v",err))
		return
	}

	apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC()
		UpdatedAt : time.Now().UTC()
		Name: params.Name
	}

	) 

	if err != nil{
		respondWithError(w,400,fmt.Sprintf("couldnt create user: %v",err))
		return
	}
	respondWithJSON(w,200,struct{})
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	apiKey, err := auth.GetAPIKey(r.header)
	if err != nil{
		respodWithError(w,403,fmt,Sprintf("Auth error:%v",err))
		return
	}
	user,err = apiCfg.DB.GetUserByApiKey(r.context(),apiKey)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("couldnt get user: %v",err))
		return
	}
	respondWithJSON(w,200,struct{})
}