package cron

import (
	"time"
	"wifer/cron/dump"
	"wifer/server/structs"

	gocron "github.com/go-co-op/gocron/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Start(props *structs.Props) {
	location, _ := time.LoadLocation("Europe/Moscow")
	cron, _ := gocron.NewScheduler(gocron.WithLocation(location))
	defer func() { _ = cron.Shutdown() }()

	// Обновить всех юзеров у которых просрочился премиум
	_, _ = cron.NewJob(
		gocron.CronJob("0 2 * * *", false), // каждый день в 2 ночи
		gocron.NewTask(func() {
			props.DB["users"].UpdateMany(props.Ctx, bson.M{"premium": bson.M{"$lt": time.Now().Unix()}}, bson.D{{Key: "$set", Value: bson.D{{Key: "premium", Value: int64(0)}}}})
		}),
	)

	// Ежемесечно деактивировать всех пользователей, которых не было в сети от 1 года
	_, _ = cron.NewJob(
		gocron.CronJob("10 2 1 * *", false), // каждое 1 число месяца в 2:10 ночи
		gocron.NewTask(func() {
			props.DB["users"].UpdateMany(props.Ctx, bson.M{"last_time": bson.M{"$lt": time.Now().Unix() - int64(31536000)}}, bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: false}}}})
		}),
	)

	// Выгрузить базу данных в облако
	_, _ = cron.NewJob(
		gocron.CronJob("0 3 * * *", false), // каждый день в 3 ночи
		gocron.NewTask(func() { dump.PrepareDB(props) }),
	)

	// Выгрузить фотки в облако
	_, _ = cron.NewJob(
		gocron.CronJob("30 3 * * *", false), // каждый день в 3:30 ночи
		gocron.NewTask(func() { dump.Start(props, "/images", "images") }),
	)

	cron.Start()
}
