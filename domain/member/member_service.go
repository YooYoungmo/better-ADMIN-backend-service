package member

import (
	"better-admin-backend-service/dtos"
	"context"
)

type MemberService struct {
}

func (MemberService) GetMemberBySignId(ctx context.Context, signId string) (MemberEntity, error) {
	return memberRepository{}.FindBySignId(ctx, signId)
}

func (MemberService) GetMemberByDoorayId(ctx context.Context, doorayId string) (MemberEntity, error) {
	return memberRepository{}.FindByDoorayId(ctx, doorayId)
}

func (MemberService) CreateMember(ctx context.Context, entity *MemberEntity) error {
	return memberRepository{}.Create(ctx, entity)
}

func (MemberService) GetMemberById(ctx context.Context, id uint) (MemberEntity, error) {
	return memberRepository{}.FindById(ctx, id)
}

func (MemberService) GetMembers(ctx context.Context, filters map[string]interface{}, pageable dtos.Pageable) ([]MemberEntity, int64, error) {
	return memberRepository{}.FindAll(ctx, filters, pageable)
}

func (MemberService) AssignRole(ctx context.Context, memberId uint, assignRole dtos.MemberAssignRole) error {
	repository := memberRepository{}
	memberEntity, err := repository.FindById(ctx, memberId)
	if err != nil {
		return err
	}

	err = memberEntity.AssignRole(ctx, assignRole)
	if err != nil {
		return err
	}

	return repository.Save(ctx, &memberEntity)
}

func (MemberService) GetMember(ctx context.Context, memberId uint) (MemberEntity, error) {
	return memberRepository{}.FindById(ctx, memberId)
}
