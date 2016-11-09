package mapreduce

import (
    "fmt"
    "net/rpc"
)


// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO

        // Use buffered channel as semaphore. Buffered channel is filled 
        // when all tasks done.
        done := make(chan bool, ntasks)

        for i := 0; i < ntasks; i += 1 {
            worker := <-mr.registerChannel
            args := DoTaskArgs{
                JobName: mr.jobName,
                File: func() string {
                    if phase == mapPhase {
                        return mr.files[i]
                    }
                    return ""
                    }(),
                Phase: phase,
                TaskNumber: i,
                NumOtherPhase : nios,
            }
            client, err := rpc.Dial("unix", worker)
            // Handling worker failures, just redo the task.
            if err != nil {
                i -= 1
                continue
            }
            called := client.Go("Worker.DoTask", args, new(struct{}), nil)
	    go func() {
                <-called.Done
                go func() {
	            mr.registerChannel <- worker
                }()
                done <- true
	    }()
        }

        // Wait for all tasks done.
        for i := 0; i < ntasks; i += 1 {
            <-done
        }
	fmt.Printf("Schedule: %v phase done\n", phase)
}
