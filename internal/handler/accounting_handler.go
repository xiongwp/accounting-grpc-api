package handler

import (
	"context"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	accountingv1 "github.com/xiongwp/accounting-grpc-api/gen/accounting/v1"
	"github.com/xiongwp/accounting-system/internal/domain/model"
	"github.com/xiongwp/accounting-system/internal/service"
)

// AccountingHandler gRPC处理器
type AccountingHandler struct {
	accountingv1.UnimplementedAccountingServiceServer
	accountingService service.AccountingService
	facadeService     service.AccountingFacadeService
	flowEngine        service.MoneyFlowEngine
	dayCutService     service.DayCutService
	adjustmentService service.AdjustmentService
	logger            *zap.Logger
}

// NewAccountingHandler 创建处理器
func NewAccountingHandler(
	accountingService service.AccountingService,
	facadeService service.AccountingFacadeService,
	flowEngine service.MoneyFlowEngine,
	dayCutService service.DayCutService,
	adjustmentService service.AdjustmentService,
	logger *zap.Logger,
) *AccountingHandler {
	return &AccountingHandler{
		accountingService: accountingService,
		facadeService:     facadeService,
		flowEngine:        flowEngine,
		dayCutService:     dayCutService,
		adjustmentService: adjustmentService,
		logger:            logger,
	}
}

// CreateAccount 创建账户
func (h *AccountingHandler) CreateAccount(ctx context.Context, req *accountingv1.CreateAccountRequest) (*accountingv1.CreateAccountResponse, error) {
	h.logger.Info("CreateAccount called", zap.Int64("userId", req.UserId))

	account, err := h.accountingService.CreateAccount(
		ctx,
		req.UserId,
		convertAccountType(req.AccountType),
		convertAccountCategory(req.Category),
		req.Currency,
	)

	if err != nil {
		h.logger.Error("CreateAccount failed", zap.Error(err))
		return &accountingv1.CreateAccountResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &accountingv1.CreateAccountResponse{
		Code:    200,
		Message: "success",
		Account: convertAccountToProto(account),
	}, nil
}

// GetAccount 查询账户
func (h *AccountingHandler) GetAccount(ctx context.Context, req *accountingv1.GetAccountRequest) (*accountingv1.GetAccountResponse, error) {
	var account *model.Account
	var err error

	switch identifier := req.Identifier.(type) {
	case *accountingv1.GetAccountRequest_AccountNo:
		account, err = h.accountingService.GetAccount(ctx, identifier.AccountNo)
	case *accountingv1.GetAccountRequest_UserId:
		// 需要在service层实现GetAccountByUserID
		return &accountingv1.GetAccountResponse{
			Code:    400,
			Message: "query by user_id not implemented yet",
		}, nil
	default:
		return &accountingv1.GetAccountResponse{
			Code:    400,
			Message: "invalid identifier",
		}, nil
	}

	if err != nil {
		h.logger.Error("GetAccount failed", zap.Error(err))
		return &accountingv1.GetAccountResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	if account == nil {
		return &accountingv1.GetAccountResponse{
			Code:    404,
			Message: "account not found",
		}, nil
	}

	return &accountingv1.GetAccountResponse{
		Code:    200,
		Message: "success",
		Account: convertAccountToProto(account),
	}, nil
}

// FreezeAccount 冻结账户
func (h *AccountingHandler) FreezeAccount(ctx context.Context, req *accountingv1.FreezeAccountRequest) (*accountingv1.FreezeAccountResponse, error) {
	// TODO: 实现冻结逻辑
	return &accountingv1.FreezeAccountResponse{
		Code:    200,
		Message: "success",
	}, nil
}

// UnfreezeAccount 解冻账户
func (h *AccountingHandler) UnfreezeAccount(ctx context.Context, req *accountingv1.UnfreezeAccountRequest) (*accountingv1.UnfreezeAccountResponse, error) {
	// TODO: 实现解冻逻辑
	return &accountingv1.UnfreezeAccountResponse{
		Code:    200,
		Message: "success",
	}, nil
}

// DoubleEntryBooking 复式记账
func (h *AccountingHandler) DoubleEntryBooking(ctx context.Context, req *accountingv1.DoubleEntryBookingRequest) (*accountingv1.DoubleEntryBookingResponse, error) {
	h.logger.Info("DoubleEntryBooking called", zap.String("businessNo", req.BusinessNo))

	// 转换分录
	entries := make([]service.AccountingEntry, len(req.Entries))
	for i, entry := range req.Entries {
		debit, _ := decimal.NewFromString(entry.DebitAmount)
		credit, _ := decimal.NewFromString(entry.CreditAmount)

		entries[i] = service.AccountingEntry{
			AccountNo:    entry.AccountNo,
			DebitAmount:  debit,
			CreditAmount: credit,
			Description:  entry.Description,
		}
	}

	// 根据执行模式选择处理方式
	mode := convertExecutionMode(req.Mode)
	bookingReq := &service.BookingRequest{
		Mode:       mode,
		BusinessNo: req.BusinessNo,
		DoubleEntryReq: &service.DoubleEntryBookingRequest{
			BusinessNo:   req.BusinessNo,
			BusinessType: convertBusinessType(req.BusinessType),
			Entries:      entries,
			Currency:     req.Currency,
			Description:  req.Description,
		},
	}

	response, err := h.facadeService.ProcessBooking(ctx, bookingReq)
	if err != nil {
		h.logger.Error("DoubleEntryBooking failed", zap.Error(err))
		return &accountingv1.DoubleEntryBookingResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &accountingv1.DoubleEntryBookingResponse{
		Code:           200,
		Message:        "success",
		VoucherNo:      response.VoucherNo,
		TransactionIds: response.TransactionIDs,
		Async:          response.Async,
	}, nil
}

// BatchBooking 批量记账
func (h *AccountingHandler) BatchBooking(ctx context.Context, req *accountingv1.BatchBookingRequest) (*accountingv1.BatchBookingResponse, error) {
	h.logger.Info("BatchBooking called", zap.Int("count", len(req.Requests)))

	// 转换请求
	requests := make([]service.BookingRequest, len(req.Requests))
	for i, r := range req.Requests {
		entries := make([]service.AccountingEntry, len(r.Entries))
		for j, entry := range r.Entries {
			debit, _ := decimal.NewFromString(entry.DebitAmount)
			credit, _ := decimal.NewFromString(entry.CreditAmount)

			entries[j] = service.AccountingEntry{
				AccountNo:    entry.AccountNo,
				DebitAmount:  debit,
				CreditAmount: credit,
				Description:  entry.Description,
			}
		}

		requests[i] = service.BookingRequest{
			Mode:       service.ExecutionModeSync,
			BusinessNo: r.BusinessNo,
			DoubleEntryReq: &service.DoubleEntryBookingRequest{
				BusinessNo:   r.BusinessNo,
				BusinessType: convertBusinessType(r.BusinessType),
				Entries:      entries,
				Currency:     r.Currency,
				Description:  r.Description,
			},
		}
	}

	batchReq := &service.BatchBookingRequest{
		Requests: requests,
		Parallel: req.Parallel,
	}

	response, err := h.facadeService.ProcessBatchBooking(ctx, batchReq)
	if err != nil {
		h.logger.Error("BatchBooking failed", zap.Error(err))
		return &accountingv1.BatchBookingResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 转换响应
	results := make([]*accountingv1.DoubleEntryBookingResponse, len(response.Results))
	for i, r := range response.Results {
		results[i] = &accountingv1.DoubleEntryBookingResponse{
			Code:           getResponseCode(r.Success),
			Message:        r.ErrorMessage,
			VoucherNo:      r.VoucherNo,
			TransactionIds: r.TransactionIDs,
			Async:          r.Async,
		}
	}

	return &accountingv1.BatchBookingResponse{
		Code:    200,
		Message: "success",
		Success: int32(response.Success),
		Failed:  int32(response.Failed),
		Total:   int32(response.Total),
		Results: results,
	}, nil
}

// MoneyFlow 资金流执行
func (h *AccountingHandler) MoneyFlow(ctx context.Context, req *accountingv1.MoneyFlowRequest) (*accountingv1.MoneyFlowResponse, error) {
	h.logger.Info("MoneyFlow called",
		zap.String("productCode", req.ProductCode),
		zap.String("sceneCode", req.SceneCode),
	)

	amount, _ := decimal.NewFromString(req.Amount)

	flowReq := &service.FlowExecutionRequest{
		ProductCode:  req.ProductCode,
		SceneCode:    req.SceneCode,
		BusinessNo:   req.BusinessNo,
		Amount:       amount,
		Currency:     req.Currency,
		Participants: req.Participants,
		ExtParams:    convertExtParams(req.ExtParams),
	}

	response, err := h.flowEngine.ExecuteFlow(ctx, flowReq)
	if err != nil {
		h.logger.Error("MoneyFlow failed", zap.Error(err))
		return &accountingv1.MoneyFlowResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &accountingv1.MoneyFlowResponse{
		Code:           200,
		Message:        "success",
		VoucherNo:      response.VoucherNo,
		TransactionIds: response.TransactionIDs,
	}, nil
}

// GetTransaction 查询流水
func (h *AccountingHandler) GetTransaction(ctx context.Context, req *accountingv1.GetTransactionRequest) (*accountingv1.GetTransactionResponse, error) {
	// TODO: 实现查询流水逻辑
	return &accountingv1.GetTransactionResponse{
		Code:    200,
		Message: "success",
	}, nil
}

// GetBalanceSnapshot 查询余额快照
func (h *AccountingHandler) GetBalanceSnapshot(ctx context.Context, req *accountingv1.GetBalanceSnapshotRequest) (*accountingv1.GetBalanceSnapshotResponse, error) {
	// TODO: 实现查询快照逻辑
	return &accountingv1.GetBalanceSnapshotResponse{
		Code:    200,
		Message: "success",
	}, nil
}

// TriggerDayCut 触发日切
func (h *AccountingHandler) TriggerDayCut(ctx context.Context, req *accountingv1.TriggerDayCutRequest) (*accountingv1.TriggerDayCutResponse, error) {
	h.logger.Info("TriggerDayCut called", zap.String("cutDate", req.CutDate))

	err := h.dayCutService.TriggerDayCut(ctx, req.CutDate)
	if err != nil {
		h.logger.Error("TriggerDayCut failed", zap.Error(err))
		return &accountingv1.TriggerDayCutResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &accountingv1.TriggerDayCutResponse{
		Code:    200,
		Message: "success",
	}, nil
}

// AdjustBalance 调账
func (h *AccountingHandler) AdjustBalance(ctx context.Context, req *accountingv1.AdjustBalanceRequest) (*accountingv1.AdjustBalanceResponse, error) {
	h.logger.Info("AdjustBalance called", zap.String("accountNo", req.AccountNo))

	amount, _ := decimal.NewFromString(req.Amount)

	adjustReq := &service.AdjustmentRequest{
		AccountNo:      req.AccountNo,
		AdjustmentType: service.AdjustmentType(req.AdjustmentType),
		Amount:         amount,
		IsIncrease:     req.IsIncrease,
		Reason:         req.Reason,
		Operator:       req.Operator,
		ApprovalNo:     req.ApprovalNo,
	}

	response, err := h.adjustmentService.AdjustBalance(ctx, adjustReq)
	if err != nil {
		h.logger.Error("AdjustBalance failed", zap.Error(err))
		return &accountingv1.AdjustBalanceResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &accountingv1.AdjustBalanceResponse{
		Code:          200,
		Message:       "success",
		TransactionId: response.TransactionID,
		VoucherNo:     response.VoucherNo,
		BalanceBefore: response.BalanceBefore.String(),
		BalanceAfter:  response.BalanceAfter.String(),
	}, nil
}

// ==================== 辅助函数 ====================

func convertAccountToProto(account *model.Account) *accountingv1.Account {
	return &accountingv1.Account{
		AccountNo:        account.AccountNo,
		UserId:           account.UserID,
		AccountType:      convertAccountTypeToProto(account.AccountType),
		Category:         convertAccountCategoryToProto(account.AccountCategory),
		Currency:         account.Currency,
		Balance:          account.Balance.String(),
		FrozenBalance:    account.FrozenBalance.String(),
		AvailableBalance: account.AvailableBalance.String(),
		Status:           convertAccountStatusToProto(account.Status),
		Version:          account.Version,
		CreatedAt:        timestamppb.New(account.CreatedAt),
		UpdatedAt:        timestamppb.New(account.UpdatedAt),
	}
}

func convertAccountType(t accountingv1.AccountType) model.AccountType {
	switch t {
	case accountingv1.AccountType_ACCOUNT_TYPE_USER:
		return model.AccountTypeUser
	case accountingv1.AccountType_ACCOUNT_TYPE_MERCHANT:
		return model.AccountTypeMerchant
	case accountingv1.AccountType_ACCOUNT_TYPE_PLATFORM:
		return model.AccountTypePlatform
	case accountingv1.AccountType_ACCOUNT_TYPE_TRANSIT:
		return model.AccountTypeTransit
	default:
		return model.AccountTypeUser
	}
}

func convertAccountTypeToProto(t model.AccountType) accountingv1.AccountType {
	switch t {
	case model.AccountTypeUser:
		return accountingv1.AccountType_ACCOUNT_TYPE_USER
	case model.AccountTypeMerchant:
		return accountingv1.AccountType_ACCOUNT_TYPE_MERCHANT
	case model.AccountTypePlatform:
		return accountingv1.AccountType_ACCOUNT_TYPE_PLATFORM
	case model.AccountTypeTransit:
		return accountingv1.AccountType_ACCOUNT_TYPE_TRANSIT
	default:
		return accountingv1.AccountType_ACCOUNT_TYPE_UNSPECIFIED
	}
}

func convertAccountCategory(c accountingv1.AccountCategory) model.AccountCategory {
	switch c {
	case accountingv1.AccountCategory_ACCOUNT_CATEGORY_ASSET:
		return model.AccountCategoryAsset
	case accountingv1.AccountCategory_ACCOUNT_CATEGORY_LIABILITY:
		return model.AccountCategoryLiability
	case accountingv1.AccountCategory_ACCOUNT_CATEGORY_EQUITY:
		return model.AccountCategoryEquity
	case accountingv1.AccountCategory_ACCOUNT_CATEGORY_REVENUE:
		return model.AccountCategoryRevenue
	case accountingv1.AccountCategory_ACCOUNT_CATEGORY_EXPENSE:
		return model.AccountCategoryExpense
	default:
		return model.AccountCategoryAsset
	}
}

func convertAccountCategoryToProto(c model.AccountCategory) accountingv1.AccountCategory {
	switch c {
	case model.AccountCategoryAsset:
		return accountingv1.AccountCategory_ACCOUNT_CATEGORY_ASSET
	case model.AccountCategoryLiability:
		return accountingv1.AccountCategory_ACCOUNT_CATEGORY_LIABILITY
	case model.AccountCategoryEquity:
		return accountingv1.AccountCategory_ACCOUNT_CATEGORY_EQUITY
	case model.AccountCategoryRevenue:
		return accountingv1.AccountCategory_ACCOUNT_CATEGORY_REVENUE
	case model.AccountCategoryExpense:
		return accountingv1.AccountCategory_ACCOUNT_CATEGORY_EXPENSE
	default:
		return accountingv1.AccountCategory_ACCOUNT_CATEGORY_UNSPECIFIED
	}
}

func convertAccountStatusToProto(s model.AccountStatus) accountingv1.AccountStatus {
	switch s {
	case model.AccountStatusDisabled:
		return accountingv1.AccountStatus_ACCOUNT_STATUS_DISABLED
	case model.AccountStatusActive:
		return accountingv1.AccountStatus_ACCOUNT_STATUS_ACTIVE
	case model.AccountStatusFrozen:
		return accountingv1.AccountStatus_ACCOUNT_STATUS_FROZEN
	default:
		return accountingv1.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
}

func convertBusinessType(t accountingv1.BusinessType) model.BusinessType {
	switch t {
	case accountingv1.BusinessType_BUSINESS_TYPE_TRANSFER:
		return model.BusinessTypeTransfer
	case accountingv1.BusinessType_BUSINESS_TYPE_PAYMENT:
		return model.BusinessTypePayment
	case accountingv1.BusinessType_BUSINESS_TYPE_REFUND:
		return model.BusinessTypeRefund
	case accountingv1.BusinessType_BUSINESS_TYPE_WITHDRAW:
		return model.BusinessTypeWithdraw
	case accountingv1.BusinessType_BUSINESS_TYPE_DEPOSIT:
		return model.BusinessTypeDeposit
	case accountingv1.BusinessType_BUSINESS_TYPE_COMMISSION:
		return model.BusinessTypeCommission
	default:
		return model.BusinessTypeTransfer
	}
}

func convertExecutionMode(m accountingv1.ExecutionMode) service.ExecutionMode {
	switch m {
	case accountingv1.ExecutionMode_EXECUTION_MODE_SYNC:
		return service.ExecutionModeSync
	case accountingv1.ExecutionMode_EXECUTION_MODE_ASYNC:
		return service.ExecutionModeAsync
	case accountingv1.ExecutionMode_EXECUTION_MODE_BATCH:
		return service.ExecutionModeBatch
	default:
		return service.ExecutionModeSync
	}
}

func convertExtParams(params map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range params {
		result[k] = v
	}
	return result
}

func getResponseCode(success bool) int32 {
	if success {
		return 200
	}
	return 500
}
