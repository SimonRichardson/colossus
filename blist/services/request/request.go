package request

import (
	"io"
	"sync"
	"time"

	"net/http"

	"bytes"

	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	sv "github.com/SimonRichardson/colossus/blist/services"
	"github.com/SimonRichardson/colossus/typex"
)

func Request(s *Service, t Tactic) selectors.Request {
	return request{s, t}
}

type request struct {
	service *Service
	tactic  Tactic
}

func (s request) Request(req *http.Request) (*http.Response, error) {
	return s.write(req, func(c Cluster, r *http.Request) <-chan sv.Element {
		return c.Request(r)
	})
}

func (s request) write(req *http.Request,
	fn func(Cluster, *http.Request) <-chan sv.Element,
) (*http.Response, error) {
	var (
		service      = s.service
		fitting, err = selectClusters(service.fitting, service.clusters)

		numOfWrites    = len(fitting.Clusters())
		numOfResponses = fitting.Required().Len()

		retrieved = 0
		returned  = 0
	)
	if err != nil {
		return nil, err
	}

	began := beforeWrite(service.instrumentation, numOfWrites)
	defer afterWrite(service.instrumentation, began, retrieved, returned)

	var (
		elements = make(chan sv.Element, numOfResponses)
		errs     = []error{}
		changes  = []*http.Response{}

		wg = &sync.WaitGroup{}
	)

	wg.Add(numOfResponses)
	go func() { wg.Wait(); close(elements) }()

	scatterReads(s.tactic, fitting, req, fn, wg, elements)

	for element := range elements {
		retrieved++

		if err := sv.ErrorFromElement(element); err != nil {
			errs = append(errs, err)
			continue
		}

		index := sv.IndexFromElement(element)
		if !fitting.Required().Contains(index) {
			errs = append(errs, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Unexpected result from cluster (%d)", index))
			continue
		}

		response := sv.ResponseFromElement(element)
		changes = append(changes, response)

		returned++
	}

	if len(errs) > 0 {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure").With(errs...)
	}

	return head(changes)
}

func selectClusters(fitting Provider, clusters []Cluster) (Fitting, error) {
	if len(clusters) < 1 {
		return nil, typex.Errorf(errors.Source, errors.UnexpectedArgument, "No clusters to locate.")
	}
	return fitting(clusters), nil
}

func head(x []*http.Response) (*http.Response, error) {
	if len(x) < 1 {
		return nil, typex.Errorf(errors.Source, errors.Complete,
			"Complete failure: no valid changes")
	}
	return x[0], nil
}

func scatterReads(tactic Tactic,
	fitting Fitting,
	req *http.Request,
	fn func(Cluster, *http.Request) <-chan sv.Element,
	wg *sync.WaitGroup,
	dst chan sv.Element,
) error {
	var (
		clusters      = fitting.Clusters()
		requests, err = duplicateRequests(req, len(clusters))
	)

	if err != nil {
		return err
	}

	return tactic(clusters, func(index int, n Cluster) {
		request := requests[index]

		if req := fitting.Required(); req.Len() > 0 && req.Contains(n.Index()) {
			defer wg.Done()

			for e := range fn(n, request) {
				dst <- e
			}
			return
		}

		// We don't care about the result, so just send it!
		fn(n, request)
	})
}

func beforeWrite(instr Instrumentation, numSends int) time.Time {
	began := time.Now()
	go func() {
		instr.RequestCall()
		instr.RequestSendTo(numSends)
	}()
	return began
}

func afterWrite(instr Instrumentation, began time.Time, retrieved, returned int) {
	go func() {
		instr.RequestDuration(time.Since(began))
		instr.RequestRetrieved(retrieved)
		instr.RequestReturned(returned)
	}()
}

type ReadWriter struct {
	Bytes *bytes.Buffer
}

func NewReadWriter() *ReadWriter {
	return &ReadWriter{Bytes: new(bytes.Buffer)}
}

func (r *ReadWriter) Read(p []byte) (int, error) {
	return r.Bytes.Read(p)
}

func (r *ReadWriter) Write(p []byte) (int, error) {
	return r.Bytes.Write(p)
}

func (r *ReadWriter) Close() error {
	return nil
}

func duplicateRequests(req *http.Request, amount int) ([]*http.Request, error) {
	// Make new buffers for all the requests
	buffers := make([]*ReadWriter, 0, amount)
	for i := 0; i < amount; i++ {
		buffers = append(buffers, NewReadWriter())
	}

	// Setup the ability to write to all the buffers at once.
	writers := make([]io.Writer, 0, amount)
	for _, v := range buffers {
		writers = append(writers, v)
	}
	writer := io.MultiWriter(writers...)

	if _, err := io.Copy(writer, req.Body); err != nil {
		return nil, err
	}
	defer req.Body.Close()

	// Make a new request for each amount required.
	res := make([]*http.Request, 0, amount)
	for _, v := range buffers {
		res = append(res, duplicateRequest(req, v))
	}

	return res, nil
}

func duplicateRequest(req *http.Request, body io.ReadCloser) *http.Request {
	return &http.Request{
		Method:        req.Method,
		URL:           req.URL,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Header:        req.Header,
		Body:          body,
		Host:          req.Host,
		ContentLength: req.ContentLength,
		Close:         true,
	}
}
