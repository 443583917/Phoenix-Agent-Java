package api

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/agent/runner"
	"github.com/phoenix-agent-go/agent/runtime"
	"github.com/phoenix-agent-go/agent/tools/datasource"
	"github.com/phoenix-agent-go/api/handler/agent"
	"github.com/phoenix-agent-go/api/handler/chat"
	"github.com/phoenix-agent-go/api/handler/common"
	datasourceHandler "github.com/phoenix-agent-go/api/handler/datasource"
	kgHandler "github.com/phoenix-agent-go/api/handler/kg"
	knowledgeHandler "github.com/phoenix-agent-go/api/handler/knowledge"
	"github.com/phoenix-agent-go/api/handler/modelconfig"
	"github.com/phoenix-agent-go/api/handler/platform"
	"github.com/phoenix-agent-go/api/handler/privilege"
	"github.com/phoenix-agent-go/api/handler/prompt"
	ragHandler "github.com/phoenix-agent-go/api/handler/rag"
	semanticmodelHandler "github.com/phoenix-agent-go/api/handler/semanticmodel"
	"github.com/phoenix-agent-go/api/middleware"
	"github.com/phoenix-agent-go/agent/knowledge"
	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
	"github.com/phoenix-agent-go/internal/service/tracing"
	"github.com/redis/go-redis/v9"
	tmodel "trpc.group/trpc-go/trpc-agent-go/model"
)

func SetupRouter(cfg *config.AppConfig, jwtManager *jwt.JWTManager, enforcer *casbin.Enforcer, privilegeSvc *service.PrivilegeService, platformSvc *service.PlatformService, dataSvc *service.DataService, ragSvc *service.RagService, kgSvc *service.KgService, agentManager *runtime.AgentManager, hitlHandler *runner.HitlHandler, rdb *redis.Client, llmModel tmodel.Model, dbManager *datasource.DatasourceManager, retriever *knowledge.Retriever, mtcm *runtime.MultiTurnContextManager, tracingSvc *tracing.TracingService, embeddingSvc *service.EmbeddingService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Tracing())
	r.Use(middleware.CORS(&cfg.Cors))

	// 健康检查
	r.GET("/echo", func(c *gin.Context) {
		response.Success(c, "ok")
	})

	// 静态文件（头像等）
	r.Static("/api/upload", "./storage/upload")

	// 认证路由（无需 JWT）
	authHandler := privilege.NewAuthHandler(privilegeSvc, jwtManager, rdb)
	auth := r.Group("/api/privilege/auth")
	auth.Use(middleware.RateLimit())
	auth.GET("/captcha", authHandler.Captcha)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/menus", middleware.Auth(jwtManager), authHandler.Menus)
	auth.GET("/getLoginUserInfo", middleware.Auth(jwtManager), authHandler.GetLoginUserInfo)

	// API 路由（需 JWT）
	api := r.Group("/api")
	api.Use(middleware.Auth(jwtManager))
	api.Use(middleware.RBAC(enforcer))
	{
		// 用户管理
		userHandler := privilege.NewUserHandler(privilegeSvc)
		userGroup := api.Group("/privilege/user")
		{
			userGroup.GET("/page", userHandler.Page)
			userGroup.GET("/:id", userHandler.GetByID)
			userGroup.GET("/username/:username", userHandler.GetByUsername)
			userGroup.GET("/code/:code", userHandler.GetByCode)
			userGroup.POST("", userHandler.Create)
			userGroup.PUT("", userHandler.Update)
			userGroup.DELETE("/:id", userHandler.Delete)
			userGroup.PUT("/password", userHandler.UpdatePassword)
			userGroup.PUT("/setPassword", userHandler.SetPassword)
			userGroup.PUT("/reset-password/:id", userHandler.ResetPassword)
		}

		// 角色管理
		roleHandler := privilege.NewRoleHandler(privilegeSvc)
		roleGroup := api.Group("/privilege/role")
		{
			roleGroup.POST("/page", roleHandler.Page)
			roleGroup.GET("/:id", roleHandler.GetByID)
			roleGroup.GET("/company/:companyId", roleHandler.GetByCompany)
			roleGroup.POST("", roleHandler.Create)
			roleGroup.PUT("", roleHandler.Update)
			roleGroup.DELETE("/:id", roleHandler.Delete)
			roleGroup.GET("/:id/acls", roleHandler.GetAcls)
		}

		// 用户-角色关联管理
		userRoleHandler := privilege.NewUserRoleHandler(privilegeSvc)
		userRoleGroup := api.Group("/privilege/user-role")
		{
			userRoleGroup.GET("/user/:userId", userRoleHandler.GetByUser)
			userRoleGroup.GET("/role/:roleId", userRoleHandler.GetByRole)
			userRoleGroup.POST("", userRoleHandler.Create)
			userRoleGroup.DELETE("/:id", userRoleHandler.Delete)
			userRoleGroup.POST("/batch-save", userRoleHandler.BatchSave)
			userRoleGroup.DELETE("/batch-remove", userRoleHandler.BatchRemove)
		}

		// 模块管理
		moduleHandler := privilege.NewModuleHandler(privilegeSvc)
		moduleGroup := api.Group("/privilege/module")
		{
			moduleGroup.GET("/page", moduleHandler.Page)
			moduleGroup.GET("/tree", moduleHandler.Tree)
			moduleGroup.GET("/tree/acl", moduleHandler.TreeWithACL)
			moduleGroup.GET("/:id", moduleHandler.GetByID)
			moduleGroup.GET("/system/:systemId", moduleHandler.GetBySystem)
			moduleGroup.GET("/pid/:pid", moduleHandler.GetByPID)
			moduleGroup.POST("", moduleHandler.Create)
			moduleGroup.PUT("", moduleHandler.Update)
			moduleGroup.GET("/pvalues/:moduleId", moduleHandler.GetPvalues)
			moduleGroup.PUT("/pvalue/:moduleId/:position/:enabled", moduleHandler.UpdatePvalue)
		}

		// ACL 管理
		aclHandler := privilege.NewACLHandler(privilegeSvc)
		aclGroup := api.Group("/privilege/acl")
		{
			aclGroup.GET("/page", aclHandler.Page)
			aclGroup.GET("/:id", aclHandler.GetByID)
			aclGroup.GET("/release/:releaseId", aclHandler.GetByRelease)
			aclGroup.GET("/release/module/:releaseId/:moduleId", aclHandler.GetByReleaseAndModule)
			aclGroup.PUT("", aclHandler.Update)
			aclGroup.DELETE("/:id", aclHandler.Delete)
			aclGroup.POST("/saveAll/:releaseId/:checkStatus", aclHandler.SaveAll)
			aclGroup.POST("/saveModule", aclHandler.SaveModule)
		}

		// 部门管理
		departmentHandler := privilege.NewDepartmentHandler(privilegeSvc)
		departmentGroup := api.Group("/privilege/department")
		{
			departmentGroup.GET("/orgTree", departmentHandler.OrgTree)
			departmentGroup.GET("/page", departmentHandler.Page)
			departmentGroup.GET("/tree", departmentHandler.Tree)
			departmentGroup.GET("/:id", departmentHandler.GetByID)
			departmentGroup.GET("/pid/:pid", departmentHandler.GetByPID)
			departmentGroup.GET("/company/:companyId", departmentHandler.GetByCompany)
			departmentGroup.GET("/code/:code", departmentHandler.GetByCode)
			departmentGroup.POST("", departmentHandler.Create)
			departmentGroup.PUT("", departmentHandler.Update)
			departmentGroup.POST("/sync", departmentHandler.Sync)
			departmentGroup.POST("/sync-children/:deptId", departmentHandler.SyncChildren)
		}

		// 公司管理
		companyHandler := privilege.NewCompanyHandler(privilegeSvc)
		companyGroup := api.Group("/privilege/company")
		{
			companyGroup.POST("/page", companyHandler.Page)
			companyGroup.GET("/:id", companyHandler.GetByID)
			companyGroup.GET("/code/:code", companyHandler.GetByCode)
			companyGroup.POST("", companyHandler.Create)
			companyGroup.PUT("", companyHandler.Update)
			companyGroup.DELETE("/:id", companyHandler.Delete)
		}

		// 员工绑定管理
		employeeHandler := privilege.NewEmployeeHandler(privilegeSvc)
		employeeGroup := api.Group("/privilege/employee")
		{
			employeeGroup.GET("/page", employeeHandler.Page)
			employeeGroup.GET("/:id", employeeHandler.GetByID)
			employeeGroup.GET("/emp-code/:empCode", employeeHandler.GetByEmpCode)
			employeeGroup.POST("", employeeHandler.Create)
			employeeGroup.PUT("", employeeHandler.Update)
			employeeGroup.DELETE("/:id", employeeHandler.Delete)
			employeeGroup.POST("/sync", employeeHandler.Sync)
			employeeGroup.POST("/sync-by-dept/:deptId", employeeHandler.SyncByDept)
		}

		// 字典管理
		dictionaryHandler := privilege.NewDictionaryHandler(privilegeSvc)
		dictionaryGroup := api.Group("/privilege/dictionary")
		{
			dictionaryGroup.GET("/page", dictionaryHandler.Page)
			dictionaryGroup.GET("/:id", dictionaryHandler.GetByID)
			dictionaryGroup.GET("/system/:systemSn", dictionaryHandler.GetBySystem)
			dictionaryGroup.GET("/pcode/:pcode", dictionaryHandler.GetByPCode)
			dictionaryGroup.POST("", dictionaryHandler.Create)
			dictionaryGroup.PUT("", dictionaryHandler.Update)
			dictionaryGroup.DELETE("/:id", dictionaryHandler.Delete)
		}

		// 权限值管理
		pvalueHandler := privilege.NewPvalueHandler(privilegeSvc)
		pvalueGroup := api.Group("/privilege/pvalue")
		{
			pvalueGroup.POST("/page", pvalueHandler.Page)
			pvalueGroup.GET("/:id", pvalueHandler.GetByID)
			pvalueGroup.GET("/system", pvalueHandler.GetBySystem)
			pvalueGroup.POST("", pvalueHandler.Create)
			pvalueGroup.PUT("", pvalueHandler.Update)
			pvalueGroup.DELETE("/:id", pvalueHandler.Delete)
		}

		// 登录日志管理
		loginLogHandler := privilege.NewLoginLogHandler(privilegeSvc)
		loginLogGroup := api.Group("/privilege/login-log")
		{
			loginLogGroup.GET("/page", loginLogHandler.Page)
			loginLogGroup.GET("/:id", loginLogHandler.GetByID)
			loginLogGroup.POST("", loginLogHandler.Create)
			loginLogGroup.DELETE("/:id", loginLogHandler.Delete)
		}
	}

	// 共用 handler（需提前创建以跨路由组复用）
	reactAgentHandler := agent.NewReactAgentHandler(agentManager)
	harnessHandler := agent.NewHarnessHandler(agentManager, hitlHandler)
	graphHandler := chat.NewGraphHandler(dataSvc, llmModel, dbManager, retriever, mtcm, tracingSvc)

	// Agent 路由（需 JWT）
	admin := r.Group("/api/admin")
	admin.Use(middleware.Auth(jwtManager))
	{
		agentGroup := admin.Group("/agent")
		{
			agentGroup.POST("/chat", reactAgentHandler.Chat)
			agentGroup.GET("/stream/chatsql", reactAgentHandler.StreamChatSQL)
		}

		harnessGroup := admin.Group("/harness")
		{
			harnessGroup.POST("/chat", harnessHandler.Chat)
			harnessGroup.POST("/confirm", harnessHandler.Confirm)
		}
	}

	// ================ 前端 C 端路由 (/api/front) ================
	front := r.Group("/api/front")
	front.Use(middleware.Auth(jwtManager))
	{
		front.POST("/stream/chat", reactAgentHandler.Chat)
		front.GET("/stream/chatsql", graphHandler.StreamSearch)
		front.POST("/harness/chat", harnessHandler.Chat)
		front.POST("/harness/confirm", harnessHandler.Confirm)
		if dataSvc != nil {
			pqFrontHandler := agent.NewAgentPresetQuestionHandler(dataSvc)
			front.GET("/:id/preset-questions", pqFrontHandler.List)
			front.POST("/addPresetQuestion", pqFrontHandler.Create)
			front.DELETE("/deletePresetQuestion/:id", pqFrontHandler.Delete)
		}
	}

	// 数据管理路由（需 JWT）
	dataAPI := r.Group("/api")
	dataAPI.Use(middleware.Auth(jwtManager))
	dataAPI.Use(middleware.RBAC(enforcer))
	if dataSvc != nil {
		// ---- Agent ----
		agentDataHandler := agent.NewAgentHandler(dataSvc)
		agentDataGroup := dataAPI.Group("/agent")
		{
			agentDataGroup.GET("/page", agentDataHandler.Page)
			agentDataGroup.GET("/list", agentDataHandler.List)
			agentDataGroup.GET("/:id", agentDataHandler.GetByID)
			agentDataGroup.POST("", agentDataHandler.Create)
			agentDataGroup.PUT("", agentDataHandler.Update)
			agentDataGroup.DELETE("/:id", agentDataHandler.Delete)
			agentDataGroup.POST("/:id/publish", agentDataHandler.Publish)
			agentDataGroup.POST("/:id/offline", agentDataHandler.Offline)
			agentDataGroup.POST("/:id/api-key/generate", agentDataHandler.GenerateAPIKey)
			agentDataGroup.POST("/:id/api-key/reset", agentDataHandler.ResetAPIKey)
			agentDataGroup.DELETE("/:id/api-key", agentDataHandler.DeleteAPIKey)
			agentDataGroup.PUT("/:id/api-key/toggle", agentDataHandler.ToggleAPIKeyEnabled)
			agentDataGroup.GET("/:id/api-key", agentDataHandler.GetAPIKeyMasked)
		}

		// ---- AgentCategory ----
		categoryHandler := agent.NewAgentCategoryHandler(dataSvc)
		categoryGroup := dataAPI.Group("/agent-category")
		{
			categoryGroup.GET("/page", categoryHandler.Page)
			categoryGroup.GET("/list", categoryHandler.List)
			categoryGroup.GET("/tree", categoryHandler.Tree)
			categoryGroup.GET("/pid/:pid", categoryHandler.GetByPID)
			categoryGroup.GET("/:id", categoryHandler.GetByID)
			categoryGroup.POST("", categoryHandler.Create)
			categoryGroup.PUT("", categoryHandler.Update)
			categoryGroup.DELETE("/:id", categoryHandler.Delete)
		}

		// ---- AgentDatasource (scoped under /api/agent/:agentId/datasource) ----
		dsHandler := datasourceHandler.NewAgentDatasourceHandler(dataSvc)
		dsGroup := dataAPI.Group("/agent/:id/datasource")
		{
			dsGroup.GET("/page", dsHandler.Page)
			dsGroup.GET("/:id", dsHandler.GetByID)
			dsGroup.POST("", dsHandler.Create)
			dsGroup.PUT("", dsHandler.Update)
			dsGroup.DELETE("/:id", dsHandler.Delete)
			dsGroup.PUT("/:id/toggle", dsHandler.ToggleActive)
			dsGroup.GET("/:id/tables", dsHandler.GetTables)
			dsGroup.POST("/:id/tables", dsHandler.SaveTables)
		}
		dataAPI.GET("/agent/:id/datasources", dsHandler.List)
		dataAPI.GET("/agent/:id/datasources/active", dsHandler.GetActive)
		dataAPI.POST("/agent/:id/datasources/init", dsHandler.InitSchema)

		// ---- AgentKnowledge ----
		knHandler := knowledgeHandler.NewAgentKnowledgeHandler(dataSvc, embeddingSvc)
		knGroup := dataAPI.Group("/agent-knowledge")
		{
			knGroup.GET("/page", knHandler.Page)
			knGroup.POST("/query/page", knHandler.QueryPage)
			knGroup.GET("/:id", knHandler.GetByID)
			knGroup.POST("", knHandler.Create)
			knGroup.PUT("", knHandler.Update)
			knGroup.DELETE("/:id", knHandler.Delete)
			knGroup.PUT("/:id/recall-toggle", knHandler.ToggleRecall)
			knGroup.POST("/:id/retry-embedding", knHandler.RetryEmbedding)
		}

		// ---- AgentPresetQuestion (scoped under /api/agent/:agentId/preset-question) ----
		pqHandler := agent.NewAgentPresetQuestionHandler(dataSvc)
		pqGroup := dataAPI.Group("/agent/:id/preset-question")
		{
			pqGroup.GET("/page", pqHandler.Page)
			pqGroup.GET("/", pqHandler.List)
			pqGroup.GET("/:id", pqHandler.GetByID)
			pqGroup.POST("", pqHandler.Create)
			pqGroup.PUT("", pqHandler.Update)
			pqGroup.DELETE("/:id", pqHandler.Delete)
		}

		// ---- AgentPresetQuestion batch + account filter ----
			dataAPI.GET("/agent/:id/preset-questions", pqHandler.List)
		dataAPI.GET("/agent/:id/:accountId/preset-questions", pqHandler.ListByAccount)
		dataAPI.POST("/agent/:id/preset-questions", pqHandler.BatchSave)
		// ---- Chat - Session & Message Management ----
		chatHandler := chat.NewChatHandler(dataSvc)
		dataAPI.GET("/agent/:id/sessions", chatHandler.ListSessions)
		dataAPI.POST("/agent/:id/sessions", chatHandler.CreateSession)
		dataAPI.DELETE("/agent/:id/sessions", chatHandler.DeleteAllSessions)
		dataAPI.GET("/sessions/:id/messages", chatHandler.GetMessages)
		dataAPI.POST("/sessions/:id/messages", chatHandler.CreateMessage)
		dataAPI.PUT("/sessions/:id/pin", chatHandler.PinSession)
		dataAPI.PUT("/sessions/:id/rename", chatHandler.RenameSession)
		dataAPI.DELETE("/sessions/:id", chatHandler.DeleteSession)
		dataAPI.POST("/sessions/:id/reports/html", chatHandler.GenerateReportHTML)

		// ---- Graph Search (SSE) ----
		dataAPI.GET("/stream/search", graphHandler.StreamSearch)

		// ---- Session Events (SSE) ----
		sessionEventHandler := chat.NewSessionEventHandler(dataSvc)
		dataAPI.GET("/agent/:id/sessions/stream", sessionEventHandler.StreamSessions)

		// ---- Datasource Management ----
		dsDataHandler := datasourceHandler.NewDatasourceHandler(dataSvc)
		dsDataGroup := dataAPI.Group("/datasource")
		{
			dsDataGroup.GET("/types", dsDataHandler.GetTypes)
			dsDataGroup.GET("/", dsDataHandler.List)
			dsDataGroup.GET("/:id", dsDataHandler.GetByID)
			dsDataGroup.GET("/:id/tables", dsDataHandler.GetTables)
			dsDataGroup.GET("/:id/tables/:tableName/columns", dsDataHandler.GetColumns)
			dsDataGroup.POST("", dsDataHandler.Create)
			dsDataGroup.PUT("/:id", dsDataHandler.Update)
			dsDataGroup.DELETE("/:id", dsDataHandler.Delete)
			dsDataGroup.POST("/:id/test", dsDataHandler.TestConnection)
			dsDataGroup.GET("/:id/logical-relations", dsDataHandler.ListLogicalRelations)
			dsDataGroup.POST("/:id/logical-relations", dsDataHandler.CreateLogicalRelation)
			dsDataGroup.PUT("/:id/logical-relations", dsDataHandler.UpdateLogicalRelation)
			dsDataGroup.PUT("/:id/logical-relations/:relationId", dsDataHandler.UpdateSingleLogicalRelation)
			dsDataGroup.DELETE("/:id/logical-relations", dsDataHandler.DeleteLogicalRelation)
			dsDataGroup.DELETE("/:id/logical-relations/:relationId", dsDataHandler.DeleteSingleLogicalRelation)
		}

		// ---- Model Configuration ----
		mcHandler := modelconfig.NewModelConfigHandler(dataSvc)
		mcGroup := dataAPI.Group("/model-config")
		{
			mcGroup.GET("/list", mcHandler.List)
			mcGroup.POST("/add", mcHandler.Add)
			mcGroup.PUT("/update", mcHandler.Update)
			mcGroup.DELETE("/:id", mcHandler.Delete)
			mcGroup.POST("/activate/:id", mcHandler.Activate)
			mcGroup.POST("/test", mcHandler.Test)
			mcGroup.GET("/check-ready", mcHandler.CheckReady)
		}

		// ---- Prompt Configuration ----
		pcHandler := prompt.NewPromptConfigHandler(dataSvc)
		pcGroup := dataAPI.Group("/prompt-config")
		{
			pcGroup.POST("/save", pcHandler.Save)
			pcGroup.GET("/:id", pcHandler.GetByID)
			pcGroup.GET("/list", pcHandler.List)
			pcGroup.GET("/list-by-type/:type", pcHandler.ListByType)
			pcGroup.GET("/active/:type", pcHandler.GetActiveByType)
			pcGroup.GET("/active-all/:type", pcHandler.GetActiveAllByType)
			pcGroup.DELETE("/:id", pcHandler.Delete)
			pcGroup.POST("/:id/enable", pcHandler.Enable)
			pcGroup.POST("/:id/disable", pcHandler.Disable)
			pcGroup.GET("/types", pcHandler.GetTypes)
			pcGroup.POST("/batch-enable", pcHandler.BatchEnable)
			pcGroup.POST("/batch-disable", pcHandler.BatchDisable)
			pcGroup.POST("/:id/priority", pcHandler.SetPriority)
			pcGroup.POST("/:id/display-order", pcHandler.SetDisplayOrder)
		}

		// ---- Semantic Model ----
		smHandler := semanticmodelHandler.NewSemanticModelHandler(dataSvc)
		smGroup := dataAPI.Group("/semantic-model")
		{
			smGroup.GET("/", smHandler.List)
			smGroup.GET("/:id", smHandler.GetByID)
			smGroup.POST("", smHandler.Create)
			smGroup.PUT("/:id", smHandler.Update)
			smGroup.DELETE("/:id", smHandler.Delete)
			smGroup.DELETE("/batch", smHandler.BatchDelete)
			smGroup.PUT("/enable", smHandler.Enable)
			smGroup.PUT("/disable", smHandler.Disable)
			smGroup.POST("/batch-import", smHandler.BatchImport)
			smGroup.GET("/template/download", smHandler.DownloadTemplate)
			smGroup.POST("/import/excel", smHandler.ImportExcel)
		}

		// ---- Business Knowledge ----
		bkHandler := knowledgeHandler.NewBusinessKnowledgeHandler(dataSvc)
		bkGroup := dataAPI.Group("/business-knowledge")
		{
			bkGroup.GET("/", bkHandler.List)
			bkGroup.GET("/:id", bkHandler.GetByID)
			bkGroup.POST("", bkHandler.Create)
			bkGroup.PUT("/:id", bkHandler.Update)
			bkGroup.DELETE("/:id", bkHandler.Delete)
			bkGroup.POST("/recall/:id", bkHandler.Recall)
			bkGroup.POST("/refresh-vector-store", bkHandler.RefreshVectorStore)
			bkGroup.POST("/retry-embedding/:id", bkHandler.RetryEmbedding)
		}

		// ---- File Upload ----
		uploadHandler := common.NewFileUploadHandler()
		uploadGroup := dataAPI.Group("/upload")
		{
			uploadGroup.POST("/avatar", uploadHandler.UploadAvatar)
		}

		// ---- 前端兼容别名路由（无 /api 前缀） ----
		// 前端请求经 Vite 开发代理 rewrite 后 /api 前缀被剥除，
		// 以下接口需以无前缀形式暴露，以兼容前端现有请求路径。
		compat := r.Group("/")
		compat.Use(middleware.Auth(jwtManager))
		compat.Use(middleware.RBAC(enforcer))
		{
			compat.POST("/upload/avatar", uploadHandler.UploadAvatar)
			compat.GET("/semantic-model/template/download", smHandler.DownloadTemplate)

			dsCompat := compat.Group("/datasource")
			{
				dsCompat.GET("/:id/logical-relations", dsDataHandler.ListLogicalRelations)
				dsCompat.POST("/:id/logical-relations", dsDataHandler.CreateLogicalRelation)
				dsCompat.PUT("/:id/logical-relations", dsDataHandler.UpdateLogicalRelation)
				dsCompat.PUT("/:id/logical-relations/:relationId", dsDataHandler.UpdateSingleLogicalRelation)
				dsCompat.DELETE("/:id/logical-relations", dsDataHandler.DeleteLogicalRelation)
				dsCompat.DELETE("/:id/logical-relations/:relationId", dsDataHandler.DeleteSingleLogicalRelation)
				dsCompat.GET("/:id/tables/:tableName/columns", dsDataHandler.GetColumns)
			}
		}
	}

	// 平台认证路由（仅 updatePassword 需要 JWT）
	platformAuthHandler := platform.NewAccountLoginHandler(platformSvc, jwtManager)
	platformAuth := r.Group("/auth")
	{
		platformAuth.POST("/login", platformAuthHandler.Login)
		platformAuth.POST("/logout", platformAuthHandler.Logout)
		platformAuth.POST("/thirdLogin", platformAuthHandler.ThirdLogin)
		platformAuth.PUT("/updatePassword", middleware.Auth(jwtManager), platformAuthHandler.UpdatePassword)
	}

	// 平台管理路由
	plat := r.Group("/platform")
	plat.Use(middleware.Auth(jwtManager))
	{
		// ── GroupInfo ──
		groupInfoHandler := platform.NewGroupInfoHandler(platformSvc)
		group := plat.Group("/group-info")
		{
			group.GET("/page", groupInfoHandler.Page)
			group.GET("/:id", groupInfoHandler.GetByID)
			group.GET("/sn/:sn", groupInfoHandler.GetBySN)
			group.POST("", groupInfoHandler.Create)
			group.PUT("", groupInfoHandler.Update)
			group.DELETE("/:id", groupInfoHandler.Delete)
			group.PUT("/:id/toggle-status", groupInfoHandler.ToggleStatus)
			group.DELETE("/remove-agent/:groupId/:agentId", groupInfoHandler.RemoveAgent)
		}

		// ── GroupAgentInfo ──
		gaHandler := platform.NewGroupAgentInfoHandler(platformSvc)
		ga := plat.Group("/group-agent-info")
		{
			ga.GET("/page", gaHandler.Page)
			ga.GET("/:id", gaHandler.GetByID)
			ga.GET("/group/:groupId", gaHandler.GetByGroupID)
			ga.GET("/agent/:agentId", gaHandler.GetByAgentID)
			ga.POST("", gaHandler.Create)
			ga.PUT("", gaHandler.Update)
			ga.DELETE("/:id", gaHandler.Delete)
			ga.DELETE("/group/:groupId/agent/:agentId", gaHandler.DeleteByGroupAndAgent)
			ga.GET("/list", gaHandler.List)
		}

		// ── AccountInfo ──
		acctHandler := platform.NewAccountInfoHandler(platformSvc)
		acct := plat.Group("/account-info")
		{
			acct.GET("/page", acctHandler.Page)
			acct.GET("/:id", acctHandler.GetByID)
			acct.GET("/username/:username", acctHandler.GetByUsername)
			acct.GET("/code/:code", acctHandler.GetByCode)
			acct.GET("/third-party/:thirdPartyId", acctHandler.GetByThirdPartyID)
			acct.GET("/status/:status", acctHandler.GetByStatus)
			acct.GET("/list", acctHandler.List)
			acct.GET("/getMyAgents", acctHandler.GetMyAgents)
			acct.GET("/getUnGroupPageByGroupId", acctHandler.GetUnGroupPageByGroupId)
			acct.POST("", acctHandler.Create)
			acct.PUT("", acctHandler.Update)
			acct.DELETE("/:id", acctHandler.Delete)
			acct.PUT("/batch-status", acctHandler.BatchStatus)
		}

		// ── AccountGroupInfo ──
		agHandler := platform.NewAccountGroupInfoHandler(platformSvc)
		ag := plat.Group("/account-group-info")
		{
			ag.GET("/page", agHandler.Page)
			ag.GET("/:id", agHandler.GetByID)
			ag.GET("/group/:groupId", agHandler.GetByGroupID)
			ag.GET("/account/:accountId", agHandler.GetByAccountID)
			ag.POST("", agHandler.Create)
			ag.PUT("", agHandler.Update)
			ag.DELETE("/:id", agHandler.Delete)
		}

		// ── AccountTenantInfo ──
		atHandler := platform.NewAccountTenantInfoHandler(platformSvc)
		at := plat.Group("/account-tenant-info")
		{
			at.GET("/page", atHandler.Page)
			at.GET("/:id", atHandler.GetByID)
			at.GET("/account/:accountId", atHandler.GetByAccountID)
			at.POST("", atHandler.Create)
			at.PUT("", atHandler.Update)
			at.DELETE("/:id", atHandler.Delete)
		}

		// ── TenantInfo ──
		tenantHandler := platform.NewTenantInfoHandler(platformSvc)
		tenant := plat.Group("/tenant-info")
		{
			tenant.GET("/page", tenantHandler.Page)
			tenant.GET("/:id", tenantHandler.GetByID)
			tenant.POST("", tenantHandler.Create)
			tenant.PUT("", tenantHandler.Update)
			tenant.DELETE("/:id", tenantHandler.Delete)
		}

		// ── PlatformInfo (common) ──
		piHandler := common.NewPlatformInfoHandler(platformSvc)
		pi := plat.Group("/platform-info")
		{
			pi.GET("/page", piHandler.Page)
			pi.GET("/:id", piHandler.GetByID)
			pi.GET("/type/:type", piHandler.GetByType)
			pi.GET("/type/:type/enabled", piHandler.GetByTypeEnabled)
			pi.GET("/getEnabledPlatform", piHandler.GetEnabledPlatform)
			pi.POST("", piHandler.Create)
			pi.PUT("", piHandler.Update)
			pi.DELETE("/:id", piHandler.Delete)
		}

		// ── Platform Sync (common stubs) ──
		syncHandler := common.NewPlatformSyncHandler(platformSvc)
		sync := plat.Group("/sync")
		{
			sync.POST("/all", syncHandler.SyncAll)
			sync.POST("/departments", syncHandler.SyncDepartments)
			sync.POST("/users", syncHandler.SyncUsers)
			sync.POST("/depts/:deptId", syncHandler.SyncSubDepartments)
			sync.POST("/depts/users/:deptId", syncHandler.SyncUsersByDept)
			sync.POST("/users/:userId", syncHandler.SyncUser)
		}
	}

	// ================ RAG 路由 ================
	ragAPI := r.Group("/api/rag")
	ragAPI.Use(middleware.Auth(jwtManager))
	ragAPI.Use(middleware.RBAC(enforcer))
	if ragSvc != nil {
		// ---- RAG File ----
		ragFileHandler := ragHandler.NewRagFileInfoHandler(ragSvc)
		ragFileGroup := ragAPI.Group("/file")
		{
			ragFileGroup.GET("/page", ragFileHandler.Page)
			ragFileGroup.GET("/list", ragFileHandler.List)
			ragFileGroup.GET("/:id", ragFileHandler.GetByID)
			ragFileGroup.POST("", ragFileHandler.Create)
			ragFileGroup.PUT("", ragFileHandler.Update)
			ragFileGroup.DELETE("/:id", ragFileHandler.Delete)
		}

		// ---- RAG Category ----
		ragCategoryHandler := ragHandler.NewRagCategoryHandler(ragSvc)
		ragCategoryGroup := ragAPI.Group("/category")
		{
			ragCategoryGroup.GET("/page", ragCategoryHandler.Page)
			ragCategoryGroup.GET("/list", ragCategoryHandler.List)
			ragCategoryGroup.GET("/:id", ragCategoryHandler.GetByID)
			ragCategoryGroup.POST("", ragCategoryHandler.Create)
			ragCategoryGroup.PUT("", ragCategoryHandler.Update)
			ragCategoryGroup.DELETE("/:id", ragCategoryHandler.Delete)
		}
	}

	// ================ KG 路由 ================
	kgAPI := r.Group("/api/kg")
	kgAPI.Use(middleware.Auth(jwtManager))
	kgAPI.Use(middleware.RBAC(enforcer))
	if kgSvc != nil {
		// ---- KG Entity ----
		kgEntityHandler := kgHandler.NewKGEntityHandler(kgSvc)
		kgEntityGroup := kgAPI.Group("/entity")
		{
			kgEntityGroup.GET("/page", kgEntityHandler.Page)
			kgEntityGroup.GET("/list", kgEntityHandler.List)
			kgEntityGroup.GET("/:id", kgEntityHandler.GetByID)
			kgEntityGroup.POST("", kgEntityHandler.Create)
			kgEntityGroup.PUT("", kgEntityHandler.Update)
			kgEntityGroup.DELETE("/:id", kgEntityHandler.Delete)
		}

		// ---- KG Relation ----
		kgRelationHandler := kgHandler.NewKGRelationHandler(kgSvc)
		kgRelationGroup := kgAPI.Group("/relation")
		{
			kgRelationGroup.GET("/page", kgRelationHandler.Page)
			kgRelationGroup.GET("/:id", kgRelationHandler.GetByID)
			kgRelationGroup.GET("/by-entity/:entityId", kgRelationHandler.FindByEntity)
			kgRelationGroup.POST("", kgRelationHandler.Create)
			kgRelationGroup.PUT("", kgRelationHandler.Update)
			kgRelationGroup.DELETE("/:id", kgRelationHandler.Delete)
		}

		// ---- KG Domain ----
		kgDomainHandler := kgHandler.NewKGDomainHandler(kgSvc)
		kgDomainGroup := kgAPI.Group("/domain")
		{
			kgDomainGroup.GET("/page", kgDomainHandler.Page)
			kgDomainGroup.GET("/list", kgDomainHandler.List)
			kgDomainGroup.GET("/:id", kgDomainHandler.GetByID)
			kgDomainGroup.POST("", kgDomainHandler.Create)
			kgDomainGroup.PUT("", kgDomainHandler.Update)
			kgDomainGroup.DELETE("/:id", kgDomainHandler.Delete)
		}
	}

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.Response{
			Code:    "404",
			Message: "not found",
			Success: false,
		})
	})

	return r
}
