package multi

import (
	"time"

	"github.com/SimonRichardson/colossus/instrumentation"
)

type instrument struct {
	instruments []instrumentation.Instrumentation
}

func New(instruments ...instrumentation.Instrumentation) instrumentation.Instrumentation {
	return instrument{instruments}
}

func (i instrument) ClusterCall(n int) {
	for _, v := range i.instruments {
		v.ClusterCall(n)
	}
}

func (i instrument) ClusterDuration(n int, t time.Duration) {
	for _, v := range i.instruments {
		v.ClusterDuration(n, t)
	}
}

func (i instrument) AInsertCall() {
	for _, v := range i.instruments {
		v.AInsertCall()
	}
}
func (i instrument) AInsertDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.AInsertDuration(t)
	}
}
func (i instrument) AModifyCall() {
	for _, v := range i.instruments {
		v.AModifyCall()
	}
}
func (i instrument) AModifyDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.AModifyDuration(t)
	}
}
func (i instrument) AModifyWithOperationsCall() {
	for _, v := range i.instruments {
		v.AModifyWithOperationsCall()
	}
}
func (i instrument) AModifyWithOperationsDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.AModifyWithOperationsDuration(t)
	}
}
func (i instrument) ADeleteCall() {
	for _, v := range i.instruments {
		v.ADeleteCall()
	}
}
func (i instrument) ADeleteDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ADeleteDuration(t)
	}
}
func (i instrument) ARollbackCall() {
	for _, v := range i.instruments {
		v.ARollbackCall()
	}
}
func (i instrument) ARollbackDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ARollbackDuration(t)
	}
}
func (i instrument) ASelectCall() {
	for _, v := range i.instruments {
		v.ASelectCall()
	}
}
func (i instrument) ASelectDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ASelectDuration(t)
	}
}
func (i instrument) ASelectRangeCall() {
	for _, v := range i.instruments {
		v.ASelectRangeCall()
	}
}
func (i instrument) ASelectRangeDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ASelectRangeDuration(t)
	}
}
func (i instrument) AKeysCall() {
	for _, v := range i.instruments {
		v.AKeysCall()
	}
}
func (i instrument) AKeysDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.AKeysDuration(t)
	}
}
func (i instrument) ASizeCall() {
	for _, v := range i.instruments {
		v.ASizeCall()
	}
}
func (i instrument) ASizeDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ASizeDuration(t)
	}
}
func (i instrument) AMembersCall() {
	for _, v := range i.instruments {
		v.AMembersCall()
	}
}
func (i instrument) AMembersDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.AMembersDuration(t)
	}
}
func (i instrument) ARepairCall() {
	for _, v := range i.instruments {
		v.ARepairCall()
	}
}
func (i instrument) ARepairDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ARepairDuration(t)
	}
}
func (i instrument) AQueryCall() {
	for _, v := range i.instruments {
		v.AQueryCall()
	}
}
func (i instrument) AQueryDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.AQueryDuration(t)
	}
}
func (i instrument) APauseCall() {
	for _, v := range i.instruments {
		v.APauseCall()
	}
}
func (i instrument) AResumeCall() {
	for _, v := range i.instruments {
		v.AResumeCall()
	}
}
func (i instrument) ATopologyCall() {
	for _, v := range i.instruments {
		v.ATopologyCall()
	}
}
func (i instrument) ATopologyDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ATopologyDuration(t)
	}
}

func (i instrument) InsertCall() {
	for _, v := range i.instruments {
		v.InsertCall()
	}
}

func (i instrument) InsertKeys(n int) {
	for _, v := range i.instruments {
		v.InsertKeys(n)
	}
}

func (i instrument) InsertSendTo(n int) {
	for _, v := range i.instruments {
		v.InsertSendTo(n)
	}
}

func (i instrument) InsertDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.InsertDuration(t)
	}
}

func (i instrument) InsertRetrieved(n int) {
	for _, v := range i.instruments {
		v.InsertRetrieved(n)
	}
}

func (i instrument) InsertReturned(n int) {
	for _, v := range i.instruments {
		v.InsertRetrieved(n)
	}
}

func (i instrument) InsertQuorumFailure() {
	for _, v := range i.instruments {
		v.InsertQuorumFailure()
	}
}

func (i instrument) InsertRepairRequired() {
	for _, v := range i.instruments {
		v.InsertRepairRequired()
	}
}

func (i instrument) InsertPartialFailure() {
	for _, v := range i.instruments {
		v.InsertPartialFailure()
	}
}

func (i instrument) ModifyCall() {
	for _, v := range i.instruments {
		v.ModifyCall()
	}
}

func (i instrument) ModifyKeys(n int) {
	for _, v := range i.instruments {
		v.ModifyKeys(n)
	}
}

func (i instrument) ModifySendTo(n int) {
	for _, v := range i.instruments {
		v.ModifySendTo(n)
	}
}

func (i instrument) ModifyDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ModifyDuration(t)
	}
}

func (i instrument) ModifyRetrieved(n int) {
	for _, v := range i.instruments {
		v.ModifyRetrieved(n)
	}
}

func (i instrument) ModifyReturned(n int) {
	for _, v := range i.instruments {
		v.ModifyRetrieved(n)
	}
}

func (i instrument) ModifyQuorumFailure() {
	for _, v := range i.instruments {
		v.ModifyQuorumFailure()
	}
}

func (i instrument) ModifyRepairRequired() {
	for _, v := range i.instruments {
		v.ModifyRepairRequired()
	}
}

func (i instrument) DeleteCall() {
	for _, v := range i.instruments {
		v.DeleteCall()
	}
}

func (i instrument) DeleteKeys(n int) {
	for _, v := range i.instruments {
		v.DeleteKeys(n)
	}
}

func (i instrument) DeleteSendTo(n int) {
	for _, v := range i.instruments {
		v.DeleteSendTo(n)
	}
}

func (i instrument) DeleteDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.DeleteDuration(t)
	}
}

func (i instrument) DeleteRetrieved(n int) {
	for _, v := range i.instruments {
		v.DeleteRetrieved(n)
	}
}

func (i instrument) DeleteReturned(n int) {
	for _, v := range i.instruments {
		v.DeleteRetrieved(n)
	}
}

func (i instrument) DeleteQuorumFailure() {
	for _, v := range i.instruments {
		v.DeleteQuorumFailure()
	}
}

func (i instrument) DeleteRepairRequired() {
	for _, v := range i.instruments {
		v.DeleteRepairRequired()
	}
}

func (i instrument) DeletePartialFailure() {
	for _, v := range i.instruments {
		v.DeletePartialFailure()
	}
}

func (i instrument) RollbackCall() {
	for _, v := range i.instruments {
		v.RollbackCall()
	}
}

func (i instrument) RollbackKeys(n int) {
	for _, v := range i.instruments {
		v.RollbackKeys(n)
	}
}

func (i instrument) RollbackSendTo(n int) {
	for _, v := range i.instruments {
		v.RollbackSendTo(n)
	}
}

func (i instrument) RollbackDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.RollbackDuration(t)
	}
}

func (i instrument) RollbackRetrieved(n int) {
	for _, v := range i.instruments {
		v.RollbackRetrieved(n)
	}
}

func (i instrument) RollbackReturned(n int) {
	for _, v := range i.instruments {
		v.RollbackRetrieved(n)
	}
}

func (i instrument) RollbackQuorumFailure() {
	for _, v := range i.instruments {
		v.RollbackQuorumFailure()
	}
}

func (i instrument) RollbackRepairRequired() {
	for _, v := range i.instruments {
		v.RollbackRepairRequired()
	}
}

func (i instrument) RollbackPartialFailure() {
	for _, v := range i.instruments {
		v.RollbackPartialFailure()
	}
}

func (i instrument) SelectCall() {
	for _, v := range i.instruments {
		v.SelectCall()
	}
}

func (i instrument) SelectKeys(n int) {
	for _, v := range i.instruments {
		v.SelectKeys(n)
	}
}

func (i instrument) SelectSendTo(n int) {
	for _, v := range i.instruments {
		v.SelectSendTo(n)
	}
}

func (i instrument) SelectSendAllPromotion() {
	for _, v := range i.instruments {
		v.SelectSendAllPromotion()
	}
}

func (i instrument) SelectPartialError() {
	for _, v := range i.instruments {
		v.SelectPartialError()
	}
}

func (i instrument) SelectDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.SelectDuration(t)
	}
}

func (i instrument) SelectRetrieved(n int) {
	for _, v := range i.instruments {
		v.SelectRetrieved(n)
	}
}

func (i instrument) SelectReturned(n int) {
	for _, v := range i.instruments {
		v.SelectReturned(n)
	}
}

func (i instrument) SelectFirstResponseDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.SelectFirstResponseDuration(t)
	}
}

func (i instrument) SelectBlockingDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.SelectBlockingDuration(t)
	}
}

func (i instrument) SelectOverheadDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.SelectOverheadDuration(t)
	}
}

func (i instrument) SelectRepairNeeded() {
	for _, v := range i.instruments {
		v.SelectRepairNeeded()
	}
}

func (i instrument) ScanCall() {
	for _, v := range i.instruments {
		v.ScanCall()
	}
}

func (i instrument) ScanSendTo(n int) {
	for _, v := range i.instruments {
		v.ScanSendTo(n)
	}
}

func (i instrument) ScanPartialError() {
	for _, v := range i.instruments {
		v.ScanPartialError()
	}
}

func (i instrument) ScanDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.ScanDuration(t)
	}
}

func (i instrument) ScanRetrieved(n int) {
	for _, v := range i.instruments {
		v.ScanRetrieved(n)
	}
}

func (i instrument) ScanReturned(n int) {
	for _, v := range i.instruments {
		v.ScanReturned(n)
	}
}

func (i instrument) ScanRepairNeeded(n int) {
	for _, v := range i.instruments {
		v.ScanRepairNeeded(n)
	}
}

func (i instrument) RepairCall() {
	for _, v := range i.instruments {
		v.RepairCall()
	}
}

func (i instrument) RepairRequest(n int) {
	for _, v := range i.instruments {
		v.RepairRequest(n)
	}
}

func (i instrument) RepairSendTo(n int) {
	for _, v := range i.instruments {
		v.RepairSendTo(n)
	}
}

func (i instrument) RepairDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.RepairDuration(t)
	}
}

func (i instrument) RepairScoreError() {
	for _, v := range i.instruments {
		v.RepairScoreError()
	}
}

func (i instrument) RepairError(n int) {
	for _, v := range i.instruments {
		v.RepairError(n)
	}
}

func (i instrument) PerformanceDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.PerformanceDuration(t)
	}
}

func (i instrument) PerformanceNamespaceDuration(ns string, t time.Duration) {
	for _, v := range i.instruments {
		v.PerformanceNamespaceDuration(ns, t)
	}
}

func (i instrument) PublishCall() {
	for _, v := range i.instruments {
		v.PublishCall()
	}
}

func (i instrument) PublishKeys(n int) {
	for _, v := range i.instruments {
		v.PublishKeys(n)
	}
}

func (i instrument) PublishSendTo(n int) {
	for _, v := range i.instruments {
		v.PublishSendTo(n)
	}
}

func (i instrument) PublishRetrieved(n int) {
	for _, v := range i.instruments {
		v.PublishRetrieved(n)
	}
}

func (i instrument) PublishReturned(n int) {
	for _, v := range i.instruments {
		v.PublishReturned(n)
	}
}

func (i instrument) PublishDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.PublishDuration(t)
	}
}

func (i instrument) SemaphoreCall() {
	for _, v := range i.instruments {
		v.SemaphoreCall()
	}
}
func (i instrument) SemaphoreSendTo(n int) {
	for _, v := range i.instruments {
		v.SemaphoreSendTo(n)
	}
}
func (i instrument) SemaphoreDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.SemaphoreDuration(t)
	}
}
func (i instrument) SemaphoreRetrieved(n int) {
	for _, v := range i.instruments {
		v.SemaphoreRetrieved(n)
	}
}
func (i instrument) SemaphoreReturned(n int) {
	for _, v := range i.instruments {
		v.SemaphoreReturned(n)
	}
}

func (i instrument) HeartbeatCall() {
	for _, v := range i.instruments {
		v.HeartbeatCall()
	}
}
func (i instrument) HeartbeatSendTo(n int) {
	for _, v := range i.instruments {
		v.HeartbeatSendTo(n)
	}
}
func (i instrument) HeartbeatDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.HeartbeatDuration(t)
	}
}
func (i instrument) HeartbeatRetrieved(n int) {
	for _, v := range i.instruments {
		v.HeartbeatRetrieved(n)
	}
}
func (i instrument) HeartbeatReturned(n int) {
	for _, v := range i.instruments {
		v.HeartbeatReturned(n)
	}
}

func (i instrument) KeyStoreCall() {
	for _, v := range i.instruments {
		v.KeyStoreCall()
	}
}
func (i instrument) KeyStoreSendTo(n int) {
	for _, v := range i.instruments {
		v.KeyStoreSendTo(n)
	}
}
func (i instrument) KeyStoreDuration(t time.Duration) {
	for _, v := range i.instruments {
		v.KeyStoreDuration(t)
	}
}
func (i instrument) KeyStoreRetrieved(n int) {
	for _, v := range i.instruments {
		v.KeyStoreRetrieved(n)
	}
}
func (i instrument) KeyStoreReturned(n int) {
	for _, v := range i.instruments {
		v.KeyStoreReturned(n)
	}
}
