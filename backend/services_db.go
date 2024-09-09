package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Write users to file in db/users.json.
func write_users(users map[string]*CustomUser) {
	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("error marshalling users to JSON: ", err)
	}

	file, err := os.Create("db/users.json")
	if err != nil {
		log.Fatalln("Error creating users file: ", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalln("Error writing user JSON to file: ", err)
	}
}

// Read users from db/users.json.
// Returns a map where the key is the user.login and the value is the custom user struct
func read_users() map[string]*CustomUser {
	file, err := os.Open("db/users.json")
	if err != nil {
		log.Println("error reading user file: ", err)
		return nil
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading user file: ", err)
	}

	UserMap := make(map[string]*CustomUser)

	err = json.Unmarshal(jsonData, &UserMap)
	if err != nil {
		log.Fatalln("Error unmarshalling users from Json ", err)
	}

	return UserMap
}

// Write teams to file in db/teams.json
func write_teams(teams map[string]*CustomTeam, active_only bool) {
	jsonData, err := json.Marshal(teams)
	if err != nil {
		log.Fatalln("Error marshalling teams to JSON: ", err.Error())
	}

	var file *os.File

	if active_only {
		file, err = os.Create("db/teams_active.json")
	} else {
		file, err = os.Create("db/teams.json")
	}
	if err != nil {
		log.Fatalln("Error creating teams file: ", err.Error())
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalln("error writing team JSON to file: ", err)
	}
}

// Read users from db/teams.json.
// Returns a map where the team.slug is the key and the value is the custom team struct
func read_teams(active_only bool) map[string]*CustomTeam {
	var file *os.File
	var err error

	if active_only {
		file, err = os.Open("db/teams_active.json")
	} else {
		file, err = os.Open("db/teams.json")
	}
	if err != nil {
		log.Println("error reading team file", err)
		return nil
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading team file: ", err)
	}

	teamMap := make(map[string]*CustomTeam)

	err = json.Unmarshal(jsonData, &teamMap)
	if err != nil {
		log.Fatalln("Error unmarshalling teams from Json ", err)
	}

	return teamMap

}