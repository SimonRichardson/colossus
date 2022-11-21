package noop

import (
	"time"

	"github.com/SimonRichardson/colossus/instrumentation"
)

type instrument struct{}

func New() instrumentation.Instrumentation {
	return instrument{}
}

func (i instrument) ClusterCall(int)                    {}
func (i instrument) ClusterDuration(int, time.Duration) {}

func (i instrument) AInsertCall()                                {}
func (i instrument) AInsertDuration(time.Duration)               {}
func (i instrument) AModifyCall()                                {}
func (i instrument) AModifyDuration(time.Duration)               {}
func (i instrument) AModifyWithOperationsCall()                  {}
func (i instrument) AModifyWithOperationsDuration(time.Duration) {}
func (i instrument) ADeleteCall()                                {}
func (i instrument) ADeleteDuration(time.Duration)               {}
func (i instrument) ARollbackCall()                              {}
func (i instrument) ARollbackDuration(time.Duration)             {}
func (i instrument) ASelectCall()                                {}
func (i instrument) ASelectDuration(time.Duration)               {}
func (i instrument) ASelectRangeCall()                           {}
func (i instrument) ASelectRangeDuration(time.Duration)          {}
func (i instrument) AKeysCall()                                  {}
func (i instrument) AKeysDuration(time.Duration)                 {}
func (i instrument) ASizeCall()                                  {}
func (i instrument) ASizeDuration(time.Duration)                 {}
func (i instrument) AMembersCall()                               {}
func (i instrument) AMembersDuration(time.Duration)              {}
func (i instrument) ARepairCall()                                {}
func (i instrument) ARepairDuration(time.Duration)               {}
func (i instrument) AQueryCall()                                 {}
func (i instrument) AQueryDuration(time.Duration)                {}
func (i instrument) APauseCall()                                 {}
func (i instrument) AResumeCall()                                {}
func (i instrument) ATopologyCall()                              {}
func (i instrument) ATopologyDuration(time.Duration)             {}

func (i instrument) InsertCall()                  {}
func (i instrument) InsertKeys(int)               {}
func (i instrument) InsertSendTo(int)             {}
func (i instrument) InsertDuration(time.Duration) {}
func (i instrument) InsertRetrieved(int)          {}
func (i instrument) InsertReturned(int)           {}
func (i instrument) InsertQuorumFailure()         {}
func (i instrument) InsertRepairRequired()        {}
func (i instrument) InsertPartialFailure()        {}

func (i instrument) ModifyCall()                  {}
func (i instrument) ModifyKeys(int)               {}
func (i instrument) ModifySendTo(int)             {}
func (i instrument) ModifyDuration(time.Duration) {}
func (i instrument) ModifyRetrieved(int)          {}
func (i instrument) ModifyReturned(int)           {}
func (i instrument) ModifyQuorumFailure()         {}
func (i instrument) ModifyRepairRequired()        {}

func (i instrument) DeleteCall()                  {}
func (i instrument) DeleteKeys(int)               {}
func (i instrument) DeleteSendTo(int)             {}
func (i instrument) DeleteDuration(time.Duration) {}
func (i instrument) DeleteRetrieved(int)          {}
func (i instrument) DeleteReturned(int)           {}
func (i instrument) DeleteQuorumFailure()         {}
func (i instrument) DeleteRepairRequired()        {}
func (i instrument) DeletePartialFailure()        {}

func (i instrument) RollbackCall()                  {}
func (i instrument) RollbackKeys(int)               {}
func (i instrument) RollbackSendTo(int)             {}
func (i instrument) RollbackDuration(time.Duration) {}
func (i instrument) RollbackRetrieved(int)          {}
func (i instrument) RollbackReturned(int)           {}
func (i instrument) RollbackQuorumFailure()         {}
func (i instrument) RollbackRepairRequired()        {}
func (i instrument) RollbackPartialFailure()        {}

func (i instrument) SelectCall()                               {}
func (i instrument) SelectKeys(int)                            {}
func (i instrument) SelectSendTo(int)                          {}
func (i instrument) SelectSendAllPromotion()                   {}
func (i instrument) SelectPartialError()                       {}
func (i instrument) SelectDuration(time.Duration)              {}
func (i instrument) SelectRetrieved(int)                       {}
func (i instrument) SelectReturned(int)                        {}
func (i instrument) SelectFirstResponseDuration(time.Duration) {}
func (i instrument) SelectBlockingDuration(time.Duration)      {}
func (i instrument) SelectOverheadDuration(time.Duration)      {}
func (i instrument) SelectRepairNeeded()                       {}

func (i instrument) ScanCall()                  {}
func (i instrument) ScanSendTo(int)             {}
func (i instrument) ScanPartialError()          {}
func (i instrument) ScanDuration(time.Duration) {}
func (i instrument) ScanRetrieved(int)          {}
func (i instrument) ScanReturned(int)           {}
func (i instrument) ScanRepairNeeded(int)       {}

func (i instrument) RepairCall()                  {}
func (i instrument) RepairRequest(int)            {}
func (i instrument) RepairSendTo(int)             {}
func (i instrument) RepairDuration(time.Duration) {}
func (i instrument) RepairScoreError()            {}
func (i instrument) RepairError(int)              {}

func (i instrument) PerformanceDuration(t time.Duration)                     {}
func (i instrument) PerformanceNamespaceDuration(ns string, t time.Duration) {}

func (i instrument) PublishCall()                  {}
func (i instrument) PublishKeys(int)               {}
func (i instrument) PublishSendTo(int)             {}
func (i instrument) PublishRetrieved(int)          {}
func (i instrument) PublishReturned(int)           {}
func (i instrument) PublishDuration(time.Duration) {}

func (i instrument) SemaphoreCall()                  {}
func (i instrument) SemaphoreSendTo(int)             {}
func (i instrument) SemaphoreDuration(time.Duration) {}
func (i instrument) SemaphoreRetrieved(int)          {}
func (i instrument) SemaphoreReturned(int)           {}

func (i instrument) HeartbeatCall()                  {}
func (i instrument) HeartbeatSendTo(int)             {}
func (i instrument) HeartbeatDuration(time.Duration) {}
func (i instrument) HeartbeatRetrieved(int)          {}
func (i instrument) HeartbeatReturned(int)           {}

func (i instrument) KeyStoreCall()                  {}
func (i instrument) KeyStoreSendTo(int)             {}
func (i instrument) KeyStoreDuration(time.Duration) {}
func (i instrument) KeyStoreRetrieved(int)          {}
func (i instrument) KeyStoreReturned(int)           {}
