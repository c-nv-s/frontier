package v1beta1

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/shield/core/action"
	"github.com/odpf/shield/core/resource"
	"github.com/odpf/shield/internal/api/v1beta1/mocks"
	"github.com/odpf/shield/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	shieldv1beta1 "github.com/odpf/shield/proto/v1beta1"
)

func TestHandler_CheckResourcePermission(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(res *mocks.ResourceService)
		request *shieldv1beta1.CheckResourcePermissionRequest
		want    *shieldv1beta1.CheckResourcePermissionResponse
		wantErr error
	}{
		{
			name: "should return internal error if relation service's CheckAuthz function returns some error",
			setup: func(res *mocks.ResourceService) {
				res.EXPECT().CheckAuthz(mock.AnythingOfType("*context.emptyCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, errors.New("some error"))
			},
			request: &shieldv1beta1.CheckResourcePermissionRequest{
				ObjectId:        testRelationV2.Object.ID,
				ObjectNamespace: testRelationV2.Object.NamespaceID,
				Permission:      schema.EditPermission,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return true when CheckAuthz function returns true bool",
			setup: func(res *mocks.ResourceService) {
				res.EXPECT().CheckAuthz(mock.AnythingOfType("*context.emptyCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(true, nil)
			},
			request: &shieldv1beta1.CheckResourcePermissionRequest{
				ObjectId:        testRelationV2.Object.ID,
				ObjectNamespace: testRelationV2.Object.NamespaceID,
				Permission:      schema.EditPermission,
			},
			want: &shieldv1beta1.CheckResourcePermissionResponse{
				Status: true,
			},
			wantErr: nil,
		},
		{
			name: "should return false when CheckAuthz function returns false bool",
			setup: func(res *mocks.ResourceService) {
				res.EXPECT().CheckAuthz(mock.AnythingOfType("*context.emptyCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, nil)
			},
			request: &shieldv1beta1.CheckResourcePermissionRequest{
				ObjectId:        testRelationV2.Object.ID,
				ObjectNamespace: testRelationV2.Object.NamespaceID,
				Permission:      schema.EditPermission,
			},
			want: &shieldv1beta1.CheckResourcePermissionResponse{
				Status: false,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResourceSrv := new(mocks.ResourceService)
			if tt.setup != nil {
				tt.setup(mockResourceSrv)
			}

			mockDep := Handler{resourceService: mockResourceSrv}
			resp, err := mockDep.CheckResourcePermission(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}