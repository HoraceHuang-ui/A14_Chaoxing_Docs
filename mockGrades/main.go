package main

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"
)

type Stu struct {
	Clazz string

	Answer   int
	Checkin  int
	Discuss  int
	Homework int
	Mission  int
	Exam     int
}

func (s Stu) getScore() int {
	res := 0.0
	res += float64(s.Answer*20+80) * 0.05
	id := clazzId[s.Clazz]
	res += float64(s.Checkin) / float64(maxCheck[id]) * 100 * 0.2
	res += float64(74+2*s.Discuss) * 0.05
	res += float64(s.Homework) / float64(maxHomework[id]) * 100 * 0.25
	res += float64(s.Mission) / float64(maxMission[id]) * 100 * 0.15
	res += float64(s.Exam) / float64(maxExam[id]) * 100 * 0.3
	return int(res) + 1
}

var maxCheck []int
var maxHomework []int
var maxMission []int
var maxExam []int
var clazzId map[string]int

func main() {
	maxCheck = make([]int, 3)
	maxHomework = make([]int, 3)
	maxMission = make([]int, 3)
	maxExam = make([]int, 3)

	f_clazzAnswer, _ := excelize.OpenFile("excels/clazzAnswer.xlsx")
	f_clazzCheckin, _ := excelize.OpenFile("excels/clazzCheckin.xlsx")
	f_homeworkOrigin, _ := excelize.OpenFile("excels/homeworkOrigin.xlsx")
	f_discuss, _ := excelize.OpenFile("excels/discuss.xlsx")
	f_missionPoint, _ := excelize.OpenFile("excels/missionPoint.xlsx")
	f_exams, _ := excelize.OpenFile("excels/exam.xlsx")

	maps := []map[string]Stu{{}, {}, {}}
	clazzId = map[string]int{
		"50655126": 0,
		"50655128": 1,
		"50655129": 2,
	}

	clazzAnswer, _ := f_clazzAnswer.GetRows("Sheet1")
	clazzCheckin, _ := f_clazzCheckin.GetRows("Sheet1")
	homeworkOrigin, _ := f_homeworkOrigin.GetRows("Sheet1")
	discuss, _ := f_discuss.GetRows("Sheet1")
	missionPoint, _ := f_missionPoint.GetRows("Sheet1")
	exams, _ := f_exams.GetRows("Sheet1")

	for i := 1; i < len(clazzAnswer); i++ {
		var temp []string
		clazz := clazzAnswer[i][2]
		name := clazzAnswer[i][1]

		ans := atoi(clazzAnswer[i][3])
		temp = clazzCheckin[i]
		check := atoi(temp[3]) + atoi(temp[5]) - 2*atoi(temp[4])
		disc := atoi(discuss[i][3])
		homework := atoi(homeworkOrigin[i][3])
		mission := atoi(missionPoint[i][3])
		exam := atoi(exams[i][2])

		id := clazzId[clazz]
		maps[id][name] = Stu{clazz, ans, check, disc, homework, mission, exam}
		if check > maxCheck[id] {
			maxCheck[id] = check
		}
		if homework > maxHomework[id] {
			maxHomework[id] = homework
		}
		if mission > maxMission[id] {
			maxMission[id] = mission
		}
		if exam > maxExam[id] {
			maxExam[id] = exam
		}
	}

	var f *excelize.File

	// Output total score for each student
	f = excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]interface{}{
		"姓名", "班级", "总分", "抢答", "上课考勤", "讨论", "作业提交", "任务点", "考试提交",
	})
	i := 2
	for id := range []int{0, 1, 2} {
		for name, s := range maps[id] {
			f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i),
				&[]interface{}{name, s.Clazz, s.getScore(), s.Answer, s.Checkin, s.Discuss, s.Homework, s.Mission, s.Exam},
			)
			i++
		}
	}
	if err := f.SaveAs("totalScore.xlsx"); err != nil {
		fmt.Println(err)
	}
	f.Close()

	// Mock exam grades
	examPaths := []string{"excels/exam/exam_126.xlsx", "excels/exam/exam_128.xlsx", "excels/exam/exam_129.xlsx"}
	for id, path := range examPaths {
		f, _ = excelize.OpenFile(path)
		rows, _ := f.GetRows("Sheet1")
		i := 2
		for _, row := range rows[1:] {
			name := row[0]
			score := maps[id][name].getScore()
			col := 'E'
			for _, cell := range row[4:] {
				if cell != "0" {
					rge := 100 - score + 15
					mock := rand.Intn(min(30, rge)) + score - 15
					mock = max(mock, 0)
					mock = min(mock, 100)
					f.SetCellInt("Sheet1", fmt.Sprintf("%c%d", col, i), mock)
				}
				col++
			}
			i++
		}
		if err := f.Save(); err != nil {
			fmt.Print(err)
		}
		f.Close()
	}

	// Mock homework grades
	hwPaths := []string{"excels/homework/homework_126.xlsx", "excels/homework/homework_128.xlsx", "excels/homework/homework_129.xlsx"}
	for id, path := range hwPaths {
		f, _ = excelize.OpenFile(path)
		rows, _ := f.GetRows("Sheet1")
		i := 2
		for _, row := range rows[1:] {
			name := row[1]
			score := maps[id][name].getScore()
			col := 'F'
			zeros := (maxHomework[id] - maps[id][name].Homework) * 19 / maxHomework[id]
			zeroCnt := 0
			var mock int
			for range rows[0][5:] {
				if zeroCnt < zeros {
					if a := rand.Intn(2); a < 1 {
						mock = 0
						zeroCnt++
					} else {
						rge := 100 - score + 15
						mock = rand.Intn(min(50, rge)) + score - 10
						mock = max(mock, 0)
						mock = min(mock, 100)
					}
				} else {
					rge := 100 - score + 15
					mock = rand.Intn(min(50, rge)) + score - 10
					mock = max(mock, 0)
					mock = min(mock, 100)
				}
				f.SetCellInt("Sheet1", fmt.Sprintf("%c%d", col, i), mock)
				col++
			}
			i++
		}
		if err := f.Save(); err != nil {
			fmt.Print(err)
		}
		f.Close()
	}

	defer func() {
		// Close the spreadsheet.
		if err := f_clazzAnswer.Close(); err != nil {
			fmt.Println(err)
		}
		if err := f_clazzCheckin.Close(); err != nil {
			fmt.Println(err)
		}
		if err := f_homeworkOrigin.Close(); err != nil {
			fmt.Println(err)
		}
		if err := f_discuss.Close(); err != nil {
			fmt.Println(err)
		}
		if err := f_missionPoint.Close(); err != nil {
			fmt.Println(err)
		}
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
}

func atoi(s string) int {
	res := 0
	for _, c := range s {
		res = res*10 + int(c) - int('0')
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
