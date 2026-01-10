package processor

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockIn struct {
	data         []string
	errLine      error
	errAvailable bool
	status       int
	statusErr    error
	currentIndex int
}

func (m *mockIn) Get() (string, bool) {
	if m.currentIndex < len(m.data) {
		line := m.data[m.currentIndex]
		m.currentIndex++
		return line, true
	}
	return "", false
}

func (m *mockIn) GetErr() (error, bool) {
	return m.errLine, m.errAvailable
}

func (m *mockIn) GetStatus() (int, error) {
	return m.status, m.statusErr
}

type mockOut struct {
	storedData   []string
	err          error
	statusCode   int
	errRecorded  error
	statusErrMsg string
}

func (m *mockOut) Set(data string) error {
	if m.err != nil {
		return m.err
	}
	m.storedData = append(m.storedData, data)
	return nil
}

func (m *mockOut) SetErr(data error) error {
	m.errRecorded = data
	return nil
}

func (m *mockOut) SetStatus(status int) error {
	if m.statusErrMsg != "" {
		return errors.New(m.statusErrMsg)
	}
	m.statusCode = status
	return nil
}

type testEndpoint struct {
	in             *mockIn
	out            *mockOut
	expectedErr    bool
	expectedData   []string
	expectedSetErr error
	expectedStatus int
	closeContext   bool
}

var _ = Describe("Endpoint", func() {
	Context("parseUintSlice", func() {
		DescribeTable("Sunny",
			func(data *testEndpoint) {
				ctx, cancel := context.WithCancel(context.Background())
				if data.closeContext {
					cancel()
				} else {
					defer cancel()
				}

				svc := &endpoint{
					in:  data.in,
					out: data.out,
				}

				err := svc.Run(ctx)

				if (err != nil) != data.expectedErr {
					Expect(err).To(Equal(data.expectedErr))
				}
				if !data.expectedErr {
					if data.expectedSetErr != nil {
						Expect(data.out.errRecorded).To(Equal(data.expectedSetErr))
					}
					Expect(data.out.statusCode).To(Equal(data.expectedStatus))
					if data.expectedData != nil {
						Expect(data.out.storedData).To(HaveLen(len(data.expectedData)))
					}
				}

			},
			Entry("All data processed successfully", &testEndpoint{
				in: &mockIn{
					data:         []string{"line1", "line2"},
					errAvailable: false,
				},
				out:          &mockOut{},
				expectedErr:  false,
				expectedData: []string{"line1", "line2"},
			}),
			Entry("Output Set fails on data", &testEndpoint{
				in: &mockIn{
					data: []string{"line1"},
				},
				out: &mockOut{
					err: errors.New("cannot store"),
				},
				expectedErr: true,
			}),
			Entry("Input status returned successfully", &testEndpoint{
				in: &mockIn{
					status: 0,
				},
				out:            &mockOut{},
				expectedErr:    false,
				expectedStatus: 0,
			}),
			Entry("Input status fails to retrieve", &testEndpoint{
				in: &mockIn{
					statusErr: errors.New("status retrieval failed"),
				},
				out:            &mockOut{},
				expectedErr:    false,
				expectedSetErr: errors.New("status retrieval failed"),
			}),
			Entry("Output SetStatus fails", &testEndpoint{
				in: &mockIn{
					status: 2,
				},
				out: &mockOut{
					statusErrMsg: "output status failure",
				},
				expectedErr:    false,
				expectedSetErr: errors.New("output status failure"),
			}),
			Entry("Context cancellation during processing", &testEndpoint{
				in: &mockIn{
					data: []string{"line1", "line2"},
				},
				out:          &mockOut{},
				expectedErr:  false,
				expectedData: []string{},
				closeContext: true,
			}),
		)
	})

})
