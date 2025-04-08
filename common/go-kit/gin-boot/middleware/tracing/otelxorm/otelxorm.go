package otelxorm

import (
	"context"
	"strings"

	"github.com/xwb1989/sqlparser"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"xorm.io/xorm/contexts"

	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/tracing"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

const (
	instrumentationName = "yuansuan.cn/tracing/otelxorm"
)

type tracingSpanKeyType int

const tracingSpanKey tracingSpanKeyType = iota

type tracingHook struct {
	cfg    *conf_type.Tracing
	tracer trace.Tracer
}

func (th *tracingHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	span := trace.SpanFromContext(c.Ctx)
	// not tracing and ignore dangling
	if !span.IsRecording() && !th.cfg.Database.Dangling {
		return c.Ctx, nil
	}

	opts := []oteltrace.SpanStartOption{
		oteltrace.WithAttributes(
			semconv.DBSystemMySQL,
		),
		oteltrace.WithAttributes(semconv.PeerServiceKey.String("xorm")),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	}

	c.Ctx, span = th.tracer.Start(c.Ctx, spanName(c.SQL), opts...)
	c.Ctx = context.WithValue(c.Ctx, tracingSpanKey, span)

	attrs := []attribute.KeyValue{
		attribute.String("sql", c.SQL),
	}
	if th.cfg.Database.Binding {
		attrs = append(attrs, attribute.String("args", util.MustParseJSON(c.Args)))
	}
	span.AddEvent("before", trace.WithAttributes(attrs...))

	return c.Ctx, nil
}

func (th *tracingHook) AfterProcess(c *contexts.ContextHook) error {
	v := c.Ctx.Value(tracingSpanKey)
	if v == nil {
		return nil
	}

	span := v.(trace.Span)

	var lastInsertId, affectedRows int64
	if c.Result != nil {
		lastInsertId, _ = c.Result.LastInsertId()
		affectedRows, _ = c.Result.RowsAffected()
	}

	span.AddEvent("after", trace.WithAttributes(
		attribute.Int64("lastInsertId", lastInsertId),
		attribute.Int64("affectedRows", affectedRows),
		attribute.String("executeTime", c.ExecuteTime.String()),
	))

	if c.Err != nil {
		span.RecordError(c.Err)
	}

	span.End()
	return nil
}

func NewTracingHook(cfg *conf_type.Tracing) *tracingHook {
	tp := otel.GetTracerProvider()
	return &tracingHook{
		cfg:    cfg,
		tracer: tp.Tracer(instrumentationName, oteltrace.WithInstrumentationVersion(tracing.SemVersion())),
	}
}

func spanName(sql string) string {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return "UNKNOWN - UNKNOWN"
	}

	var action, table string
	switch node := stmt.(type) {
	case *sqlparser.Select:
		action = "SELECT"
		table = toString(node.From)
	case *sqlparser.Insert:
		action = "INSERT"
		table = toString(node.Table.Name)
	case *sqlparser.Update:
		action = "UPDATE"
		table = toString(node.TableExprs)
	case *sqlparser.Delete:
		action = "DELETE"
		table = toString(node.TableExprs)
	default:
		action = "UNKNOWN"
		table = "UNKNOWN"
	}

	return action + " - " + table
}

func toString(node sqlparser.SQLNode) string {
	buf := sqlparser.NewTrackedBuffer(nil)
	node.Format(buf)
	return strings.Split(buf.String(), " ")[0]
}
