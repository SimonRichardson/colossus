package prometheus

import (
	"fmt"
	"sync"
	"time"

	"github.com/SimonRichardson/colossus/instrumentation"
	"github.com/prometheus/client_golang/prometheus"
)

type instrument struct {
	mutex         *sync.Mutex
	prefix        string
	maxSummaryAge time.Duration

	clusterCall     map[int]prometheus.Counter
	clusterDuration map[int]prometheus.Summary

	aInsertCall                   prometheus.Counter
	aInsertDuration               prometheus.Summary
	aModifyCall                   prometheus.Counter
	aModifyDuration               prometheus.Summary
	aModifyWithOperationsCall     prometheus.Counter
	aModifyWithOperationsDuration prometheus.Summary
	aDeleteCall                   prometheus.Counter
	aDeleteDuration               prometheus.Summary
	aRollbackCall                 prometheus.Counter
	aRollbackDuration             prometheus.Summary
	aSelectCall                   prometheus.Counter
	aSelectDuration               prometheus.Summary
	aSelectRangeCall              prometheus.Counter
	aSelectRangeDuration          prometheus.Summary
	aKeysCall                     prometheus.Counter
	aKeysDuration                 prometheus.Summary
	aSizeCall                     prometheus.Counter
	aSizeDuration                 prometheus.Summary
	aMembersCall                  prometheus.Counter
	aMembersDuration              prometheus.Summary
	aRepairCall                   prometheus.Counter
	aRepairDuration               prometheus.Summary
	aQueryCall                    prometheus.Counter
	aQueryDuration                prometheus.Summary
	aPauseCall                    prometheus.Counter
	aResumeCall                   prometheus.Counter
	aTopologyCall                 prometheus.Counter
	aTopologyDuration             prometheus.Summary

	insertCall           prometheus.Counter
	insertKeys           prometheus.Counter
	insertSendTo         prometheus.Counter
	insertRetrieved      prometheus.Counter
	insertReturned       prometheus.Counter
	insertQuorumFailure  prometheus.Counter
	insertRepairRequired prometheus.Counter
	insertPartialFailure prometheus.Counter
	insertDuration       prometheus.Summary

	modifyCall           prometheus.Counter
	modifyKeys           prometheus.Counter
	modifySendTo         prometheus.Counter
	modifyRetrieved      prometheus.Counter
	modifyReturned       prometheus.Counter
	modifyQuorumFailure  prometheus.Counter
	modifyRepairRequired prometheus.Counter
	modifyDuration       prometheus.Summary

	deleteCall           prometheus.Counter
	deleteKeys           prometheus.Counter
	deleteSendTo         prometheus.Counter
	deleteRetrieved      prometheus.Counter
	deleteReturned       prometheus.Counter
	deleteQuorumFailure  prometheus.Counter
	deleteRepairRequired prometheus.Counter
	deletePartialFailure prometheus.Counter
	deleteDuration       prometheus.Summary

	rollbackCall           prometheus.Counter
	rollbackKeys           prometheus.Counter
	rollbackSendTo         prometheus.Counter
	rollbackRetrieved      prometheus.Counter
	rollbackReturned       prometheus.Counter
	rollbackQuorumFailure  prometheus.Counter
	rollbackRepairRequired prometheus.Counter
	rollbackPartialFailure prometheus.Counter
	rollbackDuration       prometheus.Summary

	selectCall                  prometheus.Counter
	selectKeys                  prometheus.Counter
	selectSendTo                prometheus.Counter
	selectRetrieved             prometheus.Counter
	selectReturned              prometheus.Counter
	selectQuorumFailure         prometheus.Counter
	selectRepairRequired        prometheus.Counter
	selectPartialFailure        prometheus.Counter
	selectSendAllPromotion      prometheus.Counter
	selectPartialError          prometheus.Counter
	selectDuration              prometheus.Summary
	selectFirstResponseDuration prometheus.Summary
	selectBlockingDuration      prometheus.Summary
	selectOverheadDuration      prometheus.Summary

	scanCall           prometheus.Counter
	scanSendTo         prometheus.Counter
	scanRetrieved      prometheus.Counter
	scanReturned       prometheus.Counter
	scanRepairRequired prometheus.Counter
	scanPartialFailure prometheus.Counter
	scanDuration       prometheus.Summary

	repairCall       prometheus.Counter
	repairRequest    prometheus.Counter
	repairSendTo     prometheus.Counter
	repairScoreError prometheus.Counter
	repairError      prometheus.Counter
	repairDuration   prometheus.Summary

	performanceDuration          prometheus.Summary
	performanceNamespaceDuration map[string]prometheus.Summary

	publishCall      prometheus.Counter
	publishKeys      prometheus.Counter
	publishSendTo    prometheus.Counter
	publishRetrieved prometheus.Counter
	publishReturned  prometheus.Counter
	publishDuration  prometheus.Summary

	semaphoreCall      prometheus.Counter
	semaphoreSendTo    prometheus.Counter
	semaphoreDuration  prometheus.Summary
	semaphoreRetrieved prometheus.Counter
	semaphoreReturned  prometheus.Counter

	heartbeatCall      prometheus.Counter
	heartbeatSendTo    prometheus.Counter
	heartbeatDuration  prometheus.Summary
	heartbeatRetrieved prometheus.Counter
	heartbeatReturned  prometheus.Counter

	keyStoreCall      prometheus.Counter
	keyStoreSendTo    prometheus.Counter
	keyStoreDuration  prometheus.Summary
	keyStoreRetrieved prometheus.Counter
	keyStoreReturned  prometheus.Counter
}

func New(prefix string, maxSummaryAge time.Duration) instrumentation.Instrumentation {
	i := instrument{
		mutex:         &sync.Mutex{},
		prefix:        prefix,
		maxSummaryAge: maxSummaryAge,

		clusterCall:     map[int]prometheus.Counter{},
		clusterDuration: map[int]prometheus.Summary{},

		aInsertCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_insert_call_count",
			Help:      "How many aggregate insertion calls have been made.",
		}),
		aInsertDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_insert_call_duration",
			Help:      "How long the aggregate insertion calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aModifyCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_modify_call_count",
			Help:      "How many aggregate modification calls have been made.",
		}),
		aModifyDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_modify_call_duration",
			Help:      "How long the aggregate modification calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aModifyWithOperationsCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_modify_with_operations_call_count",
			Help:      "How many aggregate modification with operations calls have been made.",
		}),
		aModifyWithOperationsDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_modify_with_operations_call_duration",
			Help:      "How long the aggregate modification with operations calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aDeleteCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_delete_call_count",
			Help:      "How many aggregate deletion calls have been made.",
		}),
		aDeleteDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_delete_call_duration",
			Help:      "How long the aggregate deletion calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aRollbackCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_rollback_call_count",
			Help:      "How many aggregate rollback calls have been made.",
		}),
		aRollbackDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_rollback_call_duration",
			Help:      "How long the aggregate rollback calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aSelectCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_select_call_count",
			Help:      "How many aggregate selection calls have been made.",
		}),
		aSelectDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_select_call_duration",
			Help:      "How long the aggregate selection calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aSelectRangeCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_select_range_call_count",
			Help:      "How many aggregate range selection calls have been made.",
		}),
		aSelectRangeDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_select_range_call_duration",
			Help:      "How long the aggregate range selection calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aKeysCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_keys_call_count",
			Help:      "How many aggregate keys calls have been made.",
		}),
		aKeysDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_keys_call_duration",
			Help:      "How long the aggregate keys calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aSizeCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_size_call_count",
			Help:      "How many aggregate size calls have been made.",
		}),
		aSizeDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_size_call_duration",
			Help:      "How long the aggregate size calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aMembersCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_members_call_count",
			Help:      "How many aggregate members calls have been made.",
		}),
		aMembersDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_members_call_duration",
			Help:      "How long the aggregate members calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aRepairCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_repair_call_count",
			Help:      "How many aggregate repair calls have been made.",
		}),
		aRepairDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_repair_call_duration",
			Help:      "How long the aggregate repair calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aQueryCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_query_call_count",
			Help:      "How many aggregate query calls have been made.",
		}),
		aQueryDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_query_call_duration",
			Help:      "How long the aggregate query calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		aResumeCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_resume_call_count",
			Help:      "How many aggregate resume calls have been made.",
		}),
		aPauseCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_pause_call_count",
			Help:      "How many aggregate pause calls have been made.",
		}),
		aTopologyCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "aggregate_topology_call_count",
			Help:      "How many aggregate topology calls have been made.",
		}),
		aTopologyDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "aggregate_topology_call_duration",
			Help:      "How long the aggregate topology calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),

		insertCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_call_count",
			Help:      "How many insertion calls have been made.",
		}),
		insertKeys: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_keys_count",
			Help:      "How many insertion keys have been made.",
		}),
		insertSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_send_to_count",
			Help:      "How many insertion send_to have been made.",
		}),
		insertRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_retrieved_count",
			Help:      "How many insertion retrieved have been made.",
		}),
		insertReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_returned_count",
			Help:      "How many insertion returned have been made.",
		}),
		insertQuorumFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_quorum_failure_count",
			Help:      "How many insertion quorum_failure have been made.",
		}),
		insertRepairRequired: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_repair_required_count",
			Help:      "How many insertion repair_required have been made.",
		}),
		insertPartialFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "insert_partial_failure_count",
			Help:      "How many insertion partial_failure have been made.",
		}),
		insertDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "insert_call_duration",
			Help:      "How long the insertion calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),

		modifyCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_call_count",
			Help:      "How many modification calls have been made.",
		}),
		modifyKeys: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_keys_count",
			Help:      "How many modification keys have been made.",
		}),
		modifySendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_send_to_count",
			Help:      "How many modification send_to have been made.",
		}),
		modifyRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_retrieved_count",
			Help:      "How many modification retrieved have been made.",
		}),
		modifyReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_returned_count",
			Help:      "How many modification returned have been made.",
		}),
		modifyQuorumFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_quorum_failure_count",
			Help:      "How many modification quorum_failure have been made.",
		}),
		modifyRepairRequired: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "modify_repair_required_count",
			Help:      "How many modification repair_required have been made.",
		}),
		modifyDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "modify_call_duration",
			Help:      "How long the modification calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),

		deleteCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_call_count",
			Help:      "How many deletion calls have been made.",
		}),
		deleteKeys: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_keys_count",
			Help:      "How many deletion keys have been made.",
		}),
		deleteSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_send_to_count",
			Help:      "How many deletion send_to have been made.",
		}),
		deleteRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_retrieved_count",
			Help:      "How many deletion retrieved have been made.",
		}),
		deleteReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_returned_count",
			Help:      "How many deletion returned have been made.",
		}),
		deleteQuorumFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_quorum_failure_count",
			Help:      "How many deletion quorum_failure have been made.",
		}),
		deleteRepairRequired: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_repair_required_count",
			Help:      "How many deletion repair_required have been made.",
		}),
		deletePartialFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "delete_partial_failure_count",
			Help:      "How many deletion partial_failure have been made.",
		}),
		deleteDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "delete_call_duration",
			Help:      "How long the deletion calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		rollbackCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_call_count",
			Help:      "How many rollback calls have been made.",
		}),
		rollbackKeys: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_keys_count",
			Help:      "How many rollback keys have been made.",
		}),
		rollbackSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_send_to_count",
			Help:      "How many rollback send_to have been made.",
		}),
		rollbackRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_retrieved_count",
			Help:      "How many rollback retrieved have been made.",
		}),
		rollbackReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_returned_count",
			Help:      "How many rollback returned have been made.",
		}),
		rollbackQuorumFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_quorum_failure_count",
			Help:      "How many rollback quorum_failure have been made.",
		}),
		rollbackRepairRequired: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_repair_required_count",
			Help:      "How many rollback repair_required have been made.",
		}),
		rollbackPartialFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "rollback_partial_failure_count",
			Help:      "How many rollback partial_failure have been made.",
		}),
		rollbackDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "rollback_call_duration",
			Help:      "How long the rollback calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		selectCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_call_count",
			Help:      "How many selection calls have been made.",
		}),
		selectKeys: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_keys_count",
			Help:      "How many selection keys have been made.",
		}),
		selectSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_send_to_count",
			Help:      "How many selection send_to have been made.",
		}),
		selectRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_retrieved_count",
			Help:      "How many selection retrieved have been made.",
		}),
		selectReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_returned_count",
			Help:      "How many selection returned have been made.",
		}),
		selectQuorumFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_quorum_failure_count",
			Help:      "How many selection quorum_failure have been made.",
		}),
		selectRepairRequired: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_repair_required_count",
			Help:      "How many selection repair_required have been made.",
		}),
		selectPartialFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_partial_failure_count",
			Help:      "How many selection partial_failure have been made.",
		}),
		selectSendAllPromotion: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_send_all_promotion_count",
			Help:      "How many selection send_all_promotion have been made.",
		}),
		selectPartialError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "select_send_partial_error_count",
			Help:      "How many selection send_partial_error have been made.",
		}),
		selectDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "select_call_duration",
			Help:      "How long the selection calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		selectFirstResponseDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "select_first_response_duration",
			Help:      "How long the selection first_response calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		selectBlockingDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "select_blocking_duration",
			Help:      "How long the selection blocking calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		selectOverheadDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "select_overhead_duration",
			Help:      "How long the selection overhead calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),

		scanCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "scan_call_count",
			Help:      "How many scan calls have been made.",
		}),
		scanSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "scan_send_to_count",
			Help:      "How many scan send_to have been made.",
		}),
		scanRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "scan_retrieved_count",
			Help:      "How many scan retrieved have been made.",
		}),
		scanReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "scan_returned_count",
			Help:      "How many scan returned have been made.",
		}),
		scanRepairRequired: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "scan_repair_required_count",
			Help:      "How many scan repair_required have been made.",
		}),
		scanPartialFailure: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "scan_partial_failure_count",
			Help:      "How many scan partial_failure have been made.",
		}),
		scanDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "scan_call_duration",
			Help:      "How long the scan calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),

		repairCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "repair_call_count",
			Help:      "How many repair calls have been made.",
		}),
		repairSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "repair_send_to_count",
			Help:      "How many repair send_to have been made.",
		}),
		repairRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "repair_request_count",
			Help:      "How many repair request have been made.",
		}),
		repairScoreError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "repair_score_error_count",
			Help:      "How many repair score_error have been made.",
		}),
		repairError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "repair_error_count",
			Help:      "How many repair error have been made.",
		}),
		repairDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "repair_call_duration",
			Help:      "How long the repair calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),

		performanceDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "performance_call_duration",
			Help:      "How long the performance calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		performanceNamespaceDuration: map[string]prometheus.Summary{},

		publishCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "publish_call_count",
			Help:      "How many publish calls have been made.",
		}),
		publishKeys: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "publish_keys_count",
			Help:      "How many publish keys have been made.",
		}),
		publishSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "publish_send_to_count",
			Help:      "How many publish send_to have been made.",
		}),
		publishRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "publish_retrieved_count",
			Help:      "How many publish retrieved have been made.",
		}),
		publishReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "publish_returned_count",
			Help:      "How many publish returned have been made.",
		}),
		publishDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "publish_call_duration",
			Help:      "How long the publish calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		semaphoreCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "semaphore_call",
			Help:      "How many semaphore_call have been made.",
		}),
		semaphoreSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "semaphore_send_to",
			Help:      "How many semaphore_send_to have been made.",
		}),
		semaphoreDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "semaphore_duration",
			Help:      "How long the semaphore_duration calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		semaphoreRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "semaphore_retrieved",
			Help:      "How many semaphore_retrieved have been made.",
		}),
		semaphoreReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "semaphore_returned",
			Help:      "How many semaphore_returned have been made.",
		}),
		heartbeatCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "heartbeat_call",
			Help:      "How many heartbeat_call have been made.",
		}),
		heartbeatSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "heartbeat_send_to",
			Help:      "How many heartbeat_send_to have been made.",
		}),
		heartbeatDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "heartbeat_duration",
			Help:      "How long the heartbeat_duration calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		heartbeatRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "heartbeat_retrieved",
			Help:      "How many heartbeat_retrieved have been made.",
		}),
		heartbeatReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "heartbeat_returned",
			Help:      "How many heartbeat_returned have been made.",
		}),
		keyStoreCall: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "keyStore_call",
			Help:      "How many keyStore_call have been made.",
		}),
		keyStoreSendTo: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "keyStore_send_to",
			Help:      "How many keyStore_send_to have been made.",
		}),
		keyStoreDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: prefix,
			Name:      "keyStore_duration",
			Help:      "How long the keyStore_duration calls took in nanoseconds.",
			MaxAge:    maxSummaryAge,
		}),
		keyStoreRetrieved: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "keyStore_retrieved",
			Help:      "How many keyStore_retrieved have been made.",
		}),
		keyStoreReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: prefix,
			Name:      "keyStore_returned",
			Help:      "How many keyStore_returned have been made.",
		}),
	}

	prometheus.MustRegister(i.aInsertCall, i.aInsertDuration)
	prometheus.MustRegister(i.aModifyCall, i.aModifyDuration)
	prometheus.MustRegister(i.aModifyWithOperationsCall, i.aModifyWithOperationsDuration)
	prometheus.MustRegister(i.aDeleteCall, i.aDeleteDuration)
	prometheus.MustRegister(i.aRollbackCall, i.aRollbackDuration)
	prometheus.MustRegister(i.aSelectCall, i.aSelectDuration)
	prometheus.MustRegister(i.aSelectRangeCall, i.aSelectRangeDuration)
	prometheus.MustRegister(i.aKeysCall, i.aKeysDuration)
	prometheus.MustRegister(i.aSizeCall, i.aSizeDuration)
	prometheus.MustRegister(i.aMembersCall, i.aMembersDuration)
	prometheus.MustRegister(i.aRepairCall, i.aRepairDuration)
	prometheus.MustRegister(i.aQueryCall, i.aQueryDuration)
	prometheus.MustRegister(i.aPauseCall, i.aResumeCall)
	prometheus.MustRegister(i.aTopologyCall, i.aTopologyDuration)

	prometheus.MustRegister(i.insertCall, i.insertDuration, i.insertKeys,
		i.insertPartialFailure, i.insertQuorumFailure, i.insertRepairRequired,
		i.insertRetrieved, i.insertReturned, i.insertSendTo,
	)

	prometheus.MustRegister(i.modifyCall, i.modifyDuration, i.modifyKeys,
		i.modifyQuorumFailure, i.modifyRepairRequired,
		i.modifyRetrieved, i.modifyReturned, i.modifySendTo,
	)

	prometheus.MustRegister(i.deleteCall, i.deleteDuration, i.deleteKeys,
		i.deletePartialFailure, i.deleteQuorumFailure, i.deleteRepairRequired,
		i.deleteRetrieved, i.deleteReturned, i.deleteSendTo,
	)

	prometheus.MustRegister(i.rollbackCall, i.rollbackDuration, i.rollbackKeys,
		i.rollbackPartialFailure, i.rollbackQuorumFailure, i.rollbackRepairRequired,
		i.rollbackRetrieved, i.rollbackReturned, i.rollbackSendTo,
	)

	prometheus.MustRegister(i.selectCall, i.selectDuration, i.selectKeys,
		i.selectPartialFailure, i.selectQuorumFailure, i.selectRepairRequired,
		i.selectRetrieved, i.selectReturned, i.selectSendTo,
		i.selectSendAllPromotion, i.selectPartialError,
		i.selectFirstResponseDuration, i.selectBlockingDuration,
		i.selectOverheadDuration,
	)

	prometheus.MustRegister(i.scanCall, i.scanDuration,
		i.scanPartialFailure, i.scanRepairRequired,
		i.scanRetrieved, i.scanReturned, i.scanSendTo,
	)

	prometheus.MustRegister(i.repairCall, i.repairDuration, i.repairError,
		i.repairRequest, i.repairScoreError, i.repairSendTo,
	)

	prometheus.MustRegister(i.performanceDuration)

	prometheus.MustRegister(i.publishCall, i.publishDuration, i.publishKeys,
		i.publishRetrieved, i.publishReturned, i.publishSendTo,
	)

	prometheus.MustRegister(i.semaphoreCall)
	prometheus.MustRegister(i.semaphoreSendTo)
	prometheus.MustRegister(i.semaphoreDuration)
	prometheus.MustRegister(i.semaphoreRetrieved)
	prometheus.MustRegister(i.semaphoreReturned)

	prometheus.MustRegister(i.heartbeatCall)
	prometheus.MustRegister(i.heartbeatSendTo)
	prometheus.MustRegister(i.heartbeatDuration)
	prometheus.MustRegister(i.heartbeatRetrieved)
	prometheus.MustRegister(i.heartbeatReturned)

	prometheus.MustRegister(i.keyStoreCall)
	prometheus.MustRegister(i.keyStoreSendTo)
	prometheus.MustRegister(i.keyStoreDuration)
	prometheus.MustRegister(i.keyStoreRetrieved)
	prometheus.MustRegister(i.keyStoreReturned)

	return i
}

func (i instrument) ClusterCall(n int) {
	var (
		counter prometheus.Counter
		ok      bool
	)
	if counter, ok = i.clusterCall[n]; !ok {
		i.mutex.Lock()

		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: i.prefix,
			Name:      fmt.Sprintf("cluster.%d.call.count", n),
			Help:      fmt.Sprintf("How many cluster calls for index %d, have been made.", n),
		})

		if err := prometheus.Register(counter); err == nil {
			i.clusterCall[n] = counter
		}

		i.mutex.Unlock()
	}
	counter.Inc()
}

func (i instrument) ClusterDuration(n int, t time.Duration) {
	var (
		summary prometheus.Summary
		ok      bool
	)
	if summary, ok = i.clusterDuration[n]; !ok {
		i.mutex.Lock()

		summary = prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: i.prefix,
			Name:      fmt.Sprintf("cluster.%d.duration", n),
			Help:      fmt.Sprintf("How long the cluster calls for index %d, took in nanoseconds.", n),
			MaxAge:    i.maxSummaryAge,
		})

		if err := prometheus.Register(summary); err == nil {
			i.clusterDuration[n] = summary
		}

		i.mutex.Unlock()
	}
	summary.Observe(float64(t.Nanoseconds()))
}

func (i instrument) AInsertCall() {
	i.aInsertCall.Inc()
}

func (i instrument) AInsertDuration(t time.Duration) {
	i.aInsertDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) AModifyCall() {
	i.aModifyCall.Inc()
}

func (i instrument) AModifyDuration(t time.Duration) {
	i.aModifyDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) AModifyWithOperationsCall() {
	i.aModifyWithOperationsCall.Inc()
}

func (i instrument) AModifyWithOperationsDuration(t time.Duration) {
	i.aModifyWithOperationsDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ADeleteCall() {
	i.aDeleteCall.Inc()
}

func (i instrument) ADeleteDuration(t time.Duration) {
	i.aDeleteDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ARollbackCall() {
	i.aRollbackCall.Inc()
}

func (i instrument) ARollbackDuration(t time.Duration) {
	i.aRollbackDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ASelectCall() {
	i.aSelectCall.Inc()
}

func (i instrument) ASelectDuration(t time.Duration) {
	i.aSelectDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ASelectRangeCall() {
	i.aSelectRangeCall.Inc()
}

func (i instrument) ASelectRangeDuration(t time.Duration) {
	i.aSelectRangeDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) AKeysCall() {
	i.aKeysCall.Inc()
}

func (i instrument) AKeysDuration(t time.Duration) {
	i.aKeysDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ASizeCall() {
	i.aSizeCall.Inc()
}

func (i instrument) ASizeDuration(t time.Duration) {
	i.aSizeDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) AMembersCall() {
	i.aMembersCall.Inc()
}

func (i instrument) AMembersDuration(t time.Duration) {
	i.aMembersDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ARepairCall() {
	i.aRepairCall.Inc()
}

func (i instrument) ARepairDuration(t time.Duration) {
	i.aRepairDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) AQueryCall() {
	i.aQueryCall.Inc()
}

func (i instrument) AQueryDuration(t time.Duration) {
	i.aQueryDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) APauseCall() {
	i.aPauseCall.Inc()
}

func (i instrument) AResumeCall() {
	i.aResumeCall.Inc()
}

func (i instrument) ATopologyCall() {
	i.aTopologyCall.Inc()
}

func (i instrument) ATopologyDuration(t time.Duration) {
	i.aTopologyDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) InsertCall() {
	i.insertCall.Inc()
}

func (i instrument) InsertKeys(n int) {
	i.insertKeys.Add(float64(n))
}

func (i instrument) InsertSendTo(n int) {
	i.insertSendTo.Add(float64(n))
}

func (i instrument) InsertDuration(t time.Duration) {
	i.insertDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) InsertRetrieved(n int) {
	i.insertRetrieved.Add(float64(n))
}

func (i instrument) InsertReturned(n int) {
	i.insertReturned.Add(float64(n))
}

func (i instrument) InsertQuorumFailure() {
	i.insertQuorumFailure.Inc()
}

func (i instrument) InsertRepairRequired() {
	i.insertRepairRequired.Inc()
}

func (i instrument) InsertPartialFailure() {
	i.insertPartialFailure.Inc()
}

func (i instrument) ModifyCall() {
	i.modifyCall.Inc()
}

func (i instrument) ModifyKeys(n int) {
	i.modifyKeys.Add(float64(n))
}

func (i instrument) ModifySendTo(n int) {
	i.modifySendTo.Add(float64(n))
}

func (i instrument) ModifyDuration(t time.Duration) {
	i.modifyDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ModifyRetrieved(n int) {
	i.modifyRetrieved.Add(float64(n))
}

func (i instrument) ModifyReturned(n int) {
	i.modifyReturned.Add(float64(n))
}

func (i instrument) ModifyQuorumFailure() {
	i.modifyQuorumFailure.Inc()
}

func (i instrument) ModifyRepairRequired() {
	i.modifyRepairRequired.Inc()
}

func (i instrument) DeleteCall() {
	i.deleteCall.Inc()
}

func (i instrument) DeleteKeys(n int) {
	i.deleteKeys.Add(float64(n))
}

func (i instrument) DeleteSendTo(n int) {
	i.deleteSendTo.Add(float64(n))
}

func (i instrument) DeleteDuration(t time.Duration) {
	i.deleteDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) DeleteRetrieved(n int) {
	i.deleteRetrieved.Add(float64(n))
}

func (i instrument) DeleteReturned(n int) {
	i.deleteReturned.Add(float64(n))
}

func (i instrument) DeleteQuorumFailure() {
	i.deleteQuorumFailure.Inc()
}

func (i instrument) DeleteRepairRequired() {
	i.deleteRepairRequired.Inc()
}

func (i instrument) DeletePartialFailure() {
	i.deletePartialFailure.Inc()
}

func (i instrument) RollbackCall() {
	i.rollbackCall.Inc()
}

func (i instrument) RollbackKeys(n int) {
	i.rollbackKeys.Add(float64(n))
}

func (i instrument) RollbackSendTo(n int) {
	i.rollbackSendTo.Add(float64(n))
}

func (i instrument) RollbackDuration(t time.Duration) {
	i.rollbackDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) RollbackRetrieved(n int) {
	i.rollbackRetrieved.Add(float64(n))
}

func (i instrument) RollbackReturned(n int) {
	i.rollbackReturned.Add(float64(n))
}

func (i instrument) RollbackQuorumFailure() {
	i.rollbackQuorumFailure.Inc()
}

func (i instrument) RollbackRepairRequired() {
	i.rollbackRepairRequired.Inc()
}

func (i instrument) RollbackPartialFailure() {
	i.rollbackPartialFailure.Inc()
}

func (i instrument) SelectCall() {
	i.selectCall.Inc()
}

func (i instrument) SelectKeys(n int) {
	i.selectKeys.Add(float64(n))
}

func (i instrument) SelectSendTo(n int) {
	i.selectSendTo.Add(float64(n))
}

func (i instrument) SelectSendAllPromotion() {
	i.selectSendAllPromotion.Inc()
}

func (i instrument) SelectPartialError() {
	i.selectPartialError.Inc()
}

func (i instrument) SelectDuration(t time.Duration) {
	i.selectDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) SelectRetrieved(n int) {
	i.selectRetrieved.Add(float64(n))
}

func (i instrument) SelectReturned(n int) {
	i.selectReturned.Add(float64(n))
}

func (i instrument) SelectFirstResponseDuration(t time.Duration) {
	i.selectFirstResponseDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) SelectBlockingDuration(t time.Duration) {
	i.selectBlockingDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) SelectOverheadDuration(t time.Duration) {
	i.selectOverheadDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) SelectRepairNeeded() {
	i.selectRepairRequired.Inc()
}

func (i instrument) ScanCall() {
	i.scanCall.Inc()
}

func (i instrument) ScanSendTo(n int) {
	i.scanSendTo.Add(float64(n))
}

func (i instrument) ScanPartialError() {
	i.scanPartialFailure.Inc()
}

func (i instrument) ScanDuration(t time.Duration) {
	i.scanDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) ScanRetrieved(n int) {
	i.scanRetrieved.Add(float64(n))
}

func (i instrument) ScanReturned(n int) {
	i.scanReturned.Add(float64(n))
}

func (i instrument) ScanRepairNeeded(n int) {
	i.scanRepairRequired.Add(float64(n))
}

func (i instrument) RepairCall() {
	i.repairCall.Inc()
}

func (i instrument) RepairRequest(n int) {
	i.repairRequest.Add(float64(n))
}

func (i instrument) RepairSendTo(n int) {
	i.repairSendTo.Add(float64(n))
}

func (i instrument) RepairDuration(t time.Duration) {
	i.repairDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) RepairScoreError() {
	i.repairScoreError.Inc()
}

func (i instrument) RepairError(n int) {
	i.repairError.Add(float64(n))
}

func (i instrument) PerformanceDuration(t time.Duration) {
	i.performanceDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) PerformanceNamespaceDuration(ns string, t time.Duration) {
	var (
		summary prometheus.Summary
		ok      bool
	)
	if summary, ok = i.performanceNamespaceDuration[ns]; !ok {
		i.mutex.Lock()

		summary = prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: i.prefix,
			Name:      fmt.Sprintf("performance.%s.duration", ns),
			Help:      fmt.Sprintf("How long the performance calls for namespace %s, took in nanoseconds.", ns),
			MaxAge:    i.maxSummaryAge,
		})

		if err := prometheus.Register(summary); err == nil {
			i.performanceNamespaceDuration[ns] = summary
		}

		i.mutex.Unlock()
	}
	summary.Observe(float64(t.Nanoseconds()))
}

func (i instrument) PublishCall() {
	i.publishCall.Inc()
}

func (i instrument) PublishKeys(n int) {
	i.publishKeys.Add(float64(n))
}

func (i instrument) PublishSendTo(n int) {
	i.publishSendTo.Add(float64(n))
}

func (i instrument) PublishRetrieved(n int) {
	i.publishRetrieved.Add(float64(n))
}

func (i instrument) PublishReturned(n int) {
	i.publishReturned.Add(float64(n))
}

func (i instrument) PublishDuration(t time.Duration) {
	i.publishDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) SemaphoreCall() {
	i.semaphoreCall.Inc()
}

func (i instrument) SemaphoreSendTo(n int) {
	i.semaphoreSendTo.Add(float64(n))
}

func (i instrument) SemaphoreDuration(t time.Duration) {
	i.semaphoreDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) SemaphoreRetrieved(n int) {
	i.semaphoreRetrieved.Add(float64(n))
}

func (i instrument) SemaphoreReturned(n int) {
	i.semaphoreReturned.Add(float64(n))
}

func (i instrument) HeartbeatCall() {
	i.heartbeatCall.Inc()
}

func (i instrument) HeartbeatSendTo(n int) {
	i.heartbeatSendTo.Add(float64(n))
}

func (i instrument) HeartbeatDuration(t time.Duration) {
	i.heartbeatDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) HeartbeatRetrieved(n int) {
	i.heartbeatRetrieved.Add(float64(n))
}

func (i instrument) HeartbeatReturned(n int) {
	i.heartbeatReturned.Add(float64(n))
}

func (i instrument) KeyStoreCall() {
	i.keyStoreCall.Inc()
}

func (i instrument) KeyStoreSendTo(n int) {
	i.keyStoreSendTo.Add(float64(n))
}

func (i instrument) KeyStoreDuration(t time.Duration) {
	i.keyStoreDuration.Observe(float64(t.Nanoseconds()))
}

func (i instrument) KeyStoreRetrieved(n int) {
	i.keyStoreRetrieved.Add(float64(n))
}

func (i instrument) KeyStoreReturned(n int) {
	i.keyStoreReturned.Add(float64(n))
}
