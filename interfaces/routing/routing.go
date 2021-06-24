package routing

import (
	"forum/application"
	"forum/infrastructure/persistence"
	"forum/interfaces/forum"
	"forum/interfaces/post"
	"forum/interfaces/service"
	"forum/interfaces/thread"
	"forum/interfaces/user"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateRouter(conn *pgxpool.Pool, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	repoUser := persistence.NewUserRepository(conn)
	repoForum := persistence.NewForumRepository(conn)
	repoPosts := persistence.NewPostRepository(conn)
	repoService := persistence.NewServiceRepository(conn)
	repoThreads := persistence.NewThreadRepository(conn)

	postsApp := application.NewPostApp(repoPosts)
	userApp := application.NewUserApp(repoUser)
	serviceApp := application.NewServiceApp(repoService)
	forumApp := application.NewForumApp(repoForum)
	threadsApp := application.NewThreadApp(repoThreads, forumApp)

	forumInfo := forum.NewForumInfo(forumApp, userApp, threadsApp, logger)
	userInfo := user.NewUserInfo(userApp, logger)
	serviceInfo := service.NewServiceInfo(serviceApp, logger)
	postsInfo := post.NewPostInfo(postsApp, userApp, threadsApp, forumApp, logger)
	threadsInfo := thread.NewThreadInfo(threadsApp, userApp, logger)

	r.HandleFunc("/api/forum/create", forumInfo.HandleCreateForum).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/create", forumInfo.HandleCreateForumThread).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", forumInfo.HandleGetForumDetails).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", forumInfo.HandleGetForumUsers).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", forumInfo.HandleGetForumThreads).Methods("GET")

	r.HandleFunc("/api/post/{id}/details", postsInfo.HandleChangePost).Methods("POST")
	r.HandleFunc("/api/post/{id}/details", postsInfo.HandleGetPostDetails).Methods("GET")

	r.HandleFunc("/api/service/clear", serviceInfo.HandleClearData).Methods("POST")
	r.HandleFunc("/api/service/status", serviceInfo.HandleGetDBStatus).Methods("GET")

	r.HandleFunc("/api/thread/{slug_or_id}/create", threadsInfo.HandleCreateThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", threadsInfo.HandleUpdateThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", threadsInfo.HandleVoteForThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", threadsInfo.HandleGetThreadDetails).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", threadsInfo.HandleGetThreadPosts).Methods("GET")

	r.HandleFunc("/api/user/{nickname}/create", userInfo.HandleCreateUser).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", userInfo.HandleUpdateUser).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", userInfo.HandleGetUser).Methods("GET")

	return r
}
