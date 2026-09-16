package directsend

import (
	"testing"
	"testing/synctest"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	synctest.Test(t, func(t *testing.T) {
		got["len_after_send_with_parked_receiver"] = BufferAfterSendToParkedReceiver()
	})
	synctest.Test(t, func(t *testing.T) {
		got["len_after_send_with_no_receiver"] = BufferAfterSendWithNoReceiver()
	})
	synctest.Test(t, func(t *testing.T) {
		got["full_buffer_and_parked_receiver_possible"] = FullBufferKeepsParkedReceiver()
	})
	synctest.Test(t, func(t *testing.T) {
		got["value_received_when_buffer_full_and_sender_parked"] = ReceiveWithFullBufferAndParkedSender()
	})

	predict.Check(t, got)
}
