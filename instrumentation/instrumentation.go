package instrumentation

import (
	"time"

	"github.com/SimonRichardson/colossus/blist/services/consul"
)

type Instrumentation interface {
	ClusterInstrumentation
	AggregateInstrumentation
	InsertInstrumentation
	ModifyInstrumentation
	DeleteInstrumentation
	RollbackInstrumentation
	SelectInstrumentation
	ScanInstrumentation
	RepairInstrumentation
	PerformanceDuration
	PublishInstrumentation
	consul.Instrumentation
}

type ClusterInstrumentation interface {
	ClusterCall(int)
	ClusterDuration(int, time.Duration)
}

type AggregateInstrumentation interface {
	AInsertCall()
	AInsertDuration(time.Duration)
	AModifyCall()
	AModifyDuration(time.Duration)
	AModifyWithOperationsCall()
	AModifyWithOperationsDuration(time.Duration)
	ADeleteCall()
	ADeleteDuration(time.Duration)
	ARollbackCall()
	ARollbackDuration(time.Duration)
	ASelectCall()
	ASelectDuration(time.Duration)
	ASelectRangeCall()
	ASelectRangeDuration(time.Duration)
	AKeysCall()
	AKeysDuration(time.Duration)
	ASizeCall()
	ASizeDuration(time.Duration)
	AMembersCall()
	AMembersDuration(time.Duration)
	ARepairCall()
	ARepairDuration(time.Duration)
	AQueryCall()
	AQueryDuration(time.Duration)
	APauseCall()
	AResumeCall()
	ATopologyCall()
	ATopologyDuration(time.Duration)
}

type InsertInstrumentation interface {
	InsertCall()
	InsertKeys(int)
	InsertSendTo(int)
	InsertDuration(time.Duration)
	InsertRetrieved(int)
	InsertReturned(int)
	InsertQuorumFailure()
	InsertRepairRequired()
	InsertPartialFailure()
}

type ModifyInstrumentation interface {
	ModifyCall()
	ModifyKeys(int)
	ModifySendTo(int)
	ModifyDuration(time.Duration)
	ModifyRetrieved(int)
	ModifyReturned(int)
	ModifyQuorumFailure()
	ModifyRepairRequired()
}

type DeleteInstrumentation interface {
	DeleteCall()
	DeleteKeys(int)
	DeleteSendTo(int)
	DeleteDuration(time.Duration)
	DeleteRetrieved(int)
	DeleteReturned(int)
	DeleteQuorumFailure()
	DeleteRepairRequired()
	DeletePartialFailure()
}

type RollbackInstrumentation interface {
	RollbackCall()
	RollbackKeys(int)
	RollbackSendTo(int)
	RollbackDuration(time.Duration)
	RollbackRetrieved(int)
	RollbackReturned(int)
	RollbackQuorumFailure()
	RollbackRepairRequired()
	RollbackPartialFailure()
}

type SelectInstrumentation interface {
	SelectCall()
	SelectKeys(int)
	SelectSendTo(int)
	SelectSendAllPromotion()
	SelectPartialError()
	SelectDuration(time.Duration)
	SelectRetrieved(int)
	SelectReturned(int)
	SelectFirstResponseDuration(time.Duration)
	SelectBlockingDuration(time.Duration)
	SelectOverheadDuration(time.Duration)
	SelectRepairNeeded()
}

type ScanInstrumentation interface {
	ScanCall()
	ScanSendTo(int)
	ScanPartialError()
	ScanDuration(time.Duration)
	ScanRetrieved(int)
	ScanReturned(int)
	ScanRepairNeeded(int)
}

type RepairInstrumentation interface {
	RepairCall()
	RepairRequest(int)
	RepairSendTo(int)
	RepairDuration(time.Duration)
	RepairScoreError()
	RepairError(int)
}

type PerformanceDuration interface {
	PerformanceDuration(time.Duration)
	PerformanceNamespaceDuration(string, time.Duration)
}

type PublishInstrumentation interface {
	PublishCall()
	PublishKeys(int)
	PublishSendTo(int)
	PublishRetrieved(int)
	PublishReturned(int)
	PublishDuration(time.Duration)
}
