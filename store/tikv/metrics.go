// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package tikv

import (
	"github.com/pingcap/kvproto/pkg/errorpb"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	txnCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "txn_total",
			Help:      "Counter of created txns.",
		})

	snapshotCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "snapshot_total",
			Help:      "Counter of snapshots.",
		})

	txnCmdCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "txn_cmd_total",
			Help:      "Counter of txn commands.",
		}, []string{"type"})

	txnCmdHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "txn_cmd_seconds",
			Help:      "Bucketed histogram of processing time of txn cmds.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 18),
		}, []string{"type"})

	backoffCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "backoff_total",
			Help:      "Counter of backoff.",
		}, []string{"type"})

	backoffHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "backoff_seconds",
			Help:      "Bucketed histogram of sleep seconds of backoff.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 18),
		}, []string{"type"})

	sendReqHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "request_seconds",
			Help:      "Bucketed histogram of sending request duration.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 18),
		}, []string{"type"})

	copBuildTaskHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "cop_buildtask_seconds",
			Help:      "Coprocessor buildTask cost time.",
		})

	copTaskLenHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "cop_task_len",
			Help:      "Coprocessor task length.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 11),
		})

	coprocessorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "coprocessor_actions_total",
			Help:      "Counter of coprocessor actions.",
		}, []string{"type"})

	gcWorkerCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "gc_worker_actions_total",
			Help:      "Counter of gc worker actions.",
		}, []string{"type"})

	gcHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "gc_seconds",
			Help:      "Bucketed histogram of gc duration.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 13),
		}, []string{"stage"})

	gcConfigGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "gc_config",
			Help:      "Gauge of GC configs.",
		}, []string{"type"},
	)

	lockResolverCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "lock_resolver_actions_total",
			Help:      "Counter of lock resolver actions.",
		}, []string{"type"})

	regionErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "region_err_total",
			Help:      "Counter of region errors.",
		}, []string{"type"})

	txnWriteKVCountHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "txn_write_kv_count",
			Help:      "Count of kv pairs to write in a transaction.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 21),
		})

	txnWriteSizeHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "txn_write_size",
			Help:      "Size of kv pairs to write in a transaction. (KB)",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 21),
		})

	rawkvCmdHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "rawkv_cmd_seconds",
			Help:      "Bucketed histogram of processing time of rawkv cmds.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 18),
		}, []string{"type"})

	rawkvSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "tikvclient",
			Name:      "rawkv_kv_size",
			Help:      "Size of key/value to put, in bytes.",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 21),
		}, []string{"type"})
)

func reportRegionError(e *errorpb.Error) {
	if e.GetNotLeader() != nil {
		regionErrorCounter.WithLabelValues("not_leader").Inc()
	} else if e.GetRegionNotFound() != nil {
		regionErrorCounter.WithLabelValues("region_not_found").Inc()
	} else if e.GetKeyNotInRegion() != nil {
		regionErrorCounter.WithLabelValues("key_not_in_region").Inc()
	} else if e.GetStaleEpoch() != nil {
		regionErrorCounter.WithLabelValues("stale_epoch").Inc()
	} else if e.GetServerIsBusy() != nil {
		regionErrorCounter.WithLabelValues("server_is_busy").Inc()
	} else if e.GetStaleCommand() != nil {
		regionErrorCounter.WithLabelValues("stale_command").Inc()
	} else if e.GetStoreNotMatch() != nil {
		regionErrorCounter.WithLabelValues("store_not_match").Inc()
	} else {
		regionErrorCounter.WithLabelValues("unknown").Inc()
	}
}

func init() {
	prometheus.MustRegister(txnCounter)
	prometheus.MustRegister(snapshotCounter)
	prometheus.MustRegister(txnCmdCounter)
	prometheus.MustRegister(txnCmdHistogram)
	prometheus.MustRegister(backoffCounter)
	prometheus.MustRegister(backoffHistogram)
	prometheus.MustRegister(sendReqHistogram)
	prometheus.MustRegister(copBuildTaskHistogram)
	prometheus.MustRegister(copTaskLenHistogram)
	prometheus.MustRegister(coprocessorCounter)
	prometheus.MustRegister(gcWorkerCounter)
	prometheus.MustRegister(gcConfigGauge)
	prometheus.MustRegister(gcHistogram)
	prometheus.MustRegister(lockResolverCounter)
	prometheus.MustRegister(regionErrorCounter)
	prometheus.MustRegister(txnWriteKVCountHistogram)
	prometheus.MustRegister(txnWriteSizeHistogram)
	prometheus.MustRegister(rawkvCmdHistogram)
	prometheus.MustRegister(rawkvSizeHistogram)
}
