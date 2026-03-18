package integration

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

const baseURL = "http://localhost:18181"

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

func assertMatches(t *testing.T, body string, regex *regexp.Regexp) {
	t.Helper()
	if !regex.MatchString(body) {
		t.Errorf("expected response to match %q, got:\n%s", regex, body)
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

func assertLength[T any](t *testing.T, slice []T, expected int) {
	t.Helper()
	if len(slice) != expected {
		t.Errorf("expected to have length %d, got %d", expected, len(slice))
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

// --- Profiler tests ---

func TestProfilerNoSamples(t *testing.T) {
	body := httpGet(t, "profiler_no_samples.php")
	// getLog() should return empty string when no samples collected
	assertEndsWith(t, body, "empty")
	// The body should be exactly "empty" — no log lines before it
	if body != "empty" {
		t.Errorf("expected body to be exactly %q, got:\n%s", "empty", body)
	}
}

func TestProfilerSingleStack(t *testing.T) {
	body := httpGet(t, "profiler_single_stack.php")
	assertEndsWith(t, body, "---END---")

	// Extract the log portion (before ---END---)
	log := strings.TrimSuffix(body, "---END---")

	// Should contain the function names from the call stack
	assertContains(t, log, "inner")
	assertContains(t, log, "outer")

	// Should have at least one line in flamegraph format
	lines := filterNonEmpty(strings.Split(log, "\n"))
	if len(lines) < 1 {
		t.Fatalf("expected at least 1 log line, got %d", len(lines))
	}
}

func TestProfilerMultipleStacks(t *testing.T) {
	body := httpGet(t, "profiler_multiple_stacks.php")
	assertEndsWith(t, body, "---END---")

	log := strings.TrimSuffix(body, "---END---")

	// Should contain both work functions
	assertContains(t, log, "work_a")
	assertContains(t, log, "work_b")

	// Should have at least 2 distinct stack lines
	lines := filterNonEmpty(strings.Split(log, "\n"))
	if len(lines) < 2 {
		t.Errorf("expected at least 2 distinct log lines, got %d:\n%s", len(lines), log)
	}
}

func TestProfilerStop(t *testing.T) {
	body := httpGet(t, "profiler_stop.php")
	assertContains(t, body, "stop_works")
	assertNotContains(t, body, "FAIL")
}

func TestProfilerLogFormat(t *testing.T) {
	body := httpGet(t, "profiler_log_format.php")

	lines := strings.Split(strings.Trim(body, "\n"), "\n")
	assertLength(t, lines, 3)
	assertMatches(t, lines[0], regexp.MustCompile("top_call;middle_call;deep_call [234]"))
	assertMatches(t, lines[1], regexp.MustCompile("TestClass::method;middle_call;deep_call [234]"))
	assertMatches(t, lines[2], regexp.MustCompile("TestClass::staticMethod;deep_call [234]"))
}

func filterNonEmpty(ss []string) []string {
	var result []string
	for _, s := range ss {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
