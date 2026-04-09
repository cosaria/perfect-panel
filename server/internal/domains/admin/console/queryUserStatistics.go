package console

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type QueryUserStatisticsOutput struct {
	Body *types.UserStatisticsResponse
}

func QueryUserStatisticsHandler(deps Deps) func(context.Context, *struct{}) (*QueryUserStatisticsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserStatisticsOutput, error) {
		l := NewQueryUserStatisticsLogic(ctx, deps)
		resp, err := l.QueryUserStatistics()
		if err != nil {
			return nil, err
		}
		return &QueryUserStatisticsOutput{Body: resp}, nil
	}
}

type QueryUserStatisticsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Query user statistics
func NewQueryUserStatisticsLogic(ctx context.Context, deps Deps) *QueryUserStatisticsLogic {
	return &QueryUserStatisticsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryUserStatisticsLogic) QueryUserStatistics() (resp *types.UserStatisticsResponse, err error) {
	if strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo" {
		return l.mockRevenueStatistics(), nil
	}
	resp = &types.UserStatisticsResponse{}
	now := time.Now()
	// query today user register count
	todayUserResisterCount, err := l.deps.UserModel.QueryResisterUserTotalByDate(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryResisterUserTotalByDate error", logger.Field("error", err.Error()))
	} else {
		resp.Today.Register = todayUserResisterCount
	}
	// query today user purchase count
	newToday, renewalToday, err := l.deps.OrderModel.QueryDateUserCounts(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryDateUserCounts error", logger.Field("error", err.Error()))
	} else {
		resp.Today.NewOrderUsers = newToday
		resp.Today.RenewalOrderUsers = renewalToday
	}
	// query month user register count
	monthUserResisterCount, err := l.deps.UserModel.QueryResisterUserTotalByMonthly(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryResisterUserTotalByMonthly error", logger.Field("error", err.Error()))
	} else {
		resp.Monthly.Register = monthUserResisterCount
	}
	// query month user purchase count
	newMonth, renewalMonth, err := l.deps.OrderModel.QueryMonthlyUserCounts(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryMonthlyUserCounts error", logger.Field("error", err.Error()))
	} else {
		resp.Monthly.NewOrderUsers = newMonth
		resp.Monthly.RenewalOrderUsers = renewalMonth
	}

	// Get monthly daily user statistics list for the current month (from 1st to current date)
	monthlyListData, err := l.deps.UserModel.QueryDailyUserStatisticsList(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryDailyUserStatisticsList error", logger.Field("error", err.Error()))
		// Don't return error, just log it and continue with empty list
	} else {
		monthlyList := make([]types.UserStatistics, len(monthlyListData))
		for i, data := range monthlyListData {
			monthlyList[i] = types.UserStatistics{
				Date:              data.Date,
				Register:          data.Register,
				NewOrderUsers:     data.NewOrderUsers,
				RenewalOrderUsers: data.RenewalOrderUsers,
			}
		}
		resp.Monthly.List = monthlyList
	}

	// query all user count
	allUserCount, err := l.deps.UserModel.QueryResisterUserTotal(l.ctx)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryResisterUserTotal error", logger.Field("error", err.Error()))
	} else {
		resp.All.Register = allUserCount
	}

	// query all user order counts
	allNewOrderUsers, allRenewalOrderUsers, err := l.deps.OrderModel.QueryTotalUserCounts(l.ctx)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryTotalUserCounts error", logger.Field("error", err.Error()))
	} else {
		resp.All.NewOrderUsers = allNewOrderUsers
		resp.All.RenewalOrderUsers = allRenewalOrderUsers
	}

	// Get all monthly user statistics list for the past 6 months
	allListData, err := l.deps.UserModel.QueryMonthlyUserStatisticsList(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryUserStatisticsLogic] QueryMonthlyUserStatisticsList error", logger.Field("error", err.Error()))
		// Don't return error, just log it and continue with empty list
	} else {
		allList := make([]types.UserStatistics, len(allListData))
		for i, data := range allListData {
			allList[i] = types.UserStatistics{
				Date:              data.Date,
				Register:          data.Register,
				NewOrderUsers:     data.NewOrderUsers,
				RenewalOrderUsers: data.RenewalOrderUsers,
			}
		}
		resp.All.List = allList
	}

	return
}

func (l *QueryUserStatisticsLogic) mockRevenueStatistics() *types.UserStatisticsResponse {
	now := time.Now()

	// Generate daily user statistics for the current month (from 1st to current date)
	monthlyList := make([]types.UserStatistics, 7)
	for i := 0; i < 7; i++ {
		dayDate := now.AddDate(0, 0, -(6 - i))
		baseRegister := int64(18 + ((6 - i) * 3) + ((6-i)%3)*8)
		monthlyList[i] = types.UserStatistics{
			Date:              dayDate.Format("2006-01-02"),
			Register:          baseRegister,
			NewOrderUsers:     int64(float64(baseRegister) * 0.65),
			RenewalOrderUsers: int64(float64(baseRegister) * 0.35),
		}
	}

	// Generate monthly user statistics for the past 6 months (oldest first)
	allList := make([]types.UserStatistics, 6)
	for i := 0; i < 6; i++ {
		monthDate := now.AddDate(0, -(5 - i), 0)
		baseRegister := int64(1800 + ((5 - i) * 200) + ((5-i)%2)*500)
		allList[i] = types.UserStatistics{
			Date:              monthDate.Format("2006-01"),
			Register:          baseRegister,
			NewOrderUsers:     int64(float64(baseRegister) * 0.65),
			RenewalOrderUsers: int64(float64(baseRegister) * 0.35),
		}
	}

	return &types.UserStatisticsResponse{
		Today: types.UserStatistics{
			Register:          28,
			NewOrderUsers:     18,
			RenewalOrderUsers: 10,
		},
		Monthly: types.UserStatistics{
			Register:          888,
			NewOrderUsers:     588,
			RenewalOrderUsers: 300,
			List:              monthlyList,
		},
		All: types.UserStatistics{
			Register:          18888,
			NewOrderUsers:     0, // This field is not used in All statistics
			RenewalOrderUsers: 0, // This field is not used in All statistics
			List:              allList,
		},
	}
}
