package utils

import "github.com/robfig/cron"

type JobFunc func()
var CronJobs []*cron.Cron

/*
* * * * * ?		//	"秒 分 时 日 月 周几" 6个位

字段名				是否必须	允许的值			允许的特定字符
秒(Seconds)			是		0-59			* / , -
分(Minutes)			是		0-59			* / , -
时(Hours)			是		0-23			* / , -
日(Day of month)		是		1-31			* / , – ?
月(Month)			是		1-12 or JAN-DEC	* / , -
星期(Day of week)	否		0-6 or SUM-SAT	* / , – ?

1）星号(*)
表示 cron 表达式能匹配该字段的所有值。如在第5个字段使用星号(month)，表示每个月

2）斜线(/)
表示增长间隔，如第1个字段(minutes) 值是 3-59/15，表示每小时的第3分钟开始执行一次，之后每隔 15 分钟执行一次（即 3、18、33、48 这些时间点执行），这里也可以表示为：3/15

3）逗号(,)
用于枚举值，如第6个字段值是 MON,WED,FRI，表示 星期一、三、五 执行

4）连字号(-)
表示一个范围，如第3个字段的值为 9-17 表示 9am 到 5pm 直接每个小时（包括9和17）

5）问号(?)
只用于日(Day of month)和星期(Day of week)，\表示不指定值，可以用于代替 *
 */

// examples:
// * * * * * ? 			//每秒执行
// */5 * * * * ? 		//每5秒执行
// * */2 * * * ? 		//每2分钟执行
// * * */2 * * ? 		//每2小时执行
// * * * */2 * ? 		//每2天执行
// * * * * */2 ? 		//每2个月执行
// 00 26 09 * * ? 		//每天09:26:00执行
// 00 26 09 10 * ? 		//每月10日09:26:00执行
// 00 26 09 10 02 ? 	//每年02月10日09:26:00执行
// 1-10/2 * * * * * 	//1-10秒之间，没2秒执行一次
// 00 10 23 * * 1-5		//周一到周五，23:10:00执行
func CronJob(spec string, jobfunc JobFunc )  {
	c := cron.New()
	c.AddFunc(spec, jobfunc)
	c.Start()
	CronJobs = append(CronJobs, c)
}
