package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	db, _ := InitDB("database.db")
	user, userErr := ReturnUser(db, "1")
	if userErr != nil {
		var uid, upw string
		fmt.Print("Id: ")
		_, _ = fmt.Scanln(&uid)
		fmt.Print("Password: ")
		_, _ = fmt.Scanln(&upw)
		user.SetInfo(uid, upw)
	}

	lectureNameList, lectureCodeList, lectureErr := ReturnLecture(db)
	if userErr == nil {
		fmt.Println("만약 정보가 다르다면 database.db를 삭제한 뒤 실행해주세요.")
		fmt.Println("학번: ", user.StudentCode)
	}
	if lectureErr == nil {
		fmt.Println("수강과목:", lectureNameList)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		// Set the headless flag to false to display the browser window
		//chromedp.Flag("headless", false),
		//chromedp.Flag("start-fullscreen", true),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	currentPath, _ := os.Getwd()
	err := chromedp.Run(ctx, page.SetDownloadBehavior("allow").WithDownloadPath(currentPath))
	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(ctx,
		Login(&user),
	)
	if err != nil {
		log.Fatal(err)
	}
	if userErr != nil {
		var stdCode string
		err = chromedp.Run(ctx,
			GetStudentCode(&stdCode),
		)

		user.SetStudentCode(stdCode)
		_ = AddUser(db, &user)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = chromedp.Run(ctx,
		GoLecturePage(),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	err = chromedp.Run(ctx,
		chromedp.ScrollIntoView("/html/body/div[1]/div[2]/bb-base-layout/div/main/div/div/div[1]/div[1]/div/div/div[3]/div/div[2]/div/div[5]/bb-base-course-card/div[1]/div[2]/a/h4", chromedp.BySearch),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if lectureErr != nil {
		var lectureName, lectureCode string

		ctx2, _ := context.WithTimeout(ctx, 1*time.Second)
		for i := 2; ; i++ {
			err = chromedp.Run(ctx2,
				GetLectureInfo(i, &lectureName, &lectureCode),
			)
			if err != nil {
				break
			}
			lectureNameList = append(lectureNameList, lectureName)
			lectureCodeList = append(lectureCodeList, lectureCode)

			AddLecture(db, lectureName, lectureCode)
		}
	}

	RemoveFile("xls")
	RemoveFile("xlsx")

	for i := 0; i < len(lectureNameList); i++ {
		_ = chromedp.Run(ctx,
			DownloadAttendance(&user.StudentCode, &lectureNameList[i], &lectureCodeList[i]),
		)
	}
	time.Sleep(1 * time.Second)

	RunPy()

	RemoveFile("xls")
	RemoveFile("xlsx")
}

func Login(user *User) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://learn.hanyang.ac.kr/ultra/institution-page"),
		chromedp.Click("//*[@id=\"entry-login-custom\"]", chromedp.BySearch),
		chromedp.WaitVisible("uid", chromedp.ByID),
		chromedp.SetValue("uid", user.UserId, chromedp.ByID),
		chromedp.SetValue("upw", user.Password, chromedp.ByID),
		chromedp.Click("//*[@id=\"login_btn\"]", chromedp.BySearch),
		chromedp.WaitVisible("/html/body/div[1]/div[2]/bb-base-layout/div/aside/div[1]/header/span/a/img", chromedp.BySearch),
	}
}

func GetStudentCode(stdCode *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://learn.hanyang.ac.kr/ultra/profile"),
		chromedp.WaitVisible("//*[@id=\"body-content\"]/div[1]/div[3]/div[1]/section[1]/ul/li[3]/div/div", chromedp.BySearch),
		chromedp.Text("//*[@id=\"body-content\"]/div[1]/div[3]/div[1]/section[1]/ul/li[3]/div/div", stdCode, chromedp.BySearch),
	}
}

func GoLecturePage() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://learn.hanyang.ac.kr/ultra/course"),
		chromedp.WaitVisible("/html/body/div[1]/div[2]/bb-base-layout/div/main/div/div/div[1]/div[1]/div/div/div[3]/div/div[2]/div/div[2]/bb-base-course-card/div[1]/div[2]/div[1]/div/span"),
	}
}

func GetLectureInfo(idx int, name *string, code *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Text("/html/body/div[1]/div[2]/bb-base-layout/div/main/div/div/div[1]/div[1]/div/div/div[3]/div/div[2]/div/div["+strconv.Itoa(idx)+"]/bb-base-course-card/div[1]/div[2]/a/h4", name, chromedp.BySearch),
		chromedp.Text("/html/body/div[1]/div[2]/bb-base-layout/div/main/div/div/div[1]/div[1]/div/div/div[3]/div/div[2]/div/div["+strconv.Itoa(idx)+"]/bb-base-course-card/div[1]/div[2]/div[1]/div/span", code, chromedp.BySearch),
	}
}

func DownloadAttendance(studentCode, lectureName, lectureCode *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://learn.hanyang.ac.kr/webapps/bbgs-OnlineAttendance-BB5a998b8c44671/excel?selectedUserId=" + *studentCode + "&crs_batch_uid=" + *lectureCode + "&title=" + *lectureName + "&column=%EC%82%AC%EC%9A%A9%EC%9E%90%EB%AA%85,%EC%9C%84%EC%B9%98,%EC%BB%A8%ED%85%90%EC%B8%A0%EB%AA%85,%ED%95%99%EC%8A%B5%ED%95%9C%EC%8B%9C%EA%B0%84,%ED%95%99%EC%8A%B5%EC%9D%B8%EC%A0%95%EC%8B%9C%EA%B0%84,%EC%BB%A8%ED%85%90%EC%B8%A0%EC%8B%9C%EA%B0%84,%EC%98%A8%EB%9D%BC%EC%9D%B8%EC%B6%9C%EC%84%9D%EC%A7%84%EB%8F%84%EC%9C%A8,%EC%98%A8%EB%9D%BC%EC%9D%B8%EC%B6%9C%EC%84%9D%EC%83%81%ED%83%9C(P/F)"),
	}
}
