package staticpage

import (
	"ginweibo/controllers"
	"ginweibo/middleware/auth"
	viewmodels "ginweibo/middleware/viewmodels"
	followerModel "ginweibo/models/follower"
	statusModel "ginweibo/models/status"
	userModel "ginweibo/models/user"
	"ginweibo/routes/named"
	"ginweibo/utils/pagination"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	currentUser, err := auth.GetCurrentUserFromContext(c)
	if err != nil {
		controllers.Render(c, "home/home.html", gin.H{})
		return
	}
	// 获取用户所有关注的人，包括自己
	following, _ := followerModel.Followings(int(currentUser.ID), 0, 0)
	userIDmap := make(map[uint]*userModel.User, 0)
	userIDmap[currentUser.ID] = currentUser
	followingIDList := make([]uint, 0)
	followingIDList = append(followingIDList, currentUser.ID)
	for _, v := range following {
		followingIDList = append(followingIDList, v.ID)
		userIDmap[v.ID] = v
	}
	// 获取分页参数
	statusesAllLength, _ := statusModel.GetByUsersStatusesCount(followingIDList)
	offset, limit, currentPage, pageTotalCount := controllers.GetPageQuery(c, 10, statusesAllLength)
	if currentPage > pageTotalCount {
		controllers.Redirect(c, named.G("root")+"?page=1", false)
		return
	}
	// 获取用户的微博
	statuses, _ := statusModel.GetByUsersStatuses(followingIDList, offset, limit)
	statusesViewModels := make([]interface{}, 0)
	for _, s := range statuses {
		statusesViewModels = append(statusesViewModels, gin.H{
			"status": viewmodels.NewStatusViewModelSerializer(s),
			"user":   viewmodels.NewUserViewModelSerializer(userIDmap[s.UserID]),
		})
	}
	// 获取关注/粉丝
	followingsLength, _ := followerModel.FollowingsCount(int(currentUser.ID))
	followersLength, _ := followerModel.FollowersCount(int(currentUser.ID))
	controllers.Render(c, "home/home.html",
		pagination.CreatePaginationFillToTplData(c, "page", currentPage, pageTotalCount, gin.H{
			"statuses":         statusesViewModels,
			"statusesLength":   statusesAllLength,
			"followingsLength": followingsLength,
			"followersLength":  followersLength,
			"userData":         viewmodels.NewUserViewModelSerializer(currentUser),
		}))
}

func Help(c *gin.Context) {
	controllers.Render(c, "home/help.html", gin.H{})
}

func About(c *gin.Context) {
	controllers.Render(c, "home/about.html", gin.H{})
}
