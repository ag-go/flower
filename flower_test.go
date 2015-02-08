package flower_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cheekybits/is"
	"github.com/matryer/flower"
)

func TestSimple(t *testing.T) {
	is := is.New(t)

	// make a manager
	manager := flower.New()
	is.OK(manager)

	// add a handler
	var calls []*flower.Job
	manager.On("event", func(j *flower.Job) {
		calls = append(calls, j)
	})

	// trigger it with some data
	data := map[string]interface{}{"data": true}
	job, err := manager.New(data, "event")
	is.NoErr(err)
	is.OK(job)
	is.OK(job.ID())

	// wait for the job to finish
	job.Wait()
	is.Equal(1, len(calls))
	is.Equal(job, calls[0])
	is.Equal(data, calls[0].Data)

}

func TestPath(t *testing.T) {
	is := is.New(t)

	// make a manager
	manager := flower.New()

	// add three handlers
	var calls []string
	manager.On("one", func(j *flower.Job) {
		calls = append(calls, "one")
	})
	manager.On("two", func(j *flower.Job) {
		calls = append(calls, "two")
	})
	manager.On("three", func(j *flower.Job) {
		calls = append(calls, "three")
	})

	data := map[string]interface{}{"data": true}
	job, err := manager.New(data, "one", "two", "three")

	is.NoErr(err)
	is.OK(job)

	job.Wait()
	is.Equal(len(calls), 3)
	is.Equal(calls[0], "one")
	is.Equal(calls[1], "two")
	is.Equal(calls[2], "three")

}

func TestAbort(t *testing.T) {

	is := is.New(t)

	// make a manager
	manager := flower.New()
	is.OK(manager)

	// add a handler
	var ticks []*flower.Job
	manager.On("event", func(j *flower.Job) {
		for {
			ticks = append(ticks, j)
			time.Sleep(100 * time.Millisecond)
			if j.ShouldStop() {
				break
			}
		}
	})

	// trigger it with some data
	data := map[string]interface{}{"data": true}
	job, err := manager.New(data, "event")
	is.NoErr(err)
	is.OK(job)
	is.OK(job.ID())

	// tell the job to stop in 100 milliseconds
	go func() {
		time.Sleep(1000 * time.Millisecond)
		job.Abort()
	}()

	// wait for the job to finish
	job.Wait()
	is.Equal(10, len(ticks))

}

func TestErrs(t *testing.T) {
	is := is.New(t)

	// make a manager
	manager := flower.New()
	is.OK(manager)

	// add a handler
	err := errors.New("something went wrong")
	manager.On("event", func(j *flower.Job) {
		j.Err = err
	})

	// trigger it with some data
	data := map[string]interface{}{"data": true}
	job, err := manager.New(data, "event")
	is.NoErr(err)
	is.OK(job)
	is.OK(job.ID())

	// wait for the job to finish
	job.Wait()
	is.Equal(err, job.Err)

}