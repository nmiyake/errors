package printer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// PrintStackWithMessages returns a string representation of the provided error. If the error implements causer and all
// of the stackTracer errors are part of the same stack trace, the returned string representation will print all of the
// messages and then print the longest stack trace once at the end. If the provided error is not of this form, the
// result of printing the provided error using the "%+v" formatting directive is returned.
func PrintSingleStack(err error) string {
	type causer interface {
		Cause() error
	}
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	type errWithStack struct {
		err error
		msg string
		stack errors.StackTrace
	}

	var stackErrs []errWithStack
	errCause := err
	for errCause != nil {
		if s, ok := errCause.(stackTracer); ok {
			stackErrs = append(stackErrs, errWithStack{
				err: errCause,
				msg: errCause.Error(),
				stack: s.StackTrace(),
			})
		}
		if c, ok := errCause.(causer); ok {
			errCause = c.Cause()
		} else {
			break
		}
	}

	if len(stackErrs) == 0 {
		return fmt.Sprintf("%+v", err)
	}

	singleStack := true
	for i := len(stackErrs)-1; i > 0; i-- {
		if !hasSuffix(stackErrs[i].stack, stackErrs[i-1].stack) {
			singleStack = false
			break
		}
	}

	if !singleStack {
		return fmt.Sprintf("%+v", err)
	}

	for i := 0; i < len(stackErrs) - 1; i++ {
		stackErrs[i].msg = currErrMsg(stackErrs[i].err, stackErrs[i+1].err)
	}

	var errs []string
	if errCause != nil {
		// if root cause is non-nil, print its error message if it differs from cause
		stackErrs[len(stackErrs)-1].msg = currErrMsg(stackErrs[len(stackErrs)-1].err, errCause)
		rootErr := errCause.Error()
		if rootErr != stackErrs[len(stackErrs)-1].msg {
			errs = append(errs, errCause.Error())
		}
	}
	for i := len(stackErrs)-1; i >= 0; i-- {
		errs = append(errs, fmt.Sprintf("%v", stackErrs[i].msg))
	}
	return strings.Join(errs, "\n") + fmt.Sprintf("%+v", stackErrs[len(stackErrs)-1].stack)
}

// PrintStackWithMessages returns a string representation of the provided error. If the error implements causer and
// the error and its causes alternate between stackTracer and non-stackTracer errors, the returned string representation
// is one in which the messages are interleaved at the proper positions in the stack and common stack frames in
// consecutive elements are removed. If the provided error is not of this form, the result of printing the provided
// error using the "%+v" formatting directive is returned.
func PrintStackWithMessages(err error) string {
	type causer interface {
		Cause() error
	}
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	type errWithStack struct {
		err error
		msg string
		stack errors.StackTrace
	}

	// error is expected to alternate between stackTracer and non-stackTracer errors
	expectStackTracer := true
	var stackErrs []errWithStack
	errCause := err
	for errCause != nil {
		s, isStackTracer := errCause.(stackTracer)
		if isStackTracer {
			stackErrs = append(stackErrs, errWithStack{
				err: errCause,
				msg: errCause.Error(),
				stack: s.StackTrace(),
			})
		}
		if c, ok := errCause.(causer); ok {
			if isStackTracer != expectStackTracer {
				// if current error is a causer and does not meet the expectation of whether or not it
				// should be a stackTracer, error is not supported.
				stackErrs = nil
				break
			}
			errCause = c.Cause()
		} else {
			break
		}
		expectStackTracer = !expectStackTracer
	}

	if len(stackErrs) == 0 {
		return fmt.Sprintf("%+v", err)
	}

	for i := len(stackErrs)-1; i > 0; i-- {
		if hasSuffix(stackErrs[i].stack, stackErrs[i-1].stack) {
			// if the inner stack has the outer stack as a suffix, trim the outer stack from the inner stack
			stackErrs[i].stack = stackErrs[i].stack[0:len(stackErrs[i].stack)-len(stackErrs[i-1].stack)]
		}
	}

	for i := 0; i < len(stackErrs) - 1; i++ {
		stackErrs[i].msg = currErrMsg(stackErrs[i].err, stackErrs[i+1].err)
	}

	var errs []string
	if errCause != nil {
		// if root cause is non-nil, print its error message if it differs from cause
		stackErrs[len(stackErrs)-1].msg = currErrMsg(stackErrs[len(stackErrs)-1].err, errCause)
		rootErr := errCause.Error()
		if rootErr != stackErrs[len(stackErrs)-1].msg {
			errs = append(errs, errCause.Error())
		}
	}
	for i := len(stackErrs)-1; i >= 0; i-- {
		errs = append(errs, fmt.Sprintf("%v%+v", stackErrs[i].msg, stackErrs[i].stack))
	}
	return strings.Join(errs, "\n")
}

// currErrMsg Returns the string for the current error. If curr.Error() has the suffix ": {{next.Error()}}", returns the
// content of "curr.Error()" up to that suffix.
func currErrMsg(curr error, next error) string {
	msg := curr.Error()
	if idx := strings.LastIndex(msg, ": " + next.Error()); idx != -1 {
		return msg[:idx]
	}
	return msg
}

// hasSuffix returns true if the inner stack trace ends with the outer stack trace, false otherwise.
func hasSuffix(inner errors.StackTrace, outer errors.StackTrace) bool {
	outerIndex := len(outer)-1
	innerIndex := len(inner)-1
	for outerIndex >= 0 && innerIndex >= 0 {
		if outer[outerIndex] != inner[innerIndex] {
			break
		}
		outerIndex--
		innerIndex--
	}
	return outerIndex == 0 && innerIndex >= 0
}
