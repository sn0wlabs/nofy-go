package nofy

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var privateKey = ""

func TestMain(m *testing.M) {
	fmt.Print("Init Test Environment...")
	godotenv.Load()

	privateKey = os.Getenv("NOFY_PRIVATE_KEY")
	m.Run()
}

func TestNew(t *testing.T) {
	client := New("testKey")

	if client.Backend.Key != "dGVzdEtleTo=" {
		t.Errorf("Expected dGVzdEtleTo=, got %s", client.Backend.Key)
	}
}

func TestCustomer(t *testing.T) {
	nofy := New(privateKey)

	_, err := nofy.Customer.GetAll()

	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestUpload(t *testing.T) {
	nofy := New(privateKey)

	// read image test.jpg
	image, e := os.Open("test.jpg")
	if e != nil {
		t.Errorf("Error: %s", e)
	}

	res, err := nofy.Upload.WebP(image)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	t.Log(res)
}
