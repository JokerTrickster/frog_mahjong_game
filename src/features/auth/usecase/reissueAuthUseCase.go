package usecase

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
	"main/utils"
	"time"
)

type ReissueAuthUseCase struct {
	Repository     _interface.IReissueAuthRepository
	ContextTimeout time.Duration
}

func NewReissueAuthUseCase(repo _interface.IReissueAuthRepository, timeout time.Duration) _interface.IReissueAuthUseCase {
	return &ReissueAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ReissueAuthUseCase) Reissue(c context.Context, req *request.ReqReissue) (response.ResReissue, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 1. 액세스 토큰과 리프레시 토큰 검증
	err := VerifyAccessAndRefresh(req)
	if err != nil {
		return response.ResReissue{}, err
	}
	uID, email, err := utils.ParseToken(req.AccessToken)
	if err != nil {
		return response.ResReissue{}, err
	}
	// 디비에 토큰 있는지 체크
	err = d.Repository.CheckToken(ctx, uID, req.RefreshToken)
	if err != nil {
		return response.ResReissue{}, err
	}

	// 2. 기존 토큰 삭제
	err = d.Repository.DeleteToken(ctx, uID)
	if err != nil {
		return response.ResReissue{}, err
	}

	// 3. 액세스 토큰과 리프레시 토큰 재발급
	accessToken, accessTknExpiredAt, refreshToken, refreshTknExpiredAt, err := utils.GenerateToken(email, uID)
	if err != nil {
		return response.ResReissue{}, err
	}

	// 4. Create TokenDTO
	tokenDTO := CreateTokenDTO(uID, accessToken, accessTknExpiredAt, refreshToken, refreshTknExpiredAt)

	// 5. 액세스 토큰과 리프레시 토큰 db 저장
	err = d.Repository.SaveToken(ctx, tokenDTO)
	if err != nil {
		return response.ResReissue{}, err
	}

	// 6. 응답으로 전달
	res := response.ResReissue{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
