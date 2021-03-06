package hammer

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type Employee struct {
	Name      string `json:"name"`
	JobTitle  string `json:"job_title"`
	JobTitle2 string `json:"job_title2"`
	Nested    Nested
}

type Nested struct {
	Field  string `json:"field"`
	Field2 int    `json:"field2"`
	Field3 uint32 `json:"field3"`
}

func getVerbs() []string {
	return []string{GET, PUT, POST, DELETE, PATCH, HEAD, TRACE, CONNECT, OPTIONS}
}

func TestInitRequest(t *testing.T) {
	req := initRequest()
	if len(req.headers) != 0 || len(req.requestParams) != 0 {
		t.Error("Test Failed: TestInitRequest")
	}
}

func TestCheckError(t *testing.T) {
	if err := (&Request{err: errors.New("demoError")}).checkError(); err == nil {
		t.Error("Test Failed: TestCheckError")
	}

	if err := (&Request{}).checkError(); err != nil {
		t.Error("Test Failed: TestCheckError")
	}
}

func TestRequestBuilder(t *testing.T) {
	req := RequestBuilder()
	if len(req.headers) != 0 || len(req.requestParams) != 0 {
		t.Error("Test Failed: TestInitRequest")
	}
}

func TestWithHeaders(t *testing.T) {
	key := "testKey"
	value := "testValue"
	req := RequestBuilder().WithHeaders(key, value)

	if x, found := req.headers[key]; found {
		if x != value {
			t.Error("TestFailed:TestWithHeaders")
		}
	} else {
		t.Error("TestFailed:TestWithHeaders")
	}
}

func TestWithRequestParams(t *testing.T) {
	key := "testKey"
	value := "testValue"
	req := RequestBuilder().WithRequestParams(key, value)

	if x, found := req.requestParams[key]; found {
		if x != value {
			t.Error("TestFailed:TestWithRequestParams")
		}
	} else {
		t.Error("TestFailed:TestWithRequestParams")
	}
}

func TestWithRequestBodyParams(t *testing.T) {
	key := "testKey"
	value := "testValue"
	req := RequestBuilder().WithRequestBodyParams(key, value)

	if x, found := req.requestBodyParams[key]; found {
		if x != value {
			t.Error("TestFailed:TestWithRequestBodyParams")
		}
	} else {
		t.Error("TestFailed:TestWithRequestBodyParams")
	}
}

func TestWithURL(t *testing.T) {
	url := "http://localhost:8081/"
	req := RequestBuilder().WithURL(url)

	if req.url != url {
		t.Error("Test Failed:TestWithURL")
	}
}

func TestWithBasicAuth(t *testing.T) {

	req := RequestBuilder().WithBasicAuth("user", "pass")

	if !reflect.DeepEqual(req.basicAuth, basicAuth{
		username: "user",
		password: "pass",
	}) {
		t.Error("Test Failed: TestWithBasicAuth")
	}
}

func TestWithID(t *testing.T) {
	urlx := "http://localhost:8081/"
	id := "8s09df890asd8"
	for _, verb := range getVerbs() {
		req := (&Request{url: urlx, httpVerb: verb}).WithID(id)
		switch verb {
		case PUT, GET, DELETE:
			if req.url != (urlx + "/" + id) {
				t.Error("Test Failed:TestWithID")
			}
		default:
			if req.err == nil {
				t.Error("Test Failed:TestWithID")
			}
		}
	}
}

func TestVerbs(t *testing.T) {

	for _, verb := range getVerbs() {
		req := &Request{}

		switch verb {
		case GET:
			req = req.Get()
		case HEAD:
			req = req.Head()
		case POST:
			req = req.Post()
		case PUT:
			req = req.Put()
		case PATCH:
			req = req.Patch()
		case DELETE:
			req = req.Delete()
		case CONNECT:
			req = req.Connect()
		case OPTIONS:
			req = req.Options()
		case TRACE:
			req = req.Trace()
		}

		if req.httpVerb != verb {
			t.Errorf("TestFailed:TestVerbs {%s}", verb)
		}
	}
}

func TestWithRequestBody(t *testing.T) {

	var testCase []interface{}
	emp := Employee{
		Name:      "name",
		JobTitle:  "jobTitle1",
		JobTitle2: "jobTitle2",
		Nested: Nested{
			Field: "filed1",
		},
	}
	samplemap := map[string]string{
		"90": "Dog",
		"91": "Cat",
		"92": "Cow",
	}
	testCase = append(testCase, emp, samplemap)

	for _, body := range testCase {
		req := (&Request{}).WithRequestBody(body)

		bodyData, _ := json.Marshal(body)

		if !reflect.DeepEqual(req.requestBody, bodyData) {
			t.Error("Test Failed: TestWithRequestBody")
		}
	}
}

func TestDoRequestErr(t *testing.T) {
	req := &Request{
		url:         "http:localhost:8081/",
		httpVerb:    POST,
		requestBody: []byte(`bodySample`),
		headers: map[string]string{
			"key":  "value",
			"key2": "value2",
		},
		basicAuth: basicAuth{
			username: "user",
			password: "pass",
		},
	}

	if _, err := req.doRequest(MockClient{}); err != nil {
		t.Error("Test Failed:TestDoRequestErr ")
	}
	req.requestBody = []byte(``)
	if _, err := req.doRequest(MockClient{}); err != nil {
		t.Error("Test Failed:TestDoRequestErr ")
	}

	if _, err := req.doRequest(MockClient{err: errors.New("Error")}); err == nil {
		t.Error("Test Failed:TestDoRequestErr ")
	}

	req.httpVerb = "\\\\\\"
	if _, err := req.doRequest(MockClient{}); err == nil {
		t.Error("Test Failed:TestDoRequestErr ")
	}
}

func TestBuild(t *testing.T) {
	_, err := RequestBuilder().
		Get().
		WithURL(testURL).
		WithRequestBody(map[string]string{
			"key": "value",
		}).
		WithRequestParams("newKey", "newValue").
		WithRequestParams("newKey2", "newValue2").
		Build()

	if err != nil {
		t.Error(err)
		t.Error("Test Failed: TestBuild")
	}

	_, err = RequestBuilder().
		Get().
		WithURL(testURL).
		WithRequestBodyParams("bodyParam", "bodyparam").
		Build()

	if err != nil {
		t.Error(err)
		t.Error("Test Failed: TestBuild")
	}
}

func TestBuildErr(t *testing.T) {
	var testCase []*Request

	testCase = append(testCase,
		&Request{
			err: errors.New("dummyError"),
		},
		&Request{},
		&Request{
			httpVerb: GET,
		})

	for _, test := range testCase {
		if _, err := test.Build(); err == nil {
			t.Error("Test Failed :TestBuildErr")
		}
	}
}

func TestWithTemplate(t *testing.T) {
	tempReq, err := RequestBuilder().
		Get().
		WithURL(testURL).
		WithRequestBody(map[string]string{
			"key": "value",
		}).
		WithRequestParams("newKey", "newValue").
		WithRequestParams("newKey2", "newValue2").
		Build()
	if err != nil {
		t.Error("Test Failed: TestWithTemplate")
	}

	req := RequestBuilder().WithTemplate(tempReq)
	if !reflect.DeepEqual(req, tempReq) {
		t.Error("Test Failed: TestWithTemplate")
	}
}
