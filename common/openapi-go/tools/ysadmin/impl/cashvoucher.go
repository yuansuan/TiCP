package impl

import (
	"fmt"
	"github.com/spf13/cobra"
)

type CashVoucherOptions struct {
	AccountOptions
	CID            string // 优惠券ID，用于get和list
	Name           string // 优惠券名称，用于add和list
	Amount         int64  // 优惠券金额，用于add
	ExpiredType    int64  // 过期类型，用于add
	RelExpiredTime int64  // 相对过期时间，用于add
	AbsExpiredTime string // 绝对过期时间，用于add
	Comment        string // 优惠券备注，用于add
	Offset         int64  // 分页偏移量，用于list
	Limit          int64  // 分页大小，用于list
}

func init() {
	RegisterCmd(NewCashVoucherCommand())
}

func NewCashVoucherCommand() *cobra.Command {
	o := CashVoucherOptions{}
	cmd := &cobra.Command{
		Use:   "cashvoucher",
		Short: "优惠券管理",
		Long:  "优惠券管理，用于创建相应的优惠券",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newCashVoucherAddCommand(o),
		newCashVoucherListCommand(o),
		newCashVoucherGetCommand(o),
	)

	return cmd
}

func newCashVoucherAddCommand(o CashVoucherOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "创建优惠券",
		Long:  "创建优惠券，指定优惠券的名称、金额、过期类型及过期参数。",
		Example: `-添加优惠券
  - ysadmin cashvoucher add -N 测试优惠券 -M 10000000 -T 2 -A "2023-12-31 23:59:59" -C 备注`,
	}
	cmd.Flags().StringVarP(&o.Name, "name", "N", "", "优惠券名称")
	cmd.MarkFlagRequired("name")
	cmd.Flags().Int64VarP(&o.Amount, "amount", "M", 0, "优惠金额，单位为0.00001元")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().Int64VarP(&o.ExpiredType, "expired_type", "T", 0, "优惠券过期类型：1(相对过期时间), 2(绝对过期时间)")
	cmd.MarkFlagRequired("expired_type")
	cmd.Flags().StringVarP(&o.AbsExpiredTime, "abs_expired_time", "A", "", "绝对过期时间，例如：2023-12-31 23:59:59")
	cmd.Flags().Int64VarP(&o.RelExpiredTime, "rel_expired_time", "R", 0, "相对过期时间（单位:秒)")
	cmd.Flags().StringVarP(&o.Comment, "comment", "C", "", "备注")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if o.ExpiredType == 2 { // 改为 2 表示相对过期时间
			if !cmd.Flags().Changed("rel_expired_time") || o.RelExpiredTime == 0 {
				return fmt.Errorf("当过期类型为 2（相对过期时间）时，必须设置 --rel_expired_time")
			}
			if cmd.Flags().Changed("abs_expired_time") {
				return fmt.Errorf("当过期类型为 2（相对过期时间）时，不应设置 --abs_expired_time")
			}
			o.AbsExpiredTime = ""
		} else if o.ExpiredType == 1 { // 改为 1 表示绝对过期时间
			if !cmd.Flags().Changed("abs_expired_time") || o.AbsExpiredTime == "" {
				return fmt.Errorf("当过期类型为 1（绝对过期时间）时，必须设置 --abs_expired_time")
			}
			if cmd.Flags().Changed("rel_expired_time") {
				return fmt.Errorf("当过期类型为 1（绝对过期时间）时，不应设置 --rel_expired_time")
			}
			o.RelExpiredTime = 0
		}
		res, err := GetYsClient().Account.CashVoucherAdd(
			GetYsClient().Account.CashVoucherAdd.CashVoucherName(o.Name),
			GetYsClient().Account.CashVoucherAdd.Amount(o.Amount),
			GetYsClient().Account.CashVoucherAdd.ExpiredType(o.ExpiredType),
			GetYsClient().Account.CashVoucherAdd.AbsExpiredTime(o.AbsExpiredTime),
			GetYsClient().Account.CashVoucherAdd.RelExpiredTime(o.RelExpiredTime),
			GetYsClient().Account.CashVoucherAdd.Comment(o.Comment),
		)
		PrintResp(res, err, "add CashVoucher")
		return nil
	}
	return cmd
}

func newCashVoucherListCommand(o CashVoucherOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "查询优惠券列表",
		Long:  "列出现有的优惠券",
		Example: `- 查询优惠券列表
  - ysadmin cashvoucher list -O 1 -L 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Limit < 1 || o.Limit > 1000 {
				return fmt.Errorf("分页大小必须在 [1, 1000] 范围内")
			}
			if o.Offset < 0 {
				return fmt.Errorf("分页偏移量不能为负数")
			}

			res, err := GetYsClient().Account.CashVoucherList(
				GetYsClient().Account.CashVoucherList.PageIndex(o.Offset),
				GetYsClient().Account.CashVoucherList.PageSize(o.Limit),
			)
			PrintResp(res, err, "list CashVoucher")
			return nil
		},
	}
	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "分页偏移量")
	cmd.MarkFlagRequired("offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 10, "分页大小")
	cmd.MarkFlagRequired("limit")
	return cmd
}

func newCashVoucherGetCommand(o CashVoucherOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "查询优惠券",
		Long:  "查询优惠券，根据ID查询对应的优惠券",
		Example: `- 查询优惠券
  - ysadmin cashvoucher get -I xxxxxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.CID == "" {
				return fmt.Errorf("请提供优惠券 ID")
			}
			res, err := GetYsClient().Account.CashVoucherGet(
				GetYsClient().Account.CashVoucherGet.CashVoucherID(o.CID),
			)

			PrintResp(res, err, "get CashVoucher")
			return nil
		},
	}
	cmd.Flags().StringVarP(&o.CID, "ID", "I", "", "优惠券ID")
	cmd.MarkFlagRequired("id")
	return cmd
}
