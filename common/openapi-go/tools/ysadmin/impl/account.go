package impl

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/billlist"
	"github.com/yuansuan/ticp/common/project-root-api/account_bill/v1/create"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

const (
	// Yuan 单位换算
	Yuan                      = 100000
	accountExampleFile        = "account_example.json"
	accountExampleFileContent = `{
	"AccountName":"某某企业",
	"UserID":"531u9HG1Gpb",
	"AccountType":1
}`
)

type AccountOptions struct {
	BaseOptions
	ProductName string
	Money       int64
	TradeId     string
	Comment     string
	CSVFile     string
	Merge       bool
	OnlyConsume bool
}

func init() {
	RegisterCmd(NewAccountCommand())
}

// NewAccountCommand 资金账号管理
func NewAccountCommand() *cobra.Command {
	o := AccountOptions{}
	cmd := &cobra.Command{
		Use:   "account",
		Short: "资金账号管理",
		Long:  "资金账号管理, 不同于远算账号, 用于管理用户的账户余额, 账单等",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAccountCreateCommand(o),
		newAccountListBillCommand(o),
		newAccountGetCommand(o),
		newAccountGetByYsidCommand(o),
		newAccountAddBalanceCommand(o),
		newAccountReduceCommand(o),
		newAccountRefundCommand(o),
		newAccountExampleCommand(),
	)

	return cmd
}

func newAccountCreateCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建资金账号",
		Long:  "创建资金账号, 指定一个.json文件作为参数文件, 可使用ysadmin account example生成参考参数文件",
		Args:  cobra.ExactArgs(0),
		Example: ` - 创建资金账号, 参数文件为account.json
  - ysadmin account create -F account.json`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "JSON 文件路径")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(create.Request)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		res, err := GetYsClient().Account.Create(
			GetYsClient().Account.Create.AccountName(req.AccountName),
			GetYsClient().Account.Create.UserID(req.UserID),
			GetYsClient().Account.Create.AccountType(req.AccountType),
		)
		PrintResp(res, err, "Create Account")
		return nil
	}

	return cmd
}

func newAccountListBillCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listbill",
		Short: "获取账单列表",
		Long:  "获取账单列表, 可输出到CSV文件, 可聚合消费账单",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取账单列表, 所有的条目
  - ysadmin account listbill --all
- 获取账单列表, 分页获取, 每页10条, 第1页
  - ysadmin account listbill -L 10 -O 1
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定开始时间和结束时间
  - ysadmin account listbill -L 10 -O 1 --start_time "2021-01-01 00:00:00" --end_time "2021-01-31 23:59:59"
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定产品名称, 并输出到CSV文件
  - ysadmin account listbill -L 10 -O 1 -P 3D云应用 --csv_file /tmp/bill.csv
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定账户ID, 并输出到CSV文件, 并聚合消费账单
  - ysadmin account listbill -I 5314rXEJwrf -L 10 -O 1 --csv_file /tmp/bill.csv --merge
- 获取账单列表, 分页获取, 每页10条, 第1页, 指定账户ID, 并输出到CSV文件, 并聚合消费账单, 只展示消费类型的账单
  - ysadmin account listbill -I 5314rXEJwrf -L 10 -O 1 --csv_file /tmp/bill.csv --merge --only_consume`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "账户ID")
	cmd.Flags().StringVarP(&o.ProductName, "product_name", "P", "", "产品名称")
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")
	cmd.Flags().BoolVarP(&o.All, "all", "", false, "所有的条目")
	cmd.Flags().StringVarP(&o.StartTime, "start_time", "", "", "开始时间")
	cmd.Flags().StringVarP(&o.EndTime, "end_time", "", "", "结束时间")
	cmd.Flags().StringVarP(&o.CSVFile, "csv_file", "", "", "输出CSV格式的数据到指定文件")
	cmd.Flags().BoolVarP(&o.Merge, "merge", "", false, "是否需要聚合消费账单，按资源ID展示")
	cmd.Flags().BoolVarP(&o.OnlyConsume, "only_consume", "", false, "只展示消费类型的账单")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if o.Offset == 0 {
			// 账单列表分页从1开始，offset是分页号，不是真实偏移量
			o.Offset++
		}
		var res *billlist.Response
		for {
			tmpRes, err := GetYsClient().Account.BillList(
				GetYsClient().Account.BillList.AccountID(o.Id),
				GetYsClient().Account.BillList.PageIndex(o.Offset),
				GetYsClient().Account.BillList.PageSize(o.Limit),
				GetYsClient().Account.BillList.StartTime(o.StartTime),
				GetYsClient().Account.BillList.EndTime(o.EndTime),
				GetYsClient().Account.BillList.ProductName(o.ProductName),
			)
			if err != nil {
				PrintResp(tmpRes, err, "List Bill")
				os.Exit(0)
			}
			if res == nil {
				res = tmpRes
			} else {
				res.Data.AccountBills = append(res.Data.AccountBills, tmpRes.Data.AccountBills...)
				res.Data.Total = int64(len(res.Data.AccountBills))
			}
			if !o.All || len(tmpRes.Data.AccountBills) < int(o.Limit) {
				break
			}
			o.Offset++
		}

		if o.Merge {
			res.Data.AccountBills = o.MergeBillList(res.Data.AccountBills)
			res.Data.Total = int64(len(res.Data.AccountBills))
		}
		if len(o.CSVFile) > 0 {
			o.ToCSVFile(res.Data.AccountBills)
		} else {
			PrintResp(res, nil, "List Bill")
		}
		return nil
	}

	return cmd
}

func newAccountGetCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "获取账户信息",
		Long:  "获取账户信息",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取账户信息, 指定账户ID
  - ysadmin account get -I 5314rXEJwrf`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "账户ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Account.ByIdGet(
			GetYsClient().Account.ByIdGet.AccountID(o.Id),
		)
		PrintResp(res, err, "Get Acccount")
		return nil
	}

	return cmd
}

func newAccountGetByYsidCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getbyysid",
		Short: "根据远算账号ID获取账户信息",
		Long:  "根据远算账号ID获取账户信息",
		Args:  cobra.ExactArgs(0),
		Example: `- 根据远算账号ID获取账户信息, 指定远算账号ID
  - ysadmin account getbyysid -I 5314rXEJwrf`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "远算账号ID")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Account.ByYsIDGet(
			GetYsClient().Account.ByYsIDGet.UserID(o.Id),
		)
		PrintResp(res, err, "Get Account")
		return nil
	}

	return cmd
}

func newAccountAddBalanceCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addbalance",
		Short: "充值",
		Long:  "充值, 单位为0.00001元",
		Args:  cobra.ExactArgs(0),
		Example: `- 充值, 指定账户ID, 充值金额, 交易ID, 备注
  - ysadmin account addbalance -I 5314rXEJwrf -M 10000000 -T 531jv4i44nJ -C "充值100元"`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "账户ID")
	cmd.MarkFlagRequired("id")
	cmd.Flags().Int64VarP(&o.Money, "money", "M", 0, "充值金额, 单位为0.00001元")
	cmd.MarkFlagRequired("money")
	cmd.Flags().StringVarP(&o.TradeId, "trade_id", "T", "", "交易ID")
	cmd.MarkFlagRequired("trade_id")
	cmd.Flags().StringVarP(&o.Comment, "comment", "C", "", "备注")
	cmd.MarkFlagRequired("comment")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Account.CreditAdd(
			GetYsClient().Account.CreditAdd.AccountID(o.Id),
			GetYsClient().Account.CreditAdd.DeltaNormalBalance(o.Money),
			GetYsClient().Account.CreditAdd.TradeID(o.TradeId),
			GetYsClient().Account.CreditAdd.Comment(o.Comment),
		)
		PrintResp(res, err, "Add Balance")
		return nil
	}

	return cmd
}

func newAccountReduceCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reduce",
		Short: "扣费",
		Long:  "扣费, 单位为0.00001元",
		Args:  cobra.ExactArgs(0),
		Example: `- 扣费, 指定账户ID, 扣费金额, 交易ID, 备注, 产品名称
  - ysadmin account reduce -I 5314rXEJwrf -M 10000000 -T 531jv4i44nJ -C "扣费100元" -P 3D云应用`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "账户ID")
	cmd.MarkFlagRequired("id")
	cmd.Flags().Int64VarP(&o.Money, "money", "M", 0, "扣费金额, 单位为0.00001元")
	cmd.MarkFlagRequired("money")
	cmd.Flags().StringVarP(&o.TradeId, "trade_id", "T", "", "交易ID")
	cmd.MarkFlagRequired("trade_id")
	cmd.Flags().StringVarP(&o.Comment, "comment", "C", "", "备注")
	cmd.MarkFlagRequired("comment")
	cmd.Flags().StringVarP(&o.ProductName, "product_name", "P", "", "产品名称")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Account.ByIdReduce(
			GetYsClient().Account.ByIdReduce.AccountID(o.Id),
			GetYsClient().Account.ByIdReduce.TradeID(o.TradeId),
			GetYsClient().Account.ByIdReduce.Amount(o.Money),
			GetYsClient().Account.ByIdReduce.Comment(o.Comment),
			GetYsClient().Account.ByIdReduce.ProductName(o.ProductName),
		)
		PrintResp(res, err, "Reduce Account")
		return nil
	}

	return cmd
}

func newAccountRefundCommand(o AccountOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund",
		Short: "退款",
		Long:  "退款, 单位为0.00001元",
		Args:  cobra.ExactArgs(0),
		Example: `- 退款, 指定账户ID, 退款金额, 交易ID, 备注
  - ysadmin account refund -I 5314rXEJwrf -M 10000000 -T 531jv4i44nJ -C "退款100元"`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "账户ID")
	cmd.MarkFlagRequired("id")
	cmd.Flags().Int64VarP(&o.Money, "money", "M", 0, "退款金额, 单位为0.00001元")
	cmd.MarkFlagRequired("money")
	cmd.Flags().StringVarP(&o.TradeId, "trade_id", "T", "", "交易ID")
	cmd.MarkFlagRequired("trade_id")
	cmd.Flags().StringVarP(&o.Comment, "comment", "C", "", "备注")
	cmd.MarkFlagRequired("comment")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		res, err := GetYsClient().Account.AmountRefund(
			GetYsClient().Account.AmountRefund.AccountID(o.Id),
			GetYsClient().Account.AmountRefund.Amount(o.Money),
			GetYsClient().Account.AmountRefund.Comment(o.Comment),
			GetYsClient().Account.AmountRefund.RefundId(o.TradeId),
		)
		PrintResp(res, err, "Refund Amount")
		return nil
	}

	return cmd
}

func newAccountExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "创建资金账号示例文件",
		Long:  "创建资金账号示例文件",
		Example: `- 创建资金账号示例文件
  - ysadmin account example`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if checkFileExist(accountExampleFile) {
			fmt.Println("示例文件:", accountExampleFile)
			return nil
		}

		f, err := os.Create(accountExampleFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = f.WriteString(accountExampleFileContent); err != nil {
			return err
		}
		fmt.Println("示例文件:", accountExampleFile)
		return nil
	}

	return cmd
}

func (o *AccountOptions) MergeBillList(data []*v20230530.BillListData) []*v20230530.BillListData {
	saved := map[string]*v20230530.BillListData{}
	var keyList []string
	for _, item := range data {
		if o.OnlyConsume && item.TradeType != 1 {
			// 1为消费类型
			continue
		}
		keyId := item.ResourceID
		if len(keyId) == 0 {
			// 充值账单无ResourceID
			keyId = item.ID
		}
		if _, ok := saved[keyId]; !ok {
			saved[keyId] = item
			keyList = append(keyList, keyId)
			continue
		}
		saved[keyId].Amount += item.Amount
		saved[keyId].Quantity += item.Quantity
		if item.StartTime < saved[keyId].StartTime {
			saved[keyId].StartTime = item.StartTime
		}
		if item.EndTime > saved[keyId].EndTime {
			saved[keyId].EndTime = item.EndTime
		}
	}
	var res []*v20230530.BillListData
	for _, k := range keyList {
		res = append(res, saved[k])
	}
	return res
}

func (o *AccountOptions) ToCSVFile(data []*v20230530.BillListData) {
	if len(o.CSVFile) == 0 {
		fmt.Println("CSV文件路径为空")
		os.Exit(-1)
	}
	csvFile, err := os.Create(o.CSVFile)
	if err != nil {
		fmt.Println("创建CSV文件失败：", err)
		os.Exit(-1)
	}
	defer csvFile.Close()
	// 创建CSV写入器
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// 写入CSV标题行
	writer.Write([]string{"商品类型", "商品名称", "资源ID(作业ID/3D云应用会话ID)", "订单号", "金额(元)",
		"单价(元)", "数量", "数量单位", "账单计费开始时间", "账单计费结束时间", "账户余额(元)", "账单ID", "账单类型", "备注"})

	// 总的消费金额
	totalConsume := 0
	// 总的充值金额
	totalSaving := 0

	// 总的数量（核时/3D云应用时间）
	totalQuantity := 0.0

	// 遍历数据并写入CSV
	for _, item := range data {
		tradeType := "未知"
		if item.TradeType == 1 {
			tradeType = "消费"
			totalConsume += int(item.Amount)
			totalQuantity += item.Quantity
		} else if item.TradeType == 2 {
			tradeType = "充值"
			totalSaving += int(item.Amount)
		}
		billID := item.ID
		if o.Merge {
			billID = ""
		}
		row := []string{item.ProductName, item.MerchandiseName, item.ResourceID, item.TradeID,
			fmt.Sprintf("%f", float64(item.Amount)/100000.), fmt.Sprintf("%f", float64(item.UnitPrice)/100000.),
			fmt.Sprintf("%f", item.Quantity), item.QuantityUnit, item.StartTime, item.EndTime,
			fmt.Sprintf("%f", float64(item.AccountBalance)/100000.), billID, tradeType, item.Comment}
		writer.Write(row)
	}
	totalStr := fmt.Sprintf("总充值(元): %f, 总消费(元): %f ", float64(totalSaving)/100000., float64(totalConsume)/100000.)
	row := []string{"汇总", "", "", "", totalStr, "", fmt.Sprintf("%f", totalQuantity), "", "", "", "", "", "", ""}
	writer.Write(row)
	// 刷新缓冲区，确保数据被写入文件
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("写入CSV失败:", err)
	}
	fmt.Println(totalStr)
	fmt.Printf("总数量: %f \n", totalQuantity)
}
