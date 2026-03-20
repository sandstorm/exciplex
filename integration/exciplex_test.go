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

func TestStopCancelsQueuedCallbacks(t *testing.T) {
	body := httpGet(t, "stop_cancels_queued_callbacks.php")
	assertEndsWith(t, body, "done")
	// "callback" must appear exactly once — the CancellableCallback guard must
	// suppress any already-queued invocations after stop() was called.
	count := strings.Count(body, "callback")
	if count != 1 {
		t.Errorf("expected \"callback\" to appear exactly once, got %d times in:\n%s", count, body)
	}
}

// --- Profiler tests ---

func TestProfilerNoSamples(t *testing.T) {
	body := httpGet(t, "profiler_no_samples.php")
	// getLog() should return empty string when no samples collected
	// => The body should be exactly "empty" — no log lines before it
	if body != "empty" {
		t.Errorf("expected body to be exactly %q, got:\n%s", "empty", body)
	}
}

func TestProfilerSingleStack(t *testing.T) {
	body := httpGet(t, "profiler_single_stack.php")

	assertMatches(t, body, regexp.MustCompile(`/app/Web/profiler_single_stack\.php;test 1\n`))
}

func TestProfilerMultipleStacks(t *testing.T) {
	body := httpGet(t, "profiler_multiple_stacks.php")

	log := strings.TrimSuffix(body, "\n")
	lines := strings.Split(log, "\n")

	assertLength(t, lines, 2)
	assertContains(t, lines[0], "/app/Web/profiler_multiple_stacks.php;work_a ")
	assertContains(t, lines[1], "/app/Web/profiler_multiple_stacks.php;work_b ")
}

func TestProfilerStop(t *testing.T) {
	body := httpGet(t, "profiler_stop.php")
	assertContains(t, body, "/app/Web/profiler_stop.php;wait")
	assertContains(t, body, "stop_works")
	assertNotContains(t, body, "FAIL")
}

func TestProfilerCrossFileFunctionCall(t *testing.T) {
	body := httpGet(t, "profiler_cross_file_function_call.php")

	assertContains(t, body, "/app/Web/profiler_cross_file_function_call.php;cross_file_work")
}

func TestProfilerCrossFileNoFunctions(t *testing.T) {
	body := httpGet(t, "profiler_cross_file_no_functions.php")

	assertContains(t, body, "/app/Web/profiler_cross_file_no_functions.php;/app/Web/profiler_cross_file_no_functions_helper.php")
}

func TestProfilerClosure(t *testing.T) {
	body := httpGet(t, "profiler_closure.php")

	assertMatches(t, body, regexp.MustCompile(
		`/app/Web/profiler_closure\.php;\{closure:/app/Web/profiler_closure\.php\(3\)\} \d+`,
	))
}

func TestProfilerLogFormat(t *testing.T) {
	body := httpGet(t, "profiler_log_format.php")

	lines := strings.Split(strings.Trim(body, "\n"), "\n")
	assertLength(t, lines, 3)
	assertMatches(t, lines[0], regexp.MustCompile(`/app/Web/profiler_log_format\.php;top_call;middle_call;deep_call [234]`))
	assertMatches(t, lines[1], regexp.MustCompile(`/app/Web/profiler_log_format\.php;TestClass::method;middle_call;deep_call [234]`))
	assertMatches(t, lines[2], regexp.MustCompile(`/app/Web/profiler_log_format\.php;TestClass::staticMethod;deep_call [234]`))
}
