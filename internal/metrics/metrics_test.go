package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/shoenig/test/must"
)

func TestIncProcessedRecord(t *testing.T) {
	must.Zero(t, testutil.CollectAndCount(records_processed))
	IncProcessedRecord("test_input", "test")
	must.Eq(t, 1, testutil.CollectAndCount(records_processed))

	must.Eq(t, 1.0, testutil.ToFloat64(records_processed.WithLabelValues("test_input", "test")))
}
