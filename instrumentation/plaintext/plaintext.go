package plaintext

import (
	"fmt"
	"io"
	"time"

	"github.com/SimonRichardson/colossus/instrumentation"
)

type instrument struct{ io.Writer }

func New(w io.Writer) instrumentation.Instrumentation {
	return instrument{w}
}

func (i instrument) ClusterCall(n int) {
	fmt.Fprintf(i, "cluster.%d.call.count %d\n", n, 1)
}

func (i instrument) ClusterDuration(n int, t time.Duration) {
	fmt.Fprintf(i, "cluster.%d.call.duration %d\n", n, t.Nanoseconds()/1e6)
}

func (i instrument) AInsertCall() {
	fmt.Fprintf(i, "aggregate_insert.call.count 1\n")
}

func (i instrument) AInsertDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_insert.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) AModifyCall() {
	fmt.Fprintf(i, "aggregate_modify.call.count 1\n")
}

func (i instrument) AModifyDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_modify.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) AModifyWithOperationsCall() {
	fmt.Fprintf(i, "aggregate_modify_with_operations.call.count 1\n")
}

func (i instrument) AModifyWithOperationsDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_modify_with_operations.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ADeleteCall() {
	fmt.Fprintf(i, "aggregate_delete.call.count 1\n")
}

func (i instrument) ADeleteDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_delete.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ARollbackCall() {
	fmt.Fprintf(i, "aggregate_rollback.call.count 1\n")
}

func (i instrument) ARollbackDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_rollback.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ASelectCall() {
	fmt.Fprintf(i, "aggregate_select.call.count 1\n")
}

func (i instrument) ASelectDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_select.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ASelectRangeCall() {
	fmt.Fprintf(i, "aggregate_select_range.call.count 1\n")
}

func (i instrument) ASelectRangeDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_select_range.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) AKeysCall() {
	fmt.Fprintf(i, "aggregate_keys.call.count 1\n")
}

func (i instrument) AKeysDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_keys.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ASizeCall() {
	fmt.Fprintf(i, "aggregate_size.call.count 1\n")
}

func (i instrument) ASizeDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_size.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) AMembersCall() {
	fmt.Fprintf(i, "aggregate_members.call.count 1\n")
}

func (i instrument) AMembersDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_members.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ARepairCall() {
	fmt.Fprintf(i, "aggregate_repair.call.count 1\n")
}

func (i instrument) ARepairDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_repair.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) AQueryCall() {
	fmt.Fprintf(i, "aggregate_query.call.count 1\n")
}

func (i instrument) AQueryDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_query.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) APauseCall() {
	fmt.Fprintf(i, "aggregate_pause.call.count 1\n")
}

func (i instrument) AResumeCall() {
	fmt.Fprintf(i, "aggregate_resume.call.count 1\n")
}

func (i instrument) ATopologyCall() {
	fmt.Fprintf(i, "aggregate_topology.call.count 1\n")
}

func (i instrument) ATopologyDuration(t time.Duration) {
	fmt.Fprintf(i, "aggregate_topology.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) InsertCall() {
	fmt.Fprintf(i, "insert.call.count 1\n")
}

func (i instrument) InsertKeys(n int) {
	fmt.Fprintf(i, "insert.keys.count %d\n", n)
}

func (i instrument) InsertSendTo(n int) {
	fmt.Fprintf(i, "insert.send_to.count %d\n", n)
}

func (i instrument) InsertDuration(t time.Duration) {
	fmt.Fprintf(i, "insert.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) InsertRetrieved(n int) {
	fmt.Fprintf(i, "insert.retrieved.count %d\n", n)
}

func (i instrument) InsertReturned(n int) {
	fmt.Fprintf(i, "insert.returned.count %d\n", n)
}

func (i instrument) InsertQuorumFailure() {
	fmt.Fprintf(i, "insert.quorum_failure.count 1\n")
}

func (i instrument) InsertRepairRequired() {
	fmt.Fprintf(i, "insert.repair_required.count 1\n")
}

func (i instrument) InsertPartialFailure() {
	fmt.Fprintf(i, "insert.partial_failure.count 1\n")
}

func (i instrument) ModifyCall() {
	fmt.Fprintf(i, "modify.call.count 1\n")
}

func (i instrument) ModifyKeys(n int) {
	fmt.Fprintf(i, "modify.keys.count %d\n", n)
}

func (i instrument) ModifySendTo(n int) {
	fmt.Fprintf(i, "modify.send_to.count %d\n", n)
}

func (i instrument) ModifyDuration(t time.Duration) {
	fmt.Fprintf(i, "modify.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ModifyRetrieved(n int) {
	fmt.Fprintf(i, "modify.retrieved.count %d\n", n)
}

func (i instrument) ModifyReturned(n int) {
	fmt.Fprintf(i, "modify.returned.count %d\n", n)
}

func (i instrument) ModifyQuorumFailure() {
	fmt.Fprintf(i, "modify.quorum_failure.count 1\n")
}

func (i instrument) ModifyRepairRequired() {
	fmt.Fprintf(i, "modify.repair_required.count 1\n")
}

func (i instrument) DeleteCall() {
	fmt.Fprintf(i, "delete.call.count 1\n")
}

func (i instrument) DeleteKeys(n int) {
	fmt.Fprintf(i, "delete.keys.count %d\n", n)
}

func (i instrument) DeleteSendTo(n int) {
	fmt.Fprintf(i, "delete.send_to.count %d\n", n)
}

func (i instrument) DeleteDuration(t time.Duration) {
	fmt.Fprintf(i, "delete.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) DeleteRetrieved(n int) {
	fmt.Fprintf(i, "delete.retrieved.count %d\n", n)
}

func (i instrument) DeleteReturned(n int) {
	fmt.Fprintf(i, "delete.returned.count %d\n", n)
}

func (i instrument) DeleteQuorumFailure() {
	fmt.Fprintf(i, "delete.quorum_failure.count 1\n")
}

func (i instrument) DeleteRepairRequired() {
	fmt.Fprintf(i, "delete.repair_required.count 1\n")
}

func (i instrument) DeletePartialFailure() {
	fmt.Fprintf(i, "delete.partial_failure.count 1\n")
}

func (i instrument) RollbackCall() {
	fmt.Fprintf(i, "rollback.call.count 1\n")
}

func (i instrument) RollbackKeys(n int) {
	fmt.Fprintf(i, "rollback.keys.count %d\n", n)
}

func (i instrument) RollbackSendTo(n int) {
	fmt.Fprintf(i, "rollback.send_to.count %d\n", n)
}

func (i instrument) RollbackDuration(t time.Duration) {
	fmt.Fprintf(i, "rollback.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) RollbackRetrieved(n int) {
	fmt.Fprintf(i, "rollback.retrieved.count %d\n", n)
}

func (i instrument) RollbackReturned(n int) {
	fmt.Fprintf(i, "rollback.returned.count %d\n", n)
}

func (i instrument) RollbackQuorumFailure() {
	fmt.Fprintf(i, "rollback.quorum_failure.count 1\n")
}

func (i instrument) RollbackRepairRequired() {
	fmt.Fprintf(i, "rollback.repair_required.count 1\n")
}

func (i instrument) RollbackPartialFailure() {
	fmt.Fprintf(i, "rollback.partial_failure.count 1\n")
}

func (i instrument) SelectCall() {
	fmt.Fprintf(i, "select.call.count 1\n")
}

func (i instrument) SelectKeys(n int) {
	fmt.Fprintf(i, "select.keys.count %d\n", n)
}

func (i instrument) SelectSendTo(n int) {
	fmt.Fprintf(i, "select.send_to.count %d\n", n)
}

func (i instrument) SelectSendAllPromotion() {
	fmt.Fprintf(i, "select.send_all_promotion.count 1\n")
}

func (i instrument) SelectPartialError() {
	fmt.Fprintf(i, "select.partial_error.count 1\n")
}

func (i instrument) SelectDuration(t time.Duration) {
	fmt.Fprintf(i, "select.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) SelectRetrieved(n int) {
	fmt.Fprintf(i, "select.retrieved.count %d\n", n)
}

func (i instrument) SelectReturned(n int) {
	fmt.Fprintf(i, "select.returned.count %d\n", n)
}

func (i instrument) SelectFirstResponseDuration(t time.Duration) {
	fmt.Fprintf(i, "select.first_response_duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) SelectBlockingDuration(t time.Duration) {
	fmt.Fprintf(i, "select.blocking_duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) SelectOverheadDuration(t time.Duration) {
	fmt.Fprintf(i, "select.overhead_duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) SelectRepairNeeded() {
	fmt.Fprintf(i, "scan.select_repair_needed.count 1\n")
}

func (i instrument) ScanCall() {
	fmt.Fprintf(i, "scan.call.count 1\n")
}

func (i instrument) ScanSendTo(n int) {
	fmt.Fprintf(i, "scan.send_to.count %d\n", n)
}

func (i instrument) ScanPartialError() {
	fmt.Fprintf(i, "scan.partial_error.count 1\n")
}

func (i instrument) ScanDuration(t time.Duration) {
	fmt.Fprintf(i, "scan.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) ScanRetrieved(n int) {
	fmt.Fprintf(i, "scan.retrieved.count %d\n", n)
}

func (i instrument) ScanReturned(n int) {
	fmt.Fprintf(i, "scan.returned.count %d\n", n)
}

func (i instrument) ScanRepairNeeded(n int) {
	fmt.Fprintf(i, "scan.repair_needed.count %d\n", n)
}

func (i instrument) RepairCall() {
	fmt.Fprintf(i, "repair.call.count 1\n")
}

func (i instrument) RepairRequest(n int) {
	fmt.Fprintf(i, "repair.request.count %d\n", n)
}

func (i instrument) RepairSendTo(n int) {
	fmt.Fprintf(i, "repair.send_to.count %d\n", n)
}

func (i instrument) RepairDuration(t time.Duration) {
	fmt.Fprintf(i, "repair.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) RepairScoreError() {
	fmt.Fprintf(i, "repair.score_error.count 1\n")
}

func (i instrument) RepairError(n int) {
	fmt.Fprintf(i, "repair.error.count %d\n", n)
}

func (i instrument) PerformanceDuration(t time.Duration) {
	fmt.Fprintf(i, "performance.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) PerformanceNamespaceDuration(ns string, t time.Duration) {
	fmt.Fprintf(i, "performance.%s.duration %d\n", ns, t.Nanoseconds()/1e6)
}

func (i instrument) PublishCall() {
	fmt.Fprintf(i, "publish.call.count 1\n")
}

func (i instrument) PublishKeys(n int) {
	fmt.Fprintf(i, "publish.keys.count %d\n", n)
}

func (i instrument) PublishSendTo(n int) {
	fmt.Fprintf(i, "publish.sent_to.count %d\n", n)
}

func (i instrument) PublishRetrieved(n int) {
	fmt.Fprintf(i, "publish.retrieved.count %d\n", n)
}

func (i instrument) PublishReturned(n int) {
	fmt.Fprintf(i, "publish.returned.count %d\n", n)
}

func (i instrument) PublishDuration(t time.Duration) {
	fmt.Fprintf(i, "publish.duration %d\n", t.Nanoseconds()/1e6)
}

func (i instrument) SemaphoreCall() {
	fmt.Fprintf(i, "semaphore.call.count 1\n")
}
func (i instrument) SemaphoreSendTo(n int) {
	fmt.Fprintf(i, "semaphore.send_to.count %d\n", n)
}
func (i instrument) SemaphoreDuration(t time.Duration) {
	fmt.Fprintf(i, "semaphore.duration %d\n", t.Nanoseconds()/1e6)
}
func (i instrument) SemaphoreRetrieved(n int) {
	fmt.Fprintf(i, "semaphore.retrieved.count %d\n", n)
}
func (i instrument) SemaphoreReturned(n int) {
	fmt.Fprintf(i, "semaphore.returned.count %d\n", n)
}

func (i instrument) HeartbeatCall() {
	fmt.Fprintf(i, "heartbeat.call.count 1\n")
}
func (i instrument) HeartbeatSendTo(n int) {
	fmt.Fprintf(i, "heartbeat.send_to.count %d\n", n)
}
func (i instrument) HeartbeatDuration(t time.Duration) {
	fmt.Fprintf(i, "heartbeat.duration %d\n", t.Nanoseconds()/1e6)
}
func (i instrument) HeartbeatRetrieved(n int) {
	fmt.Fprintf(i, "heartbeat.retrieved.count %d\n", n)
}
func (i instrument) HeartbeatReturned(n int) {
	fmt.Fprintf(i, "heartbeat.returned.count %d\n", n)
}

func (i instrument) KeyStoreCall() {
	fmt.Fprintf(i, "keystore.call.count 1\n")
}
func (i instrument) KeyStoreSendTo(n int) {
	fmt.Fprintf(i, "keystore.send_to.count %d\n", n)
}
func (i instrument) KeyStoreDuration(t time.Duration) {
	fmt.Fprintf(i, "keystore.duration %d\n", t.Nanoseconds()/1e6)
}
func (i instrument) KeyStoreRetrieved(n int) {
	fmt.Fprintf(i, "keystore.retrieved.count %d\n", n)
}
func (i instrument) KeyStoreReturned(n int) {
	fmt.Fprintf(i, "keystore.returned.count %d\n", n)
}
