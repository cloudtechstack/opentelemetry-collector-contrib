// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package awss3exporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awss3exporter"

import (
	"errors"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type marshaler interface {
	MarshalTraces(td ptrace.Traces) ([]byte, error)
	MarshalLogs(ld plog.Logs) ([]byte, error)
	MarshalMetrics(md pmetric.Metrics) ([]byte, error)
	format() string
}

var (
	ErrUnknownMarshaler = errors.New("unknown marshaler")
)

func newMarshaler(mType MarshalerType, logger *zap.Logger) (marshaler, error) {
	marshaler := &s3Marshaler{logger: logger}
	switch mType {
	case OtlpProtobuf:
		marshaler.logsMarshaler = &plog.ProtoMarshaler{}
		marshaler.tracesMarshaler = &ptrace.ProtoMarshaler{}
		marshaler.metricsMarshaler = &pmetric.ProtoMarshaler{}
		marshaler.fileFormat = "binpb"
	case OtlpJSON:
		marshaler.logsMarshaler = &plog.JSONMarshaler{}
		marshaler.tracesMarshaler = &ptrace.JSONMarshaler{}
		marshaler.metricsMarshaler = &pmetric.JSONMarshaler{}
		marshaler.fileFormat = "json"
	case SumoIC:
		sumomarshaler := newSumoICMarshaler()
		marshaler.logsMarshaler = &sumomarshaler
		marshaler.fileFormat = "json.gz"
	case Body:
		exportbodyMarshaler := newbodyMarshaler()
		marshaler.logsMarshaler = &exportbodyMarshaler
		marshaler.fileFormat = exportbodyMarshaler.format()
	default:
		return nil, ErrUnknownMarshaler
	}
	return marshaler, nil
}
