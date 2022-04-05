package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type UserList struct {
	Users []User `json: "users"`
}

type User struct {
	Fname string `json: "Fname"`
	Lname string `json: "Lname"`
	Email string `json: "email"`
	PASS  string `json:"PASS"`
}

const (
	SSO = ""
	ADD = ""
	EMAIL =""
	PASSWORD = ""
)

func main() {
	fmt.Println("Bot Started")
	file, _ := ioutil.ReadFile("users.json")
	data := UserList{}

	_ = json.Unmarshal([]byte(file), &data)

	path, _ := launcher.LookPath()
	fmt.Println("path: ", path)
	l := launcher.New().Bin(path).Headless(false).Leakless(true).RemoteDebuggingPort(10000)
	defer l.Cleanup()
	u, _ := l.Launch()
	b := rod.New().ControlURL(u).MustConnect().Trace(true).Timeout(60 * time.Second)
	fmt.Println("browser: ", b)
	defer b.MustClose()
	p := b.MustPage(SSO)
	p.MustWaitLoad()

	email := p.MustElement(`[id="email"]`)
	email.MustInput(EMAIL)
	time.Sleep(3 * time.Second)
	password := p.MustElement(`[id="pw"]`)
	password.MustInput(PASSWORD)
	time.Sleep(3 * time.Second)
	p.MustElement(`[id="submitAdmin"]`).MustClick()
	p.MustWaitLoad()

	for _, user := range data.Users {
		time.Sleep(2 * time.Second)
		email := strings.Split(user.Email, "@")[0]
		password := strings.Trim(user.PASS, " ")
		fullName := user.Fname + " " + user.Lname
		newP := b.MustPage(ADD)
		newP.MustWaitLoad()
		newP.MustElement(`[id="localpart"]`).MustInput(email)
		newP.MustElement(`[id="pw"]`).MustInput(password)
		newP.MustElement(`[id="pw2"]`).MustInput(password)
		newP.MustElement(`[id="displayed_name"]`).MustInput(fullName)
		newP.MustElement(`[id="submit"]`).MustClick()
		newP.MustWaitLoad()
	}

	fmt.Println("Bye")
}
