package domain

import (
	"github.com/8treenet/freedom/example/fshop/domain/dependency"
	"github.com/8treenet/freedom/example/fshop/domain/dto"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		//绑定创建领域服务函数到框架，框架会根据客户的使用做依赖倒置和依赖注入的处理。
		initiator.BindService(func() *User {
			return &User{} //创建User领域服务
		})
		//控制器客户使用需要明确使用 InjectController
		initiator.InjectController(func(ctx freedom.Context) (service *User) {
			initiator.GetService(ctx, &service)
			return
		})
	})
}

// User 用户领域服务.
type User struct {
	Worker   freedom.Worker      //运行时，一个请求绑定一个运行时
	UserRepo dependency.UserRepo //依赖倒置用户资源库
}

// ChangePassword 修改密码
func (s *User) ChangePassword(userId int, newPassword, oldPassword string) (e error) {
	//使用用户仓库读取用户实体
	userEntity, e := s.UserRepo.Get(userId)
	if e != nil {
		return
	}

	//修改密码
	if e = userEntity.ChangePassword(newPassword, oldPassword); e != nil {
		return
	}

	//使用用户仓库持久化实体
	e = s.UserRepo.Save(userEntity)
	return
}

// Register .
func (s *User) Register(user dto.RegisterUserReq) (result dto.UserInfoRes, e error) {
	userEntity, e := s.UserRepo.New(user, 10000)
	if e != nil {
		return
	}
	result.Id = userEntity.Id
	result.Money = userEntity.Money
	result.Name = userEntity.Name
	return
}

// Get .
func (s *User) Get(userId int) (result dto.UserInfoRes, e error) {
	userEntity, e := s.UserRepo.Get(userId)
	if e != nil {
		return
	}
	result.Id = userEntity.Id
	result.Money = userEntity.Money
	result.Name = userEntity.Name
	return
}
