package test

import (
	"cometosee/controller"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func Test_POST(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("error in loading env :%v", err)
	}

	//intailize fireabase app
	ctx := context.Background()

	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_ID")}, sa)

	if err != nil {
		t.Fatalf("error in intailizing firebase app :%v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		t.Fatalf("Error intailizing firestore cient :%v", err)
	}

	defer client.Close()

	body := map[string]string{
		"caption":   "testinggg",
		"image_url": "https://www.bing.com/images/search?view=detailV2&ccid=Yg8snzOE&id=9FABF87A86BCFC4AB7BE97C01CC07CC867FD2298&thid=OIP.Yg8snzOEMZcajsBowQ7gHQHaEp&mediaurl=https%3a%2f%2fimg.freepik.com%2fpremium-vector%2fcricket-standings-point-table-template-participating-countries-asia_597133-851.jpg%3fw%3d2000&cdnurl=https%3a%2f%2fth.bing.com%2fth%2fid%2fR.620f2c9f338431971a8ec068c10ee01d%3frik%3dmCL9Z8h8wBzAlw%26pid%3dImgRaw%26r%3d0&exph=1255&expw=2000&q=cricket+standings&FORM=IRPRST&ck=A3E7D388E4EC85880220BFA51F454633&selectedIndex=4&itb=0",
		"community": "cricket",
		"username":  "barun bhattarai",
	}
	bodyJson, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(string(bodyJson)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	controller.UploadPost(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Excepted status 200, got %d", rec.Code)
	}

	//taking doc id from response
	resp := rec.Body.String()
	parts := strings.Split(resp, "docId: ")

	docsid := strings.TrimSpace(parts[1])

	//check if the post is uploaded infirebase or not
	docs, err := client.Collection("posts").Doc(docsid).Get(ctx)
	if err != nil || !docs.Exists() {
		t.Fatalf("Post not found firestore")
	}

	//clear after test
	_, _ = client.Collection("posts").Doc(docsid).Delete(ctx)
}
