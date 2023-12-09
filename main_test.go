package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMethod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	reader := strings.NewReader("Document=this is a very random statement")
	req, _ := http.NewRequest("POST", "/add", reader) //BTW check for error
	router.ServeHTTP(w, req)

	result := []string{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("the result is")
	fmt.Println(result)

	assert.Equal(t, reflect.Struct, reflect.TypeOf(result).Kind())
}
