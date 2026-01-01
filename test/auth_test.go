package test

import (
	"cometosee/controller"
	"cometosee/intailizer"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestSignup_TestLogin(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading env file: %v", err)
	}
	intailizer.DatabaseConnection()

	if intailizer.DB == nil {
		t.Fatalf("DB not initialized")
	}

	email := "testing@gmail.com"
	body := `{
	 "username":"tester",
	 "email":"` + email + `",
	 "password":"password123"
	}`

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	controller.Signup(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", rec.Code)
	}

	var count int
	err = intailizer.DB.QueryRow(
		`SELECT COUNT(*) FROM cometoseeauth WHERE email=$1`, email,
	).Scan(&count)

	if err != nil || count != 1 {
		t.Fatalf("User not inserted into database")
	}

	defer intailizer.DB.Exec(`DELETE FROM cometoseeauth WHERE email=$1`, email)

	//login__test
	loginreq := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	loginreq.Header.Set("Content-Type", "application/json")

	loginrec := httptest.NewRecorder()
	controller.Login(loginrec, loginreq)

	//check http status
	if loginrec.Code != http.StatusOK {
		t.Fatalf("Expected 200,got %d", loginrec.Code)
	}

	//check if jwt  is returenrd or not
	responsestr := loginrec.Body.String()
	if !strings.Contains(responsestr, "token") {
		t.Fatalf("Jwt token not found in response :%s", responsestr)
	}
}
