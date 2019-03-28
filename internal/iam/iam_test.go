package iam_test
//
// import (
// 	"testing"
// 	"github.com/satori/go.uuid"
// 	"github.com/go-pg/pg/orm"
//
// 	"user-api/internal/iam"
// 	"user-api/internal/model"
//
// 	"user-api/mock"
// 	"user-api/mock/mockdb"
// 	"github.com/stretchr/testify/assert"
//
// 	"github.com/go-pg/pg"
// 	iampb "user-api/rpc/iam"
// )
//
// func TestAuth(t *testing.T) {
// 	cases := []struct {
// 		name     string
// 		req      *iampb.AuthReq
// 		udb      *mockdb.User
// 		sec      *mock.Secure
// 		tg       *mock.JWT
// 		wantErr  bool
// 		wantData *iampb.AuthResp
// 	}{
// 		{
// 			name: "Fail on validation",
// 			req: &iampb.AuthReq{
// 				Auth: "onlyauth",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fail on FindByAuth",
// 			req: &iampb.AuthReq{
// 				Auth:     "email@mail.com",
// 				Password: "hunter2",
// 			},
// 			udb: &mockdb.User{
// 				FindByAuthFn: func(orm.DB, string) (*model.User, error) {
// 					return nil, mock.ErrGeneric
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fail on MatchesHash",
// 			req: &iampb.AuthReq{
// 				Auth:     "juzernejm",
// 				Password: "hunter2",
// 			},
// 			udb: &mockdb.User{
// 				FindByAuthFn: func(orm.DB, string) (*model.User, error) {
// 					return &model.User{
// 						FirstName: "John",
// 						Password:  "(has*_*h3d)",
// 					}, nil
// 				},
// 			},
// 			sec: &mock.Secure{
// 				MatchesHashFn: func(string, string) bool {
// 					return false
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fail on GenerateToken",
// 			req: &iampb.AuthReq{
// 				Auth:     "juzernejm",
// 				Password: "hunter2",
// 			},
// 			udb: &mockdb.User{
// 				FindByAuthFn: func(orm.DB, string) (*model.User, error) {
// 					return &model.User{
// 						FirstName: "John",
// 						Password:  "(has*_*h3d)",
// 					}, nil
// 				},
// 			},
// 			sec: &mock.Secure{
// 				MatchesHashFn: func(string, string) bool {
// 					return true
// 				},
// 			},
// 			tg: &mock.JWT{
// 				GenerateTokenFn: func(*model.AuthUser) (string, error) {
// 					return "", mock.ErrGeneric
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fail on UpdateLastLogin",
// 			req: &iampb.AuthReq{
// 				Auth:     "juzernejm",
// 				Password: "hunter2",
// 			},
// 			udb: &mockdb.User{
// 				FindByAuthFn: func(orm.DB, string) (*model.User, error) {
// 					return &model.User{
// 						FirstName: "John",
// 						Password:  "(has*_*h3d)",
// 					}, nil
// 				},
// 				UpdateLastLoginFn: func(orm.DB, *model.User) error {
// 					return mock.ErrGeneric
// 				},
// 			},
// 			sec: &mock.Secure{
// 				MatchesHashFn: func(string, string) bool {
// 					return true
// 				},
// 			},
// 			tg: &mock.JWT{
// 				GenerateTokenFn: func(*model.AuthUser) (string, error) {
// 					return "jwttoken", nil
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Success",
// 			req: &iampb.AuthReq{
// 				Auth:     "juzernejm",
// 				Password: "hunter2",
// 			},
// 			udb: &mockdb.User{
// 				FindByAuthFn: func(orm.DB, string) (*model.User, error) {
// 					return &model.User{
// 						FirstName: "John",
// 						Password:  "(has*_*h3d)",
// 					}, nil
// 				},
// 				UpdateLastLoginFn: func(orm.DB, *model.User) error {
// 					return nil
// 				},
// 			},
// 			sec: &mock.Secure{
// 				MatchesHashFn: func(string, string) bool {
// 					return true
// 				},
// 			},
// 			tg: &mock.JWT{
// 				GenerateTokenFn: func(*model.AuthUser) (string, error) {
// 					return "jwttoken", nil
// 				},
// 			},
// 			wantData: &iampb.AuthResp{
// 				Token: "jwttoken",
// 			},
// 		},
// 	}
// 	db := &pg.DB{}
// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := iam.New(db, tt.tg, tt.udb, tt.sec)
// 			resp, err := s.Auth(nil, tt.req)
// 			if tt.wantData != nil {
// 				tt.wantData.RefreshToken = resp.RefreshToken
// 			}
// 			assert.Equal(t, tt.wantData, resp)
// 			assert.Equal(t, tt.wantErr, err != nil)
// 		})
// 	}
// }
//
// func TestRefresh(t *testing.T) {
// 	cases := []struct {
// 		name     string
// 		req      *iampb.RefreshReq
// 		udb      *mockdb.User
// 		tg       *mock.JWT
// 		wantErr  bool
// 		wantData *iampb.RefreshResp
// 	}{
// 		{
// 			name: "Fail on validation",
// 			req: &iampb.RefreshReq{
// 				Token: "tooshort",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fail on FindByToken",
// 			req: &iampb.RefreshReq{
// 				Token: "lengthis10lengthis20",
// 			},
// 			udb: &mockdb.User{
// 				FindByTokenFn: func(orm.DB, string) (*model.User, error) {
// 					return nil, mock.ErrGeneric
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Fail on GenerateToken",
// 			req: &iampb.RefreshReq{
// 				Token: "lengthis10lengthis20",
// 			},
// 			udb: &mockdb.User{
// 				FindByTokenFn: func(orm.DB, string) (*model.User, error) {
// 					return &model.User{
// 						Id:       uuid.NewV4(),
// 						TenantId: 321,
// 						Username: "johndoe",
// 						Email:    "johndoe@mail.com",
// 						RoleId:   221,
// 					}, nil
// 				},
// 			},
// 			tg: &mock.JWT{
// 				GenerateTokenFn: func(*model.AuthUser) (string, error) {
// 					return "", mock.ErrGeneric
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Success",
// 			req: &iampb.RefreshReq{
// 				Token: "lengthis10lengthis20",
// 			},
// 			udb: &mockdb.User{
// 				FindByTokenFn: func(orm.DB, string) (*model.User, error) {
// 					return &model.User{
// 						Id:       uuid.NewV4(),
// 						TenantId: 321,
// 						Username: "johndoe",
// 						Email:    "johndoe@mail.com",
// 						RoleId:   221,
// 					}, nil
// 				},
// 			},
// 			tg: &mock.JWT{
// 				GenerateTokenFn: func(*model.AuthUser) (string, error) {
// 					return "newjwttoken", nil
// 				},
// 			},
// 			wantData: &iampb.RefreshResp{
// 				Token: "newjwttoken",
// 			},
// 		},
// 	}
// 	db := &pg.DB{}
// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := iam.New(db, tt.tg, tt.udb, nil)
// 			resp, err := s.Refresh(nil, tt.req)
// 			assert.Equal(t, tt.wantData, resp)
// 			assert.Equal(t, tt.wantErr, err != nil)
// 		})
// 	}
// }
