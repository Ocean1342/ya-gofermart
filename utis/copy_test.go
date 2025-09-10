package utis

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloneRequest(t *testing.T) {
	// Создаём оригинальный запрос с телом
	originalBody := "test request body"
	originalReq, err := http.NewRequest(
		"POST",
		"https://example.com/api",
		bytes.NewBufferString(originalBody),
	)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Устанавливаем заголовки и другие параметры
	originalReq.Header.Set("Content-Type", "application/json")
	originalReq.Header.Set("X-Test", "123")
	originalReq.URL.RawQuery = "param=value"

	// Клонируем запрос
	clonedReq, err := CopyHTTPRequest(originalReq)
	if err != nil {
		t.Fatalf("cloneRequest failed: %v", err)
	}

	// --- Проверяем, что клонированный запрос корректен ---

	// 1. Проверяем, что метод и URL совпадают
	if clonedReq.Method != originalReq.Method {
		t.Errorf("Method mismatch: got %s, want %s", clonedReq.Method, originalReq.Method)
	}
	if clonedReq.URL.String() != originalReq.URL.String() {
		t.Errorf("URL mismatch: got %s, want %s", clonedReq.URL.String(), originalReq.URL.String())
	}

	// 2. Проверяем заголовки
	if clonedReq.Header.Get("Content-Type") != originalReq.Header.Get("Content-Type") {
		t.Errorf("Content-Type header mismatch")
	}
	if clonedReq.Header.Get("X-Test") != originalReq.Header.Get("X-Test") {
		t.Errorf("X-Test header mismatch")
	}

	// 3. Проверяем, что тело скопировано и независимо
	clonedBody, err := io.ReadAll(clonedReq.Body)
	if err != nil {
		t.Fatalf("Failed to read cloned body: %v", err)
	}
	if string(clonedBody) != originalBody {
		t.Errorf("Cloned body mismatch: got %s, want %s", string(clonedBody), originalBody)
	}

	// 4. Проверяем, что оригинальное тело осталось нетронутым
	originalBodyAfterClone, err := io.ReadAll(originalReq.Body)
	if err != nil {
		t.Fatalf("Failed to read original body after clone: %v", err)
	}
	if string(originalBodyAfterClone) != originalBody {
		t.Errorf("Original body was modified after clone: got %s, want %s", string(originalBodyAfterClone), originalBody)
	}

	// 5. Проверяем, что можно прочитать тело клонированного запроса повторно (если поддерживается GetBody)
	if clonedReq.GetBody != nil {
		newBody, err := clonedReq.GetBody()
		if err != nil {
			t.Fatalf("GetBody failed: %v", err)
		}
		defer newBody.Close()

		newBodyContent, err := io.ReadAll(newBody)
		if err != nil {
			t.Fatalf("Failed to read GetBody content: %v", err)
		}
		if string(newBodyContent) != originalBody {
			t.Errorf("GetBody content mismatch: got %s, want %s", string(newBodyContent), originalBody)
		}
	}
}

func TestCopyHTTPResponse(t *testing.T) {
	// Создаём тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "123")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response body"))
	}))
	defer ts.Close()

	// Делаем запрос и получаем оригинальный Response
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Клонируем Response
	clonedResp, err := CopyHTTPResponse(resp)
	if err != nil {
		t.Fatalf("CopyHTTPResponse failed: %v", err)
	}
	defer clonedResp.Body.Close()

	// Проверяем, что статус совпадает
	if clonedResp.StatusCode != resp.StatusCode {
		t.Errorf("StatusCode mismatch: got %d, want %d", clonedResp.StatusCode, resp.StatusCode)
	}

	// Проверяем заголовки
	if clonedResp.Header.Get("X-Test") != resp.Header.Get("X-Test") {
		t.Errorf("Header mismatch: got %s, want %s", clonedResp.Header.Get("X-Test"), resp.Header.Get("X-Test"))
	}

	// Проверяем тело ответа
	originalBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read original body: %v", err)
	}

	clonedBody, err := io.ReadAll(clonedResp.Body)
	if err != nil {
		t.Fatalf("Failed to read cloned body: %v", err)
	}

	if string(clonedBody) != string(originalBody) {
		t.Errorf("Body mismatch: got %s, want %s", string(clonedBody), string(originalBody))
	}

	// Проверяем, что оригинальное тело осталось доступным
	if len(originalBody) == 0 {
		t.Error("Original body was consumed by cloning")
	}
}
