package main

import (
	"net/http"
	"strconv"
	"strings"
)

// Beat struct.
type Beat struct {
	ID          int    `gorm:"column:id" json:"id"`
	Artist      string `json:"artist"`
	URL         string `gorm:"column:beat_url" json:"url"`
	Votes       int    `gorm:"column:votes" json:"votes"`
	ChallengeID int    `gorm:"column:challenge_id" json:"challenge_id,omitempty"`
	UserID      string `gorm:"column:user_id" json:"user_id,omitempty"`
	Color       string `json:"color"`
}

// SubmitBeat ...
func SubmitBeat(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	battleID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	// TODO - Change this GetBattle statement or change GetBattle, this doesn't need a * sql statement.

	battle := GetBattle(db, battleID)
	if battle.Title == "" {
		http.Redirect(w, r, "/", 301)
		return
	}

	var user = GetUser(w, r)

	m := map[string]interface{}{
		"Battle": battle,
		"User":   user,
	}

	URL := r.URL.RequestURI()

	tpl := "SubmitBeat"
	if strings.Contains(URL, "update") {
		tpl = "UpdateBeat"
	}

	print(tpl)
	tmpl.ExecuteTemplate(w, tpl, m)
}

// InsertBeat ...
func InsertBeat(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	user := GetUser(w, r)
	if !user.Authenticated {
		http.Redirect(w, r, "/auth/discord", 301)
		return
	}

	battleID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	redirectURL := "/battle/" + strconv.Itoa(battleID)

	// TODO - BATTLE ID AND DEADLINE
	isOpen := RowExists(db, "SELECT id FROM challenges WHERE id = ?", battleID)

	if !isOpen {
		http.Redirect(w, r, redirectURL, 301)
		return
	}

	if user.Authenticated && r.Method == "POST" {
		if !strings.Contains(r.FormValue("track"), "soundcloud") {
			http.Redirect(w, r, redirectURL, 301)
			return
		}

		stmt := "INSERT INTO beats(url, challenge_id, user_id) VALUES(?,?,?)"

		// If has entered, update instead.
		if RowExists(db, "SELECT challenge_id FROM beats WHERE user_id = ? AND challenge_id = ?", user.ID, battleID) {
			stmt = "UPDATE beats SET url=? WHERE challenge_id=? AND user_id=?"
		}

		ins, err := db.Prepare(stmt)
		if err != nil {
			panic(err.Error())
		}
		defer ins.Close()

		ins.Exec(r.FormValue("track"), battleID, user.ID)
	} else {
		print("Not post")
		http.Redirect(w, r, redirectURL, 301)
		return
	}
	// TODO - Redirect with alert for user.
	http.Redirect(w, r, redirectURL, 301)
	return
}

// DeleteBeat ...
func DeleteBeat(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	defer db.Close()

	user := GetUser(w, r)
	if !user.Authenticated {
		http.Redirect(w, r, "/auth/discord", 301)
		return
	}

	battleID, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	redirectURL := "/battle/" + strconv.Itoa(battleID)

	if user.Authenticated {
		stmt := "DELETE FROM beats WHERE user_id = ? AND challenge_id = ?"

		ins, err := db.Prepare(stmt)
		if err != nil {
			panic(err.Error())
		}
		defer ins.Close()

		ins.Exec(user.ID, battleID)
	} else {
		http.Redirect(w, r, redirectURL, 301)
		return
	}
	// TODO - Redirect with alert for user.
	http.Redirect(w, r, redirectURL, 301)
	return
}
