package util

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
)

func RemoveDuplicates(slice []int64) []int64 {
	// 创建map以存储唯一元素
	seen := make(map[int64]bool)
	result := []int64{}

	// 循环遍历切片，如果没看见过就将元素添加到map中
	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

// ResourceIdentityByName make rbac.ResourceIdentity from resType and resName
func ResourceIdentityByName(resType, resName, resAction string) *rbac.ResourceIdentity {
	if resAction == "" {
		resAction = common.ResourceActionNONE
	}
	resource := &rbac.ResourceName{Type: resType, Name: resName, Action: resAction}
	return &rbac.ResourceIdentity{Identity: &rbac.ResourceIdentity_Name{Name: resource}}
}

// ResourceIdentityByID make rbac.ResourceIdentity from resType and resID
func ResourceIdentityByID(resType string, resID int64) *rbac.ResourceIdentity {
	resource := &rbac.ResourceID{Type: resType, Id: resID}
	return &rbac.ResourceIdentity{Identity: &rbac.ResourceIdentity_Id{Id: resource}}
}

// ResourceIdentityBySimple ResourceIdentityBySimple
func ResourceIdentityBySimple(res *rbac.SimpleResource) *rbac.ResourceIdentity {
	if res.ResourceName == "" {
		return ResourceIdentityByID(res.ResourceType, res.ResourceId)
	}
	return ResourceIdentityByName(res.ResourceType, res.ResourceName, res.ResourceAction)
}

// ResourceIdentitiesBySimple ResourceIdentitiesBySimple
func ResourceIdentitiesBySimple(res *rbac.SimpleResources) []*rbac.ResourceIdentity {
	result := make([]*rbac.ResourceIdentity, len(res.Resources))
	for i, r := range res.Resources {
		result[i] = ResourceIdentityBySimple(r)
	}
	return result
}
