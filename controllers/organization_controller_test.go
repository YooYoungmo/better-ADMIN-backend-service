package controllers

import (
	"better-admin-backend-service/security"
	"better-admin-backend-service/testdata/testdb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOrganizationController_CreateOrganization_최상위_조직으로_추가(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"name": "테스트 조직"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/organizations", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.CreateOrganization, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_CreateOrganization_상위조직이_있는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	requestBody := `{
		"parentOrganizationId": 1,
		"name": "테스트 조직"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/organizations", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.CreateOrganization, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_GetOrganizations(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)

	// when
	handleWithFilter(OrganizationController{}.GetOrganizations, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var resp interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	expected := []interface{}{
		map[string]interface{}{
			"id":   float64(1),
			"name": "베터코드 연구소",
			"subOrganizations": []interface{}{
				map[string]interface{}{
					"id":   float64(3),
					"name": "부서B",
					"subOrganizations": []interface{}{
						map[string]interface{}{
							"id":   float64(4),
							"name": "부서C",
							"roles": []interface{}{
								map[string]interface{}{
									"id":   float64(1),
									"name": "SYSTEM MANAGER",
								},
							},
							"members": []interface{}{
								map[string]interface{}{
									"id":   float64(3),
									"name": "유영모2",
								},
							},
						},
					},
				},
			},
			"roles": []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "SYSTEM MANAGER",
				}, map[string]interface{}{
					"id":   float64(2),
					"name": "MEMBER MANAGER",
				},
			},
			"members": []interface{}{
				map[string]interface{}{
					"id":   float64(1),
					"name": "사이트 관리자",
				}, map[string]interface{}{
					"id":   float64(2),
					"name": "유영모",
				},
			},
		},
		map[string]interface{}{
			"id":   float64(5),
			"name": "베터코드 연구소2",
			"subOrganizations": []interface{}{
				map[string]interface{}{
					"id":   float64(2),
					"name": "부서A",
				},
			},
		},
	}
	assert.Equal(t, expected, resp.([]interface{}))
}

func TestOrganizationController_GetOrganization(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/:organizationId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues("1")

	// when
	handleWithFilter(OrganizationController{}.GetOrganization, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(rec.Body.String())
	var actual interface{}
	json.Unmarshal(rec.Body.Bytes(), &actual)

	expected := map[string]interface{}{
		"id":        float64(1),
		"name":      "베터코드 연구소",
		"createdAt": "1982-01-04T00:00:00Z",
		"roles": []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "SYSTEM MANAGER",
			}, map[string]interface{}{
				"id":   float64(2),
				"name": "MEMBER MANAGER",
			},
		},
		"members": []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "사이트 관리자",
			}, map[string]interface{}{
				"id":   float64(2),
				"name": "유영모",
			},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestOrganizationController_GetOrganization_ID_로_찾을수없는_경우(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/:organizationId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues("100000")

	// when
	handleWithFilter(OrganizationController{}.GetOrganization, ctx)

	// then
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOrganizationController_ChangePosition_하위로_변경(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "2"
	requestBody := `{
		"parentOrganizationId": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/:organizationId/change-position", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.ChangePosition, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_ChangePosition_최상위로_변경(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "2"
	requestBody := `{}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/:organizationId/change-position", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.ChangePosition, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_DeleteOrganization_최하위(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "4"

	req := httptest.NewRequest(http.MethodDelete, "/api/organizations/:organizationId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.DeleteOrganization, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_DeleteOrganization_최상위(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "1"

	req := httptest.NewRequest(http.MethodDelete, "/api/organizations/:organizationId", nil)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.DeleteOrganization, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_AssignRoles(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "1"
	requestBody := `{
		"roleIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/:organizationId/assign-roles", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	// when
	handleWithFilter(OrganizationController{}.AssignRoles, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_AssignMembers(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "1"
	requestBody := `{
		"memberIds": [1, 2]
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/:organizationId/assign-members", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	// when
	handleWithFilter(OrganizationController{}.AssignMembers, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestOrganizationController_ChangeOrganizationName(t *testing.T) {
	testdb.DatabaseFixture{}.SetUpDefault(gormDB)

	// given
	organizationId := "1"
	requestBody := `{
		"name": "강남 베터코드"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/organizations/:organizationId/name", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := echoApp.NewContext(req, rec)
	ctx.SetParamNames("organizationId")
	ctx.SetParamValues(organizationId)

	userClaim := security.UserClaim{
		Id: 2,
	}
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), "userClaim", &userClaim)))

	// when
	handleWithFilter(OrganizationController{}.ChangeOrganizationName, ctx)

	// then
	assert.Equal(t, http.StatusOK, rec.Code)
}
