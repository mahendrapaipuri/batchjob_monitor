package jobstats

import (
	"reflect"
	"testing"

	"github.com/go-kit/log"
)

var (
	sacctCmdOutput = `JobID|Partition|Account|Group|GID|User|UID|Submit|Start|End|Elapsed|ExitCode|State|NNodes|NodeList|JobName|WorkDir
1479763|part1|acc1|grp|1000|usr|1000|2023-02-21T14:37:02|2023-02-21T14:37:07|2023-02-21T15:26:29|00:49:22|0:0|CANCELLED by 302137|1|compute-0|test_script1|/home/usr
1481508|part1|acc1|grp|1000|usr|1000|2023-02-21T15:48:20|2023-02-21T15:49:06|2023-02-21T15:57:23|00:08:17|0:0|CANCELLED by 302137|2|compute-[0-2]|test_script2|/home/usr`
	logger            = log.NewNopLogger()
	expectedBatchJobs = []BatchJob{
		{
			Jobid:       "1479763",
			Jobuuid:     "aaaf154c-e784-9e49-2155-4aa52462782a",
			Partition:   "part1",
			Account:     "acc1",
			Grp:         "grp",
			Gid:         "1000",
			Usr:         "usr",
			Uid:         "1000",
			Submit:      "2023-02-21T14:37:02",
			Start:       "2023-02-21T14:37:07",
			End:         "2023-02-21T15:26:29",
			Elapsed:     "00:49:22",
			Exitcode:    "0:0",
			State:       "CANCELLED by 302137",
			Nnodes:      "1",
			Nodelist:    "compute-0",
			NodelistExp: "compute-0",
		},
		{
			Jobid:       "1481508",
			Jobuuid:     "69683fe8-5d89-9ec5-4b4f-8404c7cc37f2",
			Partition:   "part1",
			Account:     "acc1",
			Grp:         "grp",
			Gid:         "1000",
			Usr:         "usr",
			Uid:         "1000",
			Submit:      "2023-02-21T15:48:20",
			Start:       "2023-02-21T15:49:06",
			End:         "2023-02-21T15:57:23",
			Elapsed:     "00:08:17",
			Exitcode:    "0:0",
			State:       "CANCELLED by 302137",
			Nnodes:      "2",
			Nodelist:    "compute-[0-2]",
			NodelistExp: "compute-0|compute-1|compute-2",
		},
	}
)

func TestParseSacctCmdOutput(t *testing.T) {
	batchJobs, numJobs := parseSacctCmdOutput(sacctCmdOutput, logger)
	if !reflect.DeepEqual(batchJobs, expectedBatchJobs) {
		t.Errorf("Expected batch jobs %#v. \n\nGot %#v", expectedBatchJobs, batchJobs)
	}
	if numJobs != 2 {
		t.Errorf("Expected batch jobs num %d. Got %d", 2, numJobs)
	}
}