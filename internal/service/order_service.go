package service

import (
	"errors"
	"fmt"
	"math/rand"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/pkg/jwt"
	"photoset/internal/repository"
	"time"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	memberRepo  *repository.MembershipRepository
	photosetRepo *repository.PhotoSetRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	memberRepo *repository.MembershipRepository,
	photosetRepo *repository.PhotoSetRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		memberRepo:  memberRepo,
		photosetRepo: photosetRepo,
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID uint, orderType string, membershipID *uint, photosetID *uint) (*domain.Order, error) {
	order := &domain.Order{
		UserID: userID,
		Type:   orderType,
		Status: "pending",
	}

	switch orderType {
	case "membership":
		if membershipID == nil {
			return nil, errors.New("会员套餐ID不能为空")
		}
		plan, err := s.memberRepo.FindByID(*membershipID)
		if err != nil {
			return nil, errors.New("会员套餐不存在")
		}
		if plan.Status != 1 {
			return nil, errors.New("该套餐已下架")
		}
		order.MembershipID = membershipID
		order.Amount = plan.Price

	case "single":
		if photosetID == nil {
			return nil, errors.New("套图ID不能为空")
		}
		photoset, err := s.photosetRepo.FindByID(*photosetID)
		if err != nil {
			return nil, errors.New("套图不存在")
		}
		if photoset.IsFree == 1 {
			return nil, errors.New("该套图是免费的，无需购买")
		}
		order.PhotoSetID = photosetID
		order.Amount = photoset.Price

	default:
		return nil, errors.New("无效的订单类型")
	}

	// 生成订单号: PS + 时间戳 + 6位随机数
	order.OrderNo = fmt.Sprintf("PS%s%06d", time.Now().Format("20060102150405"), rand.Intn(1000000))
	order.ExpireSeconds = 1800

	if err := s.orderRepo.Create(order); err != nil {
		return nil, errors.New("创建订单失败")
	}

	return order, nil
}

// MockPay 模拟支付
func (s *OrderService) MockPay(userID, orderID uint) (string, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return "", errors.New("订单不存在")
	}

	// 校验：订单属于当前用户
	if order.UserID != userID {
		return "", errors.New("无权操作此订单")
	}

	// 校验：订单状态为待支付
	if order.Status != "pending" {
		return "", errors.New("订单状态异常，无法支付")
	}

	// 更新订单状态
	now := time.Now()
	order.Status = "paid"
	order.PaidAt = &now
	if err := s.orderRepo.Update(order); err != nil {
		return "", errors.New("支付失败")
	}

	// 如果是会员订单，更新用户会员状态
	if order.Type == "membership" && order.Membership != nil {
		db := database.GetMySQL()
		var user domain.User
		if err := db.First(&user, userID).Error; err != nil {
			return "", errors.New("用户不存在")
		}

		durationDays := order.Membership.Duration
		newExpiry := now.Add(time.Duration(durationDays) * 24 * time.Hour)

		// 如果当前会员未过期，在原过期时间上累加
		if user.MembershipExpires != nil && user.MembershipExpires.After(now) {
			newExpiry = user.MembershipExpires.Add(time.Duration(durationDays) * 24 * time.Hour)
		}

		user.MembershipExpires = &newExpiry

		// admin/creator 支付后不降级角色
		if user.Role != domain.RoleAdmin && user.Role != domain.RoleCreator {
			user.Role = domain.RoleMember
		}

		if err := db.Save(&user).Error; err != nil {
			return "", errors.New("更新会员状态失败")
		}
	}

	// 重新生成 JWT Token（角色可能已变化）
	newRole := "user"
	if order.Type == "membership" {
		db := database.GetMySQL()
		var user domain.User
		if err := db.First(&user, userID).Error; err == nil {
			newRole = string(user.Role)
		}
	}

	token, err := jwt.GenerateToken(userID, newRole)
	if err != nil {
		return "", errors.New("生成Token失败")
	}

	return token, nil
}

// GetOrderList 获取用户订单列表
func (s *OrderService) GetOrderList(userID uint, page, pageSize int) ([]domain.Order, int64, error) {
	return s.orderRepo.ListByUserID(userID, page, pageSize)
}

const (
	RefundWindowHours = 48 // 用户自助退款时间窗口（小时）
)

// RefundOrder 用户自助退款
func (s *OrderService) RefundOrder(userID, orderID uint) error {
	order, err := s.orderRepo.FindPaidOrderByUser(userID, orderID)
	if err != nil {
		return errors.New("订单不存在或状态不允许退款")
	}

	// 检查是否在退款时间窗口内（48小时）
	if order.PaidAt != nil {
		deadline := order.PaidAt.Add(time.Duration(RefundWindowHours) * time.Hour)
		if time.Now().After(deadline) {
			return errors.New("已超过退款时限（48小时），请联系管理员")
		}
	}

	return s.processRefund(order)
}

// AdminRefundOrder 管理员退款（无时间限制）
func (s *OrderService) AdminRefundOrder(orderID uint) error {
	order, err := s.orderRepo.FindPaidOrder(orderID)
	if err != nil {
		return errors.New("订单不存在或状态不允许退款")
	}

	return s.processRefund(order)
}

// processRefund 执行退款核心逻辑
func (s *OrderService) processRefund(order *domain.Order) error {
	// 1. 更新订单状态为 refunded
	now := time.Now()
	order.Status = "refunded"
	if err := s.orderRepo.Update(order); err != nil {
		return errors.New("退款失败: 更新订单状态失败")
	}

	// 2. 如果是会员订单，处理会员状态
	if order.Type == "membership" && order.MembershipID != nil {
		db := database.GetMySQL()
		var user domain.User
		if err := db.First(&user, order.UserID).Error; err != nil {
			return nil // 用户不存在时只退款不处理会员
		}

		if order.Membership != nil {
			durationDays := order.Membership.Duration
			// 扣减会员时长
			if user.MembershipExpires != nil {
				newExpiry := user.MembershipExpires.Add(-time.Duration(durationDays) * 24 * time.Hour)
				if newExpiry.Before(now) {
					// 会员已过期
					user.MembershipExpires = nil
					if user.Role == domain.RoleMember {
						user.Role = domain.RoleUser
					}
				} else {
					user.MembershipExpires = &newExpiry
				}
			}
			db.Save(&user)
		}
	}
	// 3. 单套图购买退款：只需要将订单状态改为 refunded，
	//    CanViewFullPhotos 查询 status='paid' 的订单会自动失效

	return nil
}
