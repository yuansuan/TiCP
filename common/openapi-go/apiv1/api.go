package apiv1

import (
	accountCashVoucherAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/accountcashvoucher/add"
	accountCashVoucherGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/accountcashvoucher/get"
	accountCashVoucherList "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/accountcashvoucher/list"
	accountCashVoucherStatusModify "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/accountcashvoucher/statusmodify"
	accountList "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/accountlist"
	accountAmountRefund "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/amountrefund"
	accountUserBillList "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/apiuser/billlist"
	accountUserIdGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/apiuser/idget"
	accountUserResourceBillList "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/apiuser/resourcebilllist"
	accountBillList "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/billlist"
	cashVoucherAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/cashvoucher/add"
	cashVoucherAvaModify "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/cashvoucher/availabilitymodify"
	cashVoucherGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/cashvoucher/get"
	cashVoucherList "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/cashvoucher/list"
	accountCreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/create"
	accountCreditAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/creditadd"
	accountCreditQuotaModify "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/creditquotamodify"
	accountFrozenModify "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/frozenmodify"
	accountIdGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/idget"
	accountIdReduce "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/idreduce"
	accountPaymentFreezeUnfreeze "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/paymentfreezeunfreeze"
	accountPaymentReduce "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/paymentreduce"
	accountYsidGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/ysidget"
	accountYsidReduce "github.com/yuansuan/ticp/common/openapi-go/apiv1/account/ysidreduce"
	remoteAppGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/api/get_url"
	sessionAdminExecScript "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/execscript"
	sessionAdminRestore "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/restore"
	sessionClose "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/close"
	sessionDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/delete"
	sessionExecScript "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/execscript"
	sessionGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/get"
	sessionList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/list"
	sessionMount "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/mount"
	sessionPowerOff "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/poweroff"
	sessionPowerOn "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/poweron"
	sessionReady "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/ready"
	sessionReboot "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/reboot"
	sessionRestore "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/restore"
	sessionStart "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/start"
	sessionUmount "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/api/umount"
	adminAddApp "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/add"
	adminAddAppAllow "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/allowadd"
	adminDeleteAppAllow "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/allowdelete"
	adminGetAppAllow "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/allowget"
	adminDeleteApp "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/delete"
	adminGetApp "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/get"
	adminListApp "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/list"
	adminAddAppQuota "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/quotaadd"
	adminDeleteAppQuota "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/quotadelete"
	adminGetAppQuota "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/quotaget"
	adminUpdateApp "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/app/update"
	adminJobCpuUsage "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobcpuusage"
	adminJobCreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobcreate"
	adminJobDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobdelete"
	adminJobGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobget"
	adminJobList "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/joblist"
	adminJobListFiltered "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/joblistfiltered"
	adminJobGetMonitorChart "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobmonitorchart"
	adminJobGetResidual "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobresidual"
	adminJobRetransmit "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobretransmit"
	adminJobGetSnapshot "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobsnapshotget"
	adminJobListSnapshot "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobsnapshotlist"
	adminJobTerminate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobterminate"
	adminJobUpdate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/admin/jobupdate"
	ListApp "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/applist"
	jobBatchGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobbatchget"
	jobCpuUsage "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobcpuusage"
	jobCreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobcreate"
	jobDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobdelete"
	jobGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobget"
	jobList "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/joblist"
	jobGetMonitorChart "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobmonitorchart"
	jobPreSchedule "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobpreschedule"
	jobGetResidual "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobresidual"
	jobResume "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobresume"
	jobGetSnapshot "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobsnapshotget"
	jobListSnapshot "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobsnapshotlist"
	jobTerminate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobterminate"
	jobTransmitResume "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobtransmitresume"
	jobTransmitSuspend "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/jobtransmitsuspend"
	systemJobNeedSyncFileList "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/system/jobneedsyncfile"
	systemJobSyncFileStateUpdate "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/system/jobsyncfilestate"
	zoneList "github.com/yuansuan/ticp/common/openapi-go/apiv1/job/zonelist"
	rdpgoInternalClean "github.com/yuansuan/ticp/common/openapi-go/apiv1/rdpgo/intern/clean"
	rdpgoInternalExecScript "github.com/yuansuan/ticp/common/openapi-go/apiv1/rdpgo/intern/execscript"
	storageAdminBatchDownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/batchDownload"
	storageAdminDownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/download"
	storageAdminUploadComplete "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/upload/complete"
	storageAdminUploadSlice "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/upload/slice"

	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"

	licenseInfoAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licenseinfo/add"
	licenseInfoDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licenseinfo/delete"
	licenseInfoPut "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licenseinfo/put"
	licenseManagerAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/add"
	licenseManagerDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/delete"
	licenseManagerGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/get"
	licenseManagerList "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/list"
	licenseManagerPut "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/put"
	moduleConfigAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/add"
	moduleConfigBatchAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/batchadd"
	moduleConfigDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/delete"
	moduleConfigList "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/list"
	moduleConfigPut "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/put"

	storageBatchDownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/batchDownload"
	storageCompressCancel "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/compress/cancel"
	storageCompressStart "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/compress/start"
	storageCompressStatus "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/compress/status"
	storageDirectoryUsageCancel "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/directory_usage/cancel"
	storageDirectoryUsageStart "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/directory_usage/start"
	storageDirectoryUsageStatus "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/directory_usage/status"

	storageCopy "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/copy"
	storageCopyRange "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/copyRange"
	storageCreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/create"
	storageDownload "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/download"
	storageLink "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/link"
	storageLsWithPage "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/ls"
	storageMkdir "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/mkdir"
	storageMv "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/mv"
	storageReadAt "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/readAt"
	storageRm "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/rm"
	storageStat "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/stat"
	storageTruncate "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/truncate"
	storageUploadComplete "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/complete"
	storageUploadFile "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/file"
	storageUploadInit "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/init"
	storageUploadSlice "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/upload/slice"
	storageWriteAt "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/writeAt"

	storageQuotaGetAdmin "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/quota/admin/get"
	storageQuotaList "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/quota/admin/list"
	storageQuotaPutAdmin "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/quota/admin/put"
	storageQuotaTotal "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/quota/admin/total"
	storageQuotaGetAPI "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/quota/api/get"

	storageOperationLogAdminList "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/operationlog/admin/list"
	storageOperationLogApiList "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/operationlog/api/list"

	storageSharedDirectoryCreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/shared_directory/api/create"
	storageSharedDirectoryDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/shared_directory/api/delete"
	storageSharedDirectoryList "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/shared_directory/api/list"

	hardwareAdminAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/add"
	hardwareAdminAddUsers "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/add_users"
	hardwareAdminDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/delete"
	hardwareAdminDeleteUsers "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/delete_users"
	hardwareAdminGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/get"
	hardwareAdminList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/list"
	hardwareAdminPatch "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/patch"
	hardwareAdminPut "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/admin/put"
	hardwareGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/api/get"
	hardwareList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/hardware/api/list"
	storageAdminCopy "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/copy"
	storageAdminCreate "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/create"
	storageAdminDirectoryUsageCancel "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/directoryUsage/cancel"
	storageAdminDirectoryUsageStart "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/directoryUsage/start"
	storageAdminDirectoryUsageStatus "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/directoryUsage/status"
	storageAdminLs "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/ls"
	storageAdminMkdir "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/mkdir"
	storageAdminMv "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/mv"
	storageAdminReadAt "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/readAt"
	storageRealpath "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/realpath"
	storageAdminRm "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/rm"
	storageAdminStat "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/stat"
	storageAdminTruncate "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/truncate"
	storageAdminUploadInit "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/upload/init"
	storageAdminWriteAt "github.com/yuansuan/ticp/common/openapi-go/apiv1/storage/admin/writeAt"

	softwareAdminAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/add"
	softwareAdminAddUsers "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/add_users"
	softwareAdminDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/delete"
	softwareAdminDeleteUsers "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/delete_users"
	softwareAdminGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/get"
	softwareAdminList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/list"
	softwareAdminPatch "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/patch"
	softwareAdminPut "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/admin/put"

	softwareGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/api/get"
	softwareList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/software/api/list"

	remoteAppAdminAdd "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/admin/add"
	remoteAppAdminDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/admin/delete"
	remoteAppAdminPatch "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/admin/patch"
	remoteAppAdminPut "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/remoteapp/admin/put"

	sessionAdminClose "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/close"
	sessionAdminList "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/list"
	sessionAdminMount "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/mount"
	sessionAdminPowerOff "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/poweroff"
	sessionAdminPowerOn "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/poweron"
	sessionAdminReboot "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/reboot"
	sessionAdminUmount "github.com/yuansuan/ticp/common/openapi-go/apiv1/cloudapp/session/admin/umount"

	hpcJobSystemCancel "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/job/system/cancel"
	hpcJobCpuUsageSystemGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/job/system/cpu_usage"
	hpcJobSystemDelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/job/system/delete"
	hpcJobSystemGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/job/system/get"
	hpcJobSystemList "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/job/system/list"
	hpcJobSystemPost "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/job/system/post"

	hpcJobSystemOutputFileSyncPause "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/filesyncaction/system/pause"
	hpcJobSystemOutputFileSyncResume "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/filesyncaction/system/resume"

	hpcCommandSystemExecute "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/command/system/post"

	hpcResourceSystemGet "github.com/yuansuan/ticp/common/openapi-go/apiv1/hpc/resource/system/get"

	iamaccountadd "github.com/yuansuan/ticp/common/openapi-go/apiv1/iam/account/add"
	iamaccountexchangecredentials "github.com/yuansuan/ticp/common/openapi-go/apiv1/iam/account/exchange"
)

type API struct {
	hc                     *xhttp.Client
	Job                    *Job
	Storage                *Storage
	License                *License
	Account                *Account
	CloudApp               *CloudApp
	HPC                    *HPC
	StorageQuota           *StorageQuota
	StorageOperationLog    *StorageOperationLog
	StorageSharedDirectory *StorageSharedDirectory
	RDPgo                  *RDPgo
	IAM                    *IAM
}

type IAM struct {
	AccountAdd          iamaccountadd.API
	ExchangeCredentials iamaccountexchangecredentials.API
}

type Account struct {
	AccountList                    accountList.API
	BillList                       accountBillList.API
	Create                         accountCreate.API
	CreditAdd                      accountCreditAdd.API
	AmountRefund                   accountAmountRefund.API
	CreditQuotaModify              accountCreditQuotaModify.API
	FrozenModify                   accountFrozenModify.API
	ByIdGet                        accountIdGet.API
	ByIdReduce                     accountIdReduce.API
	PaymentFreezeUnfreeze          accountPaymentFreezeUnfreeze.API
	PaymentReduce                  accountPaymentReduce.API
	ByYsIDGet                      accountYsidGet.API
	ByYsIDReduce                   accountYsidReduce.API
	CashVoucherAdd                 cashVoucherAdd.API
	CashVoucherGet                 cashVoucherGet.API
	CashVoucherList                cashVoucherList.API
	CashVoucherAvaModify           cashVoucherAvaModify.API
	AccountCashVoucherList         accountCashVoucherList.API
	AccountCashVoucherAdd          accountCashVoucherAdd.API
	AccountCashVoucherGet          accountCashVoucherGet.API
	AccountCashVoucherStatusModify accountCashVoucherStatusModify.API
	ByUserIDGet                    accountUserIdGet.API
	UserBillList                   accountUserBillList.API
	UserResourceList               accountUserResourceBillList.API
}

type Job struct {
	AdminJobGet                  adminJobGet.API
	AdminJobList                 adminJobList.API
	AdminJobListFiltered         adminJobListFiltered.API
	AdminJobTerminate            adminJobTerminate.API
	AdminJobRetransmit           adminJobRetransmit.API
	AdminJobDelete               adminJobDelete.API
	AdminJobCreate               adminJobCreate.API
	AdminJobGetResidual          adminJobGetResidual.API
	AdminJobGetSnapshot          adminJobGetSnapshot.API
	AdminJobListSnapshot         adminJobListSnapshot.API
	AdminJobGetMonitorchart      adminJobGetMonitorChart.API
	AdminJobCpuUsage             adminJobCpuUsage.API
	AdminJobUpdate               adminJobUpdate.API
	JobGet                       jobGet.API
	JobBatchGet                  jobBatchGet.API
	JobList                      jobList.API
	JobTerminate                 jobTerminate.API
	JobDelete                    jobDelete.API
	JobResume                    jobResume.API
	JobCreate                    jobCreate.API
	JobTransmitSuspend           jobTransmitSuspend.API
	JobTransmitResume            jobTransmitResume.API
	JobGetResidual               jobGetResidual.API
	JobGetSnapshot               jobGetSnapshot.API
	JobListSnapshot              jobListSnapshot.API
	JobGetMonitorchart           jobGetMonitorChart.API
	JobPreSchedule               jobPreSchedule.API
	JobCpuUsage                  jobCpuUsage.API
	SystemJobNeedSyncFileList    systemJobNeedSyncFileList.API
	SystemJobSyncFileStateUpdate systemJobSyncFileStateUpdate.API
	ZoneList                     zoneList.API
	AdminGetAPP                  adminGetApp.API
	AdminListAPP                 adminListApp.API
	AdminAddAPP                  adminAddApp.API
	AdminUpdateAPP               adminUpdateApp.API
	AdminDeleteAPP               adminDeleteApp.API
	AdminGetAPPQuota             adminGetAppQuota.API
	AdminAddAPPQuota             adminAddAppQuota.API
	AdminDeleteAPPQuota          adminDeleteAppQuota.API
	ListAPP                      ListApp.API
	AdminGetAppAllow             adminGetAppAllow.API
	AdminAddAppAllow             adminAddAppAllow.API
	AdminDeleteAppAllow          adminDeleteAppAllow.API
}

type Storage struct {
	LsWithPage           storageLsWithPage.API
	Mkdir                storageMkdir.API
	Mv                   storageMv.API
	Rm                   storageRm.API
	Stat                 storageStat.API
	Download             storageDownload.API
	BatchDownload        storageBatchDownload.API
	UploadInit           storageUploadInit.API
	UploadSlice          storageUploadSlice.API
	UploadComplete       storageUploadComplete.API
	UploadFile           storageUploadFile.API
	Copy                 storageCopy.API
	CopyRange            storageCopyRange.API
	WriteAt              storageWriteAt.API
	ReadAt               storageReadAt.API
	Link                 storageLink.API
	Create               storageCreate.API
	Truncate             storageTruncate.API
	CompressStart        storageCompressStart.API
	CompressStatus       storageCompressStatus.API
	CompressCancel       storageCompressCancel.API
	DirectoryUsageStart  storageDirectoryUsageStart.API
	DirectoryUsageStatus storageDirectoryUsageStatus.API
	DirectoryUsageCancel storageDirectoryUsageCancel.API

	Realpath                  storageRealpath.API
	AdminLsWithPage           storageAdminLs.API
	AdminCreate               storageAdminCreate.API
	AdminTruncate             storageAdminTruncate.API
	AdminRm                   storageAdminRm.API
	AdminMkdir                storageAdminMkdir.API
	AdminMv                   storageAdminMv.API
	AdminCopy                 storageAdminCopy.API
	AdminDownload             storageAdminDownload.API
	AdminBatchDownload        storageAdminBatchDownload.API
	AdminUploadInit           storageAdminUploadInit.API
	AdminUploadSlice          storageAdminUploadSlice.API
	AdminUploadComplete       storageAdminUploadComplete.API
	AdminReadAt               storageAdminReadAt.API
	AdminStat                 storageAdminStat.API
	AdminWriteAt              storageAdminWriteAt.API
	AdminDirectoryUsageStart  storageAdminDirectoryUsageStart.API
	AdminDirectoryUsageStatus storageAdminDirectoryUsageStatus.API
	AdminDirectoryUsageCancel storageAdminDirectoryUsageCancel.API
}

type StorageQuota struct {
	GetQuotaAPI               storageQuotaGetAPI.API
	GetQuotaAdmin             storageQuotaGetAdmin.API
	GetStorageQuotaTotalAdmin storageQuotaTotal.API
	ListStorageQuotaAdmin     storageQuotaList.API
	PutStorageQuotaAdmin      storageQuotaPutAdmin.API
}

type StorageOperationLog struct {
	ListOperationLogAdmin storageOperationLogAdminList.API
	ListOperationLogAPI   storageOperationLogApiList.API
}

type StorageSharedDirectory struct {
	List   storageSharedDirectoryList.API
	Create storageSharedDirectoryCreate.API
	Delete storageSharedDirectoryDelete.API
}

type License struct {
	AddLicenseManager    licenseManagerAdd.API
	PutLicenseManager    licenseManagerPut.API
	DeleteLicenseManager licenseManagerDelete.API
	GetLicenseManager    licenseManagerGet.API
	ListLicenseManager   licenseManagerList.API

	AddLicenseInfo    licenseInfoAdd.API
	DeleteLicenseInfo licenseInfoDelete.API
	PutLicenseInfo    licenseInfoPut.API

	ListModuleConfig      moduleConfigList.API
	AddModuleConfig       moduleConfigAdd.API
	BatchAddModuleConfigs moduleConfigBatchAdd.API
	DeleteModuleConfig    moduleConfigDelete.API
	PutModuleConfig       moduleConfigPut.API
}

type HardWareAdmin struct {
	Add    hardwareAdminAdd.API
	Get    hardwareAdminGet.API
	List   hardwareAdminList.API
	Delete hardwareAdminDelete.API
	Put    hardwareAdminPut.API
	Patch  hardwareAdminPatch.API
}

type SoftwareAdmin struct {
	Add    softwareAdminAdd.API
	Get    softwareAdminGet.API
	List   softwareAdminList.API
	Delete softwareAdminDelete.API
	Put    softwareAdminPut.API
	Patch  softwareAdminPatch.API
}

func NewAPI(hc *xhttp.Client) (*API, error) {
	return &API{
		hc: hc,
		Account: &Account{
			AccountList:                    accountList.New(hc),
			BillList:                       accountBillList.New(hc),
			Create:                         accountCreate.New(hc),
			CreditAdd:                      accountCreditAdd.New(hc),
			AmountRefund:                   accountAmountRefund.New(hc),
			CreditQuotaModify:              accountCreditQuotaModify.New(hc),
			FrozenModify:                   accountFrozenModify.New(hc),
			ByIdGet:                        accountIdGet.New(hc),
			ByIdReduce:                     accountIdReduce.New(hc),
			PaymentFreezeUnfreeze:          accountPaymentFreezeUnfreeze.New(hc),
			PaymentReduce:                  accountPaymentReduce.New(hc),
			ByYsIDGet:                      accountYsidGet.New(hc),
			ByYsIDReduce:                   accountYsidReduce.New(hc),
			AccountCashVoucherAdd:          accountCashVoucherAdd.New(hc),
			AccountCashVoucherGet:          accountCashVoucherGet.New(hc),
			AccountCashVoucherList:         accountCashVoucherList.New(hc),
			AccountCashVoucherStatusModify: accountCashVoucherStatusModify.New(hc),
			CashVoucherAdd:                 cashVoucherAdd.New(hc),
			CashVoucherAvaModify:           cashVoucherAvaModify.New(hc),
			CashVoucherGet:                 cashVoucherGet.New(hc),
			CashVoucherList:                cashVoucherList.New(hc),
			ByUserIDGet:                    accountUserIdGet.New(hc),
			UserBillList:                   accountUserBillList.New(hc),
			UserResourceList:               accountUserResourceBillList.New(hc),
		},

		Job: &Job{
			// job
			JobGet:                       jobGet.New(hc),
			JobBatchGet:                  jobBatchGet.New(hc),
			JobList:                      jobList.New(hc),
			JobTerminate:                 jobTerminate.New(hc),
			JobDelete:                    jobDelete.New(hc),
			JobResume:                    jobResume.New(hc),
			JobCreate:                    jobCreate.New(hc),
			JobTransmitSuspend:           jobTransmitSuspend.New(hc),
			JobTransmitResume:            jobTransmitResume.New(hc),
			JobGetResidual:               jobGetResidual.New(hc),
			JobGetSnapshot:               jobGetSnapshot.New(hc),
			JobListSnapshot:              jobListSnapshot.New(hc),
			JobGetMonitorchart:           jobGetMonitorChart.New(hc),
			JobPreSchedule:               jobPreSchedule.New(hc),
			JobCpuUsage:                  jobCpuUsage.New(hc),
			ZoneList:                     zoneList.New(hc),
			AdminJobGet:                  adminJobGet.New(hc),
			AdminJobList:                 adminJobList.New(hc),
			AdminJobListFiltered:         adminJobListFiltered.New(hc),
			AdminJobTerminate:            adminJobTerminate.New(hc),
			AdminJobUpdate:               adminJobUpdate.New(hc),
			AdminJobRetransmit:           adminJobRetransmit.New(hc),
			AdminJobDelete:               adminJobDelete.New(hc),
			AdminJobCreate:               adminJobCreate.New(hc),
			AdminJobGetResidual:          adminJobGetResidual.New(hc),
			AdminJobGetSnapshot:          adminJobGetSnapshot.New(hc),
			AdminJobListSnapshot:         adminJobListSnapshot.New(hc),
			AdminJobGetMonitorchart:      adminJobGetMonitorChart.New(hc),
			AdminJobCpuUsage:             adminJobCpuUsage.New(hc),
			SystemJobNeedSyncFileList:    systemJobNeedSyncFileList.New(hc),
			SystemJobSyncFileStateUpdate: systemJobSyncFileStateUpdate.New(hc),

			// app
			AdminGetAPP:         adminGetApp.New(hc),
			AdminListAPP:        adminListApp.New(hc),
			AdminAddAPP:         adminAddApp.New(hc),
			AdminUpdateAPP:      adminUpdateApp.New(hc),
			AdminDeleteAPP:      adminDeleteApp.New(hc),
			AdminGetAPPQuota:    adminGetAppQuota.New(hc),
			AdminAddAPPQuota:    adminAddAppQuota.New(hc),
			AdminDeleteAPPQuota: adminDeleteAppQuota.New(hc),
			ListAPP:             ListApp.New(hc),
			AdminGetAppAllow:    adminGetAppAllow.New(hc),
			AdminAddAppAllow:    adminAddAppAllow.New(hc),
			AdminDeleteAppAllow: adminDeleteAppAllow.New(hc),
		},

		Storage: &Storage{
			LsWithPage:                storageLsWithPage.New(hc),
			Mkdir:                     storageMkdir.New(hc),
			Mv:                        storageMv.New(hc),
			Rm:                        storageRm.New(hc),
			Stat:                      storageStat.New(hc),
			Download:                  storageDownload.New(hc),
			BatchDownload:             storageBatchDownload.New(hc),
			UploadInit:                storageUploadInit.New(hc),
			UploadSlice:               storageUploadSlice.New(hc),
			UploadComplete:            storageUploadComplete.New(hc),
			UploadFile:                storageUploadFile.New(hc),
			Copy:                      storageCopy.New(hc),
			CopyRange:                 storageCopyRange.New(hc),
			WriteAt:                   storageWriteAt.New(hc),
			ReadAt:                    storageReadAt.New(hc),
			Link:                      storageLink.New(hc),
			Create:                    storageCreate.New(hc),
			Truncate:                  storageTruncate.New(hc),
			CompressStart:             storageCompressStart.New(hc),
			CompressStatus:            storageCompressStatus.New(hc),
			CompressCancel:            storageCompressCancel.New(hc),
			DirectoryUsageStart:       storageDirectoryUsageStart.New(hc),
			DirectoryUsageStatus:      storageDirectoryUsageStatus.New(hc),
			DirectoryUsageCancel:      storageDirectoryUsageCancel.New(hc),
			Realpath:                  storageRealpath.New(hc),
			AdminLsWithPage:           storageAdminLs.New(hc),
			AdminCreate:               storageAdminCreate.New(hc),
			AdminTruncate:             storageAdminTruncate.New(hc),
			AdminRm:                   storageAdminRm.New(hc),
			AdminMkdir:                storageAdminMkdir.New(hc),
			AdminMv:                   storageAdminMv.New(hc),
			AdminCopy:                 storageAdminCopy.New(hc),
			AdminDownload:             storageAdminDownload.New(hc),
			AdminBatchDownload:        storageAdminBatchDownload.New(hc),
			AdminUploadInit:           storageAdminUploadInit.New(hc),
			AdminUploadSlice:          storageAdminUploadSlice.New(hc),
			AdminUploadComplete:       storageAdminUploadComplete.New(hc),
			AdminReadAt:               storageAdminReadAt.New(hc),
			AdminStat:                 storageAdminStat.New(hc),
			AdminWriteAt:              storageAdminWriteAt.New(hc),
			AdminDirectoryUsageStart:  storageAdminDirectoryUsageStart.New(hc),
			AdminDirectoryUsageStatus: storageAdminDirectoryUsageStatus.New(hc),
			AdminDirectoryUsageCancel: storageAdminDirectoryUsageCancel.New(hc),
		},

		StorageQuota: &StorageQuota{
			GetQuotaAPI:               storageQuotaGetAPI.New(hc),
			GetQuotaAdmin:             storageQuotaGetAdmin.New(hc),
			GetStorageQuotaTotalAdmin: storageQuotaTotal.New(hc),
			ListStorageQuotaAdmin:     storageQuotaList.New(hc),
			PutStorageQuotaAdmin:      storageQuotaPutAdmin.New(hc),
		},

		StorageOperationLog: &StorageOperationLog{
			ListOperationLogAdmin: storageOperationLogAdminList.New(hc),
			ListOperationLogAPI:   storageOperationLogApiList.New(hc),
		},

		StorageSharedDirectory: &StorageSharedDirectory{
			List:   storageSharedDirectoryList.New(hc),
			Create: storageSharedDirectoryCreate.New(hc),
			Delete: storageSharedDirectoryDelete.New(hc),
		},

		License: &License{
			AddLicenseManager:     licenseManagerAdd.New(hc),
			DeleteLicenseManager:  licenseManagerDelete.New(hc),
			PutLicenseManager:     licenseManagerPut.New(hc),
			GetLicenseManager:     licenseManagerGet.New(hc),
			ListLicenseManager:    licenseManagerList.New(hc),
			AddLicenseInfo:        licenseInfoAdd.New(hc),
			PutLicenseInfo:        licenseInfoPut.New(hc),
			DeleteLicenseInfo:     licenseInfoDelete.New(hc),
			ListModuleConfig:      moduleConfigList.New(hc),
			AddModuleConfig:       moduleConfigAdd.New(hc),
			BatchAddModuleConfigs: moduleConfigBatchAdd.New(hc),
			PutModuleConfig:       moduleConfigPut.New(hc),
			DeleteModuleConfig:    moduleConfigDelete.New(hc),
		},

		CloudApp: newCloudApp(hc),

		IAM: &IAM{
			AccountAdd:          iamaccountadd.New(hc),
			ExchangeCredentials: iamaccountexchangecredentials.New(hc),
		},

		HPC: newHPC(hc),

		RDPgo: newRDPgo(hc),
	}, nil
}

type CloudApp struct {
	Hardware  hardware
	Software  software
	Session   session
	RemoteApp remoteApp
}

func newCloudApp(hc *xhttp.Client) *CloudApp {
	return &CloudApp{
		Hardware: hardware{
			Admin: hardwareAdmin{
				Add:         hardwareAdminAdd.New(hc),
				Get:         hardwareAdminGet.New(hc),
				List:        hardwareAdminList.New(hc),
				Delete:      hardwareAdminDelete.New(hc),
				Put:         hardwareAdminPut.New(hc),
				Patch:       hardwareAdminPatch.New(hc),
				AddUsers:    hardwareAdminAddUsers.New(hc),
				DeleteUsers: hardwareAdminDeleteUsers.New(hc),
			},
			User: hardwareUser{
				Get:  hardwareGet.New(hc),
				List: hardwareList.New(hc),
			},
		},
		Software: software{
			Admin: softwareAdmin{
				Add:         softwareAdminAdd.New(hc),
				Get:         softwareAdminGet.New(hc),
				List:        softwareAdminList.New(hc),
				Delete:      softwareAdminDelete.New(hc),
				Put:         softwareAdminPut.New(hc),
				Patch:       softwareAdminPatch.New(hc),
				AddUsers:    softwareAdminAddUsers.New(hc),
				DeleteUsers: softwareAdminDeleteUsers.New(hc),
			},
			User: softwareUser{
				Get:  softwareGet.New(hc),
				List: softwareList.New(hc),
			},
		},
		Session: session{
			Admin: sessionAdmin{
				Close:      sessionAdminClose.New(hc),
				List:       sessionAdminList.New(hc),
				PowerOff:   sessionAdminPowerOff.New(hc),
				PowerOn:    sessionAdminPowerOn.New(hc),
				Reboot:     sessionAdminReboot.New(hc),
				Restore:    sessionAdminRestore.New(hc),
				ExecScript: sessionAdminExecScript.New(hc),
				Mount:      sessionAdminMount.New(hc),
				Umount:     sessionAdminUmount.New(hc),
			},
			User: sessionUser{
				Start:      sessionStart.New(hc),
				Close:      sessionClose.New(hc),
				Get:        sessionGet.New(hc),
				List:       sessionList.New(hc),
				PowerOff:   sessionPowerOff.New(hc),
				PowerOn:    sessionPowerOn.New(hc),
				Ready:      sessionReady.New(hc),
				Reboot:     sessionReboot.New(hc),
				Delete:     sessionDelete.New(hc),
				Restore:    sessionRestore.New(hc),
				ExecScript: sessionExecScript.New(hc),
				Mount:      sessionMount.New(hc),
				Umount:     sessionUmount.New(hc),
			},
		},
		RemoteApp: remoteApp{
			Admin: remoteAppAdmin{
				Put:    remoteAppAdminPut.New(hc),
				Patch:  remoteAppAdminPatch.New(hc),
				Delete: remoteAppAdminDelete.New(hc),
				Add:    remoteAppAdminAdd.New(hc),
			},
			User: remoteAppUser{
				Get: remoteAppGet.New(hc),
			},
		},
	}
}

type hardware struct {
	Admin hardwareAdmin
	User  hardwareUser
}

type hardwareAdmin struct {
	Add         hardwareAdminAdd.API
	Get         hardwareAdminGet.API
	List        hardwareAdminList.API
	Delete      hardwareAdminDelete.API
	Put         hardwareAdminPut.API
	Patch       hardwareAdminPatch.API
	AddUsers    hardwareAdminAddUsers.API
	DeleteUsers hardwareAdminDeleteUsers.API
}

type hardwareUser struct {
	Get  hardwareGet.API
	List hardwareList.API
}

type software struct {
	Admin softwareAdmin
	User  softwareUser
}

type softwareAdmin struct {
	Add         softwareAdminAdd.API
	Get         softwareAdminGet.API
	List        softwareAdminList.API
	Delete      softwareAdminDelete.API
	Put         softwareAdminPut.API
	Patch       softwareAdminPatch.API
	AddUsers    softwareAdminAddUsers.API
	DeleteUsers softwareAdminDeleteUsers.API
}

type softwareUser struct {
	Get  softwareGet.API
	List softwareList.API
}

type session struct {
	Admin sessionAdmin
	User  sessionUser
}

type sessionAdmin struct {
	Close      sessionAdminClose.API
	List       sessionAdminList.API
	PowerOff   sessionAdminPowerOff.API
	PowerOn    sessionAdminPowerOn.API
	Reboot     sessionAdminReboot.API
	Restore    sessionAdminRestore.API
	ExecScript sessionAdminExecScript.API
	Mount      sessionAdminMount.API
	Umount     sessionAdminUmount.API
}

type sessionUser struct {
	Start      sessionStart.API
	Close      sessionClose.API
	Get        sessionGet.API
	List       sessionList.API
	PowerOff   sessionPowerOff.API
	PowerOn    sessionPowerOn.API
	Ready      sessionReady.API
	Reboot     sessionReboot.API
	Delete     sessionDelete.API
	Restore    sessionRestore.API
	ExecScript sessionExecScript.API
	Mount      sessionMount.API
	Umount     sessionUmount.API
}

type remoteApp struct {
	Admin remoteAppAdmin
	User  remoteAppUser
}

type remoteAppAdmin struct {
	Put    remoteAppAdminPut.API
	Patch  remoteAppAdminPatch.API
	Delete remoteAppAdminDelete.API
	Add    remoteAppAdminAdd.API
}

type remoteAppUser struct {
	Get remoteAppGet.API
}

type HPC struct {
	Job      hpcJob
	Command  hpcCommand
	Resource hpcResource
}

func newHPC(hc *xhttp.Client) *HPC {
	return &HPC{
		Job: hpcJob{
			System: hpcJobSystem{
				Cancel:               hpcJobSystemCancel.New(hc),
				Delete:               hpcJobSystemDelete.New(hc),
				Get:                  hpcJobSystemGet.New(hc),
				List:                 hpcJobSystemList.New(hc),
				Post:                 hpcJobSystemPost.New(hc),
				PauseOutputFileSync:  hpcJobSystemOutputFileSyncPause.New(hc),
				ResumeOutputFileSync: hpcJobSystemOutputFileSyncResume.New(hc),
				GetCpuUsage:          hpcJobCpuUsageSystemGet.New(hc),
			},
		},
		Command: hpcCommand{
			System: hpcCommandSystem{
				Execute: hpcCommandSystemExecute.New(hc),
			},
		},
		Resource: hpcResource{
			System: hpcResourceSystem{
				Get: hpcResourceSystemGet.New(hc),
			},
		},
	}
}

type hpcJob struct {
	System hpcJobSystem
}

type hpcJobSystem struct {
	Cancel               hpcJobSystemCancel.API
	Delete               hpcJobSystemDelete.API
	Get                  hpcJobSystemGet.API
	List                 hpcJobSystemList.API
	Post                 hpcJobSystemPost.API
	PauseOutputFileSync  hpcJobSystemOutputFileSyncPause.API
	ResumeOutputFileSync hpcJobSystemOutputFileSyncResume.API
	GetCpuUsage          hpcJobCpuUsageSystemGet.API
}

type hpcCommand struct {
	System hpcCommandSystem
}

type hpcCommandSystem struct {
	Execute hpcCommandSystemExecute.API
}

type hpcResource struct {
	System hpcResourceSystem
}

type hpcResourceSystem struct {
	Get hpcResourceSystemGet.API
}

type RDPgo struct {
	Clean      rdpgoInternalClean.API
	ExecScript rdpgoInternalExecScript.API
}

func newRDPgo(hc *xhttp.Client) *RDPgo {
	return &RDPgo{
		Clean:      rdpgoInternalClean.New(hc),
		ExecScript: rdpgoInternalExecScript.New(hc),
	}
}
