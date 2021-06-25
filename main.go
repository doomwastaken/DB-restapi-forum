package main

import (
	"context"
	"forum/application"
	"forum/infrastructure/persistence"
	"forum/interfaces/forum"
	"forum/interfaces/post"
	"forum/interfaces/service"
	"forum/interfaces/thread"
	"forum/interfaces/user"

	"fmt"
	"log"
	"os"
	"time"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func loggerMid(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		if end.Sub(begin) > 30*time.Millisecond {
			log.Printf("%s - %s",
				string(ctx.Request.URI().FullURI()),
				end.Sub(begin).String())
		}
	})
}

func runServer(addr string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load .env file", zap.String("error", err.Error()))
	}

	postgresConnectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	postgresConn, err := pgxpool.Connect(context.Background(), postgresConnectionString)
	if err != nil {
		log.Fatal("Could not connect to postgres database", zap.String("error", err.Error()))
		return
	}

	userRepo := persistence.NewUserRepository(postgresConn)
	forumRepo := persistence.NewForumRepository(postgresConn)
	postRepo := persistence.NewPostRepository(postgresConn)
	threadRepo := persistence.NewThreadRepository(postgresConn)
	serviceRepo := persistence.NewServiceRepository(postgresConn)

	serviceApp := application.NewServiceApp(serviceRepo)
	userApp := application.NewUserApp(userRepo)
	forumApp := application.NewForumApp(forumRepo)
	postApp := application.NewPostApp(postRepo)
	threadApp := application.NewThreadApp(threadRepo, forumApp)

	forumInfo := forum.NewForumInfo(forumApp, userApp, threadApp)
	userInfo := user.NewUserInfo(userApp)
	serviceInfo := service.NewServiceInfo(serviceApp)
	postsInfo := post.NewPostInfo(postApp, userApp, threadApp, forumApp)
	threadsInfo := thread.NewThreadInfo(threadApp, userApp)

	router := router.New()

	prefix := "/api"
	router.POST(prefix+"/user/{username}/create", userInfo.HandleCreateUser)
	router.GET(prefix+"/user/{username}/profile", userInfo.HandleGetUser)
	router.POST(prefix+"/user/{username}/profile", userInfo.HandleUpdateUser)

	router.POST(prefix+"/forum/create", forumInfo.HandleCreateForum)
	router.GET(prefix+"/forum/{forumname}/details", forumInfo.HandleGetForumDetails)
	router.GET(prefix+"/forum/{forumname}/users", forumInfo.HandleGetForumUsers)
	router.GET(prefix+"/forum/{forumname}/threads", forumInfo.HandleGetForumThreads)
	router.POST(prefix+"/forum/{forumname}/create", forumInfo.HandleCreateForumThread)

	router.GET(prefix+"/thread/{threadnameOrID}/details", threadsInfo.HandleGetThreadDetails)
	router.POST(prefix+"/thread/{threadnameOrID}/details", threadsInfo.HandleUpdateThread)
	router.GET(prefix+"/thread/{threadnameOrID}/posts", threadsInfo.HandleGetThreadPosts)
	router.POST(prefix+"/thread/{threadnameOrID}/vote", threadsInfo.HandleVoteForThread)
	router.POST(prefix+"/thread/{threadnameOrID}/create", threadsInfo.HandleCreateThread)

	router.GET(prefix+"/post/{postID}/details", postsInfo.HandleGetPostDetails)
	router.POST(prefix+"/post/{postID}/details", postsInfo.HandleChangePost)

	router.GET(prefix+"/service/status", serviceInfo.HandleGetDBStatus)
	router.POST(prefix+"/service/clear", serviceInfo.HandleClearData)

	fmt.Printf("Starting server at localhost%s\n", addr)
	fasthttp.ListenAndServe(addr, loggerMid(router.Handler))
}

func main() {
	runServer(":5000")
}