package server_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	
	"github.com/gin-gonic/gin"

	"github.com/ifIMust/srsr/registry"
	"github.com/ifIMust/srsr/server"
)

var _ = Describe("Server", func() {
	var router *gin.Engine
	
	BeforeEach(func() {
		registry := registry.NewServiceRegistry()
		router = server.SetupRouter(registry)
	})
	Context("Register", func() {
		var responseRecorder *httptest.ResponseRecorder

		BeforeEach(func() {
			responseRecorder = httptest.NewRecorder()
			request := server.RegisterRequest{
				Name: "dungen",
				Address: "localhost:5000",
			}
			reqJSON, _ := json.Marshal(request)
			reqHTTP, _ := http.NewRequest("POST", "/register", strings.NewReader(string(reqJSON)))
			router.ServeHTTP(responseRecorder, reqHTTP)
		})
		
		It("returns OK", func() {
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("responds with the expected struct and ID", func() {
			r := server.RegisterResponse{}
			body, _ := io.ReadAll(responseRecorder.Body)
			err := json.Unmarshal(body, &r)
			Expect(err).To(BeNil())
			Ω(len(r.ID)).Should(BeNumerically(">", 8))
		})
	})
})
