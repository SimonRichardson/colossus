package statsd

import (
	"fmt"
	"time"

	"github.com/SimonRichardson/colossus/instrumentation"
	"github.com/peterbourgon/g2s"
)

type instrument struct {
	statter    g2s.Statter
	sampleRate float32
}

func New(statter g2s.Statter, sampleRate float32) instrumentation.Instrumentation {
	return instrument{statter, sampleRate}
}

func (i instrument) ClusterCall(n int) {
	i.statter.Counter(i.sampleRate, fmt.Sprintf("cluster.%d.call.count", n), 1)
}

func (i instrument) ClusterDuration(n int, t time.Duration) {
	i.statter.Timing(i.sampleRate, fmt.Sprintf("cluster.%d.duration", n), t)
}

func (i instrument) AInsertCall() {
	i.statter.Counter(i.sampleRate, "aggregate_insert.call.count", 1)
}

func (i instrument) AInsertDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_insert.duration", t)
}

func (i instrument) AModifyCall() {
	i.statter.Counter(i.sampleRate, "aggregate_modify.call.count", 1)
}

func (i instrument) AModifyDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_modify.duration", t)
}

func (i instrument) AModifyWithOperationsCall() {
	i.statter.Counter(i.sampleRate, "aggregate_modify_with_operations.call.count", 1)
}

func (i instrument) AModifyWithOperationsDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_modify_with_operations.duration", t)
}

func (i instrument) ADeleteCall() {
	i.statter.Counter(i.sampleRate, "aggregate_delete.call.count", 1)
}

func (i instrument) ADeleteDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_delete.duration", t)
}

func (i instrument) ARollbackCall() {
	i.statter.Counter(i.sampleRate, "aggregate_rollback.call.count", 1)
}

func (i instrument) ARollbackDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_rollback.duration", t)
}

func (i instrument) ASelectCall() {
	i.statter.Counter(i.sampleRate, "aggregate_select.call.count", 1)
}

func (i instrument) ASelectDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_select.duration", t)
}

func (i instrument) ASelectRangeCall() {
	i.statter.Counter(i.sampleRate, "aggregate_select_range.call.count", 1)
}

func (i instrument) ASelectRangeDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_select_range.duration", t)
}

func (i instrument) AKeysCall() {
	i.statter.Counter(i.sampleRate, "aggregate_keys.call.count", 1)
}

func (i instrument) AKeysDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_keys.duration", t)
}

func (i instrument) ASizeCall() {
	i.statter.Counter(i.sampleRate, "aggregate_size.call.count", 1)
}

func (i instrument) ASizeDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_size.duration", t)
}

func (i instrument) AMembersCall() {
	i.statter.Counter(i.sampleRate, "aggregate_members.call.count", 1)
}

func (i instrument) AMembersDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_members.duration", t)
}

func (i instrument) ARepairCall() {
	i.statter.Counter(i.sampleRate, "aggregate_repair.call.count", 1)
}

func (i instrument) ARepairDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_repair.duration", t)
}

func (i instrument) AQueryCall() {
	i.statter.Counter(i.sampleRate, "aggregate_query.call.count", 1)
}

func (i instrument) AQueryDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_query.duration", t)
}

func (i instrument) APauseCall() {
	i.statter.Counter(i.sampleRate, "aggregate_pause.call.count", 1)
}

func (i instrument) AResumeCall() {
	i.statter.Counter(i.sampleRate, "aggregate_resume.call.count", 1)
}

func (i instrument) ATopologyCall() {
	i.statter.Counter(i.sampleRate, "aggregate_topology.call.count", 1)
}

func (i instrument) ATopologyDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "aggregate_topology.duration", t)
}

func (i instrument) InsertCall() {
	i.statter.Counter(i.sampleRate, "insert.call.count", 1)
}

func (i instrument) InsertKeys(n int) {
	i.statter.Counter(i.sampleRate, "insert.keys.count", n)
}

func (i instrument) InsertSendTo(n int) {
	i.statter.Counter(i.sampleRate, "insert.send_to.count", n)
}

func (i instrument) InsertDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "insert.duration", t)
}

func (i instrument) InsertRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "insert.retrieved.count", n)
}

func (i instrument) InsertReturned(n int) {
	i.statter.Counter(i.sampleRate, "insert.returned.count", n)
}

func (i instrument) InsertQuorumFailure() {
	i.statter.Counter(i.sampleRate, "insert.quorum_failure.count", 1)
}

func (i instrument) InsertRepairRequired() {
	i.statter.Counter(i.sampleRate, "insert.repair_required.count", 1)
}

func (i instrument) InsertPartialFailure() {
	i.statter.Counter(i.sampleRate, "insert.partial_failure.count", 1)
}

func (i instrument) ModifyCall() {
	i.statter.Counter(i.sampleRate, "modify.call.count", 1)
}

func (i instrument) ModifyKeys(n int) {
	i.statter.Counter(i.sampleRate, "modify.keys.count", n)
}

func (i instrument) ModifySendTo(n int) {
	i.statter.Counter(i.sampleRate, "modify.send_to.count", n)
}

func (i instrument) ModifyDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "modify.duration", t)
}

func (i instrument) ModifyRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "modify.retrieved.count", n)
}

func (i instrument) ModifyReturned(n int) {
	i.statter.Counter(i.sampleRate, "modify.returned.count", n)
}

func (i instrument) ModifyQuorumFailure() {
	i.statter.Counter(i.sampleRate, "modify.quorum_failure.count", 1)
}

func (i instrument) ModifyRepairRequired() {
	i.statter.Counter(i.sampleRate, "modify.repair_required.count", 1)
}

func (i instrument) DeleteCall() {
	i.statter.Counter(i.sampleRate, "delete.call.count", 1)
}

func (i instrument) DeleteKeys(n int) {
	i.statter.Counter(i.sampleRate, "delete.keys.count", n)
}

func (i instrument) DeleteSendTo(n int) {
	i.statter.Counter(i.sampleRate, "delete.send_to.count", n)
}

func (i instrument) DeleteDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "delete.duration", t)
}

func (i instrument) DeleteRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "delete.retrieved.count", n)
}

func (i instrument) DeleteReturned(n int) {
	i.statter.Counter(i.sampleRate, "delete.returned.count", n)
}

func (i instrument) DeleteQuorumFailure() {
	i.statter.Counter(i.sampleRate, "delete.quorum_failure.count", 1)
}

func (i instrument) DeleteRepairRequired() {
	i.statter.Counter(i.sampleRate, "delete.repair_required.count", 1)
}

func (i instrument) DeletePartialFailure() {
	i.statter.Counter(i.sampleRate, "delete.partial_failure.count", 1)
}

func (i instrument) RollbackCall() {
	i.statter.Counter(i.sampleRate, "rollback.call.count", 1)
}

func (i instrument) RollbackKeys(n int) {
	i.statter.Counter(i.sampleRate, "rollback.keys.count", n)
}

func (i instrument) RollbackSendTo(n int) {
	i.statter.Counter(i.sampleRate, "rollback.send_to.count", n)
}

func (i instrument) RollbackDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "rollback.duration", t)
}

func (i instrument) RollbackRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "rollback.retrieved.count", n)
}

func (i instrument) RollbackReturned(n int) {
	i.statter.Counter(i.sampleRate, "rollback.returned.count", n)
}

func (i instrument) RollbackQuorumFailure() {
	i.statter.Counter(i.sampleRate, "rollback.quorum_failure.count", 1)
}

func (i instrument) RollbackRepairRequired() {
	i.statter.Counter(i.sampleRate, "rollback.repair_required.count", 1)
}

func (i instrument) RollbackPartialFailure() {
	i.statter.Counter(i.sampleRate, "rollback.partial_failure.count", 1)
}

func (i instrument) SelectCall() {
	i.statter.Counter(i.sampleRate, "select.call.count", 1)
}

func (i instrument) SelectKeys(n int) {
	i.statter.Counter(i.sampleRate, "select.keys.count", n)
}

func (i instrument) SelectSendTo(n int) {
	i.statter.Counter(i.sampleRate, "select.send_to.count", n)
}

func (i instrument) SelectSendAllPromotion() {
	i.statter.Counter(i.sampleRate, "select.send_all_promotion.count", 1)
}

func (i instrument) SelectPartialError() {
	i.statter.Counter(i.sampleRate, "select.partial_error.count", 1)
}

func (i instrument) SelectDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "select.duration", t)
}

func (i instrument) SelectRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "select.retrieved.count", n)
}

func (i instrument) SelectReturned(n int) {
	i.statter.Counter(i.sampleRate, "select.returned.count", n)
}

func (i instrument) SelectFirstResponseDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "select.first_response_duration", t)
}

func (i instrument) SelectBlockingDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "select.blocking_duration", t)
}

func (i instrument) SelectOverheadDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "select.overhead_duration", t)
}

func (i instrument) SelectRepairNeeded() {
	i.statter.Counter(i.sampleRate, "scan.select_repair_needed.count", 1)
}

func (i instrument) ScanCall() {
	i.statter.Counter(i.sampleRate, "scan.call.count", 1)
}

func (i instrument) ScanSendTo(n int) {
	i.statter.Counter(i.sampleRate, "scan.send_to.count", n)
}

func (i instrument) ScanPartialError() {
	i.statter.Counter(i.sampleRate, "scan.partial_error.count", 1)
}

func (i instrument) ScanDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "scan.duration", t)
}

func (i instrument) ScanRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "scan.retrieved.count", n)
}

func (i instrument) ScanReturned(n int) {
	i.statter.Counter(i.sampleRate, "scan.returned.count", n)
}

func (i instrument) ScanRepairNeeded(n int) {
	i.statter.Counter(i.sampleRate, "scan.repair_needed.count", n)
}

func (i instrument) RepairCall() {
	i.statter.Counter(i.sampleRate, "repair.call.count", 1)
}

func (i instrument) RepairRequest(n int) {
	i.statter.Counter(i.sampleRate, "repair.request.count", n)
}

func (i instrument) RepairSendTo(n int) {
	i.statter.Counter(i.sampleRate, "repair.send_to.count", n)
}

func (i instrument) RepairDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "repair.duration", t)
}

func (i instrument) RepairScoreError() {
	i.statter.Counter(i.sampleRate, "repair.score_error.count", 1)
}

func (i instrument) RepairError(n int) {
	i.statter.Counter(i.sampleRate, "repair.error.count", n)
}

func (i instrument) PerformanceDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "performance.duration", t)
}

func (i instrument) PerformanceNamespaceDuration(ns string, t time.Duration) {
	i.statter.Timing(i.sampleRate, fmt.Sprintf("performance.%s.duration", ns), t)
}

func (i instrument) PublishCall() {
	i.statter.Counter(i.sampleRate, "publish.call.count", 1)
}

func (i instrument) PublishKeys(n int) {
	i.statter.Counter(i.sampleRate, "publish.keys.count", n)
}

func (i instrument) PublishSendTo(n int) {
	i.statter.Counter(i.sampleRate, "publish.sent_to.count", n)
}

func (i instrument) PublishRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "publish.retrieved.count", n)
}

func (i instrument) PublishReturned(n int) {
	i.statter.Counter(i.sampleRate, "publish.returned.count", n)
}

func (i instrument) PublishDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "publish.duration", t)
}

func (i instrument) SemaphoreCall() {
	i.statter.Counter(i.sampleRate, "semaphore.call.count \n", 1)
}
func (i instrument) SemaphoreSendTo(n int) {
	i.statter.Counter(i.sampleRate, "semaphore.send_to.count %d\n", n)
}
func (i instrument) SemaphoreDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "semaphore.duration %d\n", t)
}
func (i instrument) SemaphoreRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "semaphore.retrieved.count %d\n", n)
}
func (i instrument) SemaphoreReturned(n int) {
	i.statter.Counter(i.sampleRate, "semaphore.returned.count %d\n", n)
}

func (i instrument) HeartbeatCall() {
	i.statter.Counter(i.sampleRate, "heartbeat.call.count \n", 1)
}
func (i instrument) HeartbeatSendTo(n int) {
	i.statter.Counter(i.sampleRate, "heartbeat.send_to.count %d\n", n)
}
func (i instrument) HeartbeatDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "heartbeat.duration %d\n", t)
}
func (i instrument) HeartbeatRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "heartbeat.retrieved.count %d\n", n)
}
func (i instrument) HeartbeatReturned(n int) {
	i.statter.Counter(i.sampleRate, "heartbeat.returned.count %d\n", n)
}

func (i instrument) KeyStoreCall() {
	i.statter.Counter(i.sampleRate, "keystore.call.count \n", 1)
}
func (i instrument) KeyStoreSendTo(n int) {
	i.statter.Counter(i.sampleRate, "keystore.send_to.count %d\n", n)
}
func (i instrument) KeyStoreDuration(t time.Duration) {
	i.statter.Timing(i.sampleRate, "keystore.duration %d\n", t)
}
func (i instrument) KeyStoreRetrieved(n int) {
	i.statter.Counter(i.sampleRate, "keystore.retrieved.count %d\n", n)
}
func (i instrument) KeyStoreReturned(n int) {
	i.statter.Counter(i.sampleRate, "keystore.returned.count %d\n", n)
}
