package impl

import (
	"fmt"

	"github.com/spf13/cobra"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

type IamOptions struct {
	BaseOptions
	Tag string
}

func (o *IamOptions) AddFlags(cmd *cobra.Command) {
	o.BaseOptions.AddBaseOptions(cmd)
	cmd.Flags().StringVarP(&o.Tag, "tag", "T", "", "AK tag，YS_ 开头的为远算云产品账号，不能随意使用")
}

func init() {
	RegisterCmd(NewIamCommand2())
}

func NewIamCommand2() *cobra.Command {
	o := IamOptions{}
	cmd := &cobra.Command{
		Use:   "iam",
		Short: "管理IAM相关资源",
		Long:  "管理IAM相关资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newIamAddCommand(o),
		newIamListCommand(o),
		newIamGetCommand(o),
		newIamUpdateCommand(o),
		newIamDeleteCommand(o),
	)

	return cmd
}

func newIamAddCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "添加IAM相关资源",
		Long:  "添加IAM相关资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newIamAddRoleCommand(o),
		newIamAddPolicyCommand(o),
		newIamAddSecretCommand(o),
		newIamAddRolePolicyRelationCommand(o),
		newIamAddUserCommand(o),
	)

	return cmd
}

func newIamAddRoleCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "为用户添加IAM角色",
		Long:  "为用户添加IAM角色",
		Args:  cobra.ExactArgs(0),
		Example: `- 为用户添加IAM角色, 参数文件为param.json
  - ysadmin iam add role -F param.json
  - param.json内容如下:
    {
        "Description": "world peace",
        "RoleName": "4T4_VIPBoxRole123",
        "TrustPolicy": {
            "Actions": null,
            "Effect": "allow",
            "Principals": [
                "4T4ZZvA2tVc"
            ],
            "Resources": null
        },
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminAddRoleRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		err = NewIamAdminClient().AddRole(req)
		PrintResp(nil, err, "Admin Add Role")

		return nil
	}

	return cmd
}

func newIamAddPolicyCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "为用户添加IAM策略",
		Long:  "为用户添加IAM策略",
		Args:  cobra.ExactArgs(0),
		Example: `- 为用户添加IAM策略, 参数文件为param.json
  - ysadmin iam add policy -F param.json
  - param.json内容如下:
    {
        "PolicyName": "4T4_VIPBoxPolicy123",
        "RoleName": "4T4_VIPBoxRole123",
        "Version": "1.0",
        "Effect": "allow",
        "Resources": [
            "*"
        ],
        "Actions": [
            "*"
        ],
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminAddPolicyRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		err = NewIamAdminClient().AddPolicy(req)
		PrintResp(nil, err, "Admin Add Policy")

		return nil
	}

	return cmd
}

func newIamAddSecretCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "为用户添加IAM密钥",
		Long:  "为用户添加IAM密钥",
		Args:  cobra.ExactArgs(0),
		Example: `- 为用户添加IAM密钥
  - ysadmin iam add secret -I 4TiSxuPtJEm -T YS_admin
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "用户id")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&o.Tag, "tag", "T", "", "AK tag, YS_ 开头的为远算云产品账号，不能随意使用")
	cmd.MarkFlagRequired("tag")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminAddSecretRequest)
		req.Tag = o.Tag
		req.UserId = o.Id
		res, err := NewIamAdminClient().AddSecret(req)
		PrintResp(res, err, "Admin Add Secret")

		return nil
	}

	return cmd
}

func newIamAddRolePolicyRelationCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rolepolicyrelation",
		Short: "为用户添加IAM角色策略关联",
		Long:  "为用户添加IAM角色策略关联",
		Args:  cobra.ExactArgs(0),
		Example: `- 为用户添加IAM角色策略关联, 参数文件为param.json
  - ysadmin iam add rolepolicyrelation -F param.json
  - param.json内容如下:
    {
        "PolicyName": "4T4_VIPBoxPolicy123",
        "RoleName": "4T4_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminAddRolePolicyRelationRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		err = NewIamAdminClient().AddRolePolicyRelation(req)
		PrintResp(nil, err, "Admin Add Role Policy Relation")

		return nil
	}

	return cmd
}

func newIamAddUserCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "添加IAM用户",
		Long:  "添加IAM用户",
		Args:  cobra.ExactArgs(0),
		Example: `- 添加IAM用户, 参数文件为param.json
  - ysadmin iam add user -F param.json
  - param.json内容如下:
    {
        "Phone": "13800138000",
        "Password": "MyPassword@1234"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminAddUserRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		res, err := NewIamAdminClient().AddUser(req)
		PrintResp(res, err, "Admin Add User")

		return nil
	}

	return cmd
}

func newIamListCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出IAM相关资源",
		Long:  "列出IAM相关资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newIamListRoleCommand(o),
		newIamListPolicyCommand(o),
		newIamListSecretCommand(o),
		newIamListSecretsCommand(o),
		newIamListUsersCommand(o),
	)

	return cmd
}

func newIamListRoleCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "列出某用户IAM角色列表",
		Long:  "列出某用户IAM角色列表,",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出某用户IAM角色列表
  - ysadmin iam list role -I 4TiSxuPtJEm
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "用户id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminListRoleRequest)
		req.UserId = o.Id
		res, err := NewIamAdminClient().ListRole(req)
		PrintResp(res, err, "Admin List Role")

		return nil
	}

	return cmd
}

func newIamListPolicyCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "列出某用户IAM策略",
		Long:  "列出某用户IAM策略",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出某用户IAM策略
  - ysadmin iam list policy -I 4TiSxuPtJEm
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "用户id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminListPolicyRequest)
		req.UserId = o.Id
		res, err := NewIamAdminClient().ListPolicy(req)
		PrintResp(res, err, "Admin List Policy")

		return nil
	}

	return cmd
}

func newIamListSecretCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "列出某用户IAM密钥",
		Long:  "列出某用户IAM密钥",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出某用户IAM密钥
  - ysadmin iam list secret -I 4TiSxuPtJEm
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "用户id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminListSecretRequest)
		req.UserId = o.Id
		res, err := NewIamAdminClient().ListSecret(req)
		PrintResp(res, err, "Admin List Secret")

		return nil
	}

	return cmd
}

func newIamListSecretsCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "列出所有IAM密钥",
		Long:  "列出所有IAM密钥",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出所有IAM密钥
  - ysadmin iam list secrets
- 列出所有IAM密钥, 带分页参数
  - ysadmin iam list secrets -O 0 -L 10
`,
	}

	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminListSecretsRequest)
		req.PageOffset = o.Offset
		req.PageSize = o.Limit
		res, err := NewIamAdminClient().ListSecrets(req)
		PrintResp(res, err, "Admin List All Secrets")

		return nil
	}

	return cmd
}

func newIamListUsersCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "列出所有IAM用户",
		Long:  "列出所有IAM用户",
		Args:  cobra.ExactArgs(0),
		Example: `- 列出所有IAM用户
  - ysadmin iam list users
- 列出所有IAM用户, 带分页参数
  - ysadmin iam list users -O 0 -L 10
- 列出所有IAM用户, 带用户id
  - ysadmin iam list users -I 4TiSxuPtJEm
`,
	}

	cmd.Flags().Int64VarP(&o.Offset, "offset", "O", 0, "offset")
	cmd.Flags().Int64VarP(&o.Limit, "limit", "L", 1000, "limit")
	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "用户id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminListUserByNameRequest)
		req.PageOffset = o.Offset
		req.PageSize = o.Limit
		req.Name = o.Id
		res, err := NewIamAdminClient().ListUserByName(req)
		PrintResp(res, err, "Admin List All Users")

		return nil
	}

	return cmd
}

func newIamGetCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "获取IAM相关资源",
		Long:  "获取IAM相关资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newIamGetSecretCommand(o),
		newIamGetRoleCommand(o),
		newIamGetPolicyCommand(o),
		newIamGetUserCommand(o),
	)

	return cmd
}

func newIamGetSecretCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "获取IAM密钥",
		Long:  "获取IAM密钥",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取IAM密钥
  - ysadmin iam get secret -I 6I02NW57S0IJADN08OXK
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "AK id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminGetSecretRequest)
		req.AccessKeyId = o.Id
		res, err := NewIamAdminClient().GetSecret(req)
		PrintResp(res, err, "Get Secret "+o.Id)

		return nil
	}

	return cmd
}

func newIamGetRoleCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "获取某用户某个具体IAM角色信息",
		Long:  "获取某用户某个具体IAM角色信息",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取某用户某个具体IAM角色信息, 参数文件为param.json
  - ysadmin iam get role -F param.json
  - param.json内容如下:
    {
        "RoleName": "4T4_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminGetRoleRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		o.Id = req.UserId
		res, err := NewIamAdminClient().GetRole(req)
		PrintResp(res, err, "Get Role "+o.Id)

		return nil
	}

	return cmd
}

func newIamGetPolicyCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "获取某用户某个具体IAM策略信息",
		Long:  "获取某用户某个具体IAM策略信息",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取某用户某个具体IAM策略信息, 参数文件为param.json
  - ysadmin iam get policy -F param.json
  - param.json内容如下:
    {
        "PolicyName": "4T4_VIPBoxPolicy123",
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminGetPolicyRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		o.Id = req.UserId
		res, err := NewIamAdminClient().GetPolicy(req)
		PrintResp(res, err, "Get Policy "+o.Id)

		return nil
	}

	return cmd
}

func newIamGetUserCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "获取IAM用户",
		Long:  "获取IAM用户",
		Args:  cobra.ExactArgs(0),
		Example: `- 获取IAM用户
  - ysadmin iam get user -I 4TiSxuPtJEm
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "用户id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminGetUserRequest)
		req.UserId = o.Id
		res, err := NewIamAdminClient().GetUserInfo(req)
		PrintResp(res, err, "Admin Get User")

		return nil
	}

	return cmd
}

func newIamUpdateCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "更新IAM相关资源",
		Long:  "更新IAM相关资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newIamUpdateRoleCommand(o),
		newIamUpdatePolicyCommand(o),
		newIamUpdateUserCommand(o),
	)

	return cmd
}

func newIamUpdateRoleCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "更新某用户IAM角色",
		Long:  "更新某用户IAM角色",
		Args:  cobra.ExactArgs(0),
		Example: `- 更新某用户IAM角色, 参数文件为param.json
  - ysadmin iam update role -F param.json
  - param.json内容如下:
    {
        "Description": "world peace",
        "RoleName": "4T4_VIPBoxRole123",
        "TrustPolicy": {
            "Actions": null,
            "Effect": "allow",
            "Principals": [
                "4T4ZZvA2tVc"
            ],
            "Resources": null
        },
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminUpdateRoleRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		err = NewIamAdminClient().UpdateRole(req)
		PrintResp(nil, err, "Admin Update Role")

		return nil
	}

	return cmd
}

func newIamUpdatePolicyCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "更新某用户IAM策略",
		Long:  "更新某用户IAM策略",
		Args:  cobra.ExactArgs(0),
		Example: `- 更新某用户IAM策略, 参数文件为param.json
  - ysadmin iam update policy -F param.json
  - param.json内容如下:
    {
        "PolicyName": "4T4_VIPBoxPolicy123",
        "RoleName": "4T4_VIPBoxRole123",
        "Version": "1.0",
        "Effect": "allow",
        "Resources": [
            "*"
        ],
        "Actions": [
            "*"
        ],
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminUpdatePolicyRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		err = NewIamAdminClient().UpdatePolicy(req)
		PrintResp(nil, err, "Admin Update Policy")

		return nil
	}

	return cmd
}

func newIamUpdateUserCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "更新IAM用户",
		Long:  "更新IAM用户",
		Args:  cobra.ExactArgs(0),
		Example: `- 更新IAM用户, 参数文件为param.json
  - ysadmin iam update user -F param.json
  - param.json内容如下:
    {
        "UserId": "4TiSxxxxxx",
        "Name": "test"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminUpdateUserRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}
		err = NewIamAdminClient().UpdateUser(req)
		PrintResp(nil, err, "Admin Update User")

		return nil
	}

	return cmd
}

func newIamDeleteCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "删除IAM相关资源",
		Long:  "删除IAM相关资源",
		RunE:  helpRun,
	}

	cmd.AddCommand(
		newIamDeleteSecretCommand(o),
		newIamDeleteRoleCommand(o),
		newIamDeletePolicyCommand(o),
		newIamDeleteRolePolicyRelationCommand(o),
	)

	return cmd
}

func newIamDeleteSecretCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "删除IAM密钥",
		Long:  "删除IAM密钥",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除IAM密钥
  - ysadmin iam delete secret -I 6I02NW57S0IJADN08OXK
`,
	}

	cmd.Flags().StringVarP(&o.Id, "id", "I", "", "AK id")
	cmd.MarkFlagRequired("id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminDeleteSecretRequest)
		req.AccessKeyId = o.Id
		err := NewIamAdminClient().DeleteSecret(req)
		PrintResp(nil, err, "Admin Delete Secret")

		return nil
	}

	return cmd
}

func newIamDeleteRoleCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "删除某用户IAM角色",
		Long:  "删除某用户IAM角色",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除某用户IAM角色, 参数文件为param.json
  - ysadmin iam delete role -F param.json
  - param.json内容如下:
    {
        "RoleName": "4T4_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminDeleteRoleRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}

		err = NewIamAdminClient().DeleteRole(req)
		PrintResp(nil, err, "Admin Delete Role")

		return nil
	}

	return cmd
}

func newIamDeletePolicyCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "删除某用户IAM策略",
		Long:  "删除某用户IAM策略",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除某用户IAM策略, 参数文件为param.json
  - ysadmin iam delete policy -F param.json
  - param.json内容如下:
    {
        "PolicyName": "4T4_VIPBoxPolicy123",
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminDeletePolicyRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}

		err = NewIamAdminClient().DeletePolicy(req)
		PrintResp(nil, err, "Admin Delete Policy")

		return nil
	}

	return cmd
}

func newIamDeleteRolePolicyRelationCommand(o IamOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rolepolicyrelation",
		Short: "删除某用户IAM角色策略关联",
		Long:  "删除某用户IAM角色策略关联",
		Args:  cobra.ExactArgs(0),
		Example: `- 删除某用户IAM角色策略关联, 参数文件为param.json
  - ysadmin iam delete rolepolicyrelation -F param.json
  - param.json内容如下:
    {
        "PolicyName": "4T4_VIPBoxPolicy123",
        "RoleName": "4T4_VIPBoxRole123",
        "UserId": "4TiSxuPtJEm"
    }
`,
	}

	cmd.Flags().StringVarP(&o.JsonFile, "file", "F", "", "json文件")
	cmd.MarkFlagRequired("file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		req := new(iam_api.AdminDeleteRolePolicyRelationRequest)
		err := ReadAndUnmarshal(o.JsonFile, req)
		if err != nil {
			fmt.Printf("read json file error: %v\n", err)
			return nil
		}

		err = NewIamAdminClient().DeleteRolePolicyRelation(req)
		PrintResp(nil, err, "Admin Delete Role Policy Relation")

		return nil
	}

	return cmd
}
