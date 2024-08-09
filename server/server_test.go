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
		})
		When("the request is valid", func() {
			BeforeEach(func() {
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
		When("the request is valid with empty address", func() {
			BeforeEach(func() {
				request := message.RegisterRequest{
					Name:    "dungen",
					Address: "",
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
		When("the request uses non-absolute URI", func() {
			BeforeEach(func() {
				request := message.RegisterRequest{
					Name:    "dungen",
					Address: "localhost",
				}
				reqJSON, _ := json.Marshal(request)
				reqHTTP, _ := http.NewRequest("POST", "/register", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("returns Bad", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("the request is malformed", func() {
			BeforeEach(func() {
				reqJSON := "{'Desc: 'dungen', 'Adress': 'localhost:5000'}"
				reqHTTP, _ := http.NewRequest("POST", "/register", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("responds Bad Request", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("the request lacks required name field", func() {
			BeforeEach(func() {
				reqJSON := "{\"Address\": \"localhost:5000\"}"
				reqHTTP, _ := http.NewRequest("POST", "/register", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("responds Bad Request", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
	Context("Deregister", func() {
		var responseRecorder *httptest.ResponseRecorder
		BeforeEach(func() {
			responseRecorder = httptest.NewRecorder()
		})
		When("the request is malformed", func() {
			BeforeEach(func() {
				reqJSON := "{}"
				reqHTTP, _ := http.NewRequest("POST", "/deregister", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("responds Bad Request", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("the request is valid, but no ID matches", func() {
			BeforeEach(func() {
				request := message.DeregisterRequest{
					ID: "2145",
				}
				reqJSON, _ := json.Marshal(request)
				reqHTTP, _ := http.NewRequest("POST", "/deregister", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("returns OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})

			It("responds with Success False", func() {
				deregResp := message.DeregisterResponse{}
				body, _ := io.ReadAll(responseRecorder.Body)
				json.Unmarshal(body, &deregResp)
				Expect(deregResp.Success).To(BeFalse())
			})
		})
		When("the request is valid, and an ID matches", func() {
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
				request := message.DeregisterRequest{
					ID: regResp.ID,
				}
				reqJSON, _ = json.Marshal(request)
				reqHTTP, _ := http.NewRequest("POST", "/deregister", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("returns OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})
		})
	})
	Context("Heartbeat", func() {
		var responseRecorder *httptest.ResponseRecorder
		var reqHTTP *http.Request

		Context("with malformed request", func() {
			BeforeEach(func() {
				responseRecorder = httptest.NewRecorder()
				reqJSON := "{}"
				reqHTTP, _ = http.NewRequest("POST", "/heartbeat", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("responds Bad Request", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

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
	Context("Lookup", func() {
		var responseRecorder *httptest.ResponseRecorder
		var reqHTTP *http.Request
		var response message.LookupResponse

		Context("with matching name", func() {
			var address string

			BeforeEach(func() {
				responseRecorder = httptest.NewRecorder()
				address = "localhost:5000"
				regRequest := message.RegisterRequest{
					Name:    "dungen",
					Address: address,
				}
				reqJSON, _ := json.Marshal(regRequest)
				registerHTTP, _ := http.NewRequest("POST", "/register", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, registerHTTP)

				regResp := message.RegisterResponse{}
				body, _ := io.ReadAll(responseRecorder.Body)
				json.Unmarshal(body, &regResp)

				responseRecorder = httptest.NewRecorder()
				request := message.LookupRequest{
					Name: "dungen",
				}
				reqJSON, _ = json.Marshal(request)
				reqHTTP, _ = http.NewRequest("POST", "/lookup", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)

				body, _ = io.ReadAll(responseRecorder.Body)
				json.Unmarshal(body, &response)
			})
			It("returns OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})
			It("was successful", func() {
				Expect(response.Success).To(BeTrue())
			})
			It("was successful", func() {
				Expect(response.Address).To(Equal(address))
			})
		})
		Context("with malformed request", func() {
			BeforeEach(func() {
				responseRecorder = httptest.NewRecorder()
				reqJSON := "{\"nombre\": \"youthere\"}"
				reqHTTP, _ = http.NewRequest("POST", "/lookup", strings.NewReader(string(reqJSON)))
				router.ServeHTTP(responseRecorder, reqHTTP)
			})
			It("responds Bad Request", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
