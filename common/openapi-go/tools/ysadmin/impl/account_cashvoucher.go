package impl

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RegisterCmd(NewAccountCashVoucherCommand())
}

func NewAccountCashVoucherCommand() *cobra.Command {
	o := CashVoucherOptions{}
	cmd := &cobra.Command{
		Use:   "accountcash",
		Short: "优惠券发放",
		Long:  "优惠券发放，用于为资金账户发放相应的优惠券、创建相应的优惠券",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newAccountCashVoucherAddCommand(&o),
		newAccountCashVoucherGetCommand(&o),
		newAccountCashVoucherListCommand(&o),
	)
	return cmd
}

func newAccountCashVoucherAddCommand(o *CashVoucherOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "为资金账户添加优惠券",
		Long:  "为资金账户添加优惠券",
		Example: `- 为资金账户添加优惠券
  - ysadmin accountcash add -A xxxxxx -C xxxxxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.CID == "" {
				return fmt.Errorf("必须提供优惠券ID，使用 -C 参数指定")
			}
			if o.Id == "" {
				return fmt.Errorf("必须提供至少一个账户ID，使用 -A 参数指定")
			}
			res, err := GetYsClient().Account.AccountCashVoucherAdd(
				GetYsClient().Account.AccountCashVoucherAdd.AccountIDs(o.Id),
				GetYsClient().Account.AccountCashVoucherAdd.CashVoucherID(o.CID),
			)
			PrintResp(res, err, "add AccountCashVoucher")
			return nil
		},
	}
	cmd.Flags().StringVarP(&o.Id, "Id", "A", "", "AccountID")
	cmd.MarkFlagRequired("Id")
	cmd.Flags().StringVarP(&o.CID, "CID", "C", "", "CashVoucherID")
	cmd.MarkFlagRequired("CID")
	return cmd
}

func newAccountCashVoucherGetCommand(o *CashVoucherOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "查询资金账户的优惠券情况",
		Long:  "查询资金账户的优惠券情况",
		Example: `- 查询资金账户的优惠券情况
  - ysadmin accountcash get -I xxxxxx`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Id == "" {
				return fmt.Errorf("必须提供账户ID，使用 -A 参数指定")
			}
			res, err := GetYsClient().Account.AccountCashVoucherGet(
				GetYsClient().Account.AccountCashVoucherGet.AccountCashVoucherID(o.Id),
			)
			PrintResp(res, err, "Get AccountCashVoucher")
			return nil
		},
	}
	cmd.Flags().StringVarP(&o.Id, "Id", "I", "", "AccountID")
	cmd.MarkFlagRequired("Id")
	return cmd
}

func newAccountCashVoucherListCommand(o *CashVoucherOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "查询资金账户的优惠券列表",
		Long:  "查询资金账户的优惠券列表",
		Example: `- 查询资金账户的优优惠券列表
  - ysadmin accountcash list -I xxxxxx -O 1 -L 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Id == "" {
				return fmt.Errorf("必须提供账户ID，使用 -A 参数指定")
			}
			res, err := GetYsClient().Account.AccountCashVoucherList(
				GetYsClient().Account.AccountCashVoucherList.AccountID(o.Id),
				GetYsClient().Account.AccountCashVoucherList.PageIndex(o.Offset),
				GetYsClient().Account.AccountCashVoucherList.PageSize(o.Limit),
			)
			PrintResp(res, err, "list AccountCashVoucher")
			return nil
		},
	}
	cmd.Flags().StringVarP(&o.Id, "Id", "I", "", "AccountID")
	cmd.MarkFlagRequired("Id")
	cmd.Flags().Int64VarP(&o.Offset, "Offset", "O", 1, "分页偏移量")
	cmd.MarkFlagRequired("Offset")
	cmd.Flags().Int64VarP(&o.Limit, "Limit", "L", 1, "分页大小")
	cmd.MarkFlagRequired("Limit")
	return cmd
}
