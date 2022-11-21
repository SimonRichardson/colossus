package redis

import (
	"fmt"
	"sync"
	"time"

	p "github.com/SimonRichardson/colossus/fusion/redis"
	"github.com/SimonRichardson/colossus/instrumentation"
	"github.com/SimonRichardson/colossus/teleprinter"
	r "github.com/garyburd/redigo/redis"
)

const (
	defaultInstrumentationKey = "15daf1eca73c2e497feefe1285cf28dd"
)

type instrument struct {
	mutex  sync.Mutex
	buffer []string
}

func New(pool *p.Pool, maxBufferDuration time.Duration) instrumentation.Instrumentation {
	instr := &instrument{
		mutex:  sync.Mutex{},
		buffer: make([]string, 0),
	}

	go func() {
		tick := time.Tick(maxBufferDuration)
		for {
			select {
			case <-tick:
				instr.mutex.Lock()

				if err := pool.With(defaultInstrumentationKey, func(conn r.Conn) error {

					for _, v := range instr.buffer {
						conn.Send("LPUSH", "turingery_logs", v)
					}

					return conn.Flush()
				}); err != nil {
					teleprinter.L.Error().Printf("Failed to send interumentation.\n")
				}

				instr.buffer = make([]string, 0)
				instr.mutex.Unlock()
			}
		}
	}()

	return instr
}

func formatLine(v string) string {
	return fmt.Sprintf("[%s] [INSTR] %s", time.Now().Format(time.RFC3339), v)
}

func (i *instrument) counter(value string, amount int) {
	i.mutex.Lock()
	i.buffer = append(i.buffer, formatLine(fmt.Sprintf("%s %d", value, amount)))
	i.mutex.Unlock()
}

func (i *instrument) duration(value string, duration time.Duration) {
	i.mutex.Lock()
	i.buffer = append(i.buffer, formatLine(fmt.Sprintf("%s %d", value, duration.Nanoseconds())))
	i.mutex.Unlock()
}

func (i instrument) ClusterCall(n int) {
	i.counter(fmt.Sprintf("cluster.%d.call.count", n), 1)
}

func (i instrument) ClusterDuration(n int, t time.Duration) {
	i.duration(fmt.Sprintf("cluster.%d.duration", n), t)
}

func (i instrument) AInsertCall() {
	i.counter("aggregate_insert.call.count", 1)
}

func (i instrument) AInsertDuration(t time.Duration) {
	i.duration("aggregate_insert.duration", t)
}

func (i instrument) AModifyCall() {
	i.counter("aggregate_modify.call.count", 1)
}

func (i instrument) AModifyDuration(t time.Duration) {
	i.duration("aggregate_modify.duration", t)
}

func (i instrument) AModifyWithOperationsCall() {
	i.counter("aggregate_modify_with_operations.call.count", 1)
}

func (i instrument) AModifyWithOperationsDuration(t time.Duration) {
	i.duration("aggregate_modify_with_operations.duration", t)
}

func (i instrument) ADeleteCall() {
	i.counter("aggregate_delete.call.count", 1)
}

func (i instrument) ADeleteDuration(t time.Duration) {
	i.duration("aggregate_delete.duration", t)
}

func (i instrument) ARollbackCall() {
	i.counter("aggregate_rollback.call.count", 1)
}

func (i instrument) ARollbackDuration(t time.Duration) {
	i.duration("aggregate_rollback.duration", t)
}

func (i instrument) ASelectCall() {
	i.counter("aggregate_select.call.count", 1)
}

func (i instrument) ASelectDuration(t time.Duration) {
	i.duration("aggregate_select.duration", t)
}

func (i instrument) ASelectRangeCall() {
	i.counter("aggregate_select_range.call.count", 1)
}

func (i instrument) ASelectRangeDuration(t time.Duration) {
	i.duration("aggregate_select_range.duration", t)
}

func (i instrument) AKeysCall() {
	i.counter("aggregate_keys.call.count", 1)
}

func (i instrument) AKeysDuration(t time.Duration) {
	i.duration("aggregate_keys.duration", t)
}

func (i instrument) ASizeCall() {
	i.counter("aggregate_size.call.count", 1)
}

func (i instrument) ASizeDuration(t time.Duration) {
	i.duration("aggregate_size.duration", t)
}

func (i instrument) AMembersCall() {
	i.counter("aggregate_members.call.count", 1)
}

func (i instrument) AMembersDuration(t time.Duration) {
	i.duration("aggregate_members.duration", t)
}

func (i instrument) ARepairCall() {
	i.counter("aggregate_repair.call.count", 1)
}

func (i instrument) ARepairDuration(t time.Duration) {
	i.duration("aggregate_repair.duration", t)
}

func (i instrument) AQueryCall() {
	i.counter("aggregate_query.call.count", 1)
}

func (i instrument) AQueryDuration(t time.Duration) {
	i.duration("aggregate_query.duration", t)
}

func (i instrument) APauseCall() {
	i.counter("aggregate_pause.call.count", 1)
}

func (i instrument) AResumeCall() {
	i.counter("aggregate_resume.call.count", 1)
}

func (i instrument) ATopologyCall() {
	i.counter("aggregate_topology.call.count", 1)
}

func (i instrument) ATopologyDuration(t time.Duration) {
	i.duration("aggregate_topology.duration", t)
}

func (i instrument) InsertCall() {
	i.counter("insert.call.count", 1)
}

func (i instrument) InsertKeys(n int) {
	i.counter("insert.keys.count", n)
}

func (i instrument) InsertSendTo(n int) {
	i.counter("insert.send_to.count", n)
}

func (i instrument) InsertDuration(t time.Duration) {
	i.duration("insert.duration", t)
}

func (i instrument) InsertRetrieved(n int) {
	i.counter("insert.retrieved.count", n)
}

func (i instrument) InsertReturned(n int) {
	i.counter("insert.returned.count", n)
}

func (i instrument) InsertQuorumFailure() {
	i.counter("insert.quorum_failure.count", 1)
}

func (i instrument) InsertRepairRequired() {
	i.counter("insert.repair_required.count", 1)
}

func (i instrument) InsertPartialFailure() {
	i.counter("insert.partial_failure.count", 1)
}

func (i instrument) ModifyCall() {
	i.counter("modify.call.count", 1)
}

func (i instrument) ModifyKeys(n int) {
	i.counter("modify.keys.count", n)
}

func (i instrument) ModifySendTo(n int) {
	i.counter("modify.send_to.count", n)
}

func (i instrument) ModifyDuration(t time.Duration) {
	i.duration("modify.duration", t)
}

func (i instrument) ModifyRetrieved(n int) {
	i.counter("modify.retrieved.count", n)
}

func (i instrument) ModifyReturned(n int) {
	i.counter("modify.returned.count", n)
}

func (i instrument) ModifyQuorumFailure() {
	i.counter("modify.quorum_failure.count", 1)
}

func (i instrument) ModifyRepairRequired() {
	i.counter("modify.repair_required.count", 1)
}

func (i instrument) DeleteCall() {
	i.counter("delete.call.count", 1)
}

func (i instrument) DeleteKeys(n int) {
	i.counter("delete.keys.count", n)
}

func (i instrument) DeleteSendTo(n int) {
	i.counter("delete.send_to.count", n)
}

func (i instrument) DeleteDuration(t time.Duration) {
	i.duration("delete.duration", t)
}

func (i instrument) DeleteRetrieved(n int) {
	i.counter("delete.retrieved.count", n)
}

func (i instrument) DeleteReturned(n int) {
	i.counter("delete.returned.count", n)
}

func (i instrument) DeleteQuorumFailure() {
	i.counter("delete.quorum_failure.count", 1)
}

func (i instrument) DeleteRepairRequired() {
	i.counter("delete.repair_required.count", 1)
}

func (i instrument) DeletePartialFailure() {
	i.counter("delete.partial_failure.count", 1)
}

func (i instrument) RollbackCall() {
	i.counter("rollback.call.count", 1)
}

func (i instrument) RollbackKeys(n int) {
	i.counter("rollback.keys.count", n)
}

func (i instrument) RollbackSendTo(n int) {
	i.counter("rollback.send_to.count", n)
}

func (i instrument) RollbackDuration(t time.Duration) {
	i.duration("rollback.duration", t)
}

func (i instrument) RollbackRetrieved(n int) {
	i.counter("rollback.retrieved.count", n)
}

func (i instrument) RollbackReturned(n int) {
	i.counter("rollback.returned.count", n)
}

func (i instrument) RollbackQuorumFailure() {
	i.counter("rollback.quorum_failure.count", 1)
}

func (i instrument) RollbackRepairRequired() {
	i.counter("rollback.repair_required.count", 1)
}

func (i instrument) RollbackPartialFailure() {
	i.counter("rollback.partial_failure.count", 1)
}

func (i instrument) SelectCall() {
	i.counter("select.call.count", 1)
}

func (i instrument) SelectKeys(n int) {
	i.counter("select.keys.count", n)
}

func (i instrument) SelectSendTo(n int) {
	i.counter("select.send_to.count", n)
}

func (i instrument) SelectSendAllPromotion() {
	i.counter("select.send_all_promotion.count", 1)
}

func (i instrument) SelectPartialError() {
	i.counter("select.partial_error.count", 1)
}

func (i instrument) SelectDuration(t time.Duration) {
	i.duration("select.duration", t)
}

func (i instrument) SelectRetrieved(n int) {
	i.counter("select.retrieved.count", n)
}

func (i instrument) SelectReturned(n int) {
	i.counter("select.returned.count", n)
}

func (i instrument) SelectFirstResponseDuration(t time.Duration) {
	i.duration("select.first_response_duration", t)
}

func (i instrument) SelectBlockingDuration(t time.Duration) {
	i.duration("select.blocking_duration", t)
}

func (i instrument) SelectOverheadDuration(t time.Duration) {
	i.duration("select.overhead_duration", t)
}

func (i instrument) SelectRepairNeeded() {
	i.counter("scan.select_repair_needed.count", 1)
}

func (i instrument) ScanCall() {
	i.counter("scan.call.count", 1)
}

func (i instrument) ScanSendTo(n int) {
	i.counter("scan.send_to.count", n)
}

func (i instrument) ScanPartialError() {
	i.counter("scan.partial_error.count", 1)
}

func (i instrument) ScanDuration(t time.Duration) {
	i.duration("scan.duration", t)
}

func (i instrument) ScanRetrieved(n int) {
	i.counter("scan.retrieved.count", n)
}

func (i instrument) ScanReturned(n int) {
	i.counter("scan.returned.count", n)
}

func (i instrument) ScanRepairNeeded(n int) {
	i.counter("scan.repair_needed.count", n)
}

func (i instrument) RepairCall() {
	i.counter("repair.call.count", 1)
}

func (i instrument) RepairRequest(n int) {
	i.counter("repair.request.count", n)
}

func (i instrument) RepairSendTo(n int) {
	i.counter("repair.send_to.count", n)
}

func (i instrument) RepairDuration(t time.Duration) {
	i.duration("repair.duration", t)
}

func (i instrument) RepairScoreError() {
	i.counter("repair.score_error.count", 1)
}

func (i instrument) RepairError(n int) {
	i.counter("repair.error.count", n)
}

func (i instrument) PerformanceDuration(t time.Duration) {
	i.duration("performance.duration", t)
}

func (i instrument) PerformanceNamespaceDuration(ns string, t time.Duration) {
	i.duration(fmt.Sprintf("performance.%s.duration", ns), t)
}

func (i instrument) PublishCall() {
	i.counter("publish.call.count", 1)
}

func (i instrument) PublishKeys(n int) {
	i.counter("publish.keys.count", n)
}

func (i instrument) PublishSendTo(n int) {
	i.counter("publish.sent_to.count", n)
}

func (i instrument) PublishRetrieved(n int) {
	i.counter("publish.retrieved.count", n)
}

func (i instrument) PublishReturned(n int) {
	i.counter("publish.returned.count", n)
}

func (i instrument) PublishDuration(t time.Duration) {
	i.duration("publish.duration", t)
}

func (i instrument) SemaphoreCall() {
	i.counter("semaphore.call.count \n", 1)
}
func (i instrument) SemaphoreSendTo(n int) {
	i.counter("semaphore.send_to.count %d\n", n)
}
func (i instrument) SemaphoreDuration(t time.Duration) {
	i.duration("semaphore.duration %d\n", t)
}
func (i instrument) SemaphoreRetrieved(n int) {
	i.counter("semaphore.retrieved.count %d\n", n)
}
func (i instrument) SemaphoreReturned(n int) {
	i.counter("semaphore.returned.count %d\n", n)
}

func (i instrument) HeartbeatCall() {
	i.counter("heartbeat.call.count \n", 1)
}
func (i instrument) HeartbeatSendTo(n int) {
	i.counter("heartbeat.send_to.count %d\n", n)
}
func (i instrument) HeartbeatDuration(t time.Duration) {
	i.duration("heartbeat.duration %d\n", t)
}
func (i instrument) HeartbeatRetrieved(n int) {
	i.counter("heartbeat.retrieved.count %d\n", n)
}
func (i instrument) HeartbeatReturned(n int) {
	i.counter("heartbeat.returned.count %d\n", n)
}

func (i instrument) KeyStoreCall() {
	i.counter("keystore.call.count \n", 1)
}
func (i instrument) KeyStoreSendTo(n int) {
	i.counter("keystore.send_to.count %d\n", n)
}
func (i instrument) KeyStoreDuration(t time.Duration) {
	i.duration("keystore.duration %d\n", t)
}
func (i instrument) KeyStoreRetrieved(n int) {
	i.counter("keystore.retrieved.count %d\n", n)
}
func (i instrument) KeyStoreReturned(n int) {
	i.counter("keystore.returned.count %d\n", n)
}
