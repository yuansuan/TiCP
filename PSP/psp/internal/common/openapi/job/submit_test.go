package job

import (
	"fmt"
	"testing"

	apijobcreate "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
)

func TestAdminSubmit(t *testing.T) {
	config.InitConfig()
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return
	}

	var inputPath = "http://10.0.10.2:8080/4Afa3ivYikw/yskj/starccm_test1"
	//var outputPath = "http://10.0.4.55:8899/4Rsuvx4BVfq/result/4RBveK3ibGy/ut_result/"

	cores := int(1)
	memory := int(256)
	var params apijobcreate.Params = apijobcreate.Params{
		Application: apijobcreate.Application{
			Command: `srun hostname > ./hostlist; hostlist=$( cat ./hostlist | sort | uniq -c | awk '{print $2":"$1}' | tr "\n" "," | sed 's/.$//' ); $starccm -rsh /usr/bin/ssh -power -on $hostlist -mpi openmpi -batch run $file_name`,
			AppID:   "4SGaPPLPekE", // starccm+ 17.04
		},
		Resource: &apijobcreate.Resource{Cores: &cores, Memory: &memory},
		EnvVars: map[string]string{
			"starccm":            "/home/apps/siemens/starccm+_17.04_R8/17.04.008-R8/STAR-CCM+17.04.008-R8/star/bin/starccm+",
			"CDLMD_LICENSE_FILE": "29000@115.159.149.167",
			"file_name":          "Blade.sim",
			"VAR1":               "value1",
			"VAR2":               "value2",
		},
		Input: &apijobcreate.Input{
			Type:        "hpc_storage",
			Source:      inputPath,
			Destination: "",
		},
		TmpWorkdir:        false,
		SubmitWithSuspend: false,
	}

	submitParams := &SubmitParams{
		Queue:   "default",
		Name:    "test123",
		Zone:    "az-jinan",
		Comment: "hello job",
		Params:  params,
	}

	jobID, err := SubmitJob(api, submitParams)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Println("====== jobID: ", jobID)

	//assert.Equal(t, resp != nil, true)
	//assert.Equal(t, err == nil, true)

}
