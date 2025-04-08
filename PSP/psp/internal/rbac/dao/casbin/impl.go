/*
 * Copyright (C) 2019 LambdaCal Inc.
 */

package casbin

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
)

func intsToStrings(format func(int64) string, ints []int64) []string {
	result := make([]string, len(ints))
	for index, i := range ints {
		result[index] = format(i)
	}
	return result
}

func intsToMap(ints []int64) map[int64]bool {
	result := map[int64]bool{}
	for _, i := range ints {
		result[i] = true
	}
	return result
}

// GetRolePerms GetRolePerms
func (d *Dao) GetRolePerms(roleID int64) (perms []int64) {
	rows, _ := d.enforcer.GetImplicitPermissionsForUser(fmt.RoleID(roleID))
	result := make([]int64, 0, len(rows))
	for _, r := range rows {
		result = append(result, unfmt.PermID(r[1]))
	}
	return result
}

// AddRolePerms AddRolePerms
func (d *Dao) AddRolePerms(roleID int64, permIDs []int64) {
	for _, permID := range permIDs {
		d.enforcer.AddPermissionForUser(fmt.RoleID(roleID), fmt.PermID(permID), "NONE")
	}
}

// UpdateRolePerm UpdateRolePerm
func (d *Dao) UpdateRolePerms(roleID int64, permIDs []int64) {
	d.enforcer.DeletePermissionsForUser(fmt.RoleID(roleID))
	d.AddRolePerms(roleID, permIDs)
}

// RemoveRolePerm RemoveRolePerm
func (d *Dao) RemoveRolePerms(roleID int64, permIDs []int64) {
	for _, permID := range permIDs {
		d.enforcer.DeletePermissionForUser(fmt.RoleID(roleID), fmt.PermID(permID))
	}
}

// GetObjectRoles GetObjectRoles
func (d *Dao) GetObjectRoles(id *rbac.ObjectID) (roleIDs []int64) {
	roles, _ := d.enforcer.GetRolesForUser(fmt.ObjectID(id))

	roleIDs = make([]int64, 0)
	for _, id := range roles {
		roleIDs = append(roleIDs, unfmt.RoleID(id))
	}
	return
}

// GetRoleObjects GetRoleObjects
func (d *Dao) GetRoleObjects(roleID int64) (objectIDs []*rbac.ObjectID) {
	return d.roleGetObjectIDs(fmt.RoleID(roleID))
}

func (d *Dao) roleGetObjectIDs(id string) (objectIDs []*rbac.ObjectID) {
	objs, _ := d.enforcer.GetUsersForRole(id)
	result := make([]*rbac.ObjectID, 0, len(objs))
	for _, obj := range objs {
		result = append(result, unfmt.ObjectID(obj))
	}
	return result
}

// AddObjectRoles AddObjectRoles
func (d *Dao) AddObjectRoles(objectRoles *rbac.ObjectRoles) {
	for _, roleID := range objectRoles.Roles {
		d.enforcer.AddRoleForUser(fmt.ObjectID(objectRoles.Id), fmt.RoleID(roleID))
	}
}

// UpdateObjectRoles UpdateObjectRoles
func (d *Dao) UpdateObjectRoles(objectRoles *rbac.ObjectRoles) {
	// 先删除原来赋予的角色
	roleIDs := d.GetObjectRoles(objectRoles.Id)
	d.RemoveObjectRoles(&rbac.ObjectRoles{
		Id:    objectRoles.Id,
		Roles: roleIDs,
	})
	// 添加新角色
	d.AddObjectRoles(objectRoles)
}

// RemoveObjectRoles RemoveObjectRoles
func (d *Dao) RemoveObjectRoles(objectRoles *rbac.ObjectRoles) {
	for _, roleID := range objectRoles.Roles {
		d.enforcer.DeleteRoleForUser(fmt.ObjectID(objectRoles.Id), fmt.RoleID(roleID))
	}
}

// RemoveRole RemoveRole
func (d *Dao) RemoveRole(id int64) {
	d.enforcer.DeleteRole(fmt.RoleID(id))
}

// GetObjectPermissions GetObjectPermissions 查询用户具有的所有权限id
func (d *Dao) GetObjectPermissions(request *rbac.ObjectID) []int64 {
	rows, _ := d.enforcer.GetImplicitPermissionsForUser(fmt.ObjectID(request))
	result := make([]int64, 0, len(rows))
	for _, r := range rows {
		result = append(result, unfmt.PermID(r[1]))
	}
	return result
}

// CheckPermissions CheckPermissions
// todo 改成d.enforcer.Enforce()
func (d *Dao) CheckPermissions(request *rbac.CheckPermissionsRequest) bool {
	objectPerms := intsToMap(d.GetObjectPermissions(request.Id))

	for _, perm := range request.PermissionIds {
		if !objectPerms[perm] {
			return false
		}
	}
	return true
}

// RemovePerm RemovePerm
func (d *Dao) RemovePerm(id int64) {
	d.enforcer.DeletePermission(fmt.PermID(id))
}

func (d *Dao) AllEntity(ctx context.Context) (perm []int64, roleHavePerm []int64) {

	sPerm := d.enforcer.GetAllObjects()
	perm = make([]int64, len(sPerm))
	for i, p := range sPerm {
		perm[i] = unfmt.PermID(p)
	}

	sRoleHavePerm := d.enforcer.GetAllSubjects()
	roleHavePerm = make([]int64, len(sRoleHavePerm))
	for i, r := range sRoleHavePerm {
		roleHavePerm[i] = unfmt.RoleID(r)
	}

	return
}

func (d *Dao) GetObjectsForPermission(permId int64) ([]*rbac.ObjectID, error) {
	objs, err := d.enforcer.GetImplicitUsersForPermission(fmt.PermID(permId), "NONE")
	if err != nil {
		return nil, err
	}

	result := make([]*rbac.ObjectID, 0, len(objs))
	for _, obj := range objs {
		result = append(result, unfmt.ObjectID(obj))
	}
	return result, nil
}
