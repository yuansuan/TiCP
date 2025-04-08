package impl

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type OrgServiceImpl struct {
	userDao         dao.UserDao
	orgStructureDao dao.OrgStructureDao
}

func (srv *OrgServiceImpl) ListOrgMember(ctx context.Context, orgID snowflake.ID) ([]*dto.ListMemberResponse, error) {
	member, err := srv.orgStructureDao.ListOrgMember(orgID)
	if err != nil {
		return nil, err
	}

	list := make([]*dto.ListMemberResponse, 0)

	for _, user := range member {
		list = append(list, &dto.ListMemberResponse{
			UserID:   snowflake.ID(user.Id).String(),
			UserName: user.Name,
		})
	}
	return list, nil
}

func (srv *OrgServiceImpl) UpdateOrgMember(ctx context.Context, req dto.UpdateOrgMemberRequest) error {
	err := srv.orgStructureDao.DeleteOrgUserByOrgIds(snowflake.MustParseString(req.OrgID))
	if err != nil {
		return err
	}

	return srv.addOrgMember(ctx, req.OrgID, req.UserList)
}

func (srv *OrgServiceImpl) DeleteOrgMemberByUserID(ctx context.Context, userIDs []string) error {
	return srv.orgStructureDao.DeleteOrgUserByUserIds(snowflake.BatchParseStringToID(userIDs)...)
}

func (srv *OrgServiceImpl) DeleteOrgMemberByID(ctx context.Context, ids []string) error {
	return srv.orgStructureDao.DeleteOrgUserByIds(snowflake.BatchParseStringToID(ids)...)
}

func (srv *OrgServiceImpl) AddOrgMember(ctx context.Context, orgID string, userIDs []string) error {
	return srv.addOrgMember(ctx, orgID, userIDs)
}

func (srv *OrgServiceImpl) UpdateOrg(ctx context.Context, req dto.UpdateOrgRequest) error {

	orgID := snowflake.MustParseString(req.ID)
	err := srv.orgStructureDao.Update(&model.OrgStructure{
		Id:       orgID,
		Name:     req.Name,
		Comment:  req.Comment,
		ParentId: snowflake.MustParseString(req.ParentID),
	})

	if err != nil {
		return err
	}

	err = srv.orgStructureDao.DeleteOrgUserByOrgIds(orgID)

	return srv.addOrgMember(ctx, req.ID, req.UserList)
}

func (srv *OrgServiceImpl) Delete(ctx context.Context, orgID snowflake.ID) error {
	err := srv.orgStructureDao.Delete(orgID)
	if err != nil {
		return err
	}

	err = srv.orgStructureDao.DeleteOrgUserByOrgIds(orgID)
	if err != nil {
		return err
	}

	return nil
}

func NewOrgService() service.OrgService {
	return &OrgServiceImpl{
		userDao:         dao.NewUserDaoImpl(),
		orgStructureDao: dao.NewOrgStructureDaoImpl(),
	}
}

func (srv *OrgServiceImpl) CreateOrg(ctx context.Context, req dto.CreateOrgRequest) error {
	var node *snowflake.Node
	node, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("new snowflake node err: %v", err)
		return err
	}

	orgID := node.Generate()
	_, err = srv.orgStructureDao.Add(&model.OrgStructure{
		Id:       orgID,
		Name:     req.Name,
		Comment:  req.Comment,
		ParentId: snowflake.MustParseString(req.ParentID),
	})
	if err != nil {
		return err
	}

	return srv.addOrgMember(ctx, orgID.String(), req.UserList)
}

func (srv *OrgServiceImpl) addOrgMember(ctx context.Context, orgID string, userIDs []string) error {
	if len(userIDs) > 0 {
		node, err := snowflake.GetInstance()
		if err != nil {
			logging.Default().Errorf("new snowflake node err: %v", err)
			return err
		}

		orgUserList := make([]*model.OrgUserRelation, 0)
		for _, userID := range userIDs {
			orgUserList = append(orgUserList, &model.OrgUserRelation{
				Id:     node.Generate(),
				UserId: snowflake.MustParseString(userID),
				OrgId:  snowflake.MustParseString(orgID),
			})
		}
		err = srv.orgStructureDao.InsertOrgUser(orgUserList...)
		if err != nil {
			return err
		}
	}

	return nil
}
