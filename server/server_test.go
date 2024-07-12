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

	"github.com/ifIMust/srsr/message"
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
			request := message.RegisterRequest{
				Name:    "dungen",
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
			r := message.RegisterResponse{}
			body, _ := io.ReadAll(responseRecorder.Body)
			json.Unmarshal(body, &r)
			Ω(len(r.ID)).Should(BeNumerically(">", 8))
		})
	})

	Context("Heartbeat", func() {
		var responseRecorder *httptest.ResponseRecorder
		var reqHTTP *http.Request

		Context("without matching service ID", func() {
			BeforeEach(func() {
				responseRecorder = httptest.NewRecorder()
				request := message.HeartbeatRequest{
					ID: "2145",
				}
				reqJSON, _ := json.Marshal(request)
				reqHTTP, _ = http.NewRequest("POST", "/heartbeat", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})

			It("returns OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})
			It("was unsuccessful", func() {
				r := message.HeartbeatResponse{}
				body, _ := io.ReadAll(responseRecorder.Body)
				json.Unmarshal(body, &r)
				Expect(r.Success).To(BeFalse())
			})
		})
		Context("with matching service ID", func() {
			BeforeEach(func() {
				responseRecorder = httptest.NewRecorder()
				regRequest := message.RegisterRequest{
					Name:    "dungen",
					Address: "localhost:5000",
				}
				reqJSON, _ := json.Marshal(regRequest)
				registerHTTP, _ := http.NewRequest("POST", "/register", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, registerHTTP)

				regResp := message.RegisterResponse{}
				body, _ := io.ReadAll(responseRecorder.Body)
				json.Unmarshal(body, &regResp)

				responseRecorder = httptest.NewRecorder()
				request := message.HeartbeatRequest{
					ID: regResp.ID,
				}
				reqJSON, _ = json.Marshal(request)
				reqHTTP, _ = http.NewRequest("POST", "/heartbeat", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("returns OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})
			It("was successful", func() {
				r := message.HeartbeatResponse{}
				body, _ := io.ReadAll(responseRecorder.Body)
				json.Unmarshal(body, &r)
				Expect(r.Success).To(BeTrue())
			})

		})

	})
})
