package main

import (
	"encoding/json"
	"fmt"
	"github.com/cking-bot/dice_bot/models"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func rollHandler(w http.ResponseWriter, r *http.Request) {
	var rolls []int
	var value int
	var response string
	var diceType int
	var cmd models.Command

	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch cmd.Dice {
	case "4":
		diceType = 3
	case "6":
		diceType = 5
	case "8":
		diceType = 7
	case "10":
		diceType = 9
	case "12":
		diceType = 11
	case "20":
		diceType = 19
	default:
		diceType = 5
	}

	if cmd.Number > 10 || cmd.Number < 1 {
		cmd.Number = 1
	}

	if cmd.Advantage == "1" || cmd.Advantage == "2" {
		cmd.Number += 1
	}

	for r := 0; r < cmd.Number; r++ {
		rand.Seed(time.Now().UnixMicro())
		roll := rand.Int() % diceType
		rolls = append(rolls, (roll + 1))
		time.Sleep(3 * time.Microsecond)
	}

	log.Println(rolls)
	min := rolls[0]
	max := rolls[0]
	for _, v := range rolls {
		value += v
	}

	for _, v := range rolls {
		if v < min {
			min = v
		}
	}

	for _, v := range rolls {
		if v > max {
			max = v
		}
	}

	switch cmd.Advantage {
	case "1":
		response = fmt.Sprintf("Rolled %d D%d (advantage): %d | Rolls: %v", cmd.Number-1, diceType, max, rolls)
	case "2":
		response = fmt.Sprintf("Rolled %d D%d (disadvantage): %d | Rolls: %v", cmd.Number-1, diceType, min, rolls)
	default:
		response = fmt.Sprintf("Rolled %d D%d: %d | Rolls: %v", cmd.Number, diceType, value, rolls)
	}
	reply, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func main() {
	Router := mux.NewRouter()

	Router.HandleFunc("/roll", rollHandler).Methods("POST")

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", Router))
}
