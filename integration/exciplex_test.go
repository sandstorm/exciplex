package integration

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

const baseURL = "http://localhost:18181"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func httpGet(t *testing.T, path string) string {
	t.Helper()
	resp, err := http.Get(baseURL + "/" + path)
	if err != nil {
		t.Fatalf("HTTP request to %s failed: %v", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d for %s, body: %s", resp.StatusCode, path, string(body))
	}
	return string(body)
}

func assertContains(t *testing.T, body, substr string) {
	t.Helper()
	if !strings.Contains(body, substr) {
		t.Errorf("expected response to contain %q, got:\n%s", substr, body)
	}
}

func assertEndsWith(t *testing.T, body, suffix string) {
	t.Helper()
	if !strings.HasSuffix(body, suffix) {
		t.Errorf("expected response to end with %q, got:\n%s", suffix, body)
	}
}

func assertNotContains(t *testing.T, body, substr string) {
	t.Helper()
	if strings.Contains(body, substr) {
		t.Errorf("expected response NOT to contain %q, got:\n%s", substr, body)
	}
}

func TestInterruptIsCalled(t *testing.T) {
	body := httpGet(t, "interrupt_is_called.php")
	assertContains(t, body, "interrupt was called")
	assertContains(t, body, " 0, ")
	assertContains(t, body, " 9, ")
}

func TestDieInInterruptEndsExecution(t *testing.T) {
	body := httpGet(t, "die_in_interrupt_ends_execution.php")
	assertEndsWith(t, body, "last visible text, no further counting")
	assertNotContains(t, body, " 9, ")
}

func TestIntervalIsCalledRepeatedly(t *testing.T) {
	body := httpGet(t, "interval_is_called_repeatedly.php")
	assertContains(t, body, "interval 0")
	assertContains(t, body, "interval 1")
	assertContains(t, body, "interval 2")
	assertContains(t, body, " 9, ")
}

func TestStopTimeoutBeforeFired(t *testing.T) {
	body := httpGet(t, "stop_timeout_before_fired.php")
	assertNotContains(t, body, "SHOULD NOT APPEAR")
	assertContains(t, body, " 9, ")
	assertEndsWith(t, body, "done")
}

func TestStopIntervalBeforeFired(t *testing.T) {
	body := httpGet(t, "stop_interval_before_fired.php")
	assertNotContains(t, body, "SHOULD NOT APPEAR")
	assertContains(t, body, " 9, ")
	assertEndsWith(t, body, "done")
}

func TestStopIntervalAfterOneCall(t *testing.T) {
	body := httpGet(t, "stop_interval_after_one_call.php")
	assertContains(t, body, "interval 0")
	assertNotContains(t, body, "interval 1")
	assertContains(t, body, " 9, ")
	assertEndsWith(t, body, "done")
}
