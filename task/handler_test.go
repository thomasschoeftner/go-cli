package task

import "testing"

func TestJobCopy(t *testing.T) {
	job := Job{"akey" : "aval", "bkey" : "bval"}
	clone := job.Copy()
	if &job == &clone {
		t.Error("job copy failed - copied job has same address as original")
	}
	if len(job) != len(clone) {
		t.Errorf("job copy failed - expected %d key/vals, but got %d key/vals", len(job), len(clone))
	}
	for key, val := range job {
		if _, found := clone[key]; !found {
			t.Errorf("job copy failed - expected %s:%s in copy, but was not found", key, val)
		}
	}
}

func TestWithParamImplicatesJobCopy(t *testing.T) {
	job := Job{"akey" : "aval", "bkey" : "bval"}
	clone := job.WithParam("ckey", "cval")

	if &job == &clone {
		t.Error("with param failed - extended job has same address as original")
	}

	if len(job) + 1 != len(clone) {
		t.Errorf("with param failed - extended job, should have 1 element more than original, but has %d", len(clone))
	}
}