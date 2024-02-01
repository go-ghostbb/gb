package gbclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"ghostbb.io/internal/utils"
	gbtrace "ghostbb.io/net/gb_trace"
	gbstr "ghostbb.io/text/gb_str"
	gbconv "ghostbb.io/util/gb_conv"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// clientTracer is used for implementing httptrace.ClientTrace.
type clientTracer struct {
	context.Context
	span        trace.Span
	request     *http.Request
	requestBody []byte
	headers     map[string]interface{}
	mtx         sync.Mutex
}

// newClientTrace creates and returns object of newClientTrace.
func newClientTrace(ctx context.Context, span trace.Span, request *http.Request) *httptrace.ClientTrace {
	ct := &clientTracer{
		Context: ctx,
		span:    span,
		request: request,
		headers: make(map[string]interface{}),
	}

	reqBodyContent, _ := io.ReadAll(ct.request.Body)
	ct.requestBody = reqBodyContent
	ct.request.Body = utils.NewReadCloser(reqBodyContent, false)

	return &httptrace.ClientTrace{
		GetConn:              ct.getConn,
		GotConn:              ct.gotConn,
		PutIdleConn:          ct.putIdleConn,
		GotFirstResponseByte: ct.gotFirstResponseByte,
		Got100Continue:       ct.got100Continue,
		Got1xxResponse:       ct.got1xxResponse,
		DNSStart:             ct.dnsStart,
		DNSDone:              ct.dnsDone,
		ConnectStart:         ct.connectStart,
		ConnectDone:          ct.connectDone,
		TLSHandshakeStart:    ct.tlsHandshakeStart,
		TLSHandshakeDone:     ct.tlsHandshakeDone,
		WroteHeaderField:     ct.wroteHeaderField,
		WroteHeaders:         ct.wroteHeaders,
		Wait100Continue:      ct.wait100Continue,
		WroteRequest:         ct.wroteRequest,
	}
}

func (ct *clientTracer) getConn(host string) {}

func (ct *clientTracer) gotConn(info httptrace.GotConnInfo) {
	remoteAddr := ""
	if info.Conn.RemoteAddr() != nil {
		remoteAddr = info.Conn.RemoteAddr().String()
	}
	localAddr := ""
	if info.Conn.LocalAddr() != nil {
		localAddr = info.Conn.LocalAddr().String()
	}
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpAddressRemote, remoteAddr),
		attribute.String(tracingAttrHttpAddressLocal, localAddr),
	)
}

func (ct *clientTracer) putIdleConn(err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
}

func (ct *clientTracer) dnsStart(info httptrace.DNSStartInfo) {
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpDnsStart, info.Host),
	)
}

func (ct *clientTracer) dnsDone(info httptrace.DNSDoneInfo) {
	var buffer strings.Builder
	for _, v := range info.Addrs {
		if buffer.Len() != 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v.String())
	}
	if info.Err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, info.Err))
	}
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpDnsDone, buffer.String()),
	)
}

func (ct *clientTracer) connectStart(network, addr string) {
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpConnectStart, network+"@"+addr),
	)
}

func (ct *clientTracer) connectDone(network, addr string, err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	ct.span.SetAttributes(
		attribute.String(tracingAttrHttpConnectDone, network+"@"+addr),
	)
}

func (ct *clientTracer) tlsHandshakeStart() {}

func (ct *clientTracer) tlsHandshakeDone(_ tls.ConnectionState, err error) {
	if err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
}

func (ct *clientTracer) wroteHeaderField(k string, v []string) {
	if len(v) > 1 {
		ct.headers[k] = v
	} else if len(v) == 1 {
		ct.headers[k] = v[0]
	}
}

func (ct *clientTracer) wroteHeaders() {}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, info.Err))
	}

	ct.span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(
		attribute.String(tracingEventHttpRequestHeaders, gbconv.String(ct.headers)),
		attribute.String(tracingEventHttpRequestBaggage, gbtrace.GetBaggageMap(ct.Context).String()),
		attribute.String(tracingEventHttpRequestBody, gbstr.StrLimit(
			string(ct.requestBody),
			gbtrace.MaxContentLogSize(),
			"...",
		)),
	))
}

func (ct *clientTracer) got100Continue() {}

func (ct *clientTracer) wait100Continue() {}

func (ct *clientTracer) gotFirstResponseByte() {}

func (ct *clientTracer) got1xxResponse(code int, header textproto.MIMEHeader) error {
	return nil
}
